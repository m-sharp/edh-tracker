# External Integrations

**Analysis Date:** 2026-03-22

## APIs & External Services

**Google OAuth2:**
- Service: Google Identity (OAuth 2.0 + OpenID Connect)
- Purpose: User authentication — login, profile fetch, account creation
- SDK/Client: `golang.org/x/oauth2` + `golang.org/x/oauth2/google` (`lib/routers/auth.go`)
- Scopes requested: `openid`, `userinfo.email`, `userinfo.profile`
- Userinfo endpoint called directly: `https://www.googleapis.com/oauth2/v3/userinfo`
- Auth: `GOOGLE_CLIENT_ID` + `GOOGLE_CLIENT_SECRET` env vars
- Redirect: `OAUTH_REDIRECT_URL` env var (must match Google Console configuration)
- Implementation: `lib/routers/auth.go` — `AuthRouter` with `Login`, `Callback`, `Logout`, `Me` handlers

## Data Storage

**Databases:**
- Type: MySQL
- Database name: `pod_tracker` (constant in `lib/db.go`)
- Connection: `DBHOST`, `DBUSER`, `DBPASSWORD`, `DBPORT` env vars
- Client: GORM (`gorm.io/gorm` + `gorm.io/driver/mysql`) with underlying `database/sql` connection pool
- Connection pool: max 10 open connections, max 10 idle, 2-minute TTL (`lib/db.go`)
- Migrations: auto-run on startup via `lib/migrations/` (numbered sequential migrations)
- Local dev DB: `mysql/Dockerfile` with `mysql:latest` image, database `pod_tracker`

**File Storage:**
- Not applicable — no file upload or object storage integration

**Caching:**
- None — no Redis or in-memory cache layer

## Authentication & Identity

**Auth Provider: Google OAuth2**
- Flow: Authorization Code flow with CSRF nonce (`lib/routers/auth.go`)
- CSRF protection: random 16-byte hex nonce stored in `edh_csrf` short-lived cookie (5-minute TTL); validated against OAuth `state` parameter on callback
- Session: JWT stored in `edh_session` HTTP-only cookie (24-hour sliding window)
- JWT algorithm: HS256 signed with `JWT_SECRET` env var
- JWT claims: `user_id` (int), `player_id` (int), `exp`, `iat`
- Implementation: `lib/trackerHttp/auth.go` — `RequireAuth` and `OptionalAuth` middleware; `IssueJWT`
- Secure cookies: enabled in production; disabled when `DEV` env var is set
- New user flow: Google profile fetched → player + user rows inserted atomically in a single transaction (`lib/repositories/user/repo.go` `CreatePlayerAndUser`)
- Seeded player linking: if a user row exists for a matching email but no OAuth subject, `LinkOAuth` is called to associate credentials (`lib/routers/auth.go` `Callback`)

**Route-Level Auth Flags:**
- `Route.RequireAuth = true` — enforces `RequireAuth` middleware (401 if no valid JWT)
- `Route.NoAuth = true` — opts out of automatic auth on state-changing routes
- Applied per-route in `api.go` `SetupRoutes`

## Monitoring & Observability

**Error Tracking:**
- None — no Sentry or external error tracking integration

**Logs:**
- Uber Zap structured logger (`go.uber.org/zap`)
- Logger instantiated in `main.go` via `lib.GetLogger(cfg)` and passed through all layers
- Named loggers per component (e.g., `log.Named("AuthRouter")`, `log.Named("DBClient")`)
- GORM query logging at `Warn` level; `ErrRecordNotFound` suppressed via `quietLogger` wrapper (`lib/db.go`)

## CI/CD & Deployment

**Hosting:**
- Docker containers (three images: API server, React/static app server, MySQL)
- No cloud provider detected in codebase

**CI Pipeline:**
- None detected — no `.github/workflows/`, CircleCI, or similar config files found

## Environment Configuration

**Required env vars (API server):**
- `DBHOST` — MySQL host
- `DBUSER` — MySQL username
- `DBPASSWORD` — MySQL password
- `DBPORT` — MySQL port
- `GOOGLE_CLIENT_ID` — Google OAuth2 client ID
- `GOOGLE_CLIENT_SECRET` — Google OAuth2 client secret
- `OAUTH_REDIRECT_URL` — OAuth callback URL (must match Google Console)
- `JWT_SECRET` — HMAC secret for signing session JWTs
- `FRONTEND_URL` — Frontend origin; used for CORS and post-auth redirects

**Optional env vars:**
- `DEV` — any non-empty value enables development mode (disables secure flag on cookies)
- `SEED` — any non-empty value triggers `lib/seeder/seeder.go` on startup

**Secrets location:**
- All secrets injected as environment variables at container runtime; no secrets committed to the repository

## Webhooks & Callbacks

**Incoming:**
- `GET /api/auth/google/callback` — Google OAuth2 redirect callback; handled by `AuthRouter.Callback` in `lib/routers/auth.go`

**Outgoing:**
- None — no outgoing webhooks to external services

---

*Integration audit: 2026-03-22*
