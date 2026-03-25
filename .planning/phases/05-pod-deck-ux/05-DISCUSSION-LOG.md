# Phase 5: Pod & Deck UX - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-25
**Phase:** 05-pod-deck-ux
**Areas discussed:** Pod creation & onboarding, Deck creation flow, Players tab stats scope, Decks tab sort + retired decks, Record display, AppBar playing cards icon

---

## Folded Todos

| Todo | Decision |
|------|----------|
| Fix record display to default to four places | ✓ Folded into Phase 5 |
| Make playing cards icon clickable to go home on mobile | ✓ Folded into Phase 5 |
| Security review all API and frontend route authorization | Deferred to Phase 6 |

---

## Pod Creation & Onboarding

### Where should 'Create Pod' live?

| Option | Description | Selected |
|--------|-------------|----------|
| HomeView empty state only | Empty state becomes onboarding screen with Create + Join | |
| Pod page header + HomeView | Create on HomeView AND AppBar pod area for existing users | ✓ |
| Dedicated /new-pod route | Standalone route, linked from HomeView + pod header | |

**User's choice:** Pod page header + HomeView

### How should 'Create a Pod' work on HomeView empty state?

| Option | Description | Selected |
|--------|-------------|----------|
| Inline expand form | Expand inline text field + submit | |
| Modal dialog | Opens MUI Dialog with pod name field | ✓ |
| Navigate to /new-pod | Dedicated creation route | |

**User's choice:** Modal dialog

### Where on the pod page header should 'Create Pod' appear?

| Option | Description | Selected |
|--------|-------------|----------|
| + icon next to pod name | Small ⊕ icon button opens dialog | |
| Dropdown from pod name | Tapping pod name opens dropdown: all pods + Create new pod | ✓ |
| Floating action button | FAB opens create dialog | |

**User's choice:** Dropdown from pod name (doubles as pod switcher)

### Should HomeView empty state include 'Join with invite code'?

| Option | Description | Selected |
|--------|-------------|----------|
| Both on empty state | Create Pod + Join with code as parallel CTAs | |
| Create only on HomeView | Joining is via /join route | ✓ |

**User's choice:** Create only on HomeView

### Pod selector dropdown — pod switcher?

| Option | Description | Selected |
|--------|-------------|----------|
| Yes — pod switcher + create | All pods listed + Create new pod option | ✓ |
| No — create only | Dropdown offers only Create new pod | |

**User's choice:** Yes — full pod switcher

### After creating a pod from AppBar?

**User's choice:** Navigate to new pod (/pod/{newId})

### Player settings after moving Create Pod to AppBar?

**User's choice:** Keep Leave Pod in settings, remove Create Pod section

---

## Deck Creation Flow

### Where should 'Add Deck' live?

| Option | Description | Selected |
|--------|-------------|----------|
| Player view → Decks tab | Add Deck button on Player Decks tab | ✓ |
| Pod view → Decks tab | Button on Pod Decks tab | |
| Both | Button in both views | |

**User's choice:** Player view → Decks tab

### Form pattern?

| Option | Description | Selected |
|--------|-------------|----------|
| Modal dialog | MUI Dialog | |
| Inline expand | Expand below button | |
| Navigate to /deck/new route | Full-page dedicated route | ✓ |

**User's choice:** Navigate to /deck/new route

### Required fields?

**User's choice:** Name + Format always required. When Format = "Commander": Commander (required) + Partner Commander (optional). Commander fields use Autocomplete from `commander` table. New commanders can be created inline (freeSolo + `POST /api/commander`).

### After deck creation?

**User's choice:** Navigate to new deck's page (/player/{callerId}/deck/{newId})

### Commander autocomplete — how to create new?

| Option | Description | Selected |
|--------|-------------|----------|
| Autocomplete freeSolo — type + Enter | MUI freeSolo Autocomplete | ✓ |
| + button next to autocomplete | Opens New Commander dialog | |
| Fall back to deck settings | Create deck without commander, add later | |

**User's choice:** freeSolo autocomplete

### Update DeckSettingsTab commander autocomplete too?

**User's choice:** Yes — update both to support freeSolo/new commander creation

### React Router loader vs component-managed?

**User's choice:** React Router loader + navigate on submit

### Which formats trigger Commander fields?

**User's choice:** Only the "Commander" format (exact match)

### URL for /deck/new?

**User's choice:** /deck/new (top-level, not scoped to player)

### How does /deck/new know the player?

**User's choice:** Always for the logged-in caller (JWT callerPlayerID). Backend already enforces this.

### Add Deck button visibility?

**User's choice:** Only when viewer is the deck owner (isOwner check)

### Cancel/back button on /deck/new?

**User's choice:** Back button to /player/{callerId}. Also noted: add back/cancel to new game form too.

---

## Players Tab Stats Scope

### Global vs pod-scoped stats?

| Option | Description | Selected |
|--------|-------------|----------|
| Pod-scoped stats (games in this pod only) | New backend query required | ✓ |
| Global stats | Already in API response, no backend work | |

**User's choice:** Pod-scoped stats

### What stats to show?

**User's choice:** Record + Points + Kills, all pod-scoped

### Layout redesign?

**User's choice:** Card-per-player with flex layout (name + role on top row, stats below, icon buttons for Promote/Remove)

### Stats loading?

**User's choice:** Skeleton placeholder in each card while stats load

### Backend approach — new endpoint or enhance existing?

**User's choice:** Enhance existing GET /api/players?pod_id={id} — stats become pod-scoped when pod_id is present

### Replace or add pod_stats field?

**User's choice:** Replace — pod-scoped stats only when pod_id is present. Global stats unaffected on other endpoints.

### PlayerView shows?

**User's choice:** PlayerView always shows global stats (no change)

### Promote/Remove — buttons or icon buttons?

**User's choice:** TooltipIconButton (PersonAdd for Promote, PersonRemove for Remove)

### Leave Pod on Players tab card?

**User's choice:** No — Leave Pod stays in player settings

---

## Decks Tab Sort + Retired Decks

### Sort criteria for 'sorted by record'?

**User's choice:** Sort using the existing `RecordComparator` from stats.tsx (custom sort comparator, already handles dynamic place counts)

### Client-side vs server-side sort?

**User's choice:** Client-side sort on loaded page data

### Retired deck visibility?

| Option | Description | Selected |
|--------|-------------|----------|
| Hidden everywhere active, visible in game history | Current behavior extended consistently | ✓ |
| Hidden from new game form, shown with badge in tabs | Show retired with indicator | |
| Toggle on Decks tab | Show retired toggle | |

**User's choice:** Hidden everywhere active, visible in game history

### Player Decks tab — include retired?

**User's choice:** Filter retired decks client-side in the DataGrid (API may return them, but DataGrid filters by default)

### Game view result display?

**User's choice:** Always show deck name from game history (unaffected by retirement status)

---

## Record Display (Folded Todo)

### Minimum 4 places behavior?

**User's choice:** Record component always shows at least 4 places. Pads with 0 for missing places. A deck with only 2nd-place finishes shows "0 / 3 / 0 / 0" not "3". Handled internally in the Record component — no callers need to pass a prop.

---

## Playing Cards Icon (Folded Todo)

**User's choice:** Navigate to / (home). HomeView auto-redirects to first pod if user has pods.

---

## Additional Gray Areas

### New Game Form Cancel

**User's choice:** Cancel / back button navigates to /pod/:podId

### AppBar on non-pod pages

**User's choice:** AppBar works as it does today on non-pod pages (no pod selector, no changes)

### StatColumns Record column

**User's choice:** Record component handles 4-place minimum internally; StatColumns just renders `<Record record={...} />`

### Pod Decks tab Add Deck shortcut

**User's choice:** No — Add Deck lives on Player view only

---

## Claude's Discretion

- AppBar pod dropdown implementation (MUI Select, Menu, or Button + Popover)
- Card elevation, border, spacing for players tab cards
- DataGrid default sort implementation detail (initialState.sorting)
- Whether pod-scoped stats are one request or two (researcher to determine)
- Loading state granularity on /deck/new

## Deferred Ideas

- Pod Decks tab Add Deck shortcut
- Dark/light mode toggle
- Join with invite code on HomeView empty state
- Leave Pod on Players tab card
- PlayerView pod-scoped stats when visited from pod context
