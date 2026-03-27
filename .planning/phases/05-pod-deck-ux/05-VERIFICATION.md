---
phase: 05-pod-deck-ux
verified: 2026-03-25T00:00:00Z
status: passed
score: 14/14 must-haves verified
re_verification: false
---

# Phase 5: Pod/Deck UX Verification Report

**Phase Goal:** Improve pod and deck UX — pod creation accessible from AppBar and HomeView, pod players tab redesigned to show pod-scoped stats per player, deck creation accessible from player profile.
**Verified:** 2026-03-25
**Status:** passed
**Re-verification:** No — initial verification

## Goal Achievement

### Observable Truths

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | POST /api/pod returns 201 with JSON body containing new pod ID | VERIFIED | `lib/routers/pod.go` lines 180-190: captures `podID`, sets Content-Type, calls `json.NewEncoder(w).Encode` with `{ID: podID}` |
| 2 | POST /api/deck returns 201 with JSON body containing new deck ID | VERIFIED | `lib/routers/deck.go` lines 198-208: captures `deckID`, sets Content-Type, calls `json.NewEncoder(w).Encode` with `{ID: deckID}` |
| 3 | GET /api/players?pod_id=N returns pod-scoped stats | VERIFIED | `lib/business/player/functions.go` line 86: `GetAllByPod` calls `GetStatsForPlayersInPod(ctx, podID, playerIDs)` — batch, pod-scoped, single query |
| 4 | New user with no pods sees "Welcome to EDH Tracker" and "Create a Pod" button | VERIFIED | `app/src/routes/home/index.tsx` lines 64-69: exact heading and button text present |
| 5 | Clicking "Create a Pod" opens dialog with pod name field and submit | VERIFIED | `home/index.tsx` lines 72-96: full Dialog with TextField, submit, error state |
| 6 | After creating a pod, user navigates to /pod/{newPodId} | VERIFIED | `home/index.tsx` line 45: `navigate(\`/pod/${id}\`)` after `PostPod` succeeds |
| 7 | AppBar pod-name dropdown includes "Create new pod" at the bottom | VERIFIED | `app/src/routes/root.tsx` lines 156-157: Divider + `MenuItem value="create-new"` |
| 8 | Selecting "Create new pod" opens the Create Pod dialog | VERIFIED | `root.tsx` lines 120-122: `if (id === "create-new") { setCreatePodOpen(true); return; }` |
| 9 | Pod Players tab shows card-per-player layout with pod-scoped stats | VERIFIED | `app/src/routes/pod/PlayersTab.tsx`: Paper cards, `p.stats.record`, `p.stats.points`, `p.stats.kills` |
| 10 | Player can navigate to /deck/new and create a new deck | VERIFIED | `app/src/routes/deck/new/index.tsx`: full form with `newDeckLoader`, `PostDeck`, navigation after submit |
| 11 | Route /deck/new is registered and requires auth | VERIFIED | `app/src/index.tsx` lines 63-66: `{ path: "deck/new", element: <RequireAuth><NewDeckView /></RequireAuth>, loader: newDeckLoader }` |
| 12 | Player Decks tab shows "Add Deck" button only for deck owner | VERIFIED | `app/src/routes/player/DecksTab.tsx` lines 19, 40-44, 56-61: `isOwner` check, button in both empty-state and normal path |
| 13 | Pod Decks tab sorted by record descending by default | VERIFIED | `app/src/routes/pod/DecksTab.tsx` lines 58-61: `initialState.sorting.sortModel: [{ field: "record", sort: "desc" }]` |
| 14 | Retired decks hidden from Player Decks tab by default | VERIFIED | `app/src/routes/player/DecksTab.tsx` line 32: `visibleRows = (data ?? []).filter((d: Deck) => !d.retired)` |

**Score:** 14/14 truths verified

### Required Artifacts

| Artifact | Plan | Status | Details |
|----------|------|--------|---------|
| `lib/repositories/gameResult/repo.go` | 05-01 | VERIFIED | Contains `getStatsForPlayersInPod` SQL constant and `GetStatsForPlayersInPod` method |
| `lib/repositories/gameResult/stats.go` | 05-01 | VERIFIED | Contains `gameStatWithPlayer` struct |
| `lib/repositories/interfaces.go` | 05-01 | VERIFIED | `GameResultRepository` interface includes `GetStatsForPlayersInPod(ctx, podID, playerIDs)` at line 66 |
| `lib/business/player/functions.go` | 05-01 | VERIFIED | `GetAllByPod` uses `GetStatsForPlayersInPod`, no per-player `GetStatsForPlayer` loop |
| `lib/routers/pod.go` | 05-01 | VERIFIED | `json.NewEncoder` and `Content-Type: application/json` in `PodCreate` |
| `lib/routers/deck.go` | 05-01 | VERIFIED | `json.NewEncoder` and `Content-Type: application/json` in `DeckCreate` |
| `app/src/routes/home/index.tsx` | 05-03 | VERIFIED | Contains "Welcome to EDH Tracker", "Create a Pod", Dialog, PostPod import, navigate |
| `app/src/routes/root.tsx` | 05-03 | VERIFIED | Contains "create-new" sentinel, "Create new pod" MenuItem, Divider, Dialog, PostPod import |
| `app/src/http.ts` — PostPod | 05-03 | VERIFIED | `PostPod(name: string): Promise<{ id: number }>` at line 305 |
| `app/src/http.ts` — PostDeck | 05-05 | VERIFIED | `PostDeck(body: NewDeckRequest): Promise<{ id: number }>` at line 275 |
| `app/src/types.ts` | 05-05 | VERIFIED | Contains `NewDeckRequest` and `NewDeckData` interfaces |
| `app/src/routes/deck/new/index.tsx` | 05-05 | VERIFIED | `newDeckLoader`, `NewDeckView`, `freeSolo`, `PostDeck`, `PostCommander`, "Create Deck", "Discard" |
| `app/src/index.tsx` | 05-05 | VERIFIED | Contains `deck/new` path and `newDeckLoader` |
| `app/src/routes/player/DecksTab.tsx` | 05-05 | VERIFIED | "Add Deck" button, `/deck/new` link, `isOwner` check, `useAuth` import |
| `app/src/routes/deck/SettingsTab.tsx` | 05-05 | VERIFIED | `freeSolo` on both Autocomplete, `PostCommander` import, `createFilterOptions`, no `getOptionKey` |
| `app/src/routes/pod/PlayersTab.tsx` | 05-04 | VERIFIED | Paper cards, `PersonAddIcon`, `PersonOffIcon`, `TooltipIconButton`, `Record` from stats, pod-scoped stats row |
| `app/src/components/stats.tsx` | 05-02 | VERIFIED | `Math.max(...Object.keys(record).map(Number), 4)` — min 4 place columns |
| `app/src/routes/new/index.tsx` | 05-02 | VERIFIED | Contains "Discard" button |

### Key Link Verification

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `lib/business/player/functions.go` | `lib/repositories/gameResult/repo.go` | `GetStatsForPlayersInPod` call | WIRED | Line 86: `gameResultRepo.GetStatsForPlayersInPod(ctx, podID, playerIDs)` |
| `lib/repositories/gameResult/repo.go` | `lib/repositories/interfaces.go` | Interface satisfaction | WIRED | `GetStatsForPlayersInPod` in interface at line 66 |
| `app/src/routes/home/index.tsx` | `app/src/http.ts` | PostPod call | WIRED | Line 44: `const { id } = await PostPod(podName)` |
| `app/src/routes/root.tsx` | `app/src/http.ts` | PostPod call in PodSelector | WIRED | Line 132: `const { id } = await PostPod(podName)` |
| `app/src/routes/deck/new/index.tsx` | `app/src/http.ts` | PostDeck call | WIRED | Line 90: `const { id } = await PostDeck(body)` |
| `app/src/index.tsx` | `app/src/routes/deck/new/index.tsx` | Route registration with loader | WIRED | Lines 15, 63-66: import and route entry |

### Data-Flow Trace (Level 4)

| Artifact | Data Variable | Source | Produces Real Data | Status |
|----------|--------------|--------|-------------------|--------|
| `PlayersTab.tsx` | `p.stats.record`, `p.stats.points`, `p.stats.kills` | `GET /api/players?pod_id=N` → `GetAllByPod` → `GetStatsForPlayersInPod` DB query | Yes — SQL query groups by `deck.player_id` filtered to `game.pod_id = ?` | FLOWING |
| `home/index.tsx` | `createPodOpen`, navigates on submit | `PostPod` → `POST /api/pod` → returns `{"id": N}` | Yes — `podID` captured from `p.pods.Create(ctx, ...)` which returns `LastInsertId` | FLOWING |
| `deck/new/index.tsx` | navigates to `/player/{id}/deck/{newDeckId}` after submit | `PostDeck` → `POST /api/deck` → returns `{"id": N}` | Yes — `deckID` captured from `d.decks.Create(ctx, ...)` which returns `LastInsertId` | FLOWING |

### Behavioral Spot-Checks

| Behavior | Check | Status |
|----------|-------|--------|
| Go code compiles | `go vet ./lib/...` exits 0 | PASS |
| TypeScript type-checks | `tsc --noEmit` exits 0 | PASS |
| `GetStatsForPlayersInPod` in interface | `grep GetStatsForPlayersInPod lib/repositories/interfaces.go` | PASS |
| `GetAllByPod` does not call per-player stats loop | No `GetStatsForPlayer(ctx, m.PlayerID)` in `GetAllByPod` | PASS |
| Route /deck/new registered | `grep "deck/new" app/src/index.tsx` | PASS |

### Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
|-------------|------------|-------------|--------|----------|
| POD-01 | 05-03 | Pod creation accessible from AppBar, not buried in player settings | SATISFIED | "Create new pod" MenuItem in PodSelector dropdown (`root.tsx` line 157); PlayerSettingsTab has no Create Pod section |
| POD-02 | 05-03 | New users with no pods guided to create or join one | SATISFIED | HomeView empty state renders "Welcome to EDH Tracker" with "Create a Pod" CTA (`home/index.tsx` lines 63-98) |
| POD-03 | 05-02 | Pod Decks tab sorted by record (win rate or points) by default | SATISFIED | `pod/DecksTab.tsx` `initialState.sorting.sortModel: [{ field: "record", sort: "desc" }]` |
| POD-04 | 05-01, 05-04 | Pod Players tab shows each player's record and points within that pod | SATISFIED | Backend: `GetStatsForPlayersInPod` returns pod-filtered stats. Frontend: Paper cards render `p.stats.record`, `p.stats.points`, `p.stats.kills` |
| DECK-01 | 05-05 | Player can create a new deck from the UI | SATISFIED | `/deck/new` route with `NewDeckView`, "Add Deck" button on player's DecksTab |
| DECK-02 | 05-02 (pre-existing) | Commander update field has tooltip | SATISFIED (pre-existing) | `deck/SettingsTab.tsx` line 169: `TooltipIcon` with correct tooltip text |
| DECK-03 | 05-02 | Retired deck visibility consistent | SATISFIED | `player/DecksTab.tsx` filters `!d.retired`; pod DecksTab fetches active decks via paginated endpoint |

All 7 phase-5 requirement IDs from the plans are accounted for. REQUIREMENTS.md traceability table shows POD-01, POD-02 as "Pending" and DECK-01 as "Pending" — these were the open ones entering Phase 5 and are now implemented (the traceability table reflects pre-phase state; the actual code satisfies them).

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| `lib/repositories/gameResult/repo.go` | 14 | `// TODO: Better way to do deck and player stats?` | Info | Pre-existing comment, unrelated to Phase 5 work; no blocking behavior |

No other TODOs, placeholder returns, or stub patterns detected in Phase 5 files.

### Human Verification Required

#### 1. Pod creation happy path — HomeView

**Test:** Log in as a user with no pods. Verify the HomeView shows "Welcome to EDH Tracker" heading, body text "Create your first pod or ask a friend for an invite link.", and a "Create a Pod" button. Click the button, enter a pod name, click "Create Pod". Verify you are navigated to `/pod/{newId}`.
**Expected:** Full onboarding CTA renders; after submit, lands on the new pod page.
**Why human:** Real navigation and backend response required; cannot test without running app.

#### 2. AppBar "Create new pod" flow

**Test:** Log in with any user. Open the pod selector dropdown in the AppBar. Verify a "Create new pod" menu item appears at the bottom below a divider. Click it. Verify the Create Pod dialog opens. Create a pod. Verify navigation to the new pod.
**Expected:** Dropdown item present; dialog opens; successful creation navigates.
**Why human:** Requires running browser session.

#### 3. Pod Players tab card layout with pod-scoped stats

**Test:** Navigate to a pod with multiple players who have game history. Open the Players tab. Verify each player appears as a Paper card showing name (linked), Manager chip (if manager), and a stat row with record / points / kills.
**Expected:** Cards render pod-filtered stats, not global stats. A player who has played in other pods shows only results from this pod.
**Why human:** Requires real game data across multiple pods to distinguish pod-scoped from global stats.

#### 4. /deck/new form — Commander format conditional fields

**Test:** Navigate to `/deck/new`. Select "Commander" format. Verify Commander and "Partner Commander (optional)" Autocomplete fields appear. Select a non-Commander format. Verify those fields disappear.
**Expected:** Conditional rendering works; commander required for Commander format; submit disabled until commander selected.
**Why human:** Requires running React app with real format data loaded.

#### 5. freeSolo commander creation

**Test:** On `/deck/new` (or DeckSettingsTab), type a commander name that does not exist. Verify a `Create "..."` option appears in the Autocomplete dropdown. Select it. Verify the commander is created and the field is populated.
**Expected:** Inline commander creation works without leaving the form.
**Why human:** Requires live API call to POST /api/commander.

### Gaps Summary

No gaps found. All 14 must-have truths are verified in the codebase. All 7 requirement IDs (POD-01, POD-02, POD-03, POD-04, DECK-01, DECK-02, DECK-03) have implementation evidence. Go and TypeScript compile cleanly.

---

_Verified: 2026-03-25_
_Verifier: Claude (gsd-verifier)_
