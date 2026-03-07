# Phase 3d — User Repository

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
    PlayerID      int     `gorm:"column:player_id"`
    RoleID        int     `gorm:"column:role_id"`
    OAuthProvider *string `gorm:"column:oauth_provider"`
    OAuthSubject  *string `gorm:"column:oauth_subject"`
    Email         *string `gorm:"column:email"`
    DisplayName   *string `gorm:"column:display_name"`
    AvatarURL     *string `gorm:"column:avatar_url"`
}

func (Model) TableName() string { return "user" }

type RoleModel struct {
    base.GormModelBase
    Name string `gorm:"column:name"`
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
        // Insert player (using the player package model is not available here;
        // use a local struct or raw Exec)
        playerResult := tx.Exec("INSERT INTO player (name) VALUES (?)", playerName)
        if playerResult.Error != nil {
            return fmt.Errorf("failed to insert player in CreatePlayerAndUser: %w", playerResult.Error)
        }

        var playerID int64
        tx.Raw("SELECT LAST_INSERT_ID()").Scan(&playerID)

        created = Model{
            PlayerID:      int(playerID),
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

Note: The player insert inside the transaction intentionally uses a raw `Exec` rather than importing the player repository — this keeps the user package self-contained and avoids circular imports. The `created` model has its ID populated by `tx.Create`.

Alternatively, define a local `playerModel` struct scoped to this method:

```go
type txPlayerModel struct {
    ID   int    `gorm:"primaryKey;column:id"`
    Name string `gorm:"column:name"`
}
func (txPlayerModel) TableName() string { return "player" }
```

Then use `tx.Create(&txPlayerModel{Name: playerName})` to get the player ID without raw SQL.

Either approach is valid. The local struct approach avoids the `LAST_INSERT_ID()` raw query.

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

Add `testhelpers_test.go` with `newTestDB(t)`. Cleanup: truncate `user` and `player` tables (`user_role` is seeded and should not be truncated).

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/user/...` passes (or skips)
3. Smoke test: OAuth login flow (`GET /api/auth/callback`) creates user correctly; `GET /api/auth/me` returns user info
