# Phase 1: Backend Hardening - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-22
**Phase:** 01-backend-hardening
**Areas discussed:** Transaction pattern, Auth check placement, Error code discrimination

---

## Transaction Pattern

| Option | Description | Selected |
|--------|-------------|----------|
| GORM db.Transaction() in business layer | Create closure calls db.Transaction() directly, passes tx-scoped DBClient to both repos. Requires exposing DB handle to business layer. | |
| New transactional repo method | Single gameRepo.AddWithResults() method handles both inserts atomically. Hides transaction from business layer. | |
| Context-based tx propagation | Store transaction in context.Context, repos check for it. Most flexible but most complex. | |

**Revisited — user wanted to see code examples before deciding.**

| Option | Description | Selected |
|--------|-------------|----------|
| A: Context-based | lib/tx.go helpers; every repo method gains TxFromContext check boilerplate. | |
| B: Inline tx-scoped repo copies | No context changes, no repo changes; business layer creates tx-scoped repo instances inline via NewRepository(&lib.DBClient{Db: tx}). | |
| C: Pass dbClient to Create constructor only | Same as B. | ✓ |

**User's choice:** Pass `*lib.DBClient` to `game.Create` constructor only; create tx-scoped repo copies inline inside `db.Transaction()` callback. No context changes, no repo method changes.
**Notes:** User initially suggested context-based approach but after seeing concrete code examples, chose the simpler option. Concern was that context injection would be "cumbersome code-wise."

---

## Auth Check Placement

| Option | Description | Selected |
|--------|-------------|----------|
| Router layer (consistent with requirePodManager) | GameCreate calls CallerPlayerID + getPodRole before delegating. Business layer stays pure. | |
| Business layer (consistent with assertCallerOwnsDeck) | game.Create takes callerPlayerID, checks membership internally. Auth co-located with operation. | |
| Standardize both patterns on one approach | Pick one pattern, migrate the other. More work but eliminates inconsistency. | ✓ |

**Follow-up — which layer to standardize on:**

| Option | Description | Selected |
|--------|-------------|----------|
| Business layer | Auth checks alongside operations; deck pattern becomes standard; requirePodManager migrated. | |
| Router layer | Routers own all auth; business stays pure; assertCallerOwnsDeck moves up to deck router. | ✓ |

**Follow-up — scope of migration:**

| Option | Description | Selected |
|--------|-------------|----------|
| Phase 1: fix game create + flag deck inconsistency | Add SEC-01 fix now, note assertCallerOwnsDeck as future cleanup. | |
| Phase 1: fix all inconsistencies now | Move assertCallerOwnsDeck to deck router in this phase. | ✓ |

**User's choice:** Router layer owns all auth; standardize in Phase 1 — includes migrating assertCallerOwnsDeck from deck business layer to deck router.

---

## Error Code Discrimination

| Option | Description | Selected |
|--------|-------------|----------|
| Sentinel string prefix check | strings.HasPrefix(err.Error(), "forbidden:") → 403. Simple, no new types. | |
| Typed sentinel error in business layer | var ErrForbidden = errors.New("forbidden"); errors.Is() in routers. Go-idiomatic. | ✓ |
| errors.As with custom ForbiddenError type | Dedicated error type with reason field. Most expressive but most overhead. | |

**Follow-up — scope:**

| Option | Description | Selected |
|--------|-------------|----------|
| lib/business package, applied to all existing forbidden cases | ErrForbidden in shared location; all existing "forbidden: ..." calls updated to wrap it. | ✓ |
| Pod package only, just fix PromotePlayer/KickPlayer | Narrow fix, other "forbidden:" strings stay as-is. | |

**User's choice:** `var ErrForbidden` in `lib/business` package, wrapped by all existing forbidden errors throughout the business layer. Routers use `errors.Is(err, business.ErrForbidden)` consistently.

---

## Claude's Discretion

- Whether `ErrForbidden` lives in `errors.go` or inline in another file
- Whether to remove `GetAll` repo method after PERF-02 or just block the route
- Whether `assertCallerOwnsDeck` becomes a router-level helper or is inlined

## Deferred Ideas

None.
