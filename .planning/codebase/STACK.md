# Technology Stack

**Analysis Date:** 2026-03-22

## Languages

**Primary:**
- Go 1.26 - Backend API server (`main.go`, `api.go`, `lib/`)
- TypeScript 4.9.5 - Frontend React app (`app/src/`)

**Secondary:**
- SQL - Database schema migrations (`lib/migrations/`)

## Runtime

**Backend Environment:**
- Go 1.26.0 (alpine-based Docker image `golang:1.26.0-alpine3.23`)

**Frontend Environment:**
- Node.js 18 (alpine-based Docker image `node:18-alpine`) — build only; runtime is a Go static file server

**Package Manager:**
- Backend: Go modules with vendoring (`go.mod`, `go.sum`, `vendor/`)
- Frontend: npm (`app/package.json`, `app/package-lock.json`)
- Lockfile: `app/package-lock.json` present

## Frameworks

**Backend Core:**
- `gorilla/mux` v1.8.1 - HTTP router and URL parameter extraction
- `gorm.io/gorm` v1.31.1 - ORM for MySQL access (`lib/db.go`, `lib/repositories/`)
- `gorm.io/driver/mysql` v1.6.0 - GORM MySQL driver

**Frontend Core:**
- React 18.0.0 - UI component framework (`app/src/`)
- React Router DOM 6.21.1 - Client-side routing
- Material-UI (MUI) v5.15.2 + `@mui/icons-material` v5.15.3 - UI component library
- `@mui/x-data-grid` v6.18.6 - Data grid component
- `@emotion/react` v11.11.3 + `@emotion/styled` v11.11.0 - CSS-in-JS (MUI peer dependency)

**Frontend Build:**
- `react-scripts` v5.0.1 (Create React App) - Dev server, production build, test runner

**Testing:**
- `stretchr/testify` v1.11.1 - Assertions and `require` in Go tests (`assert`, `require` packages)

## Key Dependencies

**Critical:**
- `go-sql-driver/mysql` v1.8.1 - Low-level MySQL driver (used through GORM); `lib/db.go`
- `golang-jwt/jwt/v5` v5.3.1 - JWT signing and validation; `lib/trackerHttp/auth.go`
- `golang.org/x/oauth2` v0.35.0 - OAuth2 client; `lib/routers/auth.go`
- `google/uuid` v1.6.0 - UUID generation (invite codes, nonces)
- `go.uber.org/zap` v1.26.0 - Structured logging throughout all layers

**Infrastructure:**
- `cloud.google.com/go/compute/metadata` v0.3.0 - Indirect; pulled in by `golang.org/x/oauth2` for Google endpoint support

## Configuration

**Environment:**
- All config read from OS environment variables at startup via `lib/config.go`
- Required vars: `DBHOST`, `DBUSER`, `DBPASSWORD`, `DBPORT`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_REDIRECT_URL`, `JWT_SECRET`, `FRONTEND_URL`
- Optional vars: `DEV` (disables secure cookies in development), `SEED` (triggers data seeder on startup)
- Config struct: `lib.Config` with typed key constants (e.g., `lib.DBHost`, `lib.JWTSecret`)

**Build:**
- Backend: `Dockerfile` at project root — multi-stage not used, copies `vendor/` for reproducible builds
- Frontend + static server: `app/Dockerfile` — multi-stage: React build via `node:18-alpine`, Go static server via `golang:1.26.0-alpine3.23`
- MySQL: `mysql/Dockerfile` — `mysql:latest` base image with `mysql/custom.cnf`

## Platform Requirements

**Development:**
- Go 1.26+
- Node.js 18+
- MySQL instance (or Docker)
- All required env vars set

**Production:**
- Docker (three separate images: API, React app, MySQL)
- API listens on port 8081 (mapped to 8080 externally in example runs)
- Frontend static server listens on port 8081

---

*Stack analysis: 2026-03-22*
