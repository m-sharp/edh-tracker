---
created: 2026-03-26T02:11:05.421Z
title: Handle API_BASE_URL for production in http.ts
area: api
files:
  - app/src/http.ts:21
---

## Problem

`app/src/http.ts:21` has `export const API_BASE_URL = "http://localhost:8080"` with a `// TODO: How does this look in production?` comment. In production (Docker), the React app is served by a Go static file server on port 8081, and the API runs on a separate container. The hardcoded localhost URL will break in any non-local deployment.

## Solution

Determine the correct approach for production API URL resolution. Options to consider:
- Environment variable injected at build time via `REACT_APP_API_BASE_URL` (Create React App convention)
- Relative URL (empty string `""`) so API calls go to the same origin — only works if frontend and API are behind the same reverse proxy
- Nginx/proxy config that routes `/api/*` to the backend container
- Runtime config endpoint

TBD — depends on final deployment topology (same origin vs. separate domains).
