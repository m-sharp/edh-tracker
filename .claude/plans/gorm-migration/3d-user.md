# Phase 3d — User Repository

## Status
Approved

## Skill
Use the `/gorm` skill tool at the start of each implementation session for this phase.

## Scope

- Tables: `user`, `user_role`
- Files: `lib/repositories/user/model.go`, `lib/repositories/user/repo.go`
- No existing tests — write new integration tests
- Depends on: player (Phase 1a)
- Complex: transaction in `CreatePlayerAndUser`; nullable OAuth fields; two models in one package
- No business layer changes

## GORM Models

```go
package user

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

const (
    RoleAdmin  = "admin"
    RolePlayer = "player"
)

type Model struct {
    base.GormModelBase
    PlayerID      int
    RoleID        int
    OAuthProvider *string `gorm:"column:oauth_provider"` // KEEP: inferred as o_auth_provider, DB has oauth_provider
    OAuthSubject  *string `gorm:"column:oauth_subject"`  // KEEP: same reason
    Email         *string
    DisplayName   *string
    AvatarURL     *string
}

func (Model) TableName() string { return "user" }

type RoleModel struct {
    base.GormModelBase
    Name string
}

func (RoleModel) TableName() string { return "user_role" }
```

All nullable `*string` fields handled correctly by GORM (NULL ↔ nil).

## Method Mapping

| Old (sqlx) | New (GORM) | Notes |
|---|---|---|
| `GetByID` | `db.First(&m, id)` | `ErrRecordNotFound` → nil,nil |
| `GetByPlayerID` | `db.Where("player_id = ?", playerID).First(&m)` | `ErrRecordNotFound` → nil,nil |
| `GetByOAuth` | `db.Where("oauth_provider = ? AND oauth_subject = ?", ...).First(&m)` | `ErrRecordNotFound` → nil,nil |
| `GetRoleByName` | `db.Where("name = ?", name).First(&m)` | Error if not found (non-nil) |
| `Add` | `db.Create(&m)` | `m.ID` set by GORM |
| `AddWithOAuth` | `db.Create(&m)` with OAuth fields | `m.ID` set by GORM |
| `CreatePlayerAndUser` | `db.Transaction(func(tx) error {...})` | See below |
| `BulkAdd` | `db.CreateInBatches(&entries, 100)` | No return needed |
| `SoftDelete` | `db.Delete(&Model{}, id)` | Sets deleted_at |

## Repository Field

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
    return &Repository{db: client.GormDb}
}

func NewRepositoryFromDB(db *gorm.DB) *Repository {
    return &Repository{db: db}
}
```

## Special Pattern — GetRoleByName (error on not found)

Unlike other Get-by-name methods, `GetRoleByName` returns an error when not found (the role must exist). Use `First` and don't map `ErrRecordNotFound` to nil:

```go
func (r *Repository) GetRoleByName(ctx context.Context, name string) (*RoleModel, error) {
    var m RoleModel
    err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, fmt.Errorf("no UserRole found with name %q", name)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get UserRole record for name %q: %w", name, err)
    }
    return &m, nil
}
```

## Special Pattern — AddWithOAuth

Build a Model with all OAuth fields set:

```go
func (r *Repository) AddWithOAuth(
    ctx context.Context,
    playerID, roleID int,
    provider, subject, email, displayName, avatarURL string,
) (int, error) {
    m := Model{
        PlayerID:      playerID,
        RoleID:        roleID,
        OAuthProvider: &provider,
        OAuthSubject:  &subject,
        Email:         &email,
        DisplayName:   &displayName,
        AvatarURL:     &avatarURL,
    }
    if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
        return 0, fmt.Errorf("failed to insert User record with OAuth: %w", err)
    }
    return m.ID, nil
}
```

## Special Pattern — CreatePlayerAndUser (Transaction)

The current implementation uses a manual tx with `BeginTxx` / `Commit` / `Rollback`. GORM's `Transaction` helper simplifies this significantly:

```go
func (r *Repository) CreatePlayerAndUser(
    ctx context.Context,
    playerName string,
    roleID int,
    provider, subject, email, displayName, avatarURL string,
) (*Model, error) {
    var created Model

    err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        type playerRow struct {
            ID   int    `gorm:"primaryKey"`
            Name string
        }
        p := playerRow{Name: playerName}
        if err := tx.Table("player").Create(&p).Error; err != nil {
            return fmt.Errorf("failed to insert player in CreatePlayerAndUser: %w", err)
        }

        created = Model{
            PlayerID:      p.ID,
            RoleID:        roleID,
            OAuthProvider: &provider,
            OAuthSubject:  &subject,
            Email:         &email,
            DisplayName:   &displayName,
            AvatarURL:     &avatarURL,
        }
        if err := tx.Create(&created).Error; err != nil {
            return fmt.Errorf("failed to insert user in CreatePlayerAndUser: %w", err)
        }
        return nil
    })

    if err != nil {
        return nil, err
    }
    return &created, nil
}
```

The `playerRow` local type uses `tx.Table("player")` to target the correct table without a `TableName()` method. GORM sets `p.ID` from `LAST_INSERT_ID()` after insert. The `player` table's `created_at`/`updated_at` columns are absent from `playerRow`, so GORM won't try to set them — the DB's `DEFAULT NOW()` / `ON UPDATE CURRENT_TIMESTAMP` handle them.

The `created` user model has its `ID`, `CreatedAt`, and `UpdatedAt` populated by GORM's `tx.Create` directly (from `LAST_INSERT_ID()` and `time.Now()`). No re-select from DB is needed — unlike the sqlx implementation which does a `SELECT` after insert. All fields consumed by `ToEntity` (`ID`, `PlayerID`, `RoleID`, OAuth fields, `CreatedAt`, `UpdatedAt`) are populated correctly.

## BulkAdd

```go
func (r *Repository) BulkAdd(ctx context.Context, playerIDs []int, roleID int) error {
    if len(playerIDs) == 0 {
        return nil
    }
    entries := make([]Model, len(playerIDs))
    for i, id := range playerIDs {
        entries[i] = Model{PlayerID: id, RoleID: roleID}
    }
    if err := r.db.WithContext(ctx).CreateInBatches(&entries, 100).Error; err != nil {
        return fmt.Errorf("failed to bulk insert User records: %w", err)
    }
    return nil
}
```

## Behavior Changes from sqlx Migration

**`GetRoleByName` — soft-delete scope now applied**

The original sqlx query (`SELECT ... FROM user_role WHERE name = ?`) has no `deleted_at IS NULL` filter. GORM applies the soft-delete scope automatically because `RoleModel` embeds `GormModelBase` (which has `gorm.DeletedAt`). In practice this is harmless — the two seeded roles (`admin`, `player`) are never soft-deleted.

**`CreatePlayerAndUser` — no re-select after insert**

The sqlx implementation does a `SELECT ... FROM user WHERE id = ?` after inserting the user row to get a fully DB-populated model. The GORM implementation returns the `created` model directly, with `ID`, `CreatedAt`, and `UpdatedAt` populated by GORM during `Create`. The values are functionally equivalent (GORM sets timestamps from `time.Now()`; the DB uses `DEFAULT NOW()`).

**`SoftDelete` — no `RowsAffected` check**

The sqlx implementation checks `RowsAffected == 1` and returns an error if not. The GORM `db.Delete(&Model{}, id)` does not check `RowsAffected`. This is consistent with all other GORM phase implementations.

## Test Migration

No existing tests. Write new integration tests:
- `TestGetByID_Found` / `TestGetByID_NotFound`
- `TestGetByPlayerID_Found` / `TestGetByPlayerID_NotFound`
- `TestGetByOAuth_Found` / `TestGetByOAuth_NotFound`
- `TestGetRoleByName_Found` / `TestGetRoleByName_NotFound`
- `TestAdd`
- `TestAddWithOAuth`
- `TestCreatePlayerAndUser` — verify both player and user rows created atomically; verify rollback on error
- `TestBulkAdd`
- `TestSoftDelete`

**Test infrastructure:** Test file uses `package user_test`. Set up each test with:

```go
db := testHelpers.NewTestDB(t)
repo := testHelpers.NewUserRepo(db)
```

FK prerequisites: use `testHelpers.CreateTestPlayer(t, db)` for player IDs. The two seeded `user_role` rows (`admin`, `player`) can be looked up by name via `repo.GetRoleByName` or their IDs hardcoded as `1` / `2` if verified stable.

**As part of this phase**, add to `testHelpers/helpers.go`:
- `NewUserRepo(db *gorm.DB) *user.Repository` — wrapper over `user.NewRepositoryFromDB`

**Exception — `TestCreatePlayerAndUser`:** This test calls `r.db.Transaction(...)` internally. MySQL does not support true nested transactions — an inner `BEGIN` implicitly commits the outer transaction, defeating the tx rollback cleanup. For this test only, use `DELETE FROM user WHERE id = ?` and `DELETE FROM player WHERE id = ?` (using IDs from the returned model) in a dedicated `t.Cleanup` closure instead of relying on the outer rollback.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/user/...` passes (or skips)
3. Smoke test: OAuth login flow (`GET /api/auth/callback`) creates user correctly; `GET /api/auth/me` returns user info
