# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

EDH Tracker is a Magic: The Gathering Commander (EDH) game tracking app. It tracks players, decks, and game results with a points system (kills + place-based bonuses).

## Tech Stack

- **Backend**: Go 1.26, Gorilla Mux, Uber Zap logger, GORM
- **Database**: MySQL (`pod_tracker` database)
- **Frontend**: React 18 + TypeScript, React Router v6, Material-UI (MUI) v5
- **Deployment**: Docker (separate images for API, React app, and MySQL)

## Commands

### Backend
```bash
go mod vendor      # Vendor dependencies
go run main.go     # Run API server locally (requires DB env vars)
go vet ./lib/...   # Compile check (no binary output)
go test ./lib/...  # Run tests
```

### Frontend (from `app/`)
```bash
npm start        # Dev server
npm run build    # Production build
npm test         # Run tests
```

### Docker
```bash
# Build images
docker build -t edh-tracker .
docker build -f app/Dockerfile -t edh-tracker-app .

# Run API (requires DB + auth env vars)
docker run -p 8080:8081 \
  --env DBHOST=host.docker.internal \
  --env DBUSER=root \
  --env DBPASSWORD=<pass> \
  --env DBPORT=3306 \
  --env GOOGLE_CLIENT_ID=<client_id> \
  --env GOOGLE_CLIENT_SECRET=<client_secret> \
  --env OAUTH_REDIRECT_URL=<redirect_url> \
  --env JWT_SECRET=<secret> \
  --env FRONTEND_URL=http://localhost:8081 \
  --env DEV=1 edh-tracker

# Run React web app
docker run -p 8081:8081 -it edh-tracker-app:latest
```

## Architecture

See [`docs/architecture.md`](docs/architecture.md) for full detail on the 4-layer backend (routers → business → repositories → DB), functional DI pattern, and frontend structure.

### Quick Summary

```
lib/routers/       HTTP handlers
lib/business/      Domain logic + entity construction (functional DI closures)
lib/repositories/  Pure DB access returning Model types
lib/migrations/    Auto-run numbered schema migrations
lib/seeder/        Optional seed data (triggered by SEED env var)
```

`api.go` wires all routes. `lib/config.go` reads env vars.

**Required env vars**: `DBHOST`, `DBUSER`, `DBPASSWORD`, `DBPORT`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_REDIRECT_URL`, `JWT_SECRET`, `FRONTEND_URL`

**Optional env vars**: `SEED` (triggers data seeder), `DEV` (development mode — disables secure cookies)

## Data Model

See [`docs/data-model.md`](docs/data-model.md) for entity relationships, table schemas, Model/Entity split, points formula, and format validation rules.

## API Routes

See [`docs/api.md`](docs/api.md) for the full route table.

Route pattern: **plural path for GET-all, singular path for GET-one and POST** (e.g. `GET /api/players` vs `POST /api/player`).

## After Making API Changes

After modifying routes, handlers, repositories, or migrations, run the `/smoke-test` skill to rebuild the Docker image and verify that core endpoints are responding correctly.
