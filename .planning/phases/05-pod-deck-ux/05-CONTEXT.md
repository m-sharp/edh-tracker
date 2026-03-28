# Phase 5: Pod & Deck UX - Context

**Gathered:** 2026-03-25
**Status:** Ready for planning

<domain>
## Phase Boundary

Complete pod and deck feature gaps and onboarding flow. The pod creation UX moves off player settings into the AppBar and HomeView. New users get a real empty-state CTA. Players can create decks from a dedicated /deck/new route. The Pod Players tab is redesigned with pod-scoped stats. The Pod Decks tab sorts by Record by default. Retired deck visibility is defined and consistent. Two folded todos are also delivered: Record component minimum 4-place display, and playing cards icon clickable to navigate home.

**In scope:**
- POD-01: Pod creation accessible from pod page (AppBar pod-name dropdown + HomeView)
- POD-02: New user onboarding — HomeView empty state has real CTA
- POD-03: Pod Decks tab sorted by Record by default (uses RecordComparator)
- POD-04: Pod Players tab redesigned with pod-scoped stats (record + points + kills)
- DECK-01: Player can create a new deck via /deck/new route
- DECK-02: Commander update tooltip — **already implemented** in DeckSettingsTab.tsx (no work needed)
- DECK-03: Retired deck visibility defined and consistent
- Record display min 4 places (folded todo)
- Playing cards icon clickable → home (folded todo)

**Out of scope:** Auth interceptor (Phase 6), test coverage (Phase 6), production readiness (Phase 7).

</domain>

<decisions>
## Implementation Decisions

### Pod Creation Entry Point (POD-01)

- **D-01:** Pod creation moves off player settings (PlayerSettingsTab). `PlayerSettingsTab` keeps "Your Pods" list + Leave Pod, but removes the "Create New Pod" section.

- **D-02:** Pod creation is accessible from two places:
  1. **HomeView empty state** — for first-time users with no pods
  2. **AppBar pod-name dropdown** — for existing users who want to create an additional pod

- **D-03:** The AppBar pod name becomes a dropdown (MUI `Select` or `Menu`) when on `/pod/:podId`. The dropdown lists all pods the user belongs to as links (`/pod/{id}`) and includes a "Create new pod" option at the bottom. Tapping a pod name navigates to that pod. This also serves as a **pod switcher**.

- **D-04:** After creating a pod from anywhere, auto-navigate to `/pod/{newPodId}`.

- **D-05:** The pod-name dropdown only appears on `/pod/:podId`. All other views (PlayerView, DeckView, GameView) keep the AppBar as-is today.

### New User Onboarding (POD-02)

- **D-06:** HomeView empty state becomes a real onboarding screen with a "Create a Pod" button. Clicking opens a **Modal dialog** (MUI Dialog) with a pod name field + submit.

- **D-07:** HomeView does NOT include a "Join with invite link" CTA. Joining is via the /join route (following an invite link). The empty state is focused on creation.

- **D-08:** The "Create Pod" modal is also triggered from the AppBar dropdown's "Create new pod" option (same dialog component reused).

### Deck Creation (DECK-01)

- **D-09:** "Add Deck" button appears on the Player view → Decks tab, only when the viewer is the deck owner (`user.player_id === player.id`). Not shown on pod view.

- **D-10:** "Add Deck" navigates to `/deck/new` — a dedicated top-level route (not scoped in URL to the player).

- **D-11:** `/deck/new` uses React Router loader + `useNavigate` on submit. Loader fetches formats and commanders lists.

- **D-12:** After successful deck creation, navigate to `/player/{callerId}/deck/{newDeckId}`.

- **D-13:** `/deck/new` includes a Cancel / back button that navigates to `/player/{callerId}`.

- **D-14:** The deck always belongs to the logged-in caller — no player ID in the URL. Backend uses JWT `callerPlayerID` exclusively.

#### New Deck Form Fields

- **D-15:** Required: **Name** and **Format** (always required).
- **D-16:** When Format = **"Commander"**: Commander (required) and Partner Commander (optional) fields appear. For all other formats, these fields are hidden.
- **D-17:** Commander and Partner Commander are MUI Autocomplete with `freeSolo={true}`. Typing a name not in the list and pressing Enter (or selecting the "Create X" option) calls `POST /api/commander` to create the commander inline, then uses the new commander ID.
- **D-18:** After creation, navigate to `/player/{callerId}/deck/{newDeckId}`.

#### Commander Autocomplete Update (DeckSettingsTab)

- **D-19:** The existing commander Autocomplete in `DeckSettingsTab.tsx` (for updating an existing deck's commander) is also updated to support freeSolo / new commander creation. Consistent behavior across all commander Autocomplete instances in the app.

### Pod Players Tab Redesign (POD-04)

- **D-20:** Players tab switches from a `List` layout to a **card-per-player layout** (Box/Card with flex). Each card:
  ```
  ┌────────────────────────────┐
  | Mike       [Manager]        |
  | 5W/3/2 • 42pts • 12 kills  |
  |                  [👤+][👤-] |  ← icon buttons, manager-only
  └────────────────────────────┘
  ```

- **D-21:** Stats shown per player: **Record + Points + Kills**, all **pod-scoped** (games played within this specific pod only).

- **D-22:** Pod-scoped stats require **backend enhancement**: `GET /api/players?pod_id={id}` is updated to compute stats filtered by `pod_id` (join `game_result` through `game` table, filter `game.pod_id = podID`). When `pod_id` is present in the query, the `stats` field in the response reflects pod-scoped stats only (not global). Global stats are unaffected for other callers (no `pod_id` filter).

- **D-23:** Stats loading uses skeleton placeholders — player cards render immediately with name + role; the stats row shows a skeleton while the stats are being fetched. (If stats are included in the same endpoint response as the player list, skeleton is only needed if the data is fetched in a second call; researcher to determine optimal approach.)

- **D-24:** Promote/Remove action buttons become `TooltipIconButton` components (PersonAdd for Promote, PersonRemove/PersonOff for Remove). Manager-only, not shown for the caller's own card.

- **D-25:** Player names are links to `/player/{id}` (existing behavior preserved).

- **D-26:** Leave Pod stays in player settings. No Leave option on the Players tab cards.

### Pod Decks Tab Default Sort (POD-03)

- **D-27:** Pod Decks tab uses the existing `RecordComparator` from `app/src/components/stats.tsx` as the sort comparator on the Record column.

- **D-28:** DataGrid `initialState.sorting` set to sort by the Record column descending. Client-side sort on the loaded page data — no backend sort changes needed.

### Retired Deck Visibility (DECK-03)

- **D-29:** Retired decks are **hidden everywhere active**:
  - Pod Decks tab: already filters `retired=false` via backend — no change needed
  - Player Decks tab: currently returns all decks including retired; **filter client-side** in the DataGrid to exclude `retired=true` rows by default
  - New game form deck picker: retired decks are excluded (current backend behavior, no change)

- **D-30:** Retired decks **remain visible in game history** (game results store deck names as strings; historical records are unaffected by retirement status).

- **D-31:** Game view result display always shows deck names from stored game results, regardless of current retirement status. No "Retired" badge on historical game results.

### Record Display — Minimum 4 Places (Folded Todo)

- **D-32:** The `Record` component in `app/src/components/stats.tsx` is updated to always display at least 4 place columns, padding missing places with 0.
  - Example: a deck with only 2nd-place finishes currently shows `3` — after fix shows `0 / 3 / 0 / 0`
  - Example: a 3-player-only deck shows `5 / 3 / 2 / 0` (pad 4th place)
  - Example: a deck with a 5th-place result shows all 5 places: `5 / 3 / 2 / 1 / 1`
  - This is handled **internally in the Record component** — no prop changes needed in callers.
  - The `RecordComparator` already handles dynamic place counts — no change needed there.

### Playing Cards Icon → Home (Folded Todo)

- **D-33:** The `SvgIconPlayingCards` icon in the AppBar (`app/src/routes/root.tsx`) is wrapped in a `Link` or `component={Link}` pattern to navigate to `/`. HomeView then auto-redirects to the user's first pod if they have one, or shows the empty state.

### New Game Form — Cancel / Back Button

- **D-34:** The new game form (`/pod/:podId/new-game`) gets a Cancel / back button. Clicking it navigates back to `/pod/:podId`.

### Claude's Discretion

- AppBar pod dropdown implementation detail: MUI `Select`, `Menu`, or `Button + Popover` — whichever fits the AppBar layout cleanest at 375px
- Exact card elevation, border, and spacing for the player tab redesign
- Whether to use `initialState.sorting` or `onSortModelChange` for DataGrid default sort
- Whether pod-scoped stats are fetched in the same request as the player list (same endpoint, same round-trip) or as a second request with skeleton loading — researcher to determine optimal approach
- Loading state granularity on /deck/new (loader vs component mount fetch)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase Requirements
- `.planning/REQUIREMENTS.md` §Pods — POD-01, POD-02, POD-03, POD-04
- `.planning/REQUIREMENTS.md` §Decks — DECK-01, DECK-02, DECK-03
- `.planning/ROADMAP.md` §Phase 5 — goal, success criteria

### Design System (required for all UI work)
- `.planning/phases/02-design-language/02-CONTEXT.md` — color palette, typography, spacing decisions
- `.planning/phases/02-design-language/02-UI-SPEC.md` — MUI theme and component patterns

### Files to modify (frontend)
- `app/src/routes/home/index.tsx` — HomeView empty state → onboarding CTA
- `app/src/routes/root.tsx` — AppBar: pod-name dropdown (pod switcher + create), playing cards icon → link
- `app/src/routes/pod/PlayersTab.tsx` — redesign to card layout with pod-scoped stats
- `app/src/routes/pod/DecksTab.tsx` — default sort by Record using RecordComparator
- `app/src/routes/player/DecksTab.tsx` — Add Deck button (owner-only), filter retired decks client-side
- `app/src/routes/player/SettingsTab.tsx` — remove Create Pod section, keep Leave Pod
- `app/src/routes/deck/SettingsTab.tsx` — update commander Autocomplete to support freeSolo
- `app/src/routes/new/index.tsx` — add Cancel / back button to new game form
- `app/src/components/stats.tsx` — Record component: minimum 4 places
- `app/src/http.ts` — add `PostDeck` function (if not present); PostCommander already exists
- `app/src/types.ts` — add types for new deck creation request if needed
- `app/src/index.tsx` — add `/deck/new` route

### Files to create (frontend)
- `app/src/routes/deck/new/index.tsx` — new deck creation route + loader

### Files to modify (backend)
- `lib/routers/player.go` — `GetPlayers` handler: when `pod_id` is present, compute pod-scoped stats
- `lib/business/player/functions.go` — `GetAll` function: add pod-scoped stats aggregation when podID is provided
- `lib/repositories/gameResult/repo.go` — add `GetStatsForPlayerInPod` or `GetStatsForPlayersInPod` (batch) query

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `TabbedLayout` (`app/src/components/TabbedLayout.tsx`): already used in all views — no change needed
- `TooltipIcon` / `TooltipIconButton` (`app/src/components/TooltipIcon.tsx`): use TooltipIconButton for Promote/Remove in players tab
- `RecordComparator` (`app/src/components/stats.tsx`): existing sort comparator, just needs to be wired into DataGrid column definition
- `Record` component (`app/src/components/stats.tsx`): update internally for 4-place minimum
- `StatColumns` (`app/src/components/stats.tsx`): DataGrid column defs including Record column — update sortComparator on Record column
- `SvgIconPlayingCards` (`app/src/components/SvgIconPlayingCards.tsx`): already in AppBar; wrap with Link
- `PostCommander` (`app/src/http.ts`): already exists — `POST /api/commander` creates a new commander

### Established Patterns
- React Router loader pattern: use for `/deck/new` (consistent with `newGameLoader` in `routes/new/index.tsx`)
- `component={Link}` on MUI Button: established pattern for navigation buttons (Phase 3 + 4)
- `isOwner` check: used in DeckView Settings tab hiding — reuse same pattern for "Add Deck" button visibility
- `window.location.reload()`: used after mutations in several views — avoid in new code where possible; use `useNavigate` or `useLoaderData` patterns
- MUI Dialog for confirmations: established in PlayersTab (Promote/Remove), DeckSettingsTab (Retire/Delete)
- `AsyncComponentHelper`: used in PlayerSettingsTab, DeckSettingsTab — available if needed

### Integration Points
- `app/src/index.tsx`: add `/deck/new` route with loader
- `lib/routers/player.go` `GetPlayers`: currently calls `business.GetAll(ctx, podID)` — needs to be wired to pod-scoped stat computation
- `lib/repositories/gameResult/repo.go`: `getStatsForPlayer` SQL can be extended with a `pod_id` join on the `game` table; analogous to existing pattern
- Backend interface `lib/repositories/interfaces.go`: new `GetStatsForPlayersInPod` method needs to be added and implemented

### Current State Notes
- **DECK-02 is already complete**: `DeckSettingsTab.tsx` already has `<TooltipIcon title="This is for changing an existing deck's commander. To add a new deck, use the Add Deck button instead." />`. No work needed.
- `PodDecksTab.tsx`: uses DataGrid with server-side pagination — `initialState.sorting` sets default sort without breaking pagination
- `PlayersTab.tsx`: current `PlayerWithRole extends Player` includes global stats (`stats` field) — the API change will swap these for pod-scoped stats when `pod_id` is in the query
- `PlayerSettingsTab.tsx`: "Create New Pod" section is below a `<Divider>` — clean to remove; "Your Pods" list with Leave stays

</code_context>

<specifics>
## Specific Ideas

- Pod switcher dropdown doubles as pod creation entry point — tapping the pod name in the AppBar shows all the user's pods as links + "Create new pod" at the bottom
- `/deck/new` form: Commander + Partner Commander fields conditionally shown only when Format = "Commander" (exact string match against format name)
- freeSolo autocomplete for commander creation: typing a new name and pressing Enter calls `PostCommander(name)` to create it first, then uses the returned ID for the deck create call
- Record component fix: always pad to `Math.max(maxKey, 4)` places — so 5+ player games still show all places correctly
- Playing cards icon: wrap in `<Link to="/">` using `component={Link}` MUI pattern (established in Phase 3/4)
- New game form back button: navigate to `/pod/:podId` — the podId is available in the route params

</specifics>

<deferred>
## Deferred Ideas

- Pod Decks tab "Add Deck" shortcut — decided against; Add Deck lives on Player view only
- Dark/light mode toggle — out of scope for launch (persisted from Phase 2/3)
- "Join with invite code" on HomeView empty state — deferred; joining is via /join route (invite link)
- Leave Pod option on Players tab card — deferred; Leave Pod stays in Player Settings
- PlayerView showing pod-scoped stats when visited from pod context — deferred; PlayerView always shows global stats

### Reviewed Todos (not folded)
- **Security review all API and frontend route authorization** — reviewed, belongs in Phase 6 (auth/session phase). Score: 0.6, deferred.

</deferred>

---

*Phase: 05-pod-deck-ux*
*Context gathered: 2026-03-25*
