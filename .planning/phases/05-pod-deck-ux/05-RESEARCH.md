# Phase 5: Pod & Deck UX - Research

**Researched:** 2026-03-26
**Domain:** React/MUI frontend UX + Go backend API enhancement
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Pod Creation Entry Point (POD-01)**
- D-01: Pod creation moves off PlayerSettingsTab. SettingsTab keeps "Your Pods" list + Leave Pod, removes "Create New Pod" section.
- D-02: Pod creation accessible from HomeView empty state (first-time users) and AppBar pod-name dropdown (existing users).
- D-03: AppBar pod name becomes a dropdown (MUI Select or Menu) on `/pod/:podId`. Lists all pods as links + "Create new pod" at bottom.
- D-04: After creating a pod from anywhere, auto-navigate to `/pod/{newPodId}`.
- D-05: Pod-name dropdown only appears on `/pod/:podId`. Other views keep AppBar as-is.

**New User Onboarding (POD-02)**
- D-06: HomeView empty state → onboarding screen with "Create a Pod" button that opens a MUI Dialog with pod name field + submit.
- D-07: HomeView does NOT include "Join with invite link" CTA. Empty state is creation-focused.
- D-08: "Create Pod" modal is reused from AppBar dropdown's "Create new pod" option.

**Deck Creation (DECK-01)**
- D-09: "Add Deck" button on Player view → Decks tab, owner-only (viewer is deck owner).
- D-10: "Add Deck" navigates to `/deck/new` — a dedicated top-level route.
- D-11: `/deck/new` uses React Router loader + `useNavigate` on submit. Loader fetches formats and commanders.
- D-12: After successful deck creation, navigate to `/player/{callerId}/deck/{newDeckId}`.
- D-13: `/deck/new` includes a Cancel/back button navigating to `/player/{callerId}`.
- D-14: Deck always belongs to the logged-in caller — no player ID in URL. Backend uses JWT `callerPlayerID`.
- D-15: Required fields: Name and Format.
- D-16: Format = "Commander" shows Commander (required) + Partner Commander (optional). Other formats: fields hidden.
- D-17: Commander and Partner Commander use MUI Autocomplete with `freeSolo={true}`. Unknown name + Enter calls `POST /api/commander`, then uses new ID.
- D-18: After creation, navigate to `/player/{callerId}/deck/{newDeckId}`.
- D-19: DeckSettingsTab commander Autocomplete also updated to support freeSolo/new commander creation.

**Pod Players Tab Redesign (POD-04)**
- D-20: Switch from List/ListItem layout to card-per-player layout (Box/Card with flex, Paper elevation={2}).
- D-21: Stats per player: Record + Points + Kills, all pod-scoped (games in this pod only).
- D-22: Backend enhancement: `GET /api/players?pod_id={id}` computes stats filtered by pod_id via join on game table.
- D-23: Stats loading uses skeleton placeholders — player cards render immediately; stat row shows Skeleton while fetching.
- D-24: Promote/Remove become TooltipIconButton (PersonAdd/PersonOff). Manager-only, not shown for caller's own card.
- D-25: Player names are links to `/player/{id}`.
- D-26: Leave Pod stays in player settings. No Leave option on Players tab cards.

**Pod Decks Tab Default Sort (POD-03)**
- D-27: Uses existing RecordComparator from stats.tsx as sort comparator on Record column.
- D-28: DataGrid `initialState.sorting` set to sort by Record column descending. Client-side, no backend changes.

**Retired Deck Visibility (DECK-03)**
- D-29: Retired decks hidden everywhere active. Player Decks tab: filter client-side to exclude `retired=true` rows.
- D-30: Retired decks visible in game history (stored as strings, unaffected).
- D-31: Game view always shows stored deck names regardless of retirement status.

**Record Display — Minimum 4 Places**
- D-32: Record component in stats.tsx updated: `Math.max(...keys, 4)` — always display at least 4 place columns.

**Playing Cards Icon → Home**
- D-33: SvgIconPlayingCards in AppBar wrapped in `component={Link}` to navigate to `/`.

**New Game Form — Cancel / Back Button**
- D-34: New game form gets a Cancel/"Discard" button navigating back to `/pod/:podId`.

### Claude's Discretion

- AppBar pod dropdown implementation detail: MUI Select, Menu, or Button+Popover — whichever fits AppBar layout cleanest at 375px
- Exact card elevation, border, and spacing for the player tab redesign
- Whether to use `initialState.sorting` or `onSortModelChange` for DataGrid default sort
- Whether pod-scoped stats are fetched in same request as player list (one round-trip) or as a second request with skeleton loading
- Loading state granularity on /deck/new (loader vs component mount fetch)

### Deferred Ideas (OUT OF SCOPE)

- Pod Decks tab "Add Deck" shortcut — Add Deck lives on Player view only
- Dark/light mode toggle
- "Join with invite code" on HomeView empty state
- Leave Pod option on Players tab card
- PlayerView showing pod-scoped stats when visited from pod context
- Security review all API and frontend route authorization (Phase 6)
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| POD-01 | Pod creation accessible from pod page, not buried in player settings | D-02/D-03: AppBar dropdown + HomeView CTA; existing PodSelector needs "Create new pod" MenuItem |
| POD-02 | New users with no pods guided to create or join one | D-06: HomeView empty state → MUI Dialog modal; PostPod API bug must be fixed |
| POD-03 | Pod Decks tab sorted by record by default | D-27/D-28: `initialState.sorting` on DataGrid; RecordComparator already wired on Record column |
| POD-04 | Pod Players tab shows each player's record and points within that pod | D-22: backend `GetAllByPod` must call new `GetStatsForPlayersInPod` repo method |
| DECK-01 | Player can create a new deck from the UI | D-09 through D-18: `/deck/new` route; `PostDeck` missing from http.ts; backend DeckCreate must return deck ID |
| DECK-02 | Commander update tooltip already implemented | Pre-verified complete: `DeckSettingsTab.tsx` line 164 has the TooltipIcon. No work needed. |
| DECK-03 | Retired deck visibility defined and consistent | D-29: Player Decks tab client-side filter; Pod Decks tab backend already filters retired=false |
</phase_requirements>

---

## Summary

Phase 5 is entirely within the existing stack — React + MUI v5, Go + GORM + MySQL — with no new packages. The work divides into three categories: (1) frontend-only changes (DataGrid sort, Record component fix, icon link, retired filter, new game cancel button), (2) frontend + backend changes that are well-scoped (pod creation modal reused in two places, `/deck/new` route + form), and (3) a backend enhancement that requires a new SQL query (pod-scoped player stats).

The most important discovery is a **pre-existing API contract mismatch**: both `POST /api/pod` and `POST /api/deck` currently return HTTP 201 with no body. The frontend's auto-navigate-after-create features (D-04, D-12) require the new resource's ID from the response. The plan must fix both backend handlers to return `{"id": N}` and update the corresponding frontend fetch functions accordingly. The existing `PostPod` in `http.ts` already calls `res.json()` (suggesting intent to return a body), but the backend returns an empty 201 — this is a latent bug that Phase 5 must resolve.

The pod-scoped stats work (POD-04/D-22) is the most complex backend change: it requires a new `GetStatsForPlayersInPod` repo method (SQL join across game_result + game, filtered by pod_id), a new `GetAllByPodWithPodStats` business function (or extending `GetAllByPod`), and a frontend that either passes stats in the same API response or fetches them separately with skeleton loading. Research recommends the single-request approach (return pod-scoped stats in the existing `GET /api/players?pod_id=` response) since the data is already scoped to a single pod and the N+1 query can be batched with one SQL query per pod load.

**Primary recommendation:** Fix the API contract first (pod and deck create return IDs), then implement backend pod-scoped stats, then implement frontend features. The order matters because pod creation modal depends on the corrected API.

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| @mui/material | 5.15.2 | UI components (Dialog, Paper, Autocomplete, Chip, Skeleton, Select, MenuItem) | Project-locked; all required components already installed |
| @mui/x-data-grid | 6.18.6 | DataGrid with `initialState.sorting` | Already used in DecksTab and PlayerDecksTab |
| @mui/icons-material | 5.15.3 | PersonAddIcon, PersonOffIcon, ArrowBackIcon | Already installed; icons confirmed in UI-SPEC |
| react-router-dom | 6.21.1 | `useNavigate`, `useLoaderData`, `Link`, `createBrowserRouter` | Project-locked |
| Go 1.26 | — | Backend handler and SQL changes | Project-locked |
| GORM | 1.31.1 | `db.WithContext(ctx).Raw(sql, args).Scan(&result)` for new pod-scoped stats query | Established pattern from getStatsForPlayer |

### No New Packages
All components required for this phase are already installed. No `npm install` needed.

---

## Architecture Patterns

### Frontend: Shared Dialog Component (CreatePodDialog)
The "Create Pod" Dialog is triggered from two places: HomeView empty state and AppBar PodSelector dropdown. Extract it as a shared component (or inline into root.tsx and pass open state down as prop).

**Recommended:** Inline in `root.tsx` alongside the PodSelector — keeps pod creation state local to the AppBar section that already manages pod selection. HomeView opens it via a callback prop or a shared state lifted to Root.

**Pattern (established in Phase 3):**
```tsx
// MUI Dialog with controlled open state
const [createPodOpen, setCreatePodOpen] = useState(false);
// ...
<Dialog open={createPodOpen} onClose={() => setCreatePodOpen(false)} maxWidth="xs" fullWidth>
  <DialogTitle>Create a New Pod</DialogTitle>
  <DialogContent>
    <TextField label="Pod Name" value={podName} onChange={...} fullWidth />
  </DialogContent>
  <DialogActions>
    <Button onClick={() => setCreatePodOpen(false)}>Discard</Button>
    <Button variant="contained" disabled={!podName.trim()} onClick={handleCreate}>
      {submitting ? <CircularProgress size={20} /> : "Create Pod"}
    </Button>
  </DialogActions>
</Dialog>
```

### Frontend: AppBar PodSelector Extension
The existing `PodSelector` component in `root.tsx` uses MUI `Select`. Extend it to add a `Divider` + "Create new pod" MenuItem at the bottom. When the "Create new pod" value is selected, reset the Select value to the current pod and open the CreatePod dialog instead of navigating.

```tsx
// Detect the special "create-new" sentinel value
const handleChange = (e: SelectChangeEvent) => {
  if (e.target.value === "create-new") {
    // reset to current pod, open modal
    onCreatePodRequested();
    return;
  }
  // normal navigation
  navigate(`/pod/${e.target.value}`);
};
// ...
<Select value={selectedPodId} onChange={handleChange} ...>
  {pods.map(p => <MenuItem key={p.id} value={String(p.id)}>{p.name}</MenuItem>)}
  <Divider />
  <MenuItem value="create-new" sx={{ color: "primary.main" }}>Create new pod</MenuItem>
</Select>
```

**Note:** The playing cards icon wrapping is trivial — `<Box component={Link} to="/">` or `<SvgIconPlayingCards component={Link} to="/">` won't work; use a `Link` wrapper with `style={{ display: 'flex' }}`.

### Frontend: /deck/new Route

Loader fetches formats + commanders in parallel (same pattern as `newGameLoader`):

```tsx
export async function newDeckLoader(): Promise<NewDeckData> {
  const [formats, commanders] = await Promise.all([GetFormats(), GetCommanders()]);
  return { formats, commanders };
}
```

freeSolo Autocomplete pattern for commander creation:
```tsx
<Autocomplete
  options={commanders}
  freeSolo
  getOptionLabel={(opt) => typeof opt === "string" ? opt : opt.name}
  filterOptions={(options, params) => {
    const filtered = filter(options, params);
    const { inputValue } = params;
    const isExisting = options.some(opt => inputValue === opt.name);
    if (inputValue !== "" && !isExisting) {
      filtered.push({ id: -1, name: `Create "${inputValue}"` } as Commander);
    }
    return filtered;
  }}
  onChange={async (_, value) => {
    if (typeof value === "string" || (value && value.id === -1)) {
      const name = typeof value === "string" ? value : inputValue;
      const res = await PostCommander(name);
      const { id } = await res.json();
      setCommanderId(id);
    } else {
      setCommanderId(value?.id ?? null);
    }
  }}
  renderInput={(params) => <TextField {...params} label="Commander" required />}
/>
```

### Backend: Pod-Scoped Stats (POD-04/D-22)

New SQL query following the pattern of `getStatsForPlayer` — adds a `INNER JOIN game ON game_result.game_id = game.id` and `AND game.pod_id = ?` filter:

```sql
-- getStatsForPlayersInPod: batch stats for all players in a pod
SELECT deck.player_id,
       game_result.game_id,
       game_result.place,
       game_result.kill_count,
       (SELECT COUNT(*) FROM game_result gr2
          WHERE gr2.game_id = game_result.game_id
            AND gr2.deleted_at IS NULL) AS player_count
  FROM game_result
  INNER JOIN deck ON game_result.deck_id = deck.id
  INNER JOIN game ON game_result.game_id = game.id
 WHERE game.pod_id = ?
   AND deck.player_id IN ?
   AND game.deleted_at IS NULL
   AND deck.deleted_at IS NULL
   AND game_result.deleted_at IS NULL;
```

Returns rows grouped by `deck.player_id` — same aggregation as `GetStatsForDecks` but keyed by player_id and filtered to a single pod.

**Interface addition** (`lib/repositories/interfaces.go` `GameResultRepository`):
```go
GetStatsForPlayersInPod(ctx context.Context, podID int, playerIDs []int) (map[int]*Aggregate, error)
```

**Business layer** (`GetAllByPod` in `lib/business/player/functions.go`): Replace the per-player `GetStatsForPlayer` call loop with one `GetStatsForPlayersInPod(ctx, podID, allPlayerIDs)` call. Same `PlayerWithRoleEntity` response shape — `stats` field now reflects pod-scoped stats when called with a `pod_id`.

**Frontend**: `GetPlayersForPod` in `http.ts` already passes `pod_id`. When the backend returns pod-scoped stats in the same response, no separate fetch is needed and no skeleton loading for stats is required. The Skeleton in the UI-SPEC is for the two-request approach — with a single request, player cards render with all data at once.

### Backend: API Contract Fix (Create Pod + Create Deck)

**Current state (bug):**
- `POST /api/pod` → `w.WriteHeader(http.StatusCreated)` — no body
- `POST /api/deck` (DeckCreate) → `w.WriteHeader(http.StatusCreated)` — no body
- Frontend `PostPod` calls `return res.json()` → undefined result → `navigate('/pod/undefined')` (latent bug)
- `PostDeck` function does not exist in `http.ts`

**Fix required:**

Backend `PodCreate` handler: after `p.pods.Create(...)` returns the new pod ID, write `{"id": newPodID}` in the response body with 201.

Backend `DeckCreate` handler: the `d.decks.Create(...)` call already returns `(int, error)` — the ID is discarded with `_`. Capture it and write `{"id": newDeckID}` in the response with 201.

Frontend:
- `PostPod` return type changes from `Pod` to `{ id: number }` (or keep `Pod` shape if backend returns full entity — but minimal `{"id": N}` is sufficient)
- Add `PostDeck(body: NewDeckRequest): Promise<{ id: number }>` to `http.ts`
- Add `NewDeckRequest` interface to `types.ts`

### Recommended Project Structure Additions

```
app/src/routes/deck/new/
└── index.tsx              # /deck/new route component + loader
```

No new backend files needed beyond modifying existing ones.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| freeSolo autocomplete with inline create | Custom dropdown | MUI Autocomplete `freeSolo` + `filterOptions` | MUI handles keyboard events, focus, blur, aria labels |
| Pod-scoped stats batch query | Per-player N+1 loop | Single SQL with `IN ?` + GROUP BY equivalent | N+1 already a known anti-pattern in this codebase; see getStatsForDecks batch pattern |
| DataGrid default sort | `onSortModelChange` + useState | `initialState.sorting.sortModel` | `initialState` is declarative; already used in PlayerDecksTab for name sort |
| Confirm dialog for Promote/Remove | Custom confirm component | Reuse existing MUI Dialog pattern from PlayersTab | Dialog pattern already established with shared `confirmAction` state |
| Commander creation in form | Separate page/route | Inline async create within Autocomplete `onChange` | PostCommander already exists; inline creation matches DX expected by user |

---

## Common Pitfalls

### Pitfall 1: MUI Select value mismatch on "Create new pod"
**What goes wrong:** When the user selects "Create new pod" from the PodSelector dropdown, the Select's `value` momentarily becomes `"create-new"` — if the component re-renders with that value, it shows an unmatched option and MUI logs a warning.
**Why it happens:** MUI Select `value` must match one of the rendered MenuItem values.
**How to avoid:** In the `handleChange` handler, immediately call `e.preventDefault()` pattern is not available; instead, store the current selected pod in a separate `ref` or reset state before the re-render. Safest approach: keep a `selectedValue` state that always reflects the last valid pod ID, and reset it in the handleChange when "create-new" is detected before opening the modal.

### Pitfall 2: PostCommander returns Response, not Commander
**What goes wrong:** `PostCommander` in `http.ts` returns `Response`, not a parsed object. Caller must call `.json()` and handle the new ID.
**Why it happens:** PostCommander was originally written to return the raw Response for flexibility.
**How to avoid:** In the Autocomplete onChange, `const res = await PostCommander(name); const { id } = await res.json();`. Add error handling around this call and surface "Failed to create commander" inline near the field.

### Pitfall 3: DeckCreate returns 201 no-body — PostDeck must be written correctly
**What goes wrong:** The current `DeckCreate` handler discards the new deck ID (`_, err = d.decks.Create(...)`). After fixing the backend to return the ID, the frontend must parse it or navigation to `/player/{id}/deck/{newId}` breaks.
**Why it happens:** Original create endpoints follow the convention of "201 no body" but the new navigation requirement changes this.
**How to avoid:** Backend fix: capture ID, write `{"id": deckID}` with 201. Frontend: `PostDeck` returns `Promise<{ id: number }>` and the route's submit handler reads `newDeckId` from the response.

### Pitfall 4: DataGrid initialState.sorting on server-side pagination
**What goes wrong:** The Pod Decks Tab uses `paginationMode="server"`. DataGrid `initialState.sorting` affects the client-side display order but does NOT send sort parameters to the server — the underlying data is still server-ordered.
**Why it happens:** Server-side pagination only fetches pages, not sorted data — the sort must be either client-applied to the loaded page, or passed as query params to the server.
**How to avoid:** D-28 explicitly scopes this as client-side sort on loaded page data — this is acceptable for the use case. No backend sort changes needed. The DataGrid renders sorted within the loaded page only. Document this constraint in the plan as "known limitation — full sort requires backend sort params, deferred."

### Pitfall 5: pod_id filter on GetStatsForPlayersInPod must also filter game.deleted_at
**What goes wrong:** Joining game_result → game without filtering `game.deleted_at IS NULL` would include results from soft-deleted games.
**Why it happens:** GORM soft-delete auto-filter only applies to the root model being queried, not to joined tables in raw SQL.
**How to avoid:** Include `AND game.deleted_at IS NULL` explicitly in the raw SQL query constant (as shown in the SQL example above).

### Pitfall 6: freeSolo Autocomplete `getOptionLabel` with string values
**What goes wrong:** When `freeSolo={true}`, MUI Autocomplete may pass the raw string (user typed value) as the option to `getOptionLabel`. If `getOptionLabel` assumes the option is always a `Commander` object (`opt.name`), it throws.
**Why it happens:** freeSolo mode passes both `Commander` objects (from options list) and raw `string` values when the user types something new.
**How to avoid:** `getOptionLabel={(opt) => typeof opt === "string" ? opt : opt.name}`.

---

## Code Examples

### Record Component Fix (D-32)
```tsx
// Before (stats.tsx line 13):
const maxPlace = Math.max(...Object.keys(record).map(Number), 1);
// After:
const maxPlace = Math.max(...Object.keys(record).map(Number), 4);
```
Source: CONTEXT.md D-32; confirmed by reading `app/src/components/stats.tsx`.

### DataGrid Default Sort (POD-03)
```tsx
// Add to DecksTab DataGrid:
initialState={{
  sorting: {
    sortModel: [{ field: "record", sort: "desc" }],
  },
}}
```
Source: `app/src/routes/player/DecksTab.tsx` line 47 (uses same pattern for name sort).

### Retired Deck Filter (DECK-03)
```tsx
// In PlayerDecksTab, filter before passing to DataGrid:
const visibleRows = (data ?? []).filter((d) => !d.retired);
// Pass visibleRows to DataGrid instead of data
```
Source: CONTEXT.md D-29; `app/src/routes/player/DecksTab.tsx` currently shows all decks including retired.

### Backend: Pod Create Returns ID
```go
// In PodCreate handler, replace:
w.WriteHeader(http.StatusCreated)
// With:
podID, err := p.pods.Create(ctx, e.Name, callerID)
// ...error handling...
w.WriteHeader(http.StatusCreated)
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(struct{ ID int `json:"id"` }{ID: podID})
```
Source: Code review of `lib/routers/pod.go` line 180-185; `p.pods.Create` already returns `(int, error)` but ID is discarded.

### Backend: GetAllByPod Pod-Scoped Stats
```go
// Replace the per-player stat loop in GetAllByPod:
playerIDs := make([]int, len(members))
for i, m := range members {
    playerIDs[i] = m.PlayerID
}
statsMap, err := gameResultRepo.GetStatsForPlayersInPod(ctx, podID, playerIDs)
// ...then in the result loop, use statsMap[m.PlayerID] instead of calling GetStatsForPlayer
```
Source: Code review of `lib/business/player/functions.go` lines 67-112.

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Pod creation in PlayerSettingsTab | AppBar dropdown + HomeView CTA | Phase 5 | SettingsTab Create Pod section removed; PostPod API bug becomes visible |
| Global player stats in pod view | Pod-scoped stats in pod view | Phase 5 | New SQL query + interface method required |
| No deck creation in UI | /deck/new route | Phase 5 | DeckCreate must return ID; PostDeck function added to http.ts |
| Record shows only actual places | Record pads to minimum 4 places | Phase 5 | stats.tsx one-line fix |

**Pre-existing bug surface:** `PostPod` calls `res.json()` on a 201-no-body response. Currently the Create Pod feature in PlayerSettingsTab fails to navigate correctly (navigates to `/pod/undefined`). Phase 5 fixes this by changing the backend to return the new ID.

---

## Open Questions

1. **Pod-scoped stats: single request or two requests?**
   - What we know: The `GET /api/players?pod_id={id}` endpoint currently calls `GetStatsForPlayer` (global) per player. D-23 says "use skeleton if two-request approach."
   - What's unclear: CONTEXT.md defers to researcher on optimal approach.
   - Recommendation: **Single request** — modify `GetAllByPod` to call the new batch `GetStatsForPlayersInPod` instead of per-player global stats. This means the stat row has data on first render, no skeleton needed. Simpler frontend, one fewer round-trip. The CONTEXT.md D-23 skeleton guidance is a fallback if two requests are used — with single request, omit the Skeleton.

2. **PostPod return type: full Pod entity or just `{id}`?**
   - What we know: Frontend only needs `id` for navigation. Backend Create function returns `int` (new ID only).
   - Recommendation: Return minimal `{"id": N}` — consistent with DeckCreate fix. No need to re-fetch the full Pod on create; the user navigates to `/pod/{id}` which triggers a full pod loader.

3. **DeckSettingsTab freeSolo update scope (D-19)**
   - What we know: CONTEXT.md says update the existing commander Autocomplete in DeckSettingsTab to support freeSolo/inline create.
   - What's unclear: The existing Autocomplete in DeckSettingsTab does NOT have freeSolo. Adding it changes the behavior for existing users who might be selecting from the list.
   - Recommendation: Add freeSolo following the same pattern as /deck/new. The "Create X" option only appears when the typed string is not in the list — existing list-selection behavior is unchanged.

---

## Environment Availability

Step 2.6: SKIPPED (no external dependencies — all changes are within the existing Go/React/MySQL stack already verified in production).

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go: testify/assert + testify/require (go test ./lib/...) |
| Framework | Frontend: react-scripts test (npm test from app/) |
| Config file | None for Go; react-scripts test uses CRA defaults |
| Quick run command | `go vet ./lib/...` (compile check, no DB needed) |
| Full suite command | `go test ./lib/...` |
| Frontend type check | `./node_modules/.bin/tsc --noEmit` from `app/` |

### Phase Requirements → Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| POD-01 | Pod creation accessible from AppBar/HomeView | Manual (UI interaction) | n/a — no automated frontend tests | manual-only |
| POD-02 | Empty state CTA opens modal, creates pod, navigates | Manual (UI interaction) | n/a | manual-only |
| POD-03 | Pod Decks tab default sort by Record desc | Manual (visual) | n/a | manual-only |
| POD-04 | Pod Players tab shows pod-scoped stats | `go vet ./lib/...` (compile); manual for UI | `go vet ./lib/...` | partial |
| DECK-01 | /deck/new creates deck, navigates to deck view | Manual (UI interaction) | n/a | manual-only |
| DECK-02 | Commander tooltip present | Already complete — no test needed | n/a | n/a |
| DECK-03 | Retired decks hidden in Player Decks tab | Manual (visual) | n/a | manual-only |
| Backend (POD-04) | GetStatsForPlayersInPod new method | `go test ./lib/repositories/gameResult/...` | existing test file | ❌ Wave 0 for new SQL |
| TypeScript | No type errors after all frontend changes | `./node_modules/.bin/tsc --noEmit` | ✅ existing |

### Sampling Rate
- **Per task commit:** `go vet ./lib/...` and `./node_modules/.bin/tsc --noEmit`
- **Per wave merge:** `go test ./lib/...` + manual smoke-test of affected views
- **Phase gate:** Full suite green + smoke test before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `lib/repositories/gameResult/repo_test.go` — add `TestGetStatsForPlayersInPod` covering: empty player list returns empty map, players with no games in pod return zero aggregates, players with games in pod return correct pod-scoped stats. Note: existing `repo_test.go` already exists; add new test function.

*(All other requirements are frontend/UI — no automated test infrastructure needed per project convention of no frontend automated tests.)*

---

## Project Constraints (from CLAUDE.md)

- **Stack locked:** Go + Gorilla Mux + GORM + MySQL; React + MUI v5 + React Router v6 — no framework changes
- **Auth:** Google OAuth only; JWT from context; `trackerHttp.CallerPlayerID(w, r)` for handler auth
- **No breaking changes** to existing game/player/deck data
- **Compile check:** `go vet ./lib/...` (not `go build ./...` — crawls node_modules)
- **Frontend type check:** `./node_modules/.bin/tsc --noEmit` from `app/` — never `npm run build` for checking
- **After API changes:** Run `/smoke-test` skill to verify core endpoints
- **Frontend components:** Default exports from `app/src/routes/<name>.tsx`; MUI throughout; all HTTP calls via `app/src/http.ts`; all TypeScript interfaces in `app/src/types.ts`
- **JSON field names:** snake_case matching DB column names
- **POST returns 201 no body** is the established convention for creates — Phase 5 breaks this convention for pod and deck creates to enable navigation. This is intentional and scoped.
- **Functional DI pattern:** Business constructors return typed closures; new repo method added to `GameResultRepository` interface in `interfaces.go` and to compile-time check in `repositories.go`

---

## Sources

### Primary (HIGH confidence)
- Direct code review of all relevant source files — confirmed via Read tool
  - `lib/routers/pod.go` — PodCreate returns 201 no body (line 185); confirmed bug
  - `lib/routers/deck.go` — DeckCreate returns 201 no body (line 203); confirmed bug; Create returns `(int, error)` with ID discarded
  - `lib/business/player/functions.go` — GetAllByPod uses per-player GetStatsForPlayer loop; confirmed N+1 optimization opportunity
  - `lib/repositories/gameResult/repo.go` — getStatsForDecks batch pattern; direct template for getStatsForPlayersInPod
  - `lib/repositories/interfaces.go` — GameResultRepository interface; needs new method
  - `app/src/http.ts` — PostCommander returns Response; PostDeck missing; PostPod calls res.json() on no-body endpoint
  - `app/src/components/stats.tsx` — Record component; RecordComparator; StatColumns
  - `app/src/routes/pod/DecksTab.tsx` — server-side pagination DataGrid; no initialState.sorting yet
  - `app/src/routes/pod/PlayersTab.tsx` — List/ListItem layout; confirmed for redesign
  - `app/src/routes/player/DecksTab.tsx` — no retired filter; no Add Deck button
  - `app/src/routes/player/SettingsTab.tsx` — Create New Pod section to remove; PostPod navigation pattern
  - `app/src/routes/deck/SettingsTab.tsx` — DECK-02 TooltipIcon confirmed present (line 164); Autocomplete without freeSolo
  - `app/src/routes/home/index.tsx` — empty state confirmed as plain Typography
  - `app/src/routes/root.tsx` — PodSelector with MUI Select; SvgIconPlayingCards not wrapped in Link
  - `app/src/index.tsx` — current route tree; no /deck/new route
  - `.planning/phases/05-pod-deck-ux/05-CONTEXT.md` — all locked decisions
  - `.planning/phases/05-pod-deck-ux/05-UI-SPEC.md` — full design contract
  - `.claude/skills/gorm/SKILL.md` — raw SQL + batch query patterns
  - `.claude/skills/react-router/SKILL.md` — loader/action/route patterns

### Secondary (MEDIUM confidence)
- MUI Autocomplete freeSolo documentation pattern — confirmed by existing usage in `app/src/routes/deck/SettingsTab.tsx` (non-freeSolo Autocomplete) + MUI v5 docs convention

### Tertiary (LOW confidence)
- None — all findings are from direct code inspection

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — verified against installed package.json and existing usage
- Architecture: HIGH — all patterns confirmed from existing code; no speculation
- Pitfalls: HIGH — all identified from actual code bugs (PostPod/PostDeck API mismatch) and known MUI behavior

**Research date:** 2026-03-26
**Valid until:** 2026-04-25 (stable stack; no fast-moving dependencies)
