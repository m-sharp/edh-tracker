---
phase: 05-pod-deck-ux
verified: 2026-03-26T00:00:00Z
status: passed
score: 14/14 must-haves verified
re_verification:
  previous_status: passed
  previous_score: 14/14
  gaps_closed: []
  gaps_remaining: []
  regressions:
    - "Previous report incorrectly cited player/DecksTab.tsx line 32 as a visibleRows filter — the actual implementation uses DataGrid initialState.filter.filterModel (gap-closure Plan 08 superseded the original approach); behavior is correct"
---

# Phase 5: Pod/Deck UX Verification Report

**Phase Goal:** Improve Pod and Deck UX — fix creation flows, onboarding gaps, and UX friction identified during initial testing
**Verified:** 2026-03-26
**Status:** passed
**Re-verification:** Yes — independent check against actual codebase (previous: 2026-03-25, passed)

## Goal Achievement

### Observable Truths

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | POST /api/pod returns 201 with JSON body containing new pod ID | VERIFIED | `lib/routers/pod.go` lines 186-190: `Content-Type: application/json`, `WriteHeader(201)`, `json.NewEncoder(w).Encode` with `{ID: podID}` |
| 2 | POST /api/deck returns 201 with JSON body containing new deck ID | VERIFIED | `lib/routers/deck.go` lines 204-208: same pattern with `deckID` captured from `d.decks.Create(...)` |
| 3 | GET /api/players?pod_id=N returns pod-scoped stats | VERIFIED | `lib/business/player/functions.go` line 86: `GetStatsForPlayersInPod(ctx, podID, playerIDs)` — single batch query filtered to pod |
| 4 | New user with no pods sees "Welcome to EDH Tracker" and "Create a Pod" button | VERIFIED | `app/src/routes/home/index.tsx` lines 64, 69: exact strings present; dialog and PostPod wiring confirmed |
| 5 | Clicking "Create a Pod" opens dialog with pod name field and submit | VERIFIED | `home/index.tsx` lines 44-45: `const { id } = await PostPod(podName)` and navigate on success |
| 6 | After creating a pod, user navigates to /pod/{newPodId} | VERIFIED | `home/index.tsx` line 45: `navigate(\`/pod/${id}\`)` after PostPod resolves |
| 7 | AppBar pod-name dropdown includes "Create new pod" at the bottom | VERIFIED | `app/src/routes/root.tsx` line 157: `<MenuItem value="create-new">Create new pod</MenuItem>` with Divider above |
| 8 | Selecting "Create new pod" opens the Create Pod dialog | VERIFIED | `root.tsx` line 120: `if (id === "create-new") { setCreatePodOpen(true); return; }` |
| 9 | Pod Players tab shows card-per-player layout with pod-scoped stats | VERIFIED | `app/src/routes/pod/PlayersTab.tsx` line 64: `<Paper elevation={2}>` per player; line 93 renders `p.stats.record`, `p.stats.points`, `p.stats.kills` |
| 10 | Player can navigate to /deck/new and create a new deck | VERIFIED | `app/src/routes/deck/new/index.tsx`: full form wired with `PostDeck`, `PostCommander`, `freeSolo`, navigation after submit |
| 11 | Route /deck/new is registered and requires auth | VERIFIED | `app/src/index.tsx` lines 63-65: `path: "deck/new"`, `<RequireAuth><NewDeckView /></RequireAuth>`, `loader: newDeckLoader` |
| 12 | Player Decks tab shows "Add Deck" button only for deck owner | VERIFIED | `app/src/routes/player/DecksTab.tsx` line 18: `isOwner = user?.player_id === playerId`; lines 37-40, 54-58: button conditionally rendered |
| 13 | Pod Decks tab sorted by record descending by default | VERIFIED | `app/src/routes/pod/DecksTab.tsx` lines 39-41: `sortModel: [{ field: "record", sort: "desc" }]` |
| 14 | Retired decks hidden from Player Decks tab by default | VERIFIED | `app/src/routes/player/DecksTab.tsx` lines 67-71: DataGrid `initialState.filter.filterModel` with `{ field: "retired", operator: "is", value: "false" }`; "Is Retired" column at line 49 |

**Score:** 14/14 truths verified

### Required Artifacts

| Artifact | Plan | Status | Details |
|----------|------|--------|---------|
| `lib/repositories/gameResult/stats.go` | 05-01 | VERIFIED | `type gameStatWithPlayer struct` at line 28 |
| `lib/repositories/gameResult/repo.go` | 05-01 | VERIFIED | `getStatsForPlayersInPod` SQL constant at line 40; `GetStatsForPlayersInPod` method at line 184 |
| `lib/repositories/interfaces.go` | 05-01 | VERIFIED | `GetStatsForPlayersInPod(ctx, podID, playerIDs)` at line 66 |
| `lib/business/player/functions.go` | 05-01 | VERIFIED | `gameResultRepo.GetStatsForPlayersInPod(ctx, podID, playerIDs)` at line 86; no per-player loop |
| `lib/routers/pod.go` | 05-01 | VERIFIED | `json.NewEncoder(w).Encode` at line 188; `Content-Type: application/json` at line 186 |
| `lib/routers/deck.go` | 05-01 | VERIFIED | `json.NewEncoder(w).Encode` at line 206; `Content-Type: application/json` at line 204 |
| `app/src/components/stats.tsx` | 05-02 | VERIFIED | `Math.max` present for min-4-place guard |
| `app/src/routes/pod/DecksTab.tsx` | 05-02 / 05-08 | VERIFIED | DataGrid with `sortModel: [{ field: "record", sort: "desc" }]`; receives all decks from `GetAllDecksForPod` via pod/index.tsx loader |
| `app/src/routes/player/DecksTab.tsx` | 05-02 / 05-08 | VERIFIED | DataGrid filter hides retired by default; "Is Retired" column present; "Add Deck" button owner-gated |
| `app/src/routes/home/index.tsx` | 05-03 | VERIFIED | "Welcome to EDH Tracker", "Create a Pod", Dialog, PostPod import, navigate |
| `app/src/routes/root.tsx` | 05-03 | VERIFIED | "Create new pod" MenuItem, Divider, dialog wiring, PostPod import |
| `app/src/http.ts` (PostPod) | 05-03 | VERIFIED | `PostPod(name: string): Promise<{ id: number }>` at line 307 |
| `app/src/routes/pod/PlayersTab.tsx` | 05-04 | VERIFIED | Paper cards per player; `p.stats.record/points/kills`; Cancel button says "Cancel" |
| `app/src/routes/deck/new/index.tsx` | 05-05 | VERIFIED | `newDeckLoader`, `NewDeckView`, freeSolo Autocomplete, PostDeck, PostCommander, navigate after submit |
| `app/src/index.tsx` | 05-05 | VERIFIED | `deck/new` route registered with `RequireAuth` and `newDeckLoader` |
| `app/src/routes/deck/SettingsTab.tsx` | 05-05 | VERIFIED | freeSolo on both commander Autocomplete fields; `PostCommander` inline creation; `TooltipIcon` at line 169 |
| `app/src/http.ts` (PostDeck) | 05-05 | VERIFIED | `PostDeck(body: NewDeckRequest): Promise<{ id: number }>` at line 277 |
| `app/src/types.ts` | 05-05 | VERIFIED | `NewDeckRequest` at line 134; `NewDeckData` at line 141 |
| `lib/business/pod/functions.go` | 05-06 | VERIFIED | `client.GormDb.WithContext(ctx).Transaction(...)` at line 55 — pod Create wraps all 3 writes in a transaction |
| `lib/business/business.go` | 05-06 | VERIFIED | `pod.Create(r.Pods, r.PlayerPodRoles, client)` at line 91 |
| `lib/routers/commander.go` | 05-07 | VERIFIED | `json.NewEncoder(w).Encode` at line 129 — POST /api/commander returns `{id: N}` |
| `app/src/http.ts` (PostCommander) | 05-07 | VERIFIED | `PostCommander(name: string): Promise<{ id: number }>` at line 264 |
| `app/src/routes/pod/index.tsx` | 05-08 | VERIFIED | `GetAllDecksForPod(podId)` in loader at line 27; passes full deck array to `<PodDecksTab decks={decks} />` |

### Key Link Verification

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `lib/business/player/functions.go` | `lib/repositories/gameResult/repo.go` | `GetStatsForPlayersInPod` call | WIRED | Line 86: single batch call replacing N per-player calls |
| `lib/repositories/gameResult/repo.go` | `lib/repositories/interfaces.go` | Interface satisfaction | WIRED | Line 66 of interfaces.go defines method; repo.go line 184 satisfies it |
| `app/src/routes/home/index.tsx` | `app/src/http.ts` | `PostPod` call | WIRED | Line 44: `const { id } = await PostPod(podName)` |
| `app/src/routes/root.tsx` | `app/src/http.ts` | `PostPod` call in PodSelector | WIRED | Line 132: `const { id } = await PostPod(podName)` |
| `app/src/routes/deck/new/index.tsx` | `app/src/http.ts` | `PostDeck` + `PostCommander` calls | WIRED | Line 88: `const { id } = await PostDeck(body)`; line 63: `PostCommander(newName)` |
| `app/src/index.tsx` | `app/src/routes/deck/new/index.tsx` | Route registration | WIRED | Lines 15, 63-65: import and `path: "deck/new"` with loader |
| `lib/business/pod/functions.go` | `lib/business/business.go` | Constructor wiring with `client` | WIRED | `pod.Create(r.Pods, r.PlayerPodRoles, client)` at business.go line 91 |
| `lib/routers/commander.go` | `app/src/http.ts` | `PostCommander` response parsing | WIRED | Commander router encodes `{id: N}`; http.ts returns `Promise<{ id: number }>` |
| `app/src/routes/pod/index.tsx` | `app/src/routes/pod/DecksTab.tsx` | `decks` prop | WIRED | Loader calls `GetAllDecksForPod`; passes result to `<PodDecksTab decks={decks} />` |

### Data-Flow Trace (Level 4)

| Artifact | Data Variable | Source | Produces Real Data | Status |
|----------|--------------|--------|-------------------|--------|
| `pod/PlayersTab.tsx` | `p.stats.record`, `p.stats.points`, `p.stats.kills` | `GET /api/players?pod_id=N` → `GetAllByPod` → SQL at repo.go line 40 joining `game_result`, `deck`, `game` filtered to `game.pod_id = ?` | Yes | FLOWING |
| `home/index.tsx` | navigates to `/pod/${id}` | `PostPod` → `POST /api/pod` → `pods.Create(...)` returns `LastInsertId()` → JSON `{id: N}` | Yes | FLOWING |
| `deck/new/index.tsx` | navigates to `/player/{callerId}/deck/${id}` | `PostDeck` → `POST /api/deck` → `decks.Create(...)` returns `LastInsertId()` → JSON `{id: N}` | Yes | FLOWING |
| `pod/DecksTab.tsx` | `decks` prop | `GetAllDecksForPod` → `GET /api/decks?pod_id=N` → DB query | Yes | FLOWING |
| `player/DecksTab.tsx` | `data` via `GetDecksForPlayer` | `GET /api/decks?player_id=N` → DB query; DataGrid filter is client-side on real data | Yes | FLOWING |

### Behavioral Spot-Checks

| Behavior | Command | Result | Status |
|----------|---------|--------|--------|
| Go code compiles | `go vet ./lib/...` | No output (exit 0) | PASS |
| TypeScript type-checks | `./node_modules/.bin/tsc --noEmit` | No output (exit 0) | PASS |
| `GetStatsForPlayersInPod` in interface | grep `lib/repositories/interfaces.go` | Line 66 confirmed | PASS |
| `GetAllByPod` uses batch stats, no per-player loop | grep `lib/business/player/functions.go` | Line 86; no `GetStatsForPlayer` inside `GetAllByPod` | PASS |
| Route /deck/new registered | grep `app/src/index.tsx` | Lines 63-65 confirmed | PASS |
| Pod Create wraps 3 writes in transaction | grep `lib/business/pod/functions.go` | Line 55: `GormDb.WithContext(ctx).Transaction(...)` | PASS |
| POST /api/commander returns JSON ID | grep `lib/routers/commander.go` | Line 129: `json.NewEncoder(w).Encode` | PASS |

### Requirements Coverage

| Requirement | Source Plans | Description | Status | Evidence |
|-------------|-------------|-------------|--------|----------|
| POD-01 | 05-03, 05-06 | Pod creation accessible from AppBar, not buried in player settings | SATISFIED | "Create new pod" MenuItem in root.tsx line 157; pod Create in transaction so new pod immediately visible in selector |
| POD-02 | 05-03, 05-06 | New users with no pods guided to create or join one | SATISFIED | HomeView empty state: "Welcome to EDH Tracker" + "Create a Pod" CTA (home/index.tsx lines 64-69) |
| POD-03 | 05-02, 05-08 | Pod Decks tab sorted by record by default | SATISFIED | pod/DecksTab.tsx lines 39-41: `sortModel: [{ field: "record", sort: "desc" }]` |
| POD-04 | 05-01, 05-04 | Pod Players tab shows each player's record and points within that pod | SATISFIED | Backend: `GetStatsForPlayersInPod` filters by `game.pod_id`. Frontend: PlayersTab renders `p.stats.record/points/kills` per card |
| DECK-01 | 05-05, 05-07 | Player can create a new deck from the UI | SATISFIED | `/deck/new` route (index.tsx lines 63-65); "Add Deck" button in player/DecksTab; PostCommander returns `{id}` enabling inline commander creation |
| DECK-02 | 05-02 (pre-existing) | Commander update field has tooltip | SATISFIED | deck/SettingsTab.tsx line 169: `<TooltipIcon title="This is for changing an existing deck's commander. To add a new deck, use the Add Deck button instead." />` |
| DECK-03 | 05-02, 05-08 | Retired deck visibility consistent | SATISFIED | player/DecksTab.tsx lines 67-71: DataGrid filter hides retired by default; "Is Retired" column so user can clear filter |

All 7 Phase 5 requirement IDs are accounted for. REQUIREMENTS.md traceability table marks POD-01 through POD-04 and DECK-01 through DECK-03 as Phase 5 Complete. No orphaned Phase 5 requirements found.

**Note:** The task specified POD-01, POD-02, POD-03, DECK-01, DECK-02, DECK-03. The plans additionally claim POD-04, which is also satisfied. All 7 are verified.

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| `lib/repositories/gameResult/repo.go` | 14 | `// TODO: Better way to do deck and player stats?` | Info | Pre-existing comment unrelated to Phase 5; no blocking behavior |

No placeholder returns, hardcoded empty arrays feeding rendered output, or FIXME markers in any Phase 5 files.

**Correction from previous VERIFICATION.md:** The prior report at "line 32" cited `visibleRows = (data ?? []).filter((d: Deck) => !d.retired)` as the retired-deck filter mechanism. This does not exist in the current file. Gap-closure Plan 08 replaced the manual filter with a DataGrid `initialState.filter.filterModel` approach. The behavioral outcome (retired decks hidden by default, clearable by user) is the same and correctly implemented.

### Human Verification

#### 1. Pod creation happy path — HomeView

**Test:** Log in as a user with no pods. Verify the HomeView shows "Welcome to EDH Tracker" heading, body text "Create your first pod or ask a friend for an invite link.", and a "Create a Pod" button. Click the button, enter a pod name, click "Create Pod".
**Expected:** Lands on `/pod/{newId}` immediately after creation; new pod appears in the AppBar pod selector.
**Result:** PASS (verified 2026-03-27; AppBar selector update required quick task 260326-x0y fix)

#### 2. AppBar "Create new pod" flow

**Test:** Log in with any user. Open the pod selector dropdown in the AppBar. Verify a "Create new pod" menu item appears at the bottom below a divider. Click it. Verify the Create Pod dialog opens. Create a pod.
**Expected:** Dropdown item present; dialog opens; successful creation navigates to new pod; pod appears in selector.
**Result:** PASS (verified 2026-03-27; AppBar selector update required quick task 260326-x0y fix)

#### 3. Pod Players tab shows pod-scoped (not global) stats

**Test:** Navigate to a pod with players who have game history in other pods. Open the Players tab. Compare each player's stat card to what appears on their player profile (global stats). Stats should differ if the player has games outside this pod.
**Expected:** Cards show only this pod's results.
**Result:** PASS (verified 2026-03-27)

#### 4. /deck/new form — Commander conditional fields

**Test:** Navigate to `/deck/new`. Select "Commander" format. Verify Commander (required) and "Partner Commander (optional)" fields appear. Select a non-Commander format. Verify those fields disappear.
**Expected:** Conditional rendering works; submit disabled until commander selected for Commander format.
**Result:** PASS (verified 2026-03-27)

#### 5. freeSolo commander creation — deck/new and SettingsTab

**Test:** On `/deck/new`, type a commander name that does not exist. Verify a "Create ..." option appears in the Autocomplete dropdown. Select it. Verify the commander is created inline (no error, no navigation away).
**Expected:** POST /api/commander fires; field populates with the new commander; form remains on `/deck/new`.
**Result:** PASS (verified 2026-03-27)

#### 6. Retired deck filter in Player Decks tab

**Test:** Navigate to a player profile with at least one retired deck. Open the Decks tab. Verify the retired deck is hidden by default. Open the DataGrid filter panel via the toolbar and clear the "Is Retired is false" filter. Verify the retired deck appears.
**Expected:** Retired decks hidden by default; clearable via toolbar filter.
**Result:** PASS (verified 2026-03-27)

### Gaps Summary

No gaps. All 14 must-have truths verified. All 7 Phase 5 requirement IDs (POD-01 through POD-04, DECK-01 through DECK-03) have implementation evidence. Go and TypeScript compile cleanly. Eight plans executed (01-05 original plus 06-08 gap-closure) — all artifacts are substantive, wired, and data-connected. All 6 human verification items confirmed passing 2026-03-27.

---

_Verified: 2026-03-26 | Human UAT complete: 2026-03-27_
_Verifier: Claude (gsd-verifier)_
