---
phase: 01-backend-hardening
plan: "06"
subsystem: ui
tags: [typescript, react, http, error-handling, fetch]

requires: []
provides:
  - res.ok guards on all fetch calls in http.ts that consume JSON bodies
  - throw on non-2xx for all fire-and-forget mutation functions
  - try/catch with inline error display in player.tsx handleCreatePod
  - try/catch with inline error display in pod.tsx handleSaveName and handleGenerateInvite
affects: [frontend, pod, player]

tech-stack:
  added: []
  patterns:
    - "res.ok guard before res.json() on all GET and JSON-returning POST functions"
    - "Capture response from fire-and-forget mutations; throw Error on non-2xx"
    - "Inline error state (string | null) displayed via Typography color=error variant=body2"

key-files:
  created: []
  modified:
    - app/src/http.ts
    - app/src/routes/player.tsx
    - app/src/routes/pod.tsx

key-decisions:
  - "PostGame and PostCommander intentionally left unchanged — they return raw Response for caller-side handling"
  - "PostPodLeave left unchanged — it already had res.ok guard with status attachment for 403 discrimination"
  - "Error messages are simple status-code strings; no toast/snackbar infrastructure added"

patterns-established:
  - "http.ts pattern: every function either (a) guards res.ok before res.json(), (b) throws on non-2xx for mutations, or (c) returns raw Response for caller handling"
  - "Route error pattern: mutation handler state = [value, setter] with null default; clear at start, set in catch; render Typography color=error variant=body2"

requirements-completed: [SEC-04, SEC-05]

duration: 4min
completed: "2026-03-23"
---

# Phase 01 Plan 06: Frontend HTTP Error Handling Summary

**res.ok guards added to all http.ts fetch calls and try/catch with inline error display added to mutation handlers in player.tsx and pod.tsx**

## Performance

- **Duration:** 4 min
- **Started:** 2026-03-23T04:05:59Z
- **Completed:** 2026-03-23T04:10:53Z
- **Tasks:** 2
- **Files modified:** 3

## Accomplishments
- Every function in http.ts that calls `.json()` now checks `res.ok` first, preventing `SyntaxError: Unexpected token '<'` crashes when the backend returns a text/plain error body
- All fire-and-forget mutation functions now capture the response and throw on non-2xx, making upstream catch blocks active (they were dead code before)
- `handleCreatePod` in player.tsx has try/catch with `createPodError` state and inline error display
- `handleSaveName` and `handleGenerateInvite` in pod.tsx have try/catch with `nameError`/`inviteError` state and inline error display

## Task Commits

Each task was committed atomically:

1. **Task 1: Add res.ok guards and throw on non-2xx in http.ts** - `f788959` (fix)
2. **Task 2: Add missing try/catch in player.tsx and pod.tsx route handlers** - `a037a86` (fix)

## Files Created/Modified
- `app/src/http.ts` - Added res.ok guards to all GET/JSON-returning POST functions; converted fire-and-forget mutations to capture response and throw on non-2xx
- `app/src/routes/player.tsx` - Added createPodError state; wrapped handleCreatePod in try/catch; added inline error display below Create button
- `app/src/routes/pod.tsx` - Added nameError and inviteError state; wrapped handleSaveName and handleGenerateInvite in try/catch; added inline error display below respective buttons

## Decisions Made
- PostGame and PostCommander left unchanged — they return raw Response intentionally so callers can inspect status themselves
- PostPodLeave left unchanged — it already had a res.ok guard with `.status` attachment needed for 403 vs 5xx discrimination in handleLeave
- No toast/snackbar infrastructure added — inline error text only, matching existing codebase pattern

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Frontend no longer white-screens on 400/403/500 responses from backend
- String length validation errors (400) from backend will now display inline messages rather than crashing
- deck.tsx mutation catch blocks are now active code (PatchDeck/DeleteDeck now throw)

---
*Phase: 01-backend-hardening*
*Completed: 2026-03-23*
