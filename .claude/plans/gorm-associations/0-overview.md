# GORM Associations — Overview

## Goal

Eliminate N+1 query patterns in the EDH Tracker API by adding GORM `Preload` associations
to repository models. New hydrated methods are added alongside existing flat methods
(additive — no existing callers broken). The business layer is simplified as closure deps
for per-record DB lookups are replaced by preloaded data.

## Background

All repositories were migrated to GORM in the harp/gorm-implementation branch (Phases 1–5).
There are no circular import risks between repos — the cross-domain dependency graph is
acyclic (see below). This work can proceed entirely within the existing architecture.

## Constraints

- No schema changes
- No `AutoMigrate`
- No Entity layer changes beyond what preloading directly enables
- New repo methods are additive — existing flat methods remain

## Approach

1. Add association fields to `Model` structs (zero-cost until a `Preload` call is made)
2. Add `Get*WithXxx` repo methods that use `Preload`
3. Add signatures to `lib/repositories/interfaces.go`
4. Update business constructors to use the new methods and drop per-record closure deps
5. Simplify `lib/business/business.go` wiring

## Dependency Graph

```
Phase 1 — game-results     (standalone — no inter-phase deps)
Phase 2 — deck-associations (standalone — no inter-phase deps)
Phase 3 — game-result-deck  (depends on Phase 1 AND Phase 2)
Phase 4 — pod-members       (standalone — no deps)
```

Phases 1 and 2 can be implemented in parallel. Phase 3 must follow both. Phase 4 is
independent and can be done at any time.

## Aggregate Query Reduction

| Endpoint / Function | Before | After |
|---------------------|--------|-------|
| GetAllByPod (N games, P players) | 1 + N + N×P×(2–3) | ~5 queries |
| GetAllByDeck | 1 + N + N×P×(2–3) | ~5 queries |
| GetAllByPlayer | 1 + N + N×P×(2–3) | ~5 queries |
| GetByID (game) | 2 + P×(2–3) | ~5 queries |
| GetAll (decks, N decks) | N×(4–5) | ~5 batched |
| GetAllForPlayer (decks) | N×(3–4) | ~5 batched |
| GetAllByPod (decks) | N×(4–5) | ~5 batched |
| GetByID (deck) | 4–5 queries | ~5 queries |

## Skill

Use the `/gorm` skill tool at the start of each implementation session for GORM patterns,
import paths, model conventions, and test infrastructure.
