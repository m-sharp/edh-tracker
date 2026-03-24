# Phase 3: Frontend Structure - Research

**Researched:** 2026-03-24
**Domain:** React 18 + TypeScript + React Router v6 + MUI v5 frontend restructuring
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**File Organization (FEND-01)**
- D-01: Route files split into per-view subdirectories under `app/src/routes/`. Each view gets its own directory with `index.tsx` as the main view + loader, and separate files per tab sub-component. Final structure: `pod/`, `player/`, `deck/`, `game/`, `new/`, `home/` subdirectories; `login.tsx`, `join.tsx`, `RequireAuth.tsx`, `error.tsx` stay flat.
- D-02: Shared utilities (`matches.tsx`, `stats.tsx`, `common.ts`) and new shared components move to `app/src/components/`. Import paths update accordingly.
- D-03: `HomeView` moves from inline in `index.tsx` to `app/src/routes/home/index.tsx`. `index.tsx` becomes purely routing + providers.

**Shared Tab Component — TabbedLayout (FEND-02)**
- D-04: Component at `app/src/components/TabbedLayout.tsx` with `TabConfig` / `TabbedLayoutProps` interfaces as specified.
- D-05: Query string keys per view: `podTab`, `playerTab`, `deckTab`.
- D-06: Active tab resolved from query string; falls back to first non-hidden tab.
- D-07: Tab switches use `useNavigate` with `replace: true`.
- D-08: Query string stores the tab `id` string (e.g., `podTab=decks`), not a numeric index.
- D-09: When `loading=true`, render centered `<CircularProgress>` in content area.
- D-10: `variant="scrollable" scrollButtons="auto"` as TabbedLayout default.
- D-11: No `error` prop — error states handled at route/tab level.
- D-12: `hidden` tabs filtered before rendering; tab index computed from filtered list.

**Tooltip Components (FEND-03)**
- D-13: Both `TooltipIcon` and `TooltipIconButton` in `app/src/components/TooltipIcon.tsx`.
- D-14: `TooltipIcon` — `enterTouchDelay={0}` (tap toggles). `TooltipIconButton` — tap = click (tooltip hover-only).
- D-15: Mobile principle — if icon's meaning is unclear, use text button instead.
- D-16: Known use: `TooltipIcon` on commander update field in `DeckSettingsTab`.

**Bug Fixes**
- D-17 (FEND-04): Blank white screen on refresh — fix in Go static server's SPA fallback.
- D-18 (FEND-05): Loading flash in HomeView — add loading state; only show "No pods yet" after fetch resolves.

**UI-SPEC Step (DSNG-04)**
- D-19 through D-25: UI-SPEC drives per-view visual improvements. 7 views audited (Login, Home, Pod, Player, Deck, Game, Join). NewGameView excluded (Phase 4). Issues documented in `03-UI-SPEC.md`.

### Claude's Discretion
- Exact file and import path organization for `app/src/index.tsx` imports after restructure
- Whether `TabbedLayout` needs `role="tabpanel"` wrapper for accessibility
- MUI `Tooltip` props beyond `enterTouchDelay` for `TooltipIcon` mobile behavior
- Whether `AsyncComponentHelper` in `common.ts` stays or gets refactored during the move to `components/`

### Deferred Ideas (OUT OF SCOPE)
- Dark/light mode toggle
- JoinView restructure (structural split — UI audit issues from UI-SPEC ARE in scope)
- NewGameView UI improvements
- `AsyncComponentHelper` refactor beyond move-as-is
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| FEND-01 | Large route files refactored into per-view subdirectories | File-by-file move analysis below; import path mapping documented |
| FEND-02 | Shared tab component with query-string-persisted active tab | `useSearchParams` + `useNavigate` pattern confirmed; component interface fully specified in CONTEXT.md |
| FEND-03 | Shared tooltip icon and tooltip icon button components | MUI `Tooltip` + `IconButton` patterns confirmed; `InfoOutlined` already installed |
| FEND-04 | Page refresh no longer causes blank white screen | Root cause confirmed: `app/main.go`'s `spaHandler` already has SPA fallback — bug is NOT in the static server (details below) |
| FEND-05 | HomeView no longer flashes "No pods yet" before data loads | Root cause confirmed in `index.tsx`; loading-state pattern documented |
| DSNG-04 | All views individually audited against Phase 2 design language | Full per-view issue list in `03-UI-SPEC.md` (approved); implementation checklist is authoritative |
</phase_requirements>

---

## Summary

This phase is a frontend-only refactor + polish pass. No backend changes are required. The codebase has been fully read; all six requirements are well-understood with clear implementation paths.

The most important findings: (1) FEND-04's blank-screen-on-refresh root cause needs investigation — `app/main.go` already has a working SPA fallback (`spaHandler`), so the real culprit may be `RequireAuth`'s render sequence or a race condition in auth loading rather than the static server. (2) The `TabbedLayout` component is purely additive — no existing code depends on it until the three monolithic route files are migrated. (3) All new components use only already-installed packages — no new npm dependencies needed.

**Primary recommendation:** Sequence work as: Wave 0 (new shared components + HomeView fix) → Wave 1 (per-view restructuring for pod/player/deck) → Wave 2 (game/new moves + UI audit fixes per UI-SPEC). TypeScript must be verified after each wave using `cd app && ./node_modules/.bin/tsc --noEmit`.

---

## Project Constraints (from CLAUDE.md)

- Tech stack frozen: React 18 + TypeScript + MUI v5 + React Router v6. No library changes.
- No breaking changes to existing game/player/deck data.
- Frontend verify: `cd app && ./node_modules/.bin/tsc --noEmit` — never `npm run build` or `npx tsc`.
- Never `go build ./...` or `go build ./lib/...` — use `go vet ./lib/...` for compile checks.
- HTTP calls must go through `app/src/http.ts`; never call `fetch` directly from components.
- All TypeScript interfaces in `app/src/types.ts`.
- `credentials: "include"` on all fetch calls.
- MUI components from `@mui/material/styles` (sub-path), not `@mui/material` directly for theme utilities.
- Components are default exports from route files. Route components consume loader data via `useLoaderData()`.
- After API changes: run `/smoke-test` skill.
- After `app/src/` changes: run `frontend-verify` skill.

---

## Standard Stack

### Core (Already Installed — No New Dependencies)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| React | 18.0.0 | Component framework | Project constraint |
| TypeScript | 4.9.5 | Type safety | Project constraint |
| React Router DOM | 6.21.1 | Routing + query string access | Project constraint; `useSearchParams` used for FEND-02 |
| MUI `@mui/material` | 5.15.2 | UI components (Tabs, Tooltip, IconButton, CircularProgress) | Project constraint |
| `@mui/icons-material` | 5.15.3 | `InfoOutlined` for TooltipIcon default icon | Already installed |

### No New Packages Required

The UI-SPEC confirms: "No third-party registries. No new packages required. All Phase 3 UI work uses existing installed dependencies."

**Installation:** None needed.

---

## Architecture Patterns

### Recommended Project Structure (Post-Restructure)

```
app/src/
├── components/
│   ├── TabbedLayout.tsx      ← new shared component (FEND-02)
│   ├── TooltipIcon.tsx       ← new shared components (FEND-03)
│   ├── common.ts             ← moved from app/src/common.ts
│   ├── stats.tsx             ← moved from app/src/stats.tsx
│   └── matches.tsx           ← moved from app/src/matches.tsx
├── routes/
│   ├── pod/
│   │   ├── index.tsx         ← PodView + podLoader (moved from pod.tsx)
│   │   ├── DecksTab.tsx
│   │   ├── PlayersTab.tsx
│   │   ├── GamesTab.tsx
│   │   └── SettingsTab.tsx
│   ├── player/
│   │   ├── index.tsx         ← PlayerView (moved from player.tsx)
│   │   ├── OverviewTab.tsx
│   │   ├── DecksTab.tsx
│   │   ├── GamesTab.tsx
│   │   └── SettingsTab.tsx
│   ├── deck/
│   │   ├── index.tsx         ← DeckView (moved from deck.tsx)
│   │   ├── OverviewTab.tsx
│   │   ├── GamesTab.tsx
│   │   └── SettingsTab.tsx
│   ├── game/
│   │   └── index.tsx         ← GameView + gameLoader (moved from game.tsx)
│   ├── new/
│   │   └── index.tsx         ← NewGameView + newGameLoader + createGame (moved from new.tsx)
│   ├── home/
│   │   └── index.tsx         ← HomeView (extracted from index.tsx)
│   ├── login.tsx             ← stays flat
│   ├── join.tsx              ← stays flat
│   ├── RequireAuth.tsx       ← stays flat
│   └── error.tsx             ← stays flat
├── index.tsx                 ← routing config + providers only (no HomeView inline)
├── http.ts                   ← stays
├── types.ts                  ← stays
├── auth.tsx                  ← stays
└── theme.ts                  ← stays
```

### Import Path Changes After Restructure

| Old Path | New Path | Used In |
|----------|----------|---------|
| `../stats` | `../../components/stats` | pod/*, player/*, deck/* tabs |
| `../matches` | `../../components/matches` | player/GamesTab, deck/GamesTab |
| `../common` | `../../components/common` | player/*, deck/* tabs |
| `../auth` | `../../auth` | pod/index, player/index, etc. |
| `../http` | `../../http` | all tab files |
| `../types` | `../../types` | all tab files |
| `./routes/pod` | `./routes/pod/index` (or `./routes/pod`) | app/src/index.tsx |
| `./routes/game` | `./routes/game/index` (or `./routes/game`) | app/src/index.tsx |

Note: TypeScript/webpack resolves `./routes/pod` to `./routes/pod/index.tsx` automatically when the directory has an `index.tsx`. Either path form works in `app/src/index.tsx`.

### Pattern 1: TabbedLayout with Query-String Persistence

**What:** A shared component encapsulating MUI Tabs + tab content rendering, driven by the URL query string instead of local state.

**When to use:** Any view with multiple tabs (Pod, Player, Deck).

**Implementation pattern:**
```tsx
// app/src/components/TabbedLayout.tsx
import { ReactElement } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Box, CircularProgress, Tab, Tabs } from "@mui/material";

interface TabConfig {
    id: string;
    label: string;
    content: ReactElement;
    hidden?: boolean;
}

interface TabbedLayoutProps {
    queryKey: string;
    tabs: TabConfig[];
    loading?: boolean;
}

export default function TabbedLayout({ queryKey, tabs, loading }: TabbedLayoutProps): ReactElement {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    const visibleTabs = tabs.filter((t) => !t.hidden);
    const activeId = searchParams.get(queryKey);
    const activeIndex = Math.max(0, visibleTabs.findIndex((t) => t.id === activeId));
    const activeTab = visibleTabs[activeIndex];

    const handleChange = (_: React.SyntheticEvent, newIndex: number) => {
        const params = new URLSearchParams(searchParams);
        params.set(queryKey, visibleTabs[newIndex].id);
        navigate(`?${params.toString()}`, { replace: true });
    };

    return (
        <Box>
            <Tabs
                value={activeIndex}
                onChange={handleChange}
                variant="scrollable"
                scrollButtons="auto"
                sx={{ mb: 2 }}
            >
                {visibleTabs.map((t) => (
                    <Tab key={t.id} label={t.label} />
                ))}
            </Tabs>
            {loading ? (
                <Box sx={{ display: "flex", justifyContent: "center", py: 4 }}>
                    <CircularProgress />
                </Box>
            ) : (
                activeTab?.content
            )}
        </Box>
    );
}
```

**Critical detail:** `navigate(`?${params.toString()}`, { replace: true })` preserves any other query params. The `replace: true` prevents history bloat. Tab index is always derived from the filtered `visibleTabs` array — `hidden` tabs are removed before index calculation.

### Pattern 2: HomeView Loading State Fix (FEND-05)

**What:** Add `loading` boolean state initialized to `true`; only render the empty-state message after the fetch settles.

**Current broken code** (in `index.tsx`):
```tsx
useEffect(() => {
    GetPodsForPlayer(user.player_id).then((pods) => {
        if (pods.length > 0) { navigate(...) }
    });
}, [user]);
return <Typography>No pods yet...</Typography>; // always renders immediately
```

**Fixed pattern:**
```tsx
const [loading, setLoading] = useState(true);
const [pods, setPods] = useState<Pod[]>([]);

useEffect(() => {
    if (!user) return;
    GetPodsForPlayer(user.player_id).then((result) => {
        if (result.length > 0) {
            navigate(`/pod/${result[0].id}`, { replace: true });
        } else {
            setPods(result);
            setLoading(false);
        }
    }).catch(() => setLoading(false));
}, [user]);

if (loading) return <Box sx={{ display: "flex", justifyContent: "center", pt: 4 }}><CircularProgress /></Box>;
if (pods.length === 0) return <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", pt: 4 }}>
    <Typography>No pods yet. Create your first pod or ask a manager for an invite link.</Typography>
</Box>;
return null;
```

### Pattern 3: TooltipIcon and TooltipIconButton

**What:** Two small wrapper components for the two distinct tooltip+icon use cases.

```tsx
// app/src/components/TooltipIcon.tsx
import { ReactElement } from "react";
import { Tooltip } from "@mui/material";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";

interface TooltipIconProps {
    title: string;
    icon?: ReactElement;
}

export function TooltipIcon({ title, icon }: TooltipIconProps): ReactElement {
    return (
        <Tooltip title={title} enterTouchDelay={0}>
            <span style={{ display: "inline-flex", alignItems: "center", cursor: "default" }}>
                {icon ?? <InfoOutlinedIcon fontSize="small" color="action" />}
            </span>
        </Tooltip>
    );
}

interface TooltipIconButtonProps {
    title: string;
    onClick: () => void;
    icon: ReactElement;
}

export function TooltipIconButton({ title, onClick, icon }: TooltipIconButtonProps): ReactElement {
    return (
        <Tooltip title={title}>
            <IconButton size="medium" onClick={onClick}>
                {icon}
            </IconButton>
        </Tooltip>
    );
}
```

**MUI Tooltip + `forwardRef` note:** MUI `Tooltip` requires its child to accept a `ref` (forwarded ref). Raw DOM elements like `<span>` work fine. Custom components must use `React.forwardRef` if they are the direct Tooltip child. Using a `<span>` wrapper (as shown above for `TooltipIcon`) is the standard MUI pattern when the child does not forward refs. `IconButton` from MUI already forwards refs, so `TooltipIconButton` does not need a wrapper.

### Pattern 4: SvgIconPlayingCards — Reuse from root.tsx

**What:** The `SvgIconPlayingCards` function is defined in `app/src/routes/root.tsx` and needs to be reused in `app/src/routes/login.tsx` and `app/src/routes/join.tsx` (per UI-SPEC issues L-04 and J-01).

**Options:**
1. Extract to `app/src/components/SvgIconPlayingCards.tsx` and import in all three files.
2. Duplicate the ~15-line SVG in login.tsx and join.tsx (acceptable for this size).

Option 1 is cleaner; it also keeps `root.tsx` clean. This is Claude's discretion (the decision is not locked, but option 1 is recommended).

### Anti-Patterns to Avoid

- **`useState(0)` for tab index:** Replaced by query-string-driven tab state. Do not use `useState` for active tab index in any migrated view.
- **Inline `HomeView` in `index.tsx`:** Confirmed moved to `app/src/routes/home/index.tsx`. `index.tsx` must be purely routing + providers after the restructure.
- **Fixed pixel heights on DataGrid without responsive breakpoints:** Multiple views use `height: 600` or `height: 355`. Use `height: { xs: 400, sm: 600 }` or `autoHeight` per UI-SPEC.
- **Raw browser elements (`<span>`, `<em>`, `<h1>`, `<strong>`) for content:** Use MUI Typography variants. Multiple instances flagged in UI-SPEC.
- **`window.location.reload()` alternatives:** Existing pattern in save handlers — acceptable for now; not in scope to change.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Tab persistence in URL | Custom URL sync logic | `useSearchParams` + `useNavigate` | Already handles encoding, multiple params, browser back/fwd |
| Tooltip behavior | Custom hover/touch handlers | MUI `Tooltip` with `enterTouchDelay` | Handles focus, keyboard, aria-label automatically |
| Loading indicator | Custom CSS spinner | MUI `CircularProgress` | Already themed to accent gold from Phase 2 |
| Tab scroll on mobile | Custom scroll handling | MUI `variant="scrollable" scrollButtons="auto"` | Built-in, already used in pod.tsx |

---

## FEND-04 Root Cause Investigation

**Finding (HIGH confidence — code read directly):** The blank-screen-on-refresh is NOT caused by the static server. `app/main.go`'s `spaHandler` already implements SPA fallback correctly:

```go
// If the file doesn't exist on disk, serve index.html (SPA fallback)
if _, err := os.Stat(path); os.IsNotExist(err) {
    http.ServeFile(w, r, filepath.Join(absStatic, h.indexPath))
    return
}
```

The static server correctly falls through to `index.html` for unknown paths. The blank screen on refresh has a different root cause.

**Most likely root cause:** `RequireAuth` component renders `<CircularProgress>` while auth loads (correct), but on page refresh the auth state takes a moment to resolve. If `user` is null momentarily and the component renders without redirection completing, a brief white flash or blank state can occur depending on timing.

**Alternate root cause:** A CSS issue where the root `Container` in `root.tsx` uses `bgcolor: "background.default"` but the page body briefly shows white before the theme applies. Less likely since `CssBaseline` is present.

**Action for planner:** Include an investigation task that reads `app/src/routes/RequireAuth.tsx` and the auth context to determine the exact blank-screen trigger before prescribing a fix. The fix is almost certainly a one-liner but the location must be confirmed. The CONTEXT.md direction (D-17) is correct directionally — verify the actual rendering sequence first.

---

## Common Pitfalls

### Pitfall 1: Import Path Depth After Restructure

**What goes wrong:** Tab files in `app/src/routes/pod/DecksTab.tsx` are 3 levels deep. Relative imports to `app/src/components/` must use `../../components/`, not `../components/`.
**Why it happens:** Forgetting the directory nesting level when splitting files.
**How to avoid:** After every new tab file, run `tsc --noEmit` immediately. Import errors surface as TS errors.
**Warning signs:** `Cannot find module '../components/stats'` — means missing one `../` level.

### Pitfall 2: `useSearchParams` + `useNavigate` Must Be Inside Router Context

**What goes wrong:** Components using `useSearchParams` or `useNavigate` must be rendered inside a `<RouterProvider>`. If `TabbedLayout` is ever tested in isolation without a router, it will throw.
**Why it happens:** React Router hooks require Router context.
**How to avoid:** No issue in production — all views are within the router. In tests, wrap with `<MemoryRouter>`.

### Pitfall 3: Hidden Tab Index Mismatch

**What goes wrong:** If `hidden` filtering is done incorrectly, the Settings tab index shifts. E.g., if a non-manager is on `podTab=settings`, the fallback logic must map to the first visible tab, not to index 0 of the raw array.
**Why it happens:** `hidden` tabs are removed from the rendered list. MUI Tabs `value` must match the index in the *rendered* list, not the original `tabs` prop.
**How to avoid:** Always compute `activeIndex` from `visibleTabs.findIndex(...)`, not `tabs.findIndex(...)`.
**Warning signs:** MUI console warning "The `value` provided to the Tabs component is out of range."

### Pitfall 4: `useLoaderData` and Context Availability in Tab Sub-Components

**What goes wrong:** Tab sub-components extracted to separate files cannot call `useLoaderData()` — that hook only works in the component rendered directly by the route.
**Why it happens:** React Router's loader data context is scoped to the route component.
**How to avoid:** Keep `useLoaderData()` in the parent `index.tsx`. Pass data as props to tab components, or use `useParams()` in tabs that need to fetch data independently. Existing pod/player/deck tabs already follow the prop-drilling pattern — preserve it.
**Warning signs:** `useLoaderData()` returning `undefined` in a tab component.

### Pitfall 5: `AsyncComponentHelper` Pattern Limitation

**What goes wrong:** `AsyncComponentHelper` in `common.ts` calls `useEffect` + `useState` inside the function body, not inside a component. TypeScript may warn that hooks are called outside a component if the calling file does not look like a component to the linter.
**Why it happens:** It's technically a custom hook used as a helper function. The current call sites (player.tsx, deck.tsx) call it at the component function's top level — this is valid as long as it's always called at the top level of a component.
**How to avoid:** Move `common.ts` to `components/common.ts` as-is (per D-02). Don't rename or refactor the function signature during the move. Test with `tsc --noEmit` after the move.

### Pitfall 6: MUI Tooltip Requires Ref-Forwarding Child

**What goes wrong:** If `TooltipIcon`'s child does not forward refs, MUI logs a warning and the tooltip may not appear.
**Why it happens:** MUI Tooltip uses `React.cloneElement` to inject event handlers and ref into its child.
**How to avoid:** Use `<span>` as the wrapper for `TooltipIcon`'s child (native DOM elements support refs). `<IconButton>` from MUI already forwards refs — `TooltipIconButton` is safe without a wrapper.

---

## Code Examples

### TabbedLayout Usage in PodView
```tsx
// app/src/routes/pod/index.tsx (after restructure)
<TabbedLayout
    queryKey="podTab"
    tabs={[
        { id: "decks",    label: "Decks",    content: <PodDecksTab initialData={decks} /> },
        { id: "players",  label: "Players",  content: <PodPlayersTab players={players} isManager={isManager} /> },
        { id: "games",    label: "Games",    content: <PodGamesTab initialData={games} /> },
        { id: "settings", label: "Settings", content: <PodSettingsTab pod={pod} />, hidden: !isManager },
    ]}
/>
```

### Query String Navigation (Preserving Other Params)
```tsx
const [searchParams] = useSearchParams();
const navigate = useNavigate();

const switchTab = (tabId: string) => {
    const params = new URLSearchParams(searchParams); // copy existing params
    params.set("podTab", tabId);
    navigate(`?${params.toString()}`, { replace: true });
};
```

### TooltipIcon on DeckSettingsTab Commander Section (DECK-02)
```tsx
// app/src/routes/deck/SettingsTab.tsx
import { TooltipIcon } from "../../components/TooltipIcon";

<Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
    <Typography variant="h6">Commanders</Typography>
    <TooltipIcon title="This is for changing an existing deck's commander. To add a new deck, use the Add Deck button instead." />
</Box>
```

### AppBar Title Hidden on Mobile (P-07)
```tsx
// app/src/routes/root.tsx — "EDH Tracker" Typography
<Typography
    variant="h6"
    noWrap
    sx={{
        mr: 2,
        display: { xs: "none", sm: "flex" },  // ← change from "flex" to this
        fontWeight: 700,
        letterSpacing: ".3rem",
    }}
>
```

### DataGrid Responsive Height (P-01)
```tsx
// All DataGrid wrappers in Pod view
<Box sx={{ height: { xs: 400, sm: 600 }, width: "100%" }}>
    <DataGrid ... />
</Box>
```

### GameView DataGrid autoHeight (G-06)
```tsx
// game/index.tsx — GameResultsGrid
<Box sx={{ width: "100%" }}>
    <DataGrid
        rows={game.results}
        columns={columns}
        autoHeight
        // remove the fixed height: 355 wrapper Box
    />
</Box>
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `useState(0)` for tab index | Query-string-persisted tab ID via `useSearchParams` | This phase | Tabs survive page refresh; shareable URLs |
| Flat monolithic route files | Per-view subdirectories with per-tab files | This phase | Co-location, easier maintenance |
| Inline `HomeView` in router file | Dedicated `routes/home/index.tsx` | This phase | `index.tsx` stays clean |
| Raw `<span>`/`<em>`/`<h1>` for content | MUI `Typography` variants | This phase (DSNG-04) | Consistent theme-controlled typography |

---

## Environment Availability

Step 2.6: SKIPPED — this phase is purely frontend TypeScript/TSX code and configuration. No external services, databases, or CLI tools beyond the existing `tsc` binary (confirmed present in `app/node_modules/.bin/tsc` per MEMORY.md).

---

## Validation Architecture

`nyquist_validation` is `true` in `.planning/config.json`.

### Test Framework

| Property | Value |
|----------|-------|
| Framework | TypeScript compiler (`tsc`) — no Jest/Vitest configured for this project |
| Config file | `app/tsconfig.json` |
| Quick run command | `cd /mnt/d/msharp/Documents/projects/edh-tracker/app && ./node_modules/.bin/tsc --noEmit` |
| Full suite command | Same — no separate full suite; REQUIREMENTS.md explicitly excludes frontend automated tests: "Frontend automated tests: High effort for existing code; post-launch investment" |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FEND-01 | Route files split; imports resolve correctly | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ tsconfig.json |
| FEND-02 | TabbedLayout renders; query string updates on tab switch | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ |
| FEND-03 | TooltipIcon/TooltipIconButton compile with correct props | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ |
| FEND-04 | Page refresh does not produce blank screen | manual | Navigate to `/pod/1/game/1`, refresh browser | — |
| FEND-05 | HomeView does not flash "No pods yet" before data loads | manual | Load `/` while authenticated | — |
| DSNG-04 | Per-view visual issues resolved per UI-SPEC | manual | Visual inspection at 375px viewport | — |

### Sampling Rate
- **Per task commit:** `cd /mnt/d/msharp/Documents/projects/edh-tracker/app && ./node_modules/.bin/tsc --noEmit`
- **Per wave merge:** Same — run after each wave of changes
- **Phase gate:** TypeScript clean + manual verification of FEND-04, FEND-05, and DSNG-04 before `/gsd:verify-work`

### Wave 0 Gaps
None — no new test files needed. The frontend-verify skill (`tsc --noEmit`) is the sole automated check. Manual testing covers the behavioral requirements.

---

## Open Questions

1. **FEND-04 Exact Root Cause**
   - What we know: `app/main.go`'s `spaHandler` already has correct SPA fallback — the static server is not the bug.
   - What's unclear: Is the blank screen caused by `RequireAuth` rendering before auth resolves, or by a CSS timing issue, or by something else?
   - Recommendation: The planner should include a small investigation task: read `app/src/routes/RequireAuth.tsx`, reproduce the blank screen, then fix. Do not skip the investigation step.

2. **`SvgIconPlayingCards` Location**
   - What we know: Currently defined inside `root.tsx` as a private function. Needed in `login.tsx` (L-04) and `join.tsx` (J-01) per UI-SPEC.
   - What's unclear: Locked decision says extract to `components/` — not explicitly locked in CONTEXT.md (Claude's discretion for import organization).
   - Recommendation: Extract to `app/src/components/SvgIconPlayingCards.tsx`. Import in `root.tsx`, `login.tsx`, `join.tsx`.

3. **`AsyncComponentHelper` — Hook vs Helper Function**
   - What we know: The function calls `useEffect` + `useState` at the call site. It currently works in `player.tsx` and `deck.tsx`.
   - What's unclear: Whether ESLint's rules-of-hooks will flag it after the move.
   - Recommendation: Move as-is to `components/common.ts`. If ESLint flags it, rename to `useAsyncComponent` (hook naming convention). CONTEXT.md says don't modernize it in Phase 3.

---

## Sources

### Primary (HIGH confidence)
- Direct code read: `app/src/index.tsx`, `app/src/routes/pod.tsx`, `app/src/routes/player.tsx`, `app/src/routes/deck.tsx`, `app/src/routes/game.tsx`, `app/src/routes/root.tsx`, `app/src/routes/login.tsx`, `app/src/routes/join.tsx`, `app/src/common.ts`, `app/src/stats.tsx`, `app/src/matches.tsx`, `app/main.go`
- `.planning/phases/03-frontend-structure/03-CONTEXT.md` — locked decisions
- `.planning/phases/03-frontend-structure/03-UI-SPEC.md` — approved per-view audit (approved status confirmed in frontmatter)
- `.claude/skills/react-router/SKILL.md` — React Router v6 patterns for this project
- `.claude/skills/frontend-verify/SKILL.md` — TypeScript verification command

### Secondary (MEDIUM confidence)
- `.planning/REQUIREMENTS.md` — requirement descriptions
- `.planning/phases/02-design-language/02-CONTEXT.md` — established design decisions inherited by Phase 3

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all packages already installed; versions verified in CLAUDE.md
- Architecture: HIGH — all source files read directly; restructure target structure fully specified in CONTEXT.md
- FEND-04 root cause: MEDIUM — static server confirmed NOT the bug; actual root cause not yet verified (investigation task needed)
- Pitfalls: HIGH — derived from direct code reading, not from training data assumptions
- UI-SPEC implementation: HIGH — approved spec is the authoritative source; all changes are precisely specified

**Research date:** 2026-03-24
**Valid until:** 2026-04-24 (stable stack; only invalidated if MUI/React Router versions change)
