# Seeded Player → OAuth Account Linking

## Context

The seeder bulk-inserts players and creates corresponding `user` rows (player_id + role_id only — no OAuth fields). When a real person logs in via Google for the first time, the OAuth callback calls `GetByOAuth(provider, sub)`, finds no match, and calls `CreateWithOAuth` — creating a **brand-new orphan player** instead of reusing the existing seeded record and its game/deck history.

The fix: store known emails on seeded user rows, then add an email-based fallback in the OAuth callback. No new tables. No UI changes.

---

## Step 1 — User repository additions

**File**: `lib/repositories/user/repo.go`

Add three new methods:

```go
// GetByEmail returns the user with the given email, or nil if not found.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*Model, error)

// UpdateOAuth writes OAuth fields onto an existing user row.
func (r *Repository) UpdateOAuth(ctx context.Context, userID int, provider, subject, email, displayName, avatarURL string) error

// SetEmail updates the email field on the user row for a given player.
func (r *Repository) SetEmail(ctx context.Context, playerID int, email string) error
```

**File**: `lib/repositories/interfaces.go`

Add all three signatures to the `UserRepository` interface.

---

## Step 2 — User business additions

**File**: `lib/business/user/functions.go`

Add two new business functions following the existing functional-DI pattern:

```go
// GetByEmail wraps repo.GetByEmail; converts Model → Entity on hit, returns nil on miss.
func GetByEmail(userRepo repos.UserRepository) GetByEmailFunc

// LinkOAuth calls repo.UpdateOAuth then returns the updated Entity.
func LinkOAuth(userRepo repos.UserRepository) LinkOAuthFunc
```

Wire both into the `Functions` struct and `NewFunctions(...)`.

---

## Step 3 — Seeder email data + wiring

**New file**: `data/playerEmails.json`

A name → email map for the seeded players:
```json
{
  "Mike": "mike@example.com",
  "James": "james@example.com"
}
```
Replace with real email addresses. **Add this file to `.gitignore`** — it contains PII.

**File**: `lib/seeder/seeder.go`

After `repos.Users.BulkAdd(...)` creates the bare user rows, add a second pass that reads `playerEmails.json` and calls `repos.Users.SetEmail(ctx, playerID, email)` for each entry. Player IDs are already available from the bulk-insert step.

---

## Step 4 — OAuth callback fallback

**File**: `lib/routers/auth.go` — `Callback` handler

Extend the user-lookup block with an email fallback:

```go
u, err = a.usersBiz.GetByOAuth(ctx, providerGoogle, googleUser.Sub)
if err != nil { /* 500 */ }

if u == nil {
    // Fallback: seeded user row may exist for this email (no OAuth sub yet)
    u, err = a.usersBiz.GetByEmail(ctx, googleUser.Email)
    if err != nil { /* 500 */ }

    if u != nil {
        // Link OAuth credentials to the existing seeded player
        u, err = a.usersBiz.LinkOAuth(ctx, u.ID, providerGoogle, googleUser.Sub,
            googleUser.Email, googleUser.Name, googleUser.Picture)
    } else {
        // Truly new user — create player + user atomically
        u, err = a.usersBiz.CreateWithOAuth(ctx, googleUser.Name, providerGoogle,
            googleUser.Sub, googleUser.Email, googleUser.Name, googleUser.Picture)
    }
    if err != nil { /* 500 */ }
}
// JWT issuance continues as before...
```

Wire `GetByEmail` and `LinkOAuth` into `authRouter.usersBiz` (the `user.Functions` struct).

---

## Step 5 — Tests

Follow the established testing patterns (see `TESTING POLICY` in `frontend-revamp-plan.md`):

- **Repo** (`lib/repositories/user/repo_test.go`): table-driven tests for `GetByEmail` (hit + miss), `UpdateOAuth`, `SetEmail` — use `testHelpers.NewTestDB(t)`
- **Business** (`lib/business/user/functions_test.go`): `GetByEmail` and `LinkOAuth` — mock repo with `Fn` fields; table-driven
- **Router** (`lib/routers/auth_test.go`): add test case for the email-fallback path (GetByOAuth returns nil → GetByEmail returns a user → LinkOAuth called → JWT issued)

---

## Critical files

| File | Change |
|------|--------|
| `lib/repositories/user/repo.go` | Add `GetByEmail`, `UpdateOAuth`, `SetEmail` |
| `lib/repositories/interfaces.go` | Add all three to `UserRepository` interface |
| `lib/business/user/functions.go` | Add `GetByEmail`, `LinkOAuth` business functions |
| `lib/routers/auth.go` | Add email-fallback + `LinkOAuth` call in `Callback` |
| `lib/seeder/seeder.go` | Read `playerEmails.json`, call `SetEmail` for each player |
| `data/playerEmails.json` | New — name → email map; gitignored |

---

## Verification

1. Seed the DB fresh (`SEED=1`)
2. Confirm `user` rows exist for seeded players with emails set (`SELECT id, player_id, email FROM user`)
3. Log in with a Google account whose email matches one in `playerEmails.json`
4. Confirm `GET /api/auth/me` returns the **seeded** player's `player_id` (not a newly-minted one)
5. Confirm game/deck history is intact under that player

Run `go test ./lib/...` after implementation.
