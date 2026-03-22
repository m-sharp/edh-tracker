# EDH Tracker — Frontend Revamp Plan

## Overview & Decisions

- **Auth**: Google OAuth only; stateless JWT stored in HttpOnly cookie
- **User ↔ Player**: 1-to-1; first login auto-creates a linked Player record
- **Pod membership**: Invite-based via share link / invite code (UUID)
- **PodManager role**: New `player_pod_role` table; pod creator = manager; managers can promote members; seeder assigns "Mike" as manager
- **Pagination**: Server-side (limit + offset query params); returns `{"items": [...], "total": N}` when params present
- **Navigation**: Pod selector dropdown in header + link to user's own Player page
- **New game form**: Moved into pod context at `/pod/:podId/new-game`
- **Route changes**:
  - `/deck/:deckId` → `/player/:playerId/deck/:deckId`
  - `/game/:gameId` → `/pod/:podId/game/:gameId`
  - `/pod/:podId/new-game` (new)
  - Top-level `/decks`, `/players`, `/games` removed (redirect or delete)

---

## Testing Policy

Any new production code added to this project **MUST** have corresponding unit tests following the established patterns:

- **Business functions**: Mock structs with per-method `Fn` fields implementing the repository interfaces; panic on unexpected calls. See `lib/business/player/functions_test.go` and `lib/business/pod/functions_test.go`.
- **Repositories**: Real MySQL integration tests using `testHelpers.NewTestDB(t)` (wraps the connection in a transaction that auto-rolls back on `t.Cleanup`). Use fixture helpers from `lib/repositories/testHelpers/helpers.go` (e.g. `CreateTestPlayer`, `CreateTestDeck`, `CreateTestGame`). Test package naming convention: `package <domain>_test`. See `lib/repositories/player/repo_test.go`.
- **Routers**: `httptest.NewRecorder`; inject mock closures directly into the `Functions` struct. See `lib/routers/player_test.go` and `lib/routers/pod_test.go`.

---

## Dependency Graph

```
Phase 0 (Google Cloud Console Setup — manual, one-time human task)
  ↓
Phase 1 (Migrations)
  ├── Phase 2A (OAuth backend)       ← requires Phase 0 credentials
  │     └── Phase 3 (Frontend foundation: auth context + routing)
  │           ├── Phase 4A (Pod landing page)
  │           ├── Phase 4B (Player page revamp)
  │           ├── Phase 4C (Deck page revamp)
  │           └── Phase 4D (Game page revamp)
  ├── Phase 2B (Pod roles + invite backend)
  ├── Phase 2C (Edit/delete/new API endpoints)
  └── Phase 2D (Server-side pagination backend)
```

**Phase 2A–2D run in parallel after Phase 1.**
**Phase 4A–4D run in parallel after Phase 3.**
**Phase 0 can be done any time before starting Phase 2A.**

---

## Progress Tracking

Mark items `[x]` as they are completed during implementation.

### Phase 0 — Google Cloud Console Setup
- [x] Create Google Cloud project (or select existing)
- [x] Configure OAuth consent screen (app name, support email, scopes: openid/email/profile, test users)
- [x] Create OAuth 2.0 Client ID credentials (Web application type)
- [x] Add local redirect URI: `http://localhost:8080/api/auth/google/callback`
- [x] Add production redirect URI when domain is known
- [x] Record `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` for env vars

### Phase 1 — Database Migrations
- [x] Migration 17: `player_pod_role` table
- [x] Migration 18: `pod_invite` table
- [x] Migration 19: Soft-delete columns (`game` + `deck`) — **SKIPPED**: columns already existed from Migration 6; no file created
- [x] Migrations 17 and 18 registered in `lib/migrations/migrate.go` (19 not needed)

### Phase 2A — OAuth Backend
- [x] Fix CORS bugs in `lib/trackerHttp/http.go` (wildcard+credentials invalid, preflight writes body, `Add` vs `Set`)
- [x] Add dependencies (`golang.org/x/oauth2`, `github.com/golang-jwt/jwt/v5`)
- [x] Add config keys in `lib/config.go`
- [x] `lib/trackerHttp/auth.go` (`RequireAuth`, `OptionalAuth`, sliding re-issue)
- [x] User business: `GetByOAuth`, `CreateWithOAuth`
- [x] `lib/routers/auth.go` (login, callback, logout, me)
- [x] Wire `AuthRouter` + apply middleware in `api.go`

### Phase 2B — Pod Roles + Invite
- [x] `lib/repositories/playerPodRole/` (model + repo)
- [x] `lib/repositories/podInvite/` (model + repo)
- [x] Update `lib/repositories/repositories.go`
- [x] Update pod repository (`SoftDelete`, `Update`, `RemovePlayer`, `GetPlayerIDs`)
- [x] Pod business: role/invite/leave/manage functions
- [x] New pod API endpoints (PATCH, DELETE, invite, join, leave, kick, promote)
- [x] Update seeder with `player_pod_role` seed data

### Phase 2C — Edit/Delete/New Endpoints ✅
- [x] Remove `POST /api/player` route + handler (player creation is OAuth-only)
- [x] Player: `PATCH /api/player` + business + repo `Update`
- [x] Deck: `PATCH /api/deck` (full update) + `DELETE /api/deck` + business + repo
- [x] Game: `PATCH /api/game` + `DELETE /api/game` + result CRUD endpoints + business + repo
- [x] New GET filters: `GET /api/players?pod_id`, `GET /api/decks?pod_id`, `GET /api/games?player_id`

### Phase 2D — Server-Side Pagination
- [x] `PaginatedResponse[T]` type in `lib/business/pagination.go`
- [x] Game repo: `GetAllByPodPaginated`, `GetAllByDeckPaginated`
- [x] Deck repo: `GetAllByPodPaginated`, `GetAllByPlayerPaginated`
- [x] Router changes for `GET /api/games` and `GET /api/decks`

> **GORM Implementation Notes:** Pagination is straightforward with GORM. Use `.Limit(limit).Offset(offset)` chained onto existing GORM queries for data retrieval. Count queries use `.Model(&T{}).Where(...).Count(&total)` with the same WHERE clause as the data query — no raw SQL needed. Each paginated repo method must also be added to the corresponding interface in `lib/repositories/interfaces.go`. Repository tests use `testHelpers.NewTestDB(t)` with real fixtures (not mocks).

### Phase 3 — Frontend Foundation

**Sub-phase order:** 3A → (3B ∥ 3D) → 3C → (3E ∥ 3G) → 3F

#### 3A — Types
- [x] `app/src/types.ts` — add `Pod`, `PlayerWithRole`, `PaginatedResponse`, `DeckUpdateFields`, `GameResultUpdateFields`, `NewGameResultWithGame`; add `player_id` to `GameResult`

#### 3B — HTTP Client
- [x] `app/src/http.ts` — add `credentials: "include"` to all 12 existing fetches
- [x] `app/src/http.ts` — add `GetMe`, `Logout`
- [x] `app/src/http.ts` — add all remaining new API functions (pod, player, deck, game)

#### 3C — AuthContext
- [x] `app/src/auth.tsx` — `AuthUser`, `AuthContextValue`, `AuthProvider`, `useAuth`

#### 3D — Login Page (no deps; parallel with 3B)
- [x] `app/src/routes/login.tsx`

#### 3E — RequireAuth (depends on 3C)
- [x] `app/src/routes/RequireAuth.tsx`

#### 3F — Route Restructure (depends on 3C, 3D, 3E)
- [x] `app/src/index.tsx` — wrap with `AuthProvider`; new route tree with `RequireAuth`; stub components for Phase 4 routes; redirects for `/decks`, `/players`, `/games`

#### 3G — Navigation Revamp (depends on 3C; parallel with 3E)
- [x] `app/src/routes/root.tsx` — remove old nav links; logged-out/in states (avatar, pod selector with `localStorage`, logout) (pod selector + auth UI)

### Phase 4A — Pod Landing Page
- [x] `app/src/routes/pod.tsx` (tabs: Decks, Players, Games, Settings)
- [x] `app/src/routes/new.tsx` (pod context: drop global pod selector, use `useParams`)
- [x] `app/src/routes/join.tsx` (invite code landing page)

### Phase 4B — Player Page Revamp
- [x] `app/src/routes/player.tsx` — tabs (Overview, Decks, Games, Settings)
- [x] `PlayerSettingsTab` (edit name, pod list with leave button, create pod)

### Phase 4C — Deck Page Revamp
- [x] `app/src/routes/deck.tsx` — new route + tabs (Overview, Games, Settings)
- [x] Update deck links in `app/src/stats.tsx` (`CommanderColumn`)
- [x] Update deck links in `app/src/routes/game.tsx` (results grid)

### Phase 4D — Game Page Revamp
- [x] `app/src/routes/game.tsx` — new route + inline editing + result CRUD
- [x] Update game links in `app/src/matches.tsx` (include `pod_id`)

### Phase 5 - Unsorted Work
- [ ] Handle user oath token expiration mid session
---

## Phase 0 — Google Cloud Console Setup

**One-time manual setup. Must be completed before starting Phase 2A. No code changes — follow these steps in a browser.**

### Step 1: Create a Google Cloud Project

1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Click the project selector at the top → **New Project**
3. Name it (e.g., "EDH Tracker") → **Create**
4. Select the new project as active

### Step 2: Configure the OAuth Consent Screen

1. In the left sidebar: **APIs & Services → OAuth consent screen**
2. User type: **External** (allows any Google account; lets the group sign in freely)
3. Fill in:
   - App name: `EDH Tracker`
   - User support email: your email
   - Developer contact email: your email
4. Scopes: click **Add or remove scopes** → add:
   - `openid`
   - `https://www.googleapis.com/auth/userinfo.email`
   - `https://www.googleapis.com/auth/userinfo.profile`
5. Test users: add each player's Google account email address. While the app is in "Testing" mode, only these accounts can sign in (up to 100 accounts — sufficient for a pod tracker).
6. Save and continue through all steps.

> **Note:** You do not need to publish the app or go through Google's verification process. Staying in Testing mode is fine for a private group app. If you ever want to open it to new players with arbitrary Google accounts, you'd publish and verify — but for an invite-only pod tracker, Testing mode with manually-added users is the right posture.

> **Docs:** [OAuth consent screen setup](https://developers.google.com/workspace/guides/configure-oauth-consent)

### Step 3: Create OAuth 2.0 Credentials

1. **APIs & Services → Credentials → Create Credentials → OAuth 2.0 Client ID**
2. Application type: **Web application**
3. Name: `EDH Tracker Web`
4. Authorized redirect URIs — add both:
   - `http://localhost:8080/api/auth/google/callback` (local development)
   - Your production URI when the domain is known (e.g., `https://yourdomain.com/api/auth/google/callback`)
5. Click **Create**
6. Copy the **Client ID** and **Client Secret** — these become your `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` env vars

> **Note:** Google allows `http://localhost` redirect URIs without any verification — only production HTTPS URIs are subject to the consent screen publication requirements.

> **Docs:** [Setting up OAuth 2.0 for web server apps](https://developers.google.com/identity/protocols/oauth2/web-server)

### Env vars this phase unlocks

| Env Var | Source |
|---------|--------|
| `GOOGLE_CLIENT_ID` | Credentials page |
| `GOOGLE_CLIENT_SECRET` | Credentials page |
| `OAUTH_REDIRECT_URL` | Set to `http://localhost:8080/api/auth/google/callback` locally |
| `JWT_SECRET` | Generate locally: `openssl rand -hex 32` |
| `FRONTEND_URL` | `http://localhost:8081` locally |

---

## Phase 1 — Database Migrations

**Prerequisite for all other phases. Must land first.**

### Migration 17: `player_pod_role` table

```sql
CREATE TABLE player_pod_role (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    pod_id      INT NOT NULL,
    player_id   INT NOT NULL,
    role        ENUM('manager', 'member') NOT NULL DEFAULT 'member',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME NULL,
    UNIQUE KEY uq_ppr (pod_id, player_id),
    INDEX idx_ppr_pod_id    (pod_id),
    INDEX idx_ppr_player_id (player_id),
    INDEX idx_ppr_deleted_at (deleted_at),
    FOREIGN KEY (pod_id)    REFERENCES pod(id),
    FOREIGN KEY (player_id) REFERENCES player(id)
);
```

### Migration 18: `pod_invite` table

```sql
CREATE TABLE pod_invite (
    id                   INT AUTO_INCREMENT PRIMARY KEY,
    pod_id               INT NOT NULL,
    invite_code          VARCHAR(36) NOT NULL UNIQUE,
    created_by_player_id INT NOT NULL,
    expires_at           TIMESTAMP NULL,
    used_count           INT NOT NULL DEFAULT 0,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at           DATETIME NULL,
    INDEX idx_pi_pod_id      (pod_id),
    INDEX idx_pi_invite_code (invite_code),
    INDEX idx_pi_deleted_at  (deleted_at),
    FOREIGN KEY (pod_id)               REFERENCES pod(id),
    FOREIGN KEY (created_by_player_id) REFERENCES player(id)
);
```

### Migration 19: Soft-delete columns (if not present)

Check whether `game` and `deck` tables have a `deleted_at` column. If either is missing, add:

```sql
-- Only if not already present:
ALTER TABLE game ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;
ALTER TABLE deck ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;
ALTER TABLE game ADD INDEX idx_game_deleted_at (deleted_at);
ALTER TABLE deck ADD INDEX idx_deck_deleted_at (deleted_at);
```

> **Note:** `deck` already has `retired` (soft-retire). `deleted_at` is for hard-remove of a deck from the system entirely (owner deletes, not just retires).

### File changes

- Create `lib/migrations/17.go`, `18.go`, `19.go` following the existing Migration pattern
- Register all three in `lib/migrations/migrate.go`

---

## Phase 2A — Backend: Google OAuth + Session Management

**Runs in parallel with 2B, 2C, 2D. Requires Phase 1.**

### New dependencies

```
golang.org/x/oauth2
golang.org/x/oauth2/google
github.com/golang-jwt/jwt/v5
```

Add to `go.mod` / vendor.

### New config keys in `lib/config.go`

| Key | Env Var | Notes |
|-----|---------|-------|
| `GoogleClientID` | `GOOGLE_CLIENT_ID` | Required |
| `GoogleClientSecret` | `GOOGLE_CLIENT_SECRET` | Required |
| `OAuthRedirectURL` | `OAUTH_REDIRECT_URL` | e.g. `http://localhost:8080/api/auth/google/callback` |
| `JWTSecret` | `JWT_SECRET` | Required; used to sign/verify JWTs |
| `FrontendURL` | `FRONTEND_URL` | Used to redirect back to frontend after login |

### New middleware: `lib/trackerHttp/auth.go`

- `RequireAuth(jwtSecret string, secure bool) mux.MiddlewareFunc` — validates JWT from `edh_session` HttpOnly cookie; injects `userID` and `playerID` into request context; returns 401 if missing/invalid. **Sliding window**: on every valid request, re-issues a fresh `edh_session` cookie with a new `exp: now+24h` to keep sessions alive for active users. The `secure` param (derived from `cfg.Get(lib.Dev) == ""`) sets the `Secure` flag on the re-issued cookie.
- `OptionalAuth(jwtSecret string, secure bool) mux.MiddlewareFunc` — same but never rejects; injects nil if no valid token; does not re-issue cookie.
- Context helpers live in `lib/utils/context.go`: `ContextWithUserInfo` / `UserFromContext`

### User business layer additions (`lib/business/user/`)

New functions in `lib/business/user/functions.go`:

- `GetByOAuth(repo ...) GetByOAuthFunc` — find user by (provider, subject); returns `nil, nil` if not found
- `CreateWithOAuth(repo ...) CreateWithOAuthFunc` — creates Player + User in one DB transaction; returns linked Entity

New types in `lib/business/user/types.go`:
- `GetByOAuthFunc func(ctx context.Context, provider, subject string) (*Entity, error)`
- `CreateWithOAuthFunc func(ctx context.Context, playerName, provider, subject, email, displayName, avatarURL string) (*Entity, error)`

Add both to `business.Business.User` struct and wire in `NewBusiness`.

### New router: `lib/routers/auth.go`

`AuthRouter` struct holds `cfg *lib.Config`, `log *zap.Logger`, user business functions.

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/auth/google` | None | Accept optional `?redirect=` param; set CSRF state cookie (encodes state + redirect path); redirect to Google OAuth consent screen |
| `GET` | `/api/auth/google/callback` | None | Validate state; exchange code; get Google profile; find-or-create User+Player; issue JWT; redirect to `FRONTEND_URL` + decoded redirect path (or `/` if none) |
| `POST` | `/api/auth/logout` | None | Clear `edh_session` cookie; return 204 |
| `GET` | `/api/auth/me` | Required | Return current `user.Entity` as JSON |

**Callback logic:**
1. Validate `state` param matches cookie (CSRF); decode redirect path from state
2. Exchange code → Google access token
3. Fetch Google userinfo (email, sub, name, picture)
4. `user.GetByOAuth(ctx, "google", sub)`
5. If nil → `user.CreateWithOAuth(ctx, googleName, "google", sub, email, googleName, pictureURL)`
6. Build JWT: `{ user_id, player_id, exp: now+24h }`, sign with `JWT_SECRET`
7. Set `edh_session` cookie: HttpOnly, SameSite=Lax, Path=/, **Secure=true when `DEV` env var is not set** (i.e., `cfg.Get(lib.Dev) == ""`). This ensures the cookie is only sent over HTTPS in production while remaining accessible in local HTTP dev.
8. Redirect to `FRONTEND_URL` + redirect path (e.g. `/join?code=xxx`), or `FRONTEND_URL/` if none

> **Note on `/join` flow**: Frontend's `/join` route redirects unauthenticated users to `/api/auth/google?redirect=/join?code=xxx`. The state cookie encodes this path. After Google auth, the callback redirects the browser to `FRONTEND_URL/join?code=xxx` where the now-authenticated user is sent to join the pod.

### Wire into `api.go`

- Register `AuthRouter` in `SetupRoutes`
- Apply `RequireAuth` middleware to **all** state-changing sub-routers — both new ones (pod roles, invite) and existing ones (game, deck, player). This is the point where the existing `POST /api/game`, `POST /api/deck` routes become auth-protected.
- Apply `OptionalAuth` to routers that need conditional behavior (read-only endpoints that benefit from knowing who the caller is but don't require login)
- `secure` param passed to both middleware constructors: `cfg.Get(lib.Dev) == ""`

**Auth requirements for existing endpoints once middleware is applied:**

| Endpoint | Auth | Authorization check |
|----------|------|---------------------|
| `POST /api/game` | RequireAuth | Caller must be a member (any role) of the `pod_id` in the request body — call `pod.GetRole(podID, callerPlayerID)`; return 403 if empty |
| `POST /api/deck` | RequireAuth | No additional check — caller creates a deck under their own player_id (JWT player_id used, not request body) |
| `POST /api/player` | **Remove** | See Phase 2C |

### CORS fixes in `lib/trackerHttp/http.go`

Three bugs in the existing CORS middleware must be fixed before OAuth + cookies work correctly:

1. **Wildcard + Credentials is invalid**: `Access-Control-Allow-Origin: *` combined with `Access-Control-Allow-Credentials: true` is rejected by all browsers per the CORS spec. Replace `"*"` with the explicit `FRONTEND_URL` value from config.
2. **Preflight handler writes a body on 204**: `http.Error(w, "No Content", http.StatusNoContent)` sets `Content-Type: text/plain` and writes a response body. Replace with `w.WriteHeader(http.StatusNoContent)`.
3. **`Add` vs `Set` on CORS headers**: `Header.Add` appends, causing duplicate header values on repeated requests. Replace with `Header.Set` for all four CORS headers.

Updated `CORSMiddleware` should accept the allowed origin from config so it can be set explicitly:

```go
func CORSMiddleware(origin string) MiddlewareFunc {
    return func(nextHandler http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
            w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
            nextHandler(w, r)
        }
    }
}
```

Pass `cfg.Get(lib.FrontendURL)` as `origin` when registering middleware in `api.go`.

> **Note:** The long-term approach to CORS (dev proxy, production reverse proxy) is a separate concern tracked in `.claude/production-ready-cors.md`.

---

## Phase 2B — Backend: Pod Roles + Invite System

**Runs in parallel with 2A, 2C, 2D. Requires Phase 1.**

### New repository: `lib/repositories/playerPodRole/`

**`model.go`:**
```go
type Model struct {
    base.ModelBase
    PodID    int    `db:"pod_id"`
    PlayerID int    `db:"player_id"`
    Role     string `db:"role"`  // "manager" | "member"
}
```

**`repo.go`** — `Repository` struct with:
- `GetRole(ctx, podID, playerID int) (*Model, error)` — nil if not in pod
- `SetRole(ctx, podID, playerID int, role string) error` — INSERT ... ON DUPLICATE KEY UPDATE
- `GetMembersWithRoles(ctx, podID int) ([]Model, error)`
- `BulkAdd(ctx, podID int, playerIDs []int, role string) error`

### New repository: `lib/repositories/podInvite/`

**`model.go`:**
```go
type Model struct {
    base.ModelBase
    PodID             int        `db:"pod_id"`
    InviteCode        string     `db:"invite_code"`
    CreatedByPlayerID int        `db:"created_by_player_id"`
    ExpiresAt         *time.Time `db:"expires_at"`
    UsedCount         int        `db:"used_count"`
}
```

**`repo.go`** — `Repository` struct with:
- `GetByCode(ctx, code string) (*Model, error)` — nil if not found
- `Add(ctx, podID, createdByPlayerID int, code string, expiresAt *time.Time) error`
- `IncrementUsedCount(ctx, code string) error`

### Update `lib/repositories/repositories.go`

Add `PlayerPodRoles *playerPodRole.Repository` and `PodInvites *podInvite.Repository` fields + initialization in `New()`.

### Update pod repository (`lib/repositories/pod/repo.go`)

Add:
- `SoftDelete(ctx, podID int) error` — sets `deleted_at = NOW()`
- `Update(ctx, podID int, name string) error`
- `RemovePlayer(ctx, podID, playerID int) error` — soft-delete player_pod row
- `GetPlayerIDs(ctx, podID int) ([]int, error)`

### Update pod business layer (`lib/business/pod/`)

New functions in `lib/business/pod/functions.go`:

- `GetRole(roleRepo ...) GetRoleFunc` — returns "manager", "member", or "" (not in pod)
- `PromoteToManager(roleRepo ...) PromoteToManagerFunc` — verifies caller is manager; sets target to manager
- `GenerateInvite(inviteRepo ...) GenerateInviteFunc` — creates UUID code; stores in pod_invite with `expires_at = NOW() + 7 days`; returns code string
- `JoinByInvite(inviteRepo, podRepo, roleRepo ...) JoinByInviteFunc` — validates code exists and is not expired (`expires_at IS NULL OR expires_at > NOW()`); adds to player_pod + player_pod_role as member; increments used_count
- `Leave(podRepo, roleRepo ...) LeaveFunc` — removes caller from player_pod and player_pod_role; returns 403 if caller is the only manager (must promote someone else first)
- `SoftDelete(podRepo ...) SoftDeleteFunc`
- `Update(podRepo ...) UpdateFunc` — update pod name
- `GetMembersWithRoles(roleRepo, playerRepo ...) GetMembersWithRolesFunc` — returns []PlayerWithRole
- `RemovePlayer(podRepo, roleRepo ...) RemovePlayerFunc`

**Update `Create`** function: after `pod.Add()`, call `roleRepo.SetRole(podID, creatorPlayerID, "manager")`.
**Update `AddPlayer`** function: also call `roleRepo.SetRole(podID, playerID, "member")`.

### New pod API endpoints (`lib/routers/pod.go`)

| Method | Path | Auth | PodManager? | Body / Params | Description |
|--------|------|------|------------|---------------|-------------|
| `PATCH` | `/api/pod` | Required | Yes | `{name}` + `?pod_id=X` | Update pod name |
| `DELETE` | `/api/pod` | Required | Yes | `?pod_id=X` | Soft delete pod |
| `POST` | `/api/pod/invite` | Required | Yes | `{pod_id}` | Generate invite code (7-day expiry); response: `{invite_code}` |
| `POST` | `/api/pod/join` | Required | No | `{invite_code}` | Join pod; response: pod entity |
| `POST` | `/api/pod/leave` | Required | No | `{pod_id}` | Self-remove caller from pod; response: 204 |
| `PATCH` | `/api/pod/player` | Required | Yes | `{pod_id, player_id}` | Promote player_id to manager |
| `DELETE` | `/api/pod/player` | Required | Yes | `{pod_id, player_id}` | Kick player_id from pod (manager only) |

PodManager check: extract `playerID` from JWT context → `pod.GetRole(podID, playerID)` → 403 if not manager.

### Update seeder (`lib/seeder/seeder.go`)

After `BulkAddPlayers` + `repos.Pods.BulkAddPlayers(...)`, seed `player_pod_role`:
1. Call `repos.PlayerPodRoles.BulkAdd(ctx, podID, allPlayerIDs, "member")`
2. Look up Mike's playerID from `playerIDs["Mike"]`
3. Call `repos.PlayerPodRoles.SetRole(ctx, podID, mikeID, "manager")` to override

---

## Phase 2C — Backend: Edit/Delete/New API Endpoints

**Runs in parallel with 2A, 2B, 2D. Requires Phase 1.**

### New GET endpoints

| Method | Path | Query Params | Description |
|--------|------|-------------|-------------|
| `GET` | `/api/players` | `?pod_id=X` | Filter players by pod; includes role from player_pod_role |
| `GET` | `/api/decks` | `?pod_id=X` | Decks owned by any player in the pod |
| `GET` | `/api/games` | `?player_id=X` | All games any of the player's decks participated in |

**Backend notes:**
- `GET /api/players?pod_id=X`: join player_pod + player; fetch stats; fetch role per player from player_pod_role
- `GET /api/decks?pod_id=X`: get all playerIDs in pod via `pod.GetPlayerIDs`, then `deck.GetAllByPlayerIDs`
- `GET /api/games?player_id=X`: get deckIDs for player, then games where any deck_id is in those deckIDs

Player entity response for pod-scoped query should include `role` field; define `PlayerWithRoleEntity` in player business layer.

### Remove standalone player creation endpoint

Delete the `POST /api/player` route from `lib/routers/player.go` and `api.go`. Player creation now flows exclusively through the OAuth callback (`CreateWithOAuth`). The associated business function (`player.Create` or equivalent called directly) and repository method can be kept if used by `CreateWithOAuth`; only the HTTP handler and route registration are removed.

### Player edit endpoint

| Method | Path | Auth | Body | Description |
|--------|------|------|------|-------------|
| `PATCH` | `/api/player` | Required | `{name: string}` + `?player_id=X` | Update player display name; JWT playerID must match player_id |

New `Update(repo ...) UpdateFunc` → `func(ctx, playerID int, name string) error` in player business layer and repository.

### Deck endpoints

| Method | Path | Auth | Body / Params | Description |
|--------|------|------|---------------|-------------|
| `PATCH` | `/api/deck` | Required | `{name?, format_id?, commander_id?, partner_commander_id?, retired?}` + `?deck_id=X` | Edit deck fields; caller must own deck |
| `DELETE` | `/api/deck` | Required | `?deck_id=X` | Soft delete deck; caller must own deck |

Extend current `PATCH /api/deck` handler to handle all optional fields. Currently it only sets `retired=true` — that logic becomes part of the general update.

New business functions:
- `Update(ctx, deckID int, fields DeckUpdateFields) error`
- `SoftDelete(ctx, deckID int) error`

New repository functions in `lib/repositories/deck/`:
- `Update(ctx, deckID int, fields UpdateFields) error`
- `SoftDelete(ctx, deckID int) error`

Commander updates (edit `deck_commander` rows): delete existing + re-insert new ones.

**Ownership check:** fetch deck → verify `deck.PlayerID == callerPlayerID` from JWT; return 403 otherwise.

### Game endpoints

| Method | Path | Auth | Body / Params | Description |
|--------|------|------|---------------|-------------|
| `PATCH` | `/api/game` | Required | `{description}` + `?game_id=X` | Edit game description; PodManager only |
| `DELETE` | `/api/game` | Required | `?game_id=X` | Soft delete game; PodManager only |
| `POST` | `/api/game/result` | Required | `{game_id, deck_id, player_id, place, kill_count}` | Add result to existing game; PodManager only |
| `PATCH` | `/api/game/result` | Required | `{place?, kill_count?, deck_id?}` + `?result_id=X` | Edit game result; PodManager only |
| `DELETE` | `/api/game/result` | Required | `?result_id=X` | Remove game result; PodManager only |

**PodManager check for game endpoints:** fetch game → get pod_id → `pod.GetRole(podID, callerPlayerID)` → 403 if not manager.

Also ensure `GameResult` entity returned by `GET /api/game` includes `player_id` (needed for deck link `/player/:playerId/deck/:deckId`). Check current `gameResult.Entity`; add `PlayerID` field if missing.

New business functions in `lib/business/game/`:
- `Update(ctx, gameID int, description string) error`
- `SoftDelete(ctx, gameID int) error`
- `AddResult(ctx, gameID, deckID, playerID, place, killCount int) (int, error)`
- `UpdateResult(ctx, resultID, place, killCount, deckID int) error`
- `DeleteResult(ctx, resultID int) error`

New repository methods in `lib/repositories/game/`:
- `Update(ctx, gameID int, description string) error`
- `SoftDelete(ctx, gameID int) error`

New repository methods in `lib/repositories/gameResult/`:
- `Add(ctx, model Model) (int, error)` (single insert)
- `Update(ctx, resultID, place, killCount, deckID int) error`
- `Delete(ctx, resultID int) error`
- `GetByID(ctx, resultID int) (*Model, error)` (for PodManager check via game lookup)

---

## Phase 2D — Backend: Server-Side Pagination

**Runs in parallel with 2A, 2B, 2C. Requires Phase 1.**

### Response format

When `limit` and `offset` are provided, return:
```json
{ "items": [...], "total": 150, "limit": 25, "offset": 0 }
```
When pagination params are absent, return plain array (backwards compat).

Define `PaginatedResponse[T]` in a shared location (e.g. `lib/business/pagination.go`):
```go
type PaginatedResponse[T any] struct {
    Items  []T `json:"items"`
    Total  int `json:"total"`
    Limit  int `json:"limit"`
    Offset int `json:"offset"`
}
```

### Repository changes

**`lib/repositories/game/repo.go`** — add:
- `GetAllByPodPaginated(ctx, podID, limit, offset int) ([]Model, int, error)`
- `GetAllByDeckPaginated(ctx, deckID, limit, offset int) ([]Model, int, error)`

Each uses `SELECT ... LIMIT ? OFFSET ?` alongside `SELECT COUNT(*) FROM ...` with same WHERE clause.

**`lib/repositories/deck/repo.go`** — add:
- `GetAllByPodPaginated(ctx, podID, limit, offset int) ([]Model, int, error)`
- `GetAllByPlayerPaginated(ctx, playerID, limit, offset int) ([]Model, int, error)`

### Router changes

**`lib/routers/game.go` — `GET /api/games`:**
- If `limit` param present: use paginated repo method; return `PaginatedResponse`
- Supports: `?pod_id=X&limit=N&offset=M`, `?deck_id=X&limit=N&offset=M`, `?player_id=X&limit=N&offset=M`

**`lib/routers/deck.go` — `GET /api/decks`:**
- If `limit` param present: use paginated repo method; return `PaginatedResponse`
- Supports: `?pod_id=X&limit=N&offset=M`, `?player_id=X&limit=N&offset=M`

---

## Phase 3 — Frontend Foundation

**Requires Phase 2A (auth backend). Blocks Phase 4A–4D.**

### Dependency Graph

```
3A (Types)
  ↓
3B (HTTP client — credentials + GetMe/Logout + new functions)   ← depends on 3A
  ↓                                                3D (LoginPage) — no deps; run in parallel with 3B
3C (AuthContext)                                              ↓
  ├──→ 3E (RequireAuth)                    ─────────────────┐
  └──→ 3G (Navigation revamp)                               ↓
                                            3F (Route restructure) ← depends on 3C + 3D + 3E
```

**Parallel opportunities:** 3D with 3A/3B; 3E and 3G with each other (both depend only on 3C).

---

### Sub-Phase 3A — Types

**File:** `app/src/types.ts` — additions only. `AuthUser` lives in `auth.tsx`, not here.

```typescript
interface Pod {
    id: number;
    name: string;
    created_at: string;
    updated_at: string;
}

interface PlayerWithRole extends Player {
    role: "manager" | "member";
}

interface PaginatedResponse<T> {
    items: Array<T>;
    total: number;
    limit: number;
    offset: number;
}

interface DeckUpdateFields {
    name?: string;
    format_id?: number;
    commander_id?: number;
    partner_commander_id?: number | null;
    retired?: boolean;
}

interface GameResultUpdateFields {
    place?: number;
    kill_count?: number;
    deck_id?: number;
}

interface NewGameResultWithGame extends NewGameResult {
    game_id: number;
}
```

Also add `player_id: number` to the existing `GameResult` interface (needed for deck links at `/player/:playerId/deck/:deckId`).

---

### Sub-Phase 3B — HTTP Client

**File:** `app/src/http.ts`. **Depends on 3A.**

**Part 1 — Add `credentials: "include"` to all existing fetches.**
There are 11 functions / 12 `fetch` calls. `PostGame` and `PostCommander` already have options objects; the rest are bare. Example:
```typescript
// Before
const res = await fetch(`http://localhost:8080/api/players`);
// After
const res = await fetch(`http://localhost:8080/api/players`, { credentials: "include" });
```
The hardcoded `localhost:8080` base URL is **not** changed in this phase (pre-existing TODO).

**Part 2 — Add `GetMe` and `Logout`** (consumed by `AuthProvider` in 3C):
```typescript
export async function GetMe(): Promise<AuthUser> {
    const res = await fetch(`http://localhost:8080/api/auth/me`, { credentials: "include" });
    if (!res.ok) throw new Error("Unauthenticated");
    return res.json();
}

export async function Logout(): Promise<void> {
    await fetch(`http://localhost:8080/api/auth/logout`, {
        method: "POST",
        credentials: "include",
    });
}
```
`AuthUser` is imported from `./auth`.

**Part 3 — Add all remaining new functions** (all include `credentials: "include"`):
- `GetPod(podId: number): Promise<Pod>` — `GET /api/pod?pod_id=X`
- `GetPodsForPlayer(playerId: number): Promise<Array<Pod>>` — `GET /api/pod?player_id=X`
- `GetPlayersForPod(podId: number): Promise<Array<PlayerWithRole>>` — `GET /api/players?pod_id=X`
- `GetDecksForPod(podId: number, limit?: number, offset?: number): Promise<PaginatedResponse<Deck>>`
- `GetGamesForPod(podId: number, limit?: number, offset?: number): Promise<PaginatedResponse<Game>>`
- `GetGamesForPlayer(playerId: number): Promise<Array<Game>>` — `GET /api/games?player_id=X`
- `PostPod(name: string): Promise<Pod>` — `POST /api/pod`
- `PostPodInvite(podId: number): Promise<{invite_code: string}>` — `POST /api/pod/invite`
- `PostPodJoin(inviteCode: string): Promise<Pod>` — `POST /api/pod/join`
- `PostPodLeave(podId: number): Promise<void>` — `POST /api/pod/leave`
- `PatchPod(podId: number, name: string): Promise<void>` — `PATCH /api/pod?pod_id=X`
- `DeletePod(podId: number): Promise<void>` — `DELETE /api/pod?pod_id=X`
- `PatchPodPlayerRole(podId: number, playerId: number): Promise<void>` — promote to manager
- `DeletePodPlayer(podId: number, playerId: number): Promise<void>` — manager kick
- `PatchPlayer(playerId: number, name: string): Promise<void>` — `PATCH /api/player?player_id=X`
- `PatchDeck(deckId: number, fields: DeckUpdateFields): Promise<void>` — `PATCH /api/deck?deck_id=X`
- `DeleteDeck(deckId: number): Promise<void>` — `DELETE /api/deck?deck_id=X`
- `PatchGame(gameId: number, description: string): Promise<void>` — `PATCH /api/game?game_id=X`
- `DeleteGame(gameId: number): Promise<void>` — `DELETE /api/game?game_id=X`
- `PostGameResult(result: NewGameResultWithGame): Promise<void>` — `POST /api/game/result`
- `PatchGameResult(resultId: number, fields: GameResultUpdateFields): Promise<void>` — `PATCH /api/game/result?result_id=X`
- `DeleteGameResult(resultId: number): Promise<void>` — `DELETE /api/game/result?result_id=X`

`PaginatedResponse` functions append `?limit=N&offset=M` when those params are provided.

---

### Sub-Phase 3C — AuthContext

**File:** `app/src/auth.tsx` (new). **Depends on 3A, 3B (GetMe + Logout).**

```typescript
interface AuthUser {
    id: number;
    player_id: number;
    display_name: string | null;   // *string in Go → null in JSON when unset
    avatar_url: string | null;     // *string in Go → null in JSON when unset
}

interface AuthContextValue {
    user: AuthUser | null;
    loading: boolean;
    logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue>({ user: null, loading: true, logout: async () => {} });

export function AuthProvider({ children }: { children: ReactNode }): ReactElement {
    const [user, setUser] = useState<AuthUser | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        GetMe()
            .then(setUser)
            .catch(() => setUser(null))
            .finally(() => setLoading(false));
    }, []);

    const logout = async () => { await Logout(); setUser(null); };

    return <AuthContext.Provider value={{ user, loading, logout }}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue { return useContext(AuthContext); }
```

> **Note:** `user.Entity` in `lib/business/user/entity.go` has `DisplayName *string` and `AvatarURL *string` — Go pointers serialize as JSON `null` when nil. Use `string | null`, not `string`.

---

### Sub-Phase 3D — Login Page

**File:** `app/src/routes/login.tsx` (new). **No dependencies — implement any time.**

Centered MUI page: "Sign in with Google" `<Button>` pointing to `/api/auth/google`. No loader.

```typescript
export default function LoginPage(): ReactElement {
    return (
        <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", mt: 8 }}>
            <Typography variant="h4" gutterBottom>EDH Tracker</Typography>
            <Button variant="contained" href="/api/auth/google" size="large">
                Sign in with Google
            </Button>
        </Box>
    );
}
```

---

### Sub-Phase 3E — RequireAuth Wrapper

**File:** `app/src/routes/RequireAuth.tsx` (new). **Depends on 3C.**

```typescript
export default function RequireAuth({ children }: { children: ReactNode }): ReactElement {
    const { user, loading } = useAuth();
    if (loading) return <CircularProgress />;
    if (!user) return <Navigate to="/login" replace />;
    return <>{children}</>;
}
```

---

### Sub-Phase 3F — Route Restructure

**File:** `app/src/index.tsx`. **Depends on 3C (AuthProvider), 3D (LoginPage), 3E (RequireAuth).**

**Phase 4 stub approach:** Routes for `/pod/:podId`, `/join`, etc. reference components that don't exist until Phase 4. Add minimal inline stub components (`function PodView() { return <Typography>Coming soon</Typography>; }`) so the route tree is complete and navigable.

```typescript
// index.tsx — top-level render
createRoot(document.getElementById("root") as HTMLElement).render(
    <StrictMode>
        <CssBaseline enableColorScheme />
        <AuthProvider>
            <RouterProvider router={router} />
        </AuthProvider>
    </StrictMode>
);

const router = createBrowserRouter([
    {
        path: "/",
        element: <Root />,
        errorElement: <ErrorPage />,
        children: [
            { index: true,                   element: <RequireAuth><HomeView /></RequireAuth> },
            { path: "login",                 element: <LoginPage /> },
            { path: "join",                  element: <JoinView /> },           // stub → Phase 4A
            { path: "pod/:podId",            element: <RequireAuth><PodView /></RequireAuth> },     // stub → Phase 4A
            { path: "pod/:podId/new-game",   element: <RequireAuth><NewGameView /></RequireAuth> },
            { path: "pod/:podId/game/:gameId", element: <RequireAuth><GameView /></RequireAuth> },
            { path: "player/:playerId",      element: <RequireAuth><PlayerView /></RequireAuth>, loader: GetPlayer },
            { path: "player/:playerId/deck/:deckId", element: <RequireAuth><DeckView /></RequireAuth>, loader: GetDeck },
            { path: "decks",   element: <Navigate to="/" replace /> },
            { path: "players", element: <Navigate to="/" replace /> },
            { path: "games",   element: <Navigate to="/" replace /> },
        ],
    },
]);
```

`HomeView` (inline stub): calls `GetPodsForPlayer(user.player_id)` → navigates to `/pod/:firstPodId`, or shows "Create your first pod" CTA if none exist.

**Removed routes:** `/deck/:deckId`, `/game/:gameId`, `/new-game` (old top-level form).

---

### Sub-Phase 3G — Navigation Revamp

**File:** `app/src/routes/root.tsx`. **Depends on 3C (useAuth), 3A (Pod type), 3B (GetPodsForPlayer). Parallel with 3E.**

Remove hardcoded nav links (`/players`, `/decks`, `/games`, `/new-game`). Replace placeholder login text:

**Logged out:**
```typescript
<Button href="/api/auth/google">Sign in with Google</Button>
```

**Logged in:**
- Avatar + display name → link to `/player/:playerId`
- Pod selector `<Select>`:
  - Fetches `GetPodsForPlayer(user.player_id)` in `useEffect` on mount
  - Selected value: `podId` from `useParams()` if on a `/pod/*` route, else reads `localStorage.getItem("lastPodId")`
  - On change: write to `localStorage`; navigate to `/pod/:podId`
  - Shows placeholder option if no pods
- Logout button → `logout()` from `useAuth()` → navigate to `/login`

---

## Phase 4A — Pod Landing Page

**Requires Phase 3. Runs in parallel with 4B, 4C, 4D.**

### Route: `/pod/:podId`

**File:** `app/src/routes/pod.tsx`

**Loader:** Parallel fetch of `GetPod(podId)` + `GetPlayersForPod(podId)` + initial page of decks + initial page of games.

**Component structure:**

```
<PodView pod={pod} players={players} currentUserRole={role}>
  <h1>{pod.name}</h1>
  <Tabs value={tab} onChange={setTab}>
    <Tab label="Decks" />
    <Tab label="Players" />
    <Tab label="Games" />
    {isManager && <Tab label="Settings" />}
  </Tabs>
  {tab === 0 && <PodDecksTab podId={pod.id} initialData={decks} />}
  {tab === 1 && <PodPlayersTab players={players} isManager={isManager} podId={pod.id} />}
  {tab === 2 && <PodGamesTab podId={pod.id} initialData={games} />}
  {tab === 3 && isManager && <PodSettingsTab pod={pod} />}
</PodView>
```

**`<PodDecksTab>`:**
- MUI DataGrid, `paginationMode="server"`, `rowCount={total}`
- On `onPaginationModelChange` → `GetDecksForPod(podId, limit, newOffset)` → update rows
- Columns: Commander/Deck (link to `/player/${row.player_id}/deck/${row.id}`), Format, Record, Kills, Points, Games

**`<PodPlayersTab>`:**
- List of players; each row: avatar/name (link to `/player/:playerId`), role badge
- PodManager sees per-row: "Promote" button (if member) → `PatchPodPlayerRole(podId, playerId)` → reload; "Remove" button → `DeletePodPlayer(podId, playerId)` → reload

**`<PodGamesTab>`:**
- MUI DataGrid, `paginationMode="server"`, `rowCount={total}`
- Columns: Game # (link to `/pod/${podId}/game/${row.id}`), Description, Date, Participants
- "New Game" button → navigate to `/pod/${podId}/new-game`

**`<PodSettingsTab>`:**
- Edit pod name: text field + save → `PatchPod(podId, newName)` → reload
- Invite link: "Generate Invite Link" button → `PostPodInvite(podId)` → show `${window.location.origin}/join?code=xxx` with copy button (frontend constructs the full URL — no need for `FRONTEND_URL` on the client)
- Delete pod: button with `<Dialog>` confirmation → `DeletePod(podId)` → navigate to `/`

### Route: `/pod/:podId/new-game`

**File:** `app/src/routes/new.tsx` (updated)

- `pod_id` comes from `useParams().podId` — no dropdown needed
- Loader: `GetPlayersForPod(podId)` + `GetDecksForPod(podId)` (all, unpaginated) + `GetFormats()`
- On submit success: navigate to `/pod/${podId}`
- Remove unused `pod_id: 0` hardcode

### `/join` route

**File:** `app/src/routes/join.tsx`

- Reads `?code=xxx` from URL params
- If not authenticated: redirect to `/api/auth/google?redirect=/join?code=xxx` (bypasses login page entirely; Google auth will redirect back here after login)
- If authenticated: call `PostPodJoin(code)` → navigate to `/pod/${pod.id}`
- If error (expired or invalid code): show error message with link back to `/`

### Redirect old routes

In `app/src/index.tsx`: `/decks`, `/players`, `/games` → `<Navigate to="/" replace />`

---

## Phase 4B — Player Page Revamp

**Requires Phase 3. Runs in parallel with 4A, 4C, 4D.**

### File: `app/src/routes/player.tsx`

**Loader:** Unchanged — `GetPlayer({ params })` returns `Player` entity.

**Component structure:**

```
<PlayerView player={player}>
  <h1>{player.name}</h1>          // label: Display Name
  <Record record={player.stats.record} />
  <Tabs value={tab} onChange={setTab}>
    <Tab label="Overview" />
    <Tab label="Decks" />
    <Tab label="Games" />
    {isOwnProfile && <Tab label="Settings" />}
  </Tabs>
  {tab === 0 && <PlayerOverviewTab player={player} />}
  {tab === 1 && <PlayerDecksTab player={player} />}
  {tab === 2 && <PlayerGamesTab player={player} />}
  {tab === 3 && isOwnProfile && <PlayerSettingsTab player={player} />}
</PlayerView>
```

`isOwnProfile` = `useAuth().user?.player_id === player.id`

**`<PlayerOverviewTab>`:**
- Stats row (games, kills, points)
- "Pods:" list with links to `/pod/:podId` — `AsyncComponentHelper(GetPodsForPlayer(player.id))`
- Created at date

**`<PlayerDecksTab>`:**
- Existing `<DeckDisplay>` DataGrid
- Update `CommanderColumn` (in `stats.tsx`) link: `/player/${row.player_id}/deck/${row.id}`

**`<PlayerGamesTab>`:**
- `AsyncComponentHelper(GetGamesForPlayer(player.id))` → `<MatchesDisplay games={data} />`

**`<PlayerSettingsTab>` (own profile only):**
- **Edit display name**: text field + save → `PatchPlayer(player.id, newName)` → reload player data
- **Your pods**: `AsyncComponentHelper(GetPodsForPlayer(player.id))` → list each pod with link + "Leave Pod" button with `<Dialog>` confirmation → `PostPodLeave(podId)` → reload (note: if caller is the only pod manager, backend returns 403; show error "Promote another member to manager before leaving")
- **Create new pod**: text field + "Create" button → `PostPod(name)` → navigate to `/pod/${newPod.id}`

---

## Phase 4C — Deck Page Revamp

**Requires Phase 3. Runs in parallel with 4A, 4B, 4D.**

### Route: `/player/:playerId/deck/:deckId`

**Loader:** `GetDeck({ params })` — uses `deckId` from params; `player_id` is in the deck entity for validation.

**All existing links to `/deck/:deckId` must be updated:**
- `app/src/stats.tsx` — `CommanderColumn`: `/player/${row.player_id}/deck/${row.id}`
- `app/src/routes/game.tsx` — results grid Deck column: needs `player_id` on `GameResult` (Phase 2C adds it); update to `/player/${row.player_id}/deck/${row.deck_id}`

### Component structure:

```
<DeckView deck={deck}>
  <h1>{deck.name}</h1>
  <h2>{commanders}</h2>             // e.g. "Atraxa, Praetors' Voice" or "Malcolm, Keen-Eyed Navigator / Breeches, Brazen Plunderer"
  <Record record={deck.stats.record} />
  <Tabs value={tab} onChange={setTab}>
    <Tab label="Overview" />
    <Tab label="Games" />
    {isOwner && <Tab label="Settings" />}
  </Tabs>
  {tab === 0 && <DeckOverviewTab deck={deck} />}
  {tab === 1 && <DeckGamesTab deck={deck} />}
  {tab === 2 && isOwner && <DeckSettingsTab deck={deck} />}
</DeckView>
```

`isOwner` = `useAuth().user?.player_id === deck.player_id`

**`<DeckOverviewTab>`:**
- Stats row (games, kills, points)
- Owner: link to `/player/${deck.player_id}`
- Format: `{deck.format_name}`
- Retired badge if `deck.retired`
- Created at date

**`<DeckGamesTab>`:**
- Existing `<MatchUpsForDeck deck={deck} />`
- Update game links to `/pod/${game.pod_id}/game/${game.id}`

**`<DeckSettingsTab>` (owner only):**
- **Edit name**: text field (pre-filled) + save → `PatchDeck(deck.id, {name})`
- **Edit format**: `<Select>` with formats from `GetFormats()` (loaders can pre-fetch) + save → `PatchDeck(deck.id, {format_id})`
- **Edit commanders**: commander autocomplete (current value pre-filled) + optional partner autocomplete; save → `PatchDeck(deck.id, {commander_id, partner_commander_id})`
- **Retire**: button + `<Dialog>` confirmation → `PatchDeck(deck.id, {retired: true})` → navigate to `/player/${deck.player_id}`
- **Delete**: button + `<Dialog>` confirmation → `DeleteDeck(deck.id)` → navigate to `/player/${deck.player_id}`

---

## Phase 4D — Game Page Revamp

**Requires Phase 3. Runs in parallel with 4A, 4B, 4C.**

### Route: `/pod/:podId/game/:gameId`

**Loader:** Parallel fetch of `GetGame({ params })` + `GetPod(podId)`.

**All existing links to `/game/:gameId` must be updated:**
- `app/src/matches.tsx` — `MatchesDisplay` game # column: game entities include `pod_id`; update to `/pod/${row.pod_id}/game/${row.id}`

### Component structure:

```
<GameView game={game} pod={pod} isManager={isManager}>
  <h1>{pod.name} — Game #{game.id}</h1>
  <em>{game.created_at toLocaleString}</em>
  <GameDescription game={game} isManager={isManager} />
  <GameResultsGrid game={game} isManager={isManager} podId={podId} />
  {isManager && <DeleteGameButton gameId={game.id} podId={podId} />}
</GameView>
```

`isManager`: check `useAuth().user?.player_id` against pod members returned by loader (or store role in loader data).

**`<GameDescription>` (editable for PodManager):**
- Display: `Description: {game.description}` (or placeholder if empty)
- If PodManager: edit icon → text field + save/cancel → `PatchGame(game.id, newDesc)` → reload
- If PodManager and no description: show "Add description" placeholder

**`<GameResultsGrid>`:**
Columns: Place, Deck (link to `/player/${row.player_id}/deck/${row.deck_id}`), Commander, Kills, Points
If PodManager: add "Edit" icon column + "Remove" icon column

- **Edit row**: opens `<EditResultModal>` with place, kill count, deck autocomplete (from pod's decks) pre-filled; save → `PatchGameResult(result.id, fields)` → reload
- **Remove row**: `<Dialog>` confirmation → `DeleteGameResult(result.id)` → reload

**"Add Result" button (PodManager):**
- Button below grid
- Opens `<AddResultModal>`: player selector, deck autocomplete, place, kill count
- Save → `PostGameResult({game_id, deck_id, player_id, place, kill_count})` → reload

**"Delete Game" button (PodManager):**
- `<Dialog>` confirmation
- `DeleteGame(game.id)` → navigate to `/pod/${podId}`

## Phase 5 - Unsorted work

A btw session responded with this when asked out handling a user's oauth token expiring mid session:
```
 - Mid-session expiry: If the token expires while the user is actively using the app, useAuth() won't re-check. The
    user state remains set from the initial load. The user will only get bounced to /login when they make an API call
    that returns 401 — but that's in the individual fetch functions in http.ts, which currently just throw new
    Error(...). Nothing catches that throw and redirects to /login.

    So for now: graceful on load, silent failure mid-session. That's acceptable for Phase 3 — the mid-session case would
     be addressed later by adding a 401 interceptor in http.ts (e.g., call logout() from context on any 401 response).
    Not a gap in 3E itself.
```

Sounds like we need to account for this.

---

## Summary: Parallelism Guide for Subagents

| Work Item | Depends On | Can Run Parallel With |
|-----------|------------|----------------------|
| Phase 1 (Migrations) | Nothing | — |
| Phase 2A (OAuth backend) | Phase 1 | 2B, 2C, 2D |
| Phase 2B (Pod roles + invite) | Phase 1 | 2A, 2C, 2D |
| Phase 2C (Edit/delete endpoints) | Phase 1 | 2A, 2B, 2D |
| Phase 2D (Pagination backend) | Phase 1 | 2A, 2B, 2C |
| Phase 3A (Types) | Phase 2A | 3D |
| Phase 3B (HTTP client) | 3A | 3D |
| Phase 3C (AuthContext) | 3B | — |
| Phase 3D (Login page) | Phase 2A | 3A, 3B, 3C |
| Phase 3E (RequireAuth) | 3C | 3G |
| Phase 3F (Route restructure) | 3C, 3D, 3E | — |
| Phase 3G (Navigation revamp) | 3C | 3E |
| Phase 4A (Pod landing page) | Phase 3 | 4B, 4C, 4D |
| Phase 4B (Player page) | Phase 3 | 4A, 4C, 4D |
| Phase 4C (Deck page) | Phase 3 | 4A, 4B, 4D |
| Phase 4D (Game page) | Phase 3 | 4A, 4B, 4C |

---

## Verified Codebase State

The following was confirmed by inspecting the current codebase before planning — no additional work needed for these:

| Concern | Status |
|---------|--------|
| `pod` table `deleted_at` column | Already exists (Migration 8). Migration 19 correctly targets only `game` and `deck`. |
| `game.Entity` `pod_id` field | Already present — `matches.tsx` game link update is straightforward. |
| `deck.Entity` `player_id` field | Already present — no extra backend work needed for deck ownership checks. |
| `pod` repository `GetByPlayerID` | Already exists — `GET /api/pod?player_id=X` endpoint just needs wiring. |
| `PATCH /api/deck` handler | Exists as `RetireDeck` (sets `retired=true` only). Phase 2C extends it to a full update. |
| `gameResult.Entity` `player_id` field | Already present — Phase 2C note can be closed. |
| Repository test infrastructure | GORM integration tests (`testHelpers.NewTestDB`) — `go-sqlmock` no longer used. |

---

## Key Files Reference

| File | Role |
|------|------|
| `lib/migrations/migrate.go` | Register migrations 17, 18, 19 |
| `lib/config.go` | Add Google OAuth + JWT env vars |
| `api.go` | Wire AuthRouter; apply auth middleware |
| `lib/routers/auth.go` | New — OAuth callback + me + logout |
| `lib/trackerHttp/auth.go` | JWT validation middleware (`RequireAuth`, `OptionalAuth`) |
| `lib/utils/context.go` | `ContextWithUserInfo` / `UserFromContext` helpers |
| `lib/repositories/playerPodRole/` | New package |
| `lib/repositories/podInvite/` | New package |
| `lib/repositories/repositories.go` | Add new repos |
| `lib/business/user/functions.go` | Add GetByOAuth, CreateWithOAuth |
| `lib/business/pod/functions.go` | Add role + invite functions |
| `lib/business/game/functions.go` | Add edit/delete/result functions |
| `lib/business/deck/functions.go` | Add Update, SoftDelete |
| `lib/business/player/functions.go` | Add Update |
| `lib/seeder/seeder.go` | Seed player_pod_role; Mike=manager |
| `app/src/auth.tsx` | New — AuthContext + AuthProvider |
| `app/src/routes/login.tsx` | New — login page |
| `app/src/routes/join.tsx` | New — join pod via invite code |
| `app/src/routes/pod.tsx` | New — pod landing page |
| `app/src/routes/root.tsx` | Nav revamp (pod selector, login/logout) |
| `app/src/routes/player.tsx` | Add tabs + settings |
| `app/src/routes/deck.tsx` | Add tabs + settings + route change |
| `app/src/routes/game.tsx` | Add editing + route change |
| `app/src/routes/new.tsx` | Update for pod context |
| `app/src/index.tsx` | Route restructure |
| `app/src/http.ts` | Add all new HTTP functions |
| `app/src/types.ts` | Add Pod, PlayerWithRole, PaginatedResponse |
| `app/src/stats.tsx` | Update CommanderColumn link |
| `app/src/matches.tsx` | Update game link to include pod_id |
