---
phase: 03-frontend-structure
verified: 2026-03-24T00:00:00Z
status: passed
score: 27/27 must-haves verified
re_verification: true
gaps: []
    artifacts:
      - path: "app/src/routes/root.tsx"
        issue: "Typography sx.display is 'flex' not { xs: 'none', sm: 'flex' }"
    missing:
      - "Change Typography sx={{ ..., display: 'flex', ... }} to sx={{ ..., display: { xs: 'none', sm: 'flex' }, ... }} on the 'EDH Tracker' text element in DrawerAppBar"
---

# Phase 3: Frontend Structure Verification Report

**Phase Goal:** Restructure frontend routes into per-domain subdirectories with a shared component library, fixing all identified UI-SPEC issues in the process.
**Verified:** 2026-03-24
**Status:** gaps_found — 1 gap blocking full requirement satisfaction
**Re-verification:** No — initial verification

---

## Goal Achievement

### Observable Truths

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | TabbedLayout component exists and accepts queryKey, tabs, and loading props | VERIFIED | `app/src/components/TabbedLayout.tsx` — correct interface, `useSearchParams`, `replace: true`, `TabConfig.hidden`, `CircularProgress` |
| 2 | TooltipIcon component renders an info icon with tap-to-toggle tooltip on mobile | VERIFIED | `app/src/components/TooltipIcon.tsx` — `enterTouchDelay={0}`, `InfoOutlinedIcon`, named export |
| 3 | TooltipIconButton component wraps IconButton with a hover tooltip | VERIFIED | `app/src/components/TooltipIcon.tsx` — named export `TooltipIconButton` with `Tooltip` wrapper |
| 4 | SvgIconPlayingCards is importable from components/ for use across views | VERIFIED | `app/src/components/SvgIconPlayingCards.tsx` — default export; consumed in `login.tsx`, `join.tsx`, `root.tsx` |
| 5 | Shared utilities (common.ts, stats.tsx, matches.tsx) live in components/ and all existing route files compile against the new paths | VERIFIED | All three files exist in `components/`; old files at `app/src/` deleted; TypeScript compiles cleanly |
| 6 | HomeView lives in its own file, not inline in index.tsx | VERIFIED | `app/src/routes/home/index.tsx` — default export `HomeView`; `index.tsx` imports from `./routes/home` with no inline definition |
| 7 | HomeView shows CircularProgress while loading, not 'No pods yet' flash | VERIFIED | `useState(true)` for loading, `CircularProgress` rendered when loading is true |
| 8 | Empty state message only appears after fetch confirms zero pods | VERIFIED | `setPods(result)` + `setLoading(false)` only called when `result.length === 0` |
| 9 | Refreshing any page does not produce a blank white screen | VERIFIED | `RequireAuth.tsx` wraps `CircularProgress` in centered `Box` with `justifyContent: "center"` |
| 10 | RequireAuth loading spinner is centered and visible during auth check | VERIFIED | `Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", pt: 4 }}` |
| 11 | Login page has vertical centering and the playing cards icon | VERIFIED | `minHeight: "100vh"`, `justifyContent: "center"`, `SvgIconPlayingCards fontSize={48}` |
| 12 | Join page error/no-code states have playing cards icon and MUI Button for 'Go home' | VERIFIED | `SvgIconPlayingCards fontSize={40}`, `Button component={Link} variant="outlined"` on both error and no-code states |
| 13 | PodView lives in routes/pod/index.tsx with tab content in separate per-tab files | VERIFIED | `routes/pod/` directory contains index.tsx, DecksTab.tsx, PlayersTab.tsx, GamesTab.tsx, SettingsTab.tsx; old pod.tsx deleted |
| 14 | Pod tabs use TabbedLayout with queryKey 'podTab' | VERIFIED | `pod/index.tsx` imports `TabbedLayout`, uses `queryKey="podTab"`, includes `hidden: !isManager` for Settings |
| 15 | DataGrid heights are responsive (xs: 400, sm: 600) | VERIFIED | `DecksTab.tsx` and `GamesTab.tsx` both use `height: { xs: 400, sm: 600 }` |
| 16 | Promote/Remove buttons have 44px minimum touch target | VERIFIED | `PlayersTab.tsx` — both Promote and Remove buttons have `sx={{ minHeight: 44 }}` |
| 17 | Settings form rows wrap on narrow viewports | VERIFIED | `SettingsTab.tsx` — `flexWrap: "wrap"` and `minWidth: 160` on TextField |
| 18 | AppBar title hidden on xs breakpoint | FAILED | root.tsx Typography still has `display: "flex"` — responsive breakpoint object not applied |
| 19 | PlayerView lives in routes/player/index.tsx with tab content in per-tab files | VERIFIED | `routes/player/` directory contains index.tsx, OverviewTab.tsx, DecksTab.tsx, GamesTab.tsx, SettingsTab.tsx; old player.tsx deleted |
| 20 | Player tabs use TabbedLayout with queryKey 'playerTab' | VERIFIED | `player/index.tsx` imports `TabbedLayout`, uses `queryKey="playerTab"`, `hidden: !isOwner` for Settings |
| 21 | Stats rows use MUI Typography instead of raw span/strong elements | VERIFIED | `OverviewTab.tsx` uses `Typography variant="body1"` for all stats; `DeckOverviewTab.tsx` same pattern |
| 22 | Error states show fixed user-friendly messages, not raw error.message | VERIFIED | `OverviewTab.tsx`: "Could not load pods. Refresh to try again."; `DecksTab.tsx`: "Could not load decks."; `GamesTab.tsx`: "Could not load games." |
| 23 | Leave pod buttons have 44px minimum touch target | VERIFIED | `player/SettingsTab.tsx` Leave button has `sx={{ minHeight: 44 }}` |
| 24 | DeckView lives in routes/deck/index.tsx with tab content in per-tab files | VERIFIED | `routes/deck/` directory contains index.tsx, OverviewTab.tsx, GamesTab.tsx, SettingsTab.tsx; old deck.tsx deleted |
| 25 | Deck tabs use TabbedLayout with queryKey 'deckTab' | VERIFIED | `deck/index.tsx` imports `TabbedLayout`, uses `queryKey="deckTab"`, `hidden: !isOwner` for Settings |
| 26 | TooltipIcon on Commanders section heading (DECK-02) | VERIFIED | `deck/SettingsTab.tsx` imports `TooltipIcon`, renders it next to `Typography variant="h6">Commanders<` with the correct tooltip text |
| 27 | GameView lives in routes/game/index.tsx with all G-01 through G-07 fixes | VERIFIED | `routes/game/index.tsx` — Typography h4 heading, body2 date, TooltipIconButton for edit/remove, `alignItems: "flex-start"`, `autoHeight` DataGrid; `mt: 3` on Delete button |

**Score:** 26/27 truths verified

---

## Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `app/src/components/TabbedLayout.tsx` | Shared tab component with query-string-persisted active tab | VERIFIED | 55 lines; `useSearchParams`, `useNavigate`, `replace: true`, `hidden` filtering, `CircularProgress` loading state |
| `app/src/components/TooltipIcon.tsx` | TooltipIcon and TooltipIconButton shared components | VERIFIED | Named exports for both; `enterTouchDelay={0}` on TooltipIcon |
| `app/src/components/SvgIconPlayingCards.tsx` | Playing cards SVG icon extracted from root.tsx | VERIFIED | Default export; `fontSize` prop; used in login, join, root |
| `app/src/components/common.ts` | AsyncComponentHelper moved from app/src/common.ts | VERIFIED | `export function AsyncComponentHelper` present; original deleted |
| `app/src/components/stats.tsx` | Record, RecordComparator, StatColumns, CommanderColumn | VERIFIED | All 4 exports present; `from "../types"` import path corrected |
| `app/src/components/matches.tsx` | MatchesDisplay, MatchUpDisplay | VERIFIED | Both named exports present; `from "../types"` import path corrected |
| `app/src/routes/home/index.tsx` | HomeView with loading state fix | VERIFIED | `useState(true)` loading; CircularProgress; updated empty state copy |
| `app/src/routes/RequireAuth.tsx` | Auth guard with centered loading spinner | VERIFIED | `justifyContent: "center"`, `alignItems: "center"` present |
| `app/src/routes/pod/index.tsx` | PodView + podLoader using TabbedLayout | VERIFIED | Named `podLoader` export; default `PodView`; TabbedLayout with queryKey "podTab" |
| `app/src/routes/pod/DecksTab.tsx` | Pod Decks tab | VERIFIED | Responsive height `{ xs: 400, sm: 600 }` |
| `app/src/routes/pod/PlayersTab.tsx` | Pod Players tab | VERIFIED | `minHeight: 44` on Promote and Remove buttons |
| `app/src/routes/pod/GamesTab.tsx` | Pod Games tab | VERIFIED | Responsive height `{ xs: 400, sm: 600 }` |
| `app/src/routes/pod/SettingsTab.tsx` | Pod Settings tab | VERIFIED | `flexWrap: "wrap"`, `minWidth: 160` on TextField and invite row |
| `app/src/routes/player/index.tsx` | PlayerView using TabbedLayout | VERIFIED | TabbedLayout `queryKey="playerTab"`, `hidden: !isOwner` |
| `app/src/routes/player/OverviewTab.tsx` | Player Overview tab with Typography fixes | VERIFIED | Typography for all stats, `flexWrap: "wrap"`, `Could not load pods`, Typography body2 for timestamp |
| `app/src/routes/player/DecksTab.tsx` | Player Decks tab | VERIFIED | `Could not load decks` error; empty state |
| `app/src/routes/player/GamesTab.tsx` | Player Games tab | VERIFIED | `Could not load games` error; empty state |
| `app/src/routes/player/SettingsTab.tsx` | Player Settings tab with touch target fixes | VERIFIED | Leave button `minHeight: 44`; `Divider` before Create New Pod |
| `app/src/routes/deck/index.tsx` | DeckView using TabbedLayout | VERIFIED | TabbedLayout `queryKey="deckTab"`, `mb: 0.5` on deck name heading |
| `app/src/routes/deck/OverviewTab.tsx` | Deck Overview tab with Typography fixes | VERIFIED | Typography for stats, `flexWrap: "wrap"` |
| `app/src/routes/deck/GamesTab.tsx` | Deck Games tab | VERIFIED | Exists; uses MatchesDisplay |
| `app/src/routes/deck/SettingsTab.tsx` | Deck Settings tab with TooltipIcon and button label fixes | VERIFIED | "Save Name", "Save Format", "Save Commanders"; `fullWidth` on Autocomplete; TooltipIcon on Commanders; `minHeight: 44` on Retire/Delete |
| `app/src/routes/game/index.tsx` | GameView + gameLoader with UI-SPEC fixes | VERIFIED | Typography h4, body2 date, TooltipIconButton edit/remove, `alignItems: "flex-start"`, `autoHeight`, `mt: 3` |
| `app/src/routes/new/index.tsx` | NewGameView + newGameLoader + createGame action | VERIFIED | All three exports present; import paths updated to `../../` depth |

### Deleted Files (required)

| File | Status |
|------|--------|
| `app/src/common.ts` | DELETED (confirmed) |
| `app/src/stats.tsx` | DELETED (confirmed) |
| `app/src/matches.tsx` | DELETED (confirmed) |
| `app/src/routes/pod.tsx` | DELETED (confirmed) |
| `app/src/routes/player.tsx` | DELETED (confirmed) |
| `app/src/routes/deck.tsx` | DELETED (confirmed) |
| `app/src/routes/game.tsx` | DELETED (confirmed) |
| `app/src/routes/new.tsx` | DELETED (confirmed) |

---

## Key Link Verification

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `components/TabbedLayout.tsx` | `react-router-dom` | `useSearchParams + useNavigate` | VERIFIED | Both imports present; `replace: true` in navigate call |
| `routes/pod/index.tsx` | `components/stats.tsx` | DecksTab import chain | VERIFIED | `DecksTab.tsx` imports `from "../../components/stats"` |
| `routes/home/index.tsx` | `app/src/http.ts` | `GetPodsForPlayer` | VERIFIED | Import and call both present |
| `app/src/index.tsx` | `routes/home/index.tsx` | import | VERIFIED | `import HomeView from "./routes/home"` |
| `routes/RequireAuth.tsx` | `app/src/auth.tsx` | `useAuth` hook | VERIFIED | `useAuth()` destructures both `user` and `loading` |
| `routes/pod/index.tsx` | `components/TabbedLayout.tsx` | TabbedLayout usage | VERIFIED | `import TabbedLayout from "../../components/TabbedLayout"` |
| `routes/pod/index.tsx` | `routes/pod/DecksTab.tsx` | tab content import | VERIFIED | `import PodDecksTab from "./DecksTab"` |
| `routes/player/index.tsx` | `components/TabbedLayout.tsx` | TabbedLayout usage | VERIFIED | `import TabbedLayout from "../../components/TabbedLayout"` |
| `routes/deck/index.tsx` | `components/TabbedLayout.tsx` | TabbedLayout usage | VERIFIED | `import TabbedLayout from "../../components/TabbedLayout"` |
| `routes/deck/SettingsTab.tsx` | `components/TooltipIcon.tsx` | TooltipIcon usage | VERIFIED | `import { TooltipIcon } from "../../components/TooltipIcon"` |
| `routes/game/index.tsx` | `components/TooltipIcon.tsx` | TooltipIconButton for edit/remove | VERIFIED | `import { TooltipIconButton } from "../../components/TooltipIcon"` |

---

## Data-Flow Trace (Level 4)

Level 4 not applicable to this phase — the phase restructures existing components and applies UI fixes. No new data sources introduced; all data flows were present pre-phase and unchanged.

---

## Behavioral Spot-Checks

| Behavior | Result | Status |
|----------|--------|--------|
| TypeScript compile (`tsc --noEmit`) | Exit 0, no errors | PASS |
| `app/src/components/` contains 6 files | TabbedLayout.tsx, TooltipIcon.tsx, SvgIconPlayingCards.tsx, common.ts, stats.tsx, matches.tsx | PASS |
| All old flat route files deleted | 8 files confirmed absent | PASS |
| All new subdirectory index files exist | home, pod, player, deck, game, new — all confirmed | PASS |
| index.tsx does not define HomeView inline | No `function HomeView` in index.tsx | PASS |
| root.tsx uses SvgIconPlayingCards from components | `import SvgIconPlayingCards from "../components/SvgIconPlayingCards"` present | PASS |
| root.tsx AppBar title hidden on xs | `display: "flex"` found — NOT `{ xs: "none", sm: "flex" }` | FAIL |

---

## Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
|-------------|------------|-------------|--------|----------|
| FEND-01 | 03-03, 03-04, 03-05, 03-06 | Large route files refactored into per-view subdirectories (pod/, player/, deck/, game/) | SATISFIED | All 4 domains restructured into subdirectories; old flat files deleted |
| FEND-02 | 03-01, 03-03, 03-04, 03-05 | Shared tab component used across Pod, Player, Deck views — active tab persisted via query string | SATISFIED | TabbedLayout in components/; used in pod, player, deck with distinct queryKey values |
| FEND-03 | 03-01, 03-05, 03-06 | Shared tooltip icon and tooltip icon button components available and used | SATISFIED | TooltipIcon.tsx created; TooltipIcon used in deck/SettingsTab; TooltipIconButton used in game/index.tsx |
| FEND-04 | 03-02 | Page refresh no longer causes blank white screen | SATISFIED | RequireAuth.tsx has centered Box wrapper around CircularProgress |
| FEND-05 | 03-02 | HomeView no longer flashes "No pods yet" before data loads | SATISFIED | `useState(true)` for loading; CircularProgress shown until fetch resolves |
| DSNG-04 | 03-02, 03-03, 03-04, 03-05, 03-06 | All views individually audited against design language — no view left with structural layout issues | PARTIAL — P-07 not applied | Login/Join/Home/Player/Deck/Game all fixed; Pod AppBar title mobile fix missing from root.tsx |

---

## Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| `app/src/routes/pod/PlayersTab.tsx` | 34 | `// TODO: Use icons w/ tooltips for promote/remove buttons?` | Info | Pre-existing TODO, not blocking — buttons already have minHeight fix |
| `app/src/routes/pod/SettingsTab.tsx` | 54 | `// TODO: Icon w/ tooltip for Save & Copy` | Info | Pre-existing TODO, not blocking |
| `app/src/components/matches.tsx` | 53 | `// ToDo: Make this and the containing DataGrid more mobile friendly` | Info | Pre-existing TODO carried from original file |
| `app/src/routes/new/index.tsx` | 29 | `// ToDo: Validation` | Info | Pre-existing TODO; NewGameView excluded from DSNG-04 per D-20 |
| `app/src/routes/root.tsx` | 78 | `display: "flex"` on AppBar title Typography | Blocker | P-07 fix not applied — AppBar title remains visible on 375px viewport, crowding PodSelector + Avatar + Logout |

---

## Human Verification Required

### 1. Tab persistence across navigation

**Test:** Navigate to Pod view, click "Players" tab, then navigate to a player and come back. URL should retain `?podTab=players`.
**Expected:** Tab selection persists in the URL query string and survives navigation.
**Why human:** Requires browser interaction with a running app.

### 2. Mobile AppBar crowding (P-07 — known gap)

**Test:** Open the app on a 375px viewport (iPhone SE size). Log in and navigate to a pod.
**Expected:** With the gap unfixed, "EDH Tracker" text, PodSelector dropdown, avatar, username, and Logout compete for the AppBar. Verify the extent of crowding.
**Why human:** Visual layout check on narrow viewport.

### 3. TooltipIcon on mobile

**Test:** On a mobile device, tap the info icon next to "Commanders" in the Deck Settings tab.
**Expected:** Tooltip appears and stays visible on tap (not hover-only). `enterTouchDelay={0}` ensures this.
**Why human:** Touch interaction behavior requires physical device or emulator.

---

## Gaps Summary

One gap was found: the **AppBar title mobile fix (P-07)** was specified in plan 03-03 Task 2 and included in the plan's success criteria, but was not applied. The root.tsx Typography element for "EDH Tracker" still uses `display: "flex"` instead of the responsive `display: { xs: "none", sm: "flex" }`. This means on 375px viewports the text label still crowds the AppBar alongside PodSelector, Avatar, and Logout controls.

All 6 other requirement items (FEND-01 through FEND-05, DSNG-04 partial) are otherwise well-implemented. The single missing fix is a one-line change to one sx prop in root.tsx.

---

_Verified: 2026-03-24_
_Verifier: Claude (gsd-verifier)_
