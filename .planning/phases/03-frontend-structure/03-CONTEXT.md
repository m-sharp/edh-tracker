# Phase 3: Frontend Structure - Context

**Gathered:** 2026-03-23
**Status:** Ready for UI-SPEC then planning

<domain>
## Phase Boundary

Refactor the frontend from flat route files into per-view subdirectories with per-tab file extraction. Build shared components (TabbedLayout, TooltipIcon, TooltipIconButton) and move shared utilities to `app/src/components/`. Fix structural bugs (blank screen on refresh, HomeView loading flash). Run `/gsd:ui-phase 03` to produce a per-view UI-SPEC before planning — that spec drives view-level visual improvements (DSNG-04).

**In scope:**
- FEND-01: Route file restructuring into per-view subdirectories with per-tab files
- FEND-02: `TabbedLayout` component with query-string-persisted active tab
- FEND-03: `TooltipIcon` and `TooltipIconButton` shared components
- FEND-04: Fix blank white screen on page refresh
- FEND-05: Fix "No pods yet" flash before HomeView data loads
- DSNG-04: Per-view visual audit (Login, Home, Pod, Player, Deck, Deck, Game, Join views) — driven by UI-SPEC

**Out of scope:** Game form redesign (Phase 4), pod/deck functional gaps (Phase 5), auth interceptor (Phase 6).

</domain>

<decisions>
## Implementation Decisions

### File Organization (FEND-01)

- **D-01:** Route files split into per-view subdirectories under `app/src/routes/`. Each view gets its own directory with `index.tsx` as the main view + loader, and separate files per tab sub-component.

  Final structure:
  ```
  app/src/routes/
    pod/
      index.tsx         ← PodView + podLoader
      DecksTab.tsx
      PlayersTab.tsx
      GamesTab.tsx
      SettingsTab.tsx
    player/
      index.tsx
      OverviewTab.tsx
      DecksTab.tsx
      GamesTab.tsx
      SettingsTab.tsx
    deck/
      index.tsx
      OverviewTab.tsx
      GamesTab.tsx
      SettingsTab.tsx
    game/
      index.tsx         ← GameView + gameLoader
    new/
      index.tsx         ← NewGameView + newGameLoader + createGame action
    home/
      index.tsx         ← HomeView (moved out of index.tsx)
    login.tsx           ← stays flat (15 lines, no substructure needed)
    join.tsx            ← stays flat (73 lines, no substructure needed)
    RequireAuth.tsx     ← stays flat
    error.tsx           ← stays flat
  ```

- **D-02:** Shared utilities (`matches.tsx`, `stats.tsx`, `common.ts`) and new shared components (`TabbedLayout.tsx`, `TooltipIcon.tsx`) all move to `app/src/components/`. Import paths update accordingly.

- **D-03:** `HomeView` moves from inline in `index.tsx` to `app/src/routes/home/index.tsx`. `index.tsx` becomes purely routing + providers.

### Shared Tab Component — TabbedLayout (FEND-02)

- **D-04:** Component: `app/src/components/TabbedLayout.tsx`

  Interface:
  ```tsx
  interface TabConfig {
    id: string;          // used as query string value (e.g., 'decks', 'players')
    label: string;       // display text
    content: ReactElement;
    hidden?: boolean;    // omit from rendered tabs when true
  }

  interface TabbedLayoutProps {
    queryKey: string;    // query string key (e.g., 'podTab', 'playerTab', 'deckTab')
    tabs: TabConfig[];
    loading?: boolean;   // when true, shows CircularProgress in content area
  }
  ```

  Usage example:
  ```tsx
  <TabbedLayout
    queryKey="podTab"
    tabs={[
      { id: 'decks',    label: 'Decks',    content: <PodDecksTab /> },
      { id: 'players',  label: 'Players',  content: <PodPlayersTab /> },
      { id: 'games',    label: 'Games',    content: <PodGamesTab /> },
      { id: 'settings', label: 'Settings', content: <PodSettingsTab />, hidden: !isManager },
    ]}
  />
  ```

- **D-05:** Query string keys per view: `podTab`, `playerTab`, `deckTab`. Each is distinct — no shared key across views.

- **D-06:** Active tab resolved from query string: `?podTab=decks`. Falls back to first non-hidden tab if key is absent or value doesn't match any tab `id`. Out-of-range index also falls back to first tab.

- **D-07:** Tab switches use `useNavigate` with `replace: true` — no history stacking. Back/forward navigate between pages, not between tabs.

- **D-08:** Query string stores the tab `id` string (e.g., `podTab=decks`), not a numeric index. Changing a display label doesn't break bookmarks.

- **D-09:** When loading=true, render a centered `<CircularProgress>` in the content area instead of the active tab content.

- **D-10:** Scrollable tabs: `variant="scrollable" scrollButtons="auto"` — already established in `pod.tsx`, carried forward as TabbedLayout default.

- **D-11:** No `error` prop — error states handled at route level (React Router ErrorPage) or inline within tab components.

- **D-12:** `hidden` tabs are filtered out before rendering. Tab index is computed from the filtered list.

### Tooltip Components (FEND-03)

- **D-13:** Both components in `app/src/components/TooltipIcon.tsx`.

  **TooltipIcon** — informational icon with hover tooltip (tap toggles on mobile):
  ```tsx
  interface TooltipIconProps {
    title: string;
    icon?: ReactElement;  // defaults to <InfoOutlined />
  }
  ```
  Use for: help/info text next to fields. Not a button — clicking does nothing, only shows tooltip.

  **TooltipIconButton** — tappable icon button with hover tooltip (tap = click on mobile):
  ```tsx
  interface TooltipIconButtonProps {
    title: string;
    onClick: () => void;
    icon: ReactElement;
  }
  ```
  Use for: icon buttons whose tooltip provides additional context on hover.

- **D-14:** Mobile behavior:
  - `TooltipIcon`: tap toggles tooltip (MUI `enterTouchDelay={0}`)
  - `TooltipIconButton`: tap = click action (standard button behavior; tooltip is hover-only)

- **D-15:** Mobile principle (carried into UI-SPEC): if an icon button's purpose cannot be understood from the icon alone, use a text button instead. Icon buttons are for universally-understood actions.

- **D-16:** Known use case: `TooltipIcon` on the commander update field in `DeckView` (`deck/SettingsTab.tsx` after restructure): "This is for changing an existing deck's commander. To add a new deck, use the Add Deck button instead." The UI-SPEC audit may surface additional uses — apply wherever appropriate.

### Bug Fixes (FEND-04, FEND-05)

- **D-17 (FEND-04):** Blank white screen on page refresh — investigate root cause first. Most likely: the Go static file server doesn't catch-all to `index.html` for unknown routes. Fix: configure the static server to serve `index.html` for any path not matching a static file.

- **D-18 (FEND-05):** "No pods yet" flash — `HomeView` renders before the `GetPodsForPlayer` fetch resolves. Fix: add a loading state (or use React Router loader) so the message only appears after data has loaded and confirmed empty.

### UI-SPEC Step (DSNG-04)

- **D-19:** Run `/gsd:ui-phase 03` before `/gsd:plan-phase 03`. The UI-SPEC drives view-level visual improvements.

- **D-20:** Views to audit: Login, Home, Pod, Player, Deck, Game, Join (7 views). NewGameView excluded — Phase 4 redesigns it entirely.

- **D-21:** UI-SPEC format: structured issue list per view. Each issue tagged with type (layout, spacing, typography, interaction). Each entry: current behavior → desired behavior. Actionable by planner.

- **D-22:** Audit scope: mobile-first (375px primary). Full visual audit: layout, spacing, typography, interactions (hover states, focus states, dialog styling, touch targets).

- **D-23:** UI-SPEC researcher must read Phase 2 artifacts before auditing:
  - `.planning/phases/02-design-language/02-CONTEXT.md` — design decisions
  - `.planning/phases/02-design-language/02-UI-SPEC.md` — established design system

- **D-24:** UI-SPEC focuses on visual issues only — not structural changes (TabbedLayout migration, file restructuring). Those are captured here in CONTEXT.md.

- **D-25:** Issues only — no "things working well" section. Implementers preserve working patterns by not changing them.

### Claude's Discretion

- Exact file and import path organization for `app/src/index.tsx` imports after restructure — keep it clean
- Whether `TabbedLayout` needs a `TabPanel`-style wrapper with `role="tabpanel"` for accessibility
- MUI `Tooltip` props beyond `enterTouchDelay` (e.g., `leaveDelay`, `placement`) for TooltipIcon mobile behavior
- Whether `AsyncComponentHelper` in `common.ts` stays or gets refactored during the move to `components/`

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase Requirements
- `.planning/REQUIREMENTS.md` §Frontend Structure — FEND-01 through FEND-05
- `.planning/REQUIREMENTS.md` §Design — DSNG-04
- `.planning/ROADMAP.md` §Phase 3 — goal, success criteria, `UI hint: yes`

### Phase 2 Design System (required for UI-SPEC and view polish)
- `.planning/phases/02-design-language/02-CONTEXT.md` — color palette, typography, spacing decisions
- `.planning/phases/02-design-language/02-UI-SPEC.md` — established MUI theme and component patterns

### Existing Frontend
- `app/src/index.tsx` — current router definition + HomeView inline; restructure target
- `app/src/routes/pod.tsx` — 355-line monolith; split into pod/ + tab files
- `app/src/routes/player.tsx` — 253-line monolith; split into player/ + tab files
- `app/src/routes/deck.tsx` — 333-line monolith; split into deck/ + tab files
- `app/src/routes/game.tsx` — 410-line view; moves to game/index.tsx
- `app/src/routes/new.tsx` — 240-line new game form; moves to new/index.tsx
- `app/src/routes/root.tsx` — AppBar/nav; no structural change needed
- `app/src/common.ts` — `AsyncComponentHelper`; moves to `app/src/components/`
- `app/src/stats.tsx` — shared stat components; moves to `app/src/components/`
- `app/src/matches.tsx` — match display; moves to `app/src/components/`
- `app/src/theme.ts` — MUI theme; stays in `app/src/`

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `AsyncComponentHelper` (`common.ts`): data-fetch helper returning `{ data, loading, error }` — used in `player.tsx` and `deck.tsx`. Move to `components/` without changing the API.
- `StatColumns`, `Record`, `CommanderColumn` (`stats.tsx`): shared DataGrid column defs and stat display — move to `components/`.
- `MatchesDisplay` (`matches.tsx`): game result display — move to `components/`.

### Established Patterns
- Tab state in all 3 views currently uses `useState(0)` — replace with `TabbedLayout` across all views.
- Tab scrollability already established: `pod.tsx` uses `variant="scrollable" scrollButtons="auto"` — carry forward as TabbedLayout default.
- `useLoaderData()` pattern: data fetching via React Router loaders (not useEffect) — preserve in restructured files.
- `LoaderFunctionArgs` imported from `@remix-run/router/utils` in game.tsx and pod.tsx — note this import source for consistency.

### Integration Points
- `app/src/index.tsx`: Import paths for all routes update after restructure (e.g., `./routes/pod` → `./routes/pod/index` or `./routes/pod`).
- All tab sub-components in the new per-tab files import their shared dependencies from `../../../components/` (3 levels up from routes/pod/DecksTab.tsx).
- `app/src/http.ts` and `app/src/types.ts` stay in `app/src/` — no move needed.
- `app/src/auth.tsx` stays in `app/src/` — no move needed.

</code_context>

<specifics>
## Specific Ideas

- TabbedLayout URL hygiene: use `replace: true` to prevent history bloat; query string key is only written on explicit tab switch, not on initial render.
- `?podTab=decks` not `?tab=decks` — distinct keys per view to prevent contamination if query strings are ever preserved during navigation.
- `TooltipIconButton` tap behavior: tap = click (not toggle). On mobile, the tooltip hover context is secondary — the action is primary.
- Icon button mobile principle: if the icon isn't self-explanatory, use a text button. TooltipIconButton is for icons whose meaning is clear; the tooltip is supplemental.

</specifics>

<deferred>
## Deferred Ideas

- Dark/light mode toggle — out of scope for launch (was Phase 2 deferred)
- JoinView restructure — small file, no splitting needed, but the UI-SPEC audit will catch visual issues
- NewGameView UI improvements — excluded from UI-SPEC audit; Phase 4 redesigns it entirely
- `AsyncComponentHelper` refactor — move it as-is; whether to modernize it is Phase 3+ discretion

</deferred>

---

*Phase: 03-frontend-structure*
*Context gathered: 2026-03-23*
