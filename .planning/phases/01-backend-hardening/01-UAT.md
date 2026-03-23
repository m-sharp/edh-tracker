---
status: diagnosed
phase: 01-backend-hardening
source: 01-01-SUMMARY.md, 01-02-SUMMARY.md, 01-03-SUMMARY.md, 01-04-SUMMARY.md, 01-05-SUMMARY.md
started: 2026-03-22T04:00:00Z
updated: 2026-03-22T04:30:00Z
---

## Current Test

## Current Test

[testing complete]

## Tests

### 1. Cold Start Smoke Test
expected: Kill any running server/service. Clear ephemeral state. Start the API from scratch (e.g., `go run main.go` with valid env vars). Server boots without errors, any migrations complete, and a basic API call (e.g., GET /api/players) returns a valid response.
result: pass

### 2. Non-Pod-Member Cannot Create Game
expected: Attempt to POST /api/game as a player who is NOT a member of the target pod. The API should return 403 Forbidden. A player who IS a pod member can still create a game successfully (201).
result: pass

### 3. 403 vs 500 Discrimination on Pod Actions
expected: When a non-manager tries to promote/kick a pod member, they get 403 Forbidden (not 500). When a DB error occurs on those same endpoints, 500 is returned instead. The error codes correctly distinguish authorization failures from infrastructure failures.
result: pass

### 4. Deck Create Uses JWT Identity (Not Body player_id)
expected: POST /api/deck with a body that includes a `player_id` field for a different player should create the deck owned by the authenticated caller, not the player_id in the body. The body player_id is silently ignored — the JWT identity wins.
result: pass

### 5. Deck Update/Delete Blocked for Non-Owners
expected: Attempt to PATCH or DELETE a deck that belongs to a different player. The API should return 403 Forbidden. The actual owner can still update/delete their deck successfully.
result: pass

### 6. Unfiltered Deck Endpoint Returns 400
expected: GET /api/decks with no query parameters returns 400 Bad Request with a message like "pod_id or player_id query param is required". GET /api/decks?pod_id=1 (or player_id=1) returns 200 with results.
result: pass

### 7. String Field Max-Length Validation
expected: Submitting a player name longer than 256 characters, a pod name longer than 255 characters, or a deck name longer than 255 characters returns 400 Bad Request. A game description longer than 256 chars also returns 400. Inputs at or under the limit succeed normally.
result: issue
reported: "Server raises an error as expected - trying to send a name that is too long via the frontend edit results in a white screen and a console error of `Uncaught SyntaxError: Unexpected token '<'`. This will need to be captured in some later phase - appropriate error handling of requests within the frontend."
severity: major

### 8. Pod Invite Use Limit
expected: Attempting to join a pod using an invite code that has already been used 25 times returns an error (400 or 403) with a descriptive message. An invite used fewer than 25 times still works normally.
result: pass

## Summary

total: 8
passed: 7
issues: 1
pending: 0
skipped: 0
blocked: 0

## Gaps

- truth: "Submitting a field value that exceeds max length (e.g. player/pod/deck name) returns 400 and the frontend handles it gracefully — no crash or white screen"
  status: failed
  reason: "User reported: Server raises an error as expected - trying to send a name that is too long via the frontend edit results in a white screen and a console error of `Uncaught SyntaxError: Unexpected token '<'`. This will need to be captured in some later phase - appropriate error handling of requests within the frontend."
  severity: major
  test: 7
  root_cause: "Backend WriteError returns text/plain bodies; http.ts calls .json() unconditionally on responses without checking res.ok first, causing SyntaxError on non-2xx responses. Additionally, all fire-and-forget mutation functions (PatchPlayer, PatchDeck, PatchPod, etc.) never check res.ok or throw, making their call sites' try/catch blocks dead code."
  artifacts:
    - path: "app/src/http.ts"
      issue: "PostPod, PostPodJoin, PostPodInvite call res.json() with no res.ok guard; all PATCH/DELETE functions discard the Response entirely without status checks"
    - path: "app/src/routes/player.tsx"
      issue: "handleCreatePod has no try/catch; handleSaveName catch block is dead code because PatchPlayer never throws"
    - path: "app/src/routes/pod.tsx"
      issue: "handleSaveName and handleGenerateInvite propagate errors uncaught to React error boundary"
    - path: "app/src/routes/deck.tsx"
      issue: "All mutation handlers (save name/format/commanders, retire, delete) have dead catch blocks because PatchDeck/DeleteDeck never throw"
    - path: "lib/trackerHttp/http.go"
      issue: "WriteError uses http.Error() which sets Content-Type: text/plain — inconsistent with frontend expecting JSON"
  missing:
    - "app/src/http.ts: add res.ok guards throwing before .json() calls in PostPod, PostPodJoin, PostPodInvite"
    - "app/src/http.ts: all fire-and-forget mutations must check res.ok and throw on non-2xx"
    - "app/src/routes/player.tsx: wrap handleCreatePod in try/catch"
    - "app/src/routes/pod.tsx: wrap handleSaveName and handleGenerateInvite in try/catch"
    - "Optional: make WriteError return JSON body for consistent frontend parsing"
