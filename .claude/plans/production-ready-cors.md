# Production-Ready CORS Plan

**Status: Discovery / Ideation — not scheduled. Revisit after frontend-revamp lands.**

---

## Problem

The current CORS setup in `lib/http.go` uses backend middleware to add `Access-Control-Allow-*` headers to every response. This works as a dev workaround but is not the right long-term architecture:

- Backend CORS headers are a symptom of the frontend and API running on different origins
- The cleanest fix is to eliminate the cross-origin situation entirely rather than papering over it
- Even if CORS headers are kept, the current implementation has bugs (see `frontend-revamp-plan.md` Phase 2A for the immediate bug fixes)

---

## Context

The app currently runs as two separate servers:
- **API**: Go server on port `8081`
- **Frontend**: React dev server (port `3000` in dev) or nginx (port `8081` in Docker)

In development the two origins differ, causing CORS preflight requests. The backend works around this by injecting `Access-Control-Allow-Origin: *` on every response.

The frontend uses HttpOnly cookies for session management (post OAuth revamp). This makes the wildcard origin approach not just hacky — it is **broken** for cookies. The `* + credentials: true` combination is rejected by browsers, so proper CORS or same-origin is required regardless.

---

## Recommended Approaches

### Option 1 (Dev): React Dev Server Proxy

Add a proxy entry to `app/package.json` (Create React App):

```json
"proxy": "http://localhost:8081"
```

Or in `app/vite.config.ts` if Vite is ever adopted:

```ts
server: {
  proxy: {
    '/api': 'http://localhost:8081'
  }
}
```

The React dev server forwards all `/api` requests to the Go backend. From the browser's perspective, everything comes from `localhost:3000` — no cross-origin requests, no CORS headers needed at all.

**Pros:** Zero backend changes required; cookies work naturally; no CORS configuration to maintain.
**Cons:** Dev-only; doesn't address production.

### Option 2 (Production): Reverse Proxy (nginx or Caddy)

Serve both the frontend static files and the API from the same origin behind a reverse proxy.

Example nginx config:

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # Serve frontend static files
    location / {
        root /usr/share/nginx/html;
        try_files $uri /index.html;
    }

    # Proxy API to Go backend
    location /api/ {
        proxy_pass http://api:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

Same origin → no CORS headers needed at all. Cookies (SameSite=Lax, HttpOnly) work without `Access-Control-Allow-Credentials`.

**Pros:** Industry-standard; eliminates CORS entirely; simplifies cookie security model.
**Cons:** Requires nginx (or equivalent) in the Docker Compose setup; adds a deployment artifact.

### Option 3: Remove CORS middleware entirely (after proxy is in place)

Once Option 1 (dev proxy) and Option 2 (production reverse proxy) are in place, `CORSMiddleware` and `CORSPreflightHandler` in `lib/http.go` can be deleted along with the route registration in `api.go`. This is the end state.

---

## Questions to Resolve Before Implementing

1. **Does the Docker Compose setup currently include nginx?** If not, adding it is the main work item.
2. **Will there ever be a legitimate need for external/third-party clients to call the API directly?** If yes, CORS headers are still needed for those callers — but scoped to known origins, not `*`.
3. **CRA vs Vite**: The proxy config syntax differs. Confirm which bundler is in use before touching `app/package.json`.

---

## Work Items (not yet scheduled)

- [ ] Confirm Docker Compose topology (does nginx exist?)
- [ ] Add `"proxy"` field to `app/package.json` for local dev
- [ ] Add nginx service + config to Docker Compose (or confirm existing setup handles it)
- [ ] Remove `CORSMiddleware` and `CORSPreflightHandler` from `lib/http.go` and deregister from `api.go`
- [ ] Verify cookies work end-to-end in both dev (proxy) and production (nginx) environments
