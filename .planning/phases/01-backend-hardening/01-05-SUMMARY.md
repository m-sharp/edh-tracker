---
phase: 01-backend-hardening
plan: 05
subsystem: backend
tags: [security, validation, auth, jwt, invite]
dependency_graph:
  requires: []
  provides: [JWT-length-guard, input-max-length-validation, invite-use-limit]
  affects: [main.go, lib/business/player, lib/business/pod, lib/business/deck, lib/routers/game]
tech_stack:
  added: []
  patterns: [startup-guard, entity-validate, router-inline-validation]
key_files:
  created: []
  modified:
    - main.go
    - lib/business/player/entity.go
    - lib/business/pod/entity.go
    - lib/business/deck/entity.go
    - lib/routers/game.go
    - lib/business/pod/functions.go
decisions:
  - "maxInviteUses hardcoded to 25 — pod_invite table has no max_used_count column; constant placed in functions.go with a comment explaining why"
metrics:
  duration: "3 minutes"
  completed: "2026-03-23"
  tasks_completed: 3
  files_modified: 6
---

# Phase 01 Plan 05: Security Guards and Input Validation Summary

JWT startup guard, string field max-length validation, and invite use-count limit to prevent weak secrets, DB varchar overflow, and unlimited invite reuse.

## Tasks Completed

| # | Name | Commit | Files |
|---|------|--------|-------|
| 1 | Add JWT secret length guard at startup (AUTH-02) | a05f373 | main.go |
| 2 | Add max length validation to entity Validate methods and game description (SEC-05) | d6ee675 | lib/business/player/entity.go, lib/business/pod/entity.go, lib/business/deck/entity.go, lib/routers/game.go |
| 3 | Add max use count check to JoinByInvite (SEC-04) | aa4048e | lib/business/pod/functions.go |

## What Was Built

**Task 1 — JWT Startup Guard (AUTH-02):**
After `lib.NewConfig()` succeeds, `main.go` now reads the `JWT_SECRET` value and calls `log.Fatalf` if `len(jwtSecret) < 32`. This prevents the server from starting with a weak secret. The error message includes the actual byte count to aid debugging.

**Task 2 — String Field Max-Length Validation (SEC-05):**
Added length checks to four locations matching DB VARCHAR column sizes:
- `lib/business/player/entity.go` `Validate()`: player name capped at 256 chars
- `lib/business/pod/entity.go` `Entity.Validate()`: pod name capped at 255 chars
- `lib/business/pod/entity.go` `UpdatePodInputEntity.Validate()`: pod name update capped at 255 chars
- `lib/business/deck/entity.go` `ValidateCreate()`: deck name capped at 255 chars
- `lib/routers/game.go` `GameCreate`: game description capped at 256 chars (inline in router since no `game.Entity.Validate()` exists)

All return 400 with descriptive messages.

**Task 3 — Invite Max Use Count (SEC-04):**
Added `const maxInviteUses = 25` to `lib/business/pod/functions.go` (package-level, with explanatory comment). `JoinByInvite` now checks `invite.UsedCount >= maxInviteUses` after the nil/expiry checks and before `AddPlayerToPod`, returning a descriptive error if the limit is reached.

## Decisions Made

| Decision | Rationale |
|----------|-----------|
| maxInviteUses hardcoded to 25 | The `pod_invite` table has no `max_used_count` column; per plan research a hardcoded constant is the correct approach; 25 is sufficient for any realistic friend-group pod |

## Deviations from Plan

None - plan executed exactly as written.

## Known Stubs

None.

## Self-Check: PASSED

Files exist:
- main.go: FOUND
- lib/business/player/entity.go: FOUND
- lib/business/pod/entity.go: FOUND
- lib/business/deck/entity.go: FOUND
- lib/routers/game.go: FOUND
- lib/business/pod/functions.go: FOUND

Commits:
- a05f373: FOUND
- d6ee675: FOUND
- aa4048e: FOUND
