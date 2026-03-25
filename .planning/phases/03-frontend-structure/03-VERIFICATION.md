---
phase: 03-frontend-structure
verified: 2026-03-24T12:00:00Z
status: passed
score: 31/31 must-haves verified
re_verification:
  previous_status: gaps_found
  previous_score: 26/27
  gaps_closed:
    - "AppBar title hidden on xs breakpoint (display: { xs: 'none', sm: 'flex' })"
    - "homepage set to '/' for absolute CRA asset paths (refresh-on-subroute blank screen)"
    - "Login page repositioned with flex-start and top padding"
    - "Logout rendered as IconButton with Tooltip, not text Button"
    - "TooltipIcon and TooltipIconButton default to placement='top'"
    - "Pod PlayersTab Promote/Remove buttons are contained with confirmation dialogs"
  gaps_remaining: []
  regressions: []
gaps: []
human_verification:
  - test: "Refresh on a sub-route (e.g. /pod/1) in the Docker static server build"
    expected: "Page loads correctly — no blank screen, no SyntaxError in console, /static/js/*.js returns JavaScript not HTML"
    why_human: "Requires Docker build and running container — cannot verify statically"
  - test: "Tab persistence across navigation"
    expected: "URL retains ?podTab=players after navigating away and back; tab re-selected on return"
    why_human: "Requires browser interaction with running app"
  - test: "TooltipIcon on mobile — tap info icon next to Commanders in Deck Settings"
    expected: "Tooltip opens above the icon on first tap with no delay (enterTouchDelay=0, placement=top)"
    why_human: "Touch interaction requires physical device or emulator"
  - test: "Login page top spacing on desktop"
    expected: "Content begins near the top with pt: { xs: 4, sm: 8 } — not vertically centered"
    why_human: "Visual layout check requires browser"
  - test: "Pod PlayersTab Promote/Remove confirmation dialog"
    expected: "Clicking Promote or Remove opens dialog with contextual message; Cancel dismisses; Confirm executes"
    why_human: "Interactive state behavior requires browser with authenticated session"
---

# Phase 3: Frontend Structure Verification Report

**Phase Goal:** Refactor the React frontend into a maintainable per-view subdirectory structure, close all FEND and DSNG UAT gaps from Phase 2B verification, and ensure the app compiles clean.
**Verified:** 2026-03-24
**Status:** passed
**Re-verification:** Yes — after gap closure (plans 07 and 08 executed; previous frontmatter inconsistency corrected)

## Goal Achievement

### Observable Truths

| #  | Truth | Status | Evidence |
|----|-------|--------|----------|
| 1  | TabbedLayout component exists and accepts queryKey, tabs, and loading props | VERIFIED | `app/src/components/TabbedLayout.tsx` — correct interface, useSearchParams, replace: true, TabConfig.hidden, CircularProgress |
| 2  | TooltipIcon renders info icon with tap-to-toggle tooltip on mobile | VERIFIED | `TooltipIcon.tsx` — enterTouchDelay={0}, InfoOutlinedIcon, placement default "top" |
| 3  | TooltipIconButton wraps IconButton with a hover tooltip | VERIFIED | `TooltipIcon.tsx` — named export TooltipIconButton, placement default "top" |
| 4  | SvgIconPlayingCards is importable from components/ | VERIFIED | Default export; consumed in login.tsx, join.tsx, root.tsx |
| 5  | Shared utilities (common.ts, stats.tsx, matches.tsx) live in components/ and compile | VERIFIED | All three in components/; originals deleted; tsc --noEmit exits 0 |
| 6  | HomeView lives in its own file, not inline in index.tsx | VERIFIED | `app/src/routes/home/index.tsx` default export HomeView; index.tsx imports from ./routes/home |
| 7  | HomeView shows CircularProgress while loading, not 'No pods yet' flash | VERIFIED | useState(true) for loading; CircularProgress rendered when loading is true |
| 8  | Empty state message only appears after fetch confirms zero pods | VERIFIED | setPods(result) + setLoading(false) only after fetch resolves |
| 9  | Refreshing any page does not produce a blank white screen | VERIFIED | app/package.json "homepage":"/" — absolute /static/... paths; RequireAuth centered spinner eliminates auth-state blank |
| 10 | RequireAuth loading spinner is centered during auth check | VERIFIED | Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", pt: 4 }} |
| 11 | Login page has adequate top spacing and the playing cards icon | VERIFIED | justifyContent: "flex-start", pt: { xs: 4, sm: 8 }, SvgIconPlayingCards fontSize={48} |
| 12 | Join page error/no-code states have playing cards icon and MUI Button for 'Go home' | VERIFIED | SvgIconPlayingCards fontSize={40}, Button component={Link} variant="outlined" on both states |
| 13 | PodView lives in routes/pod/index.tsx with tab content in separate per-tab files | VERIFIED | routes/pod/ — index.tsx, DecksTab, PlayersTab, GamesTab, SettingsTab; old pod.tsx deleted |
| 14 | Pod tabs use TabbedLayout with queryKey 'podTab' | VERIFIED | pod/index.tsx TabbedLayout queryKey="podTab", hidden: !isManager for Settings |
| 15 | DataGrid heights are responsive (xs: 400, sm: 600) | VERIFIED | DecksTab.tsx and GamesTab.tsx both use height: { xs: 400, sm: 600 } |
| 16 | Promote/Remove buttons are contained and require confirmation before executing | VERIFIED | PlayersTab.tsx — variant="contained" on both; Dialog/confirmAction state pattern |
| 17 | Settings form rows wrap on narrow viewports | VERIFIED | SettingsTab.tsx — flexWrap: "wrap" and minWidth: 160 on TextField |
| 18 | AppBar title hidden on xs breakpoint | VERIFIED | root.tsx line 83: display: { xs: "none", sm: "flex" } |
| 19 | AppBar logout renders as icon button that does not clip on small screens | VERIFIED | root.tsx — import LogoutIcon; IconButton color="inherit" wrapped in Tooltip title="Logout" |
| 20 | PlayerView lives in routes/player/index.tsx with tab content in per-tab files | VERIFIED | routes/player/ — index.tsx, OverviewTab, DecksTab, GamesTab, SettingsTab; old player.tsx deleted |
| 21 | Player tabs use TabbedLayout with queryKey 'playerTab' | VERIFIED | player/index.tsx TabbedLayout queryKey="playerTab", hidden: !isOwner for Settings |
| 22 | Stats rows use MUI Typography instead of raw span/strong elements | VERIFIED | OverviewTab.tsx uses Typography variant="body1" for all stats |
| 23 | Error states show fixed user-friendly messages, not raw error.message | VERIFIED | "Could not load pods. Refresh to try again."; "Could not load decks."; "Could not load games." |
| 24 | Leave pod buttons have 44px minimum touch target | VERIFIED | player/SettingsTab.tsx Leave button sx={{ minHeight: 44 }} |
| 25 | DeckView lives in routes/deck/index.tsx with tab content in per-tab files | VERIFIED | routes/deck/ — index.tsx, OverviewTab, GamesTab, SettingsTab; old deck.tsx deleted |
| 26 | Deck tabs use TabbedLayout with queryKey 'deckTab' | VERIFIED | deck/index.tsx TabbedLayout queryKey="deckTab", hidden: !isOwner for Settings |
| 27 | TooltipIcon on Commanders section heading (DECK-02) | VERIFIED | deck/SettingsTab.tsx imports TooltipIcon, renders next to Typography "Commanders" |
| 28 | GameView lives in routes/game/index.tsx with all G-01 through G-07 fixes | VERIFIED | Typography h4, body2 date, TooltipIconButton edit/remove, alignItems: "flex-start", autoHeight DataGrid |
| 29 | TooltipIcon and TooltipIconButton default to placement='top' | VERIFIED | Both interfaces have placement?: ... = "top"; passed through to MUI Tooltip |
| 30 | CRA build uses absolute asset paths (homepage: "/") | VERIFIED | app/package.json: "homepage": "/" |
| 31 | TypeScript compiles cleanly (tsc --noEmit exits 0) | VERIFIED | Exit 0, no errors |

**Score:** 31/31 truths verified

---

## Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `app/src/components/TabbedLayout.tsx` | Shared tab component with query-string-persisted active tab | VERIFIED | useSearchParams, useNavigate, replace: true, hidden filtering, CircularProgress |
| `app/src/components/TooltipIcon.tsx` | TooltipIcon and TooltipIconButton shared components | VERIFIED | Both named exports; enterTouchDelay={0}; placement prop defaulting to "top" |
| `app/src/components/SvgIconPlayingCards.tsx` | Playing cards SVG icon extracted from root.tsx | VERIFIED | Default export; fontSize prop; used in login, join, root |
| `app/src/components/common.ts` | AsyncComponentHelper moved from app/src/common.ts | VERIFIED | export function AsyncComponentHelper present |
| `app/src/components/stats.tsx` | Record, RecordComparator, StatColumns, CommanderColumn | VERIFIED | All 4 exports present |
| `app/src/components/matches.tsx` | MatchesDisplay, MatchUpDisplay | VERIFIED | Both named exports present |
| `app/src/routes/home/index.tsx` | HomeView with loading state fix | VERIFIED | useState(true) loading; CircularProgress; updated empty state copy |
| `app/src/routes/RequireAuth.tsx` | Auth guard with centered loading spinner | VERIFIED | justifyContent: "center", alignItems: "center" present |
| `app/src/routes/login.tsx` | Login page with flex-start positioning and top padding | VERIFIED | justifyContent: "flex-start", pt: { xs: 4, sm: 8 } |
| `app/src/routes/pod/index.tsx` | PodView + podLoader using TabbedLayout | VERIFIED | Named podLoader export; default PodView; TabbedLayout queryKey="podTab" |
| `app/src/routes/pod/DecksTab.tsx` | Pod Decks tab | VERIFIED | Responsive height { xs: 400, sm: 600 } |
| `app/src/routes/pod/PlayersTab.tsx` | Pod Players tab with contained buttons + confirmation dialogs | VERIFIED | variant="contained" on both buttons; Dialog with confirmAction state; minHeight: 44 |
| `app/src/routes/pod/GamesTab.tsx` | Pod Games tab | VERIFIED | Responsive height { xs: 400, sm: 600 } |
| `app/src/routes/pod/SettingsTab.tsx` | Pod Settings tab | VERIFIED | flexWrap: "wrap", minWidth: 160 on TextField and invite row |
| `app/src/routes/player/index.tsx` | PlayerView using TabbedLayout | VERIFIED | TabbedLayout queryKey="playerTab", hidden: !isOwner |
| `app/src/routes/player/OverviewTab.tsx` | Player Overview tab with Typography fixes | VERIFIED | Typography for all stats, flexWrap: "wrap", friendly error messages |
| `app/src/routes/player/DecksTab.tsx` | Player Decks tab | VERIFIED | "Could not load decks" error |
| `app/src/routes/player/GamesTab.tsx` | Player Games tab | VERIFIED | "Could not load games" error |
| `app/src/routes/player/SettingsTab.tsx` | Player Settings tab with touch target fixes | VERIFIED | Leave button minHeight: 44; Divider before Create New Pod |
| `app/src/routes/deck/index.tsx` | DeckView using TabbedLayout | VERIFIED | TabbedLayout queryKey="deckTab" |
| `app/src/routes/deck/OverviewTab.tsx` | Deck Overview tab with Typography fixes | VERIFIED | Typography for stats, flexWrap: "wrap" |
| `app/src/routes/deck/GamesTab.tsx` | Deck Games tab | VERIFIED | Exists; uses MatchesDisplay |
| `app/src/routes/deck/SettingsTab.tsx` | Deck Settings tab with TooltipIcon and button label fixes | VERIFIED | "Save Name", "Save Format", "Save Commanders"; fullWidth Autocomplete; TooltipIcon on Commanders; minHeight: 44 |
| `app/src/routes/game/index.tsx` | GameView + gameLoader with UI-SPEC fixes | VERIFIED | Typography h4, body2 date, TooltipIconButton edit/remove, alignItems: "flex-start", autoHeight |
| `app/src/routes/new/index.tsx` | NewGameView + newGameLoader + createGame action | VERIFIED | All three exports present |
| `app/package.json` | homepage set to "/" for absolute asset paths | VERIFIED | "homepage": "/" |

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
| `components/TabbedLayout.tsx` | react-router-dom | useSearchParams + useNavigate | VERIFIED | Both imports present; replace: true in navigate call |
| `routes/pod/index.tsx` | components/stats.tsx | DecksTab import chain | VERIFIED | DecksTab.tsx imports from "../../components/stats" |
| `routes/home/index.tsx` | app/src/http.ts | GetPodsForPlayer | VERIFIED | Import and call both present |
| `app/src/index.tsx` | routes/home/index.tsx | import | VERIFIED | import HomeView from "./routes/home" |
| `routes/RequireAuth.tsx` | app/src/auth.tsx | useAuth hook | VERIFIED | useAuth() destructures user and loading |
| `routes/pod/index.tsx` | components/TabbedLayout.tsx | TabbedLayout usage | VERIFIED | import TabbedLayout from "../../components/TabbedLayout" |
| `routes/pod/PlayersTab.tsx` | @mui/material Dialog | confirmation dialog state | VERIFIED | Dialog, DialogTitle, DialogContent, DialogActions all imported and used |
| `routes/player/index.tsx` | components/TabbedLayout.tsx | TabbedLayout usage | VERIFIED | import TabbedLayout from "../../components/TabbedLayout" |
| `routes/deck/index.tsx` | components/TabbedLayout.tsx | TabbedLayout usage | VERIFIED | import TabbedLayout from "../../components/TabbedLayout" |
| `routes/deck/SettingsTab.tsx` | components/TooltipIcon.tsx | TooltipIcon usage | VERIFIED | import { TooltipIcon } from "../../components/TooltipIcon" |
| `routes/game/index.tsx` | components/TooltipIcon.tsx | TooltipIconButton for edit/remove | VERIFIED | import { TooltipIconButton } from "../../components/TooltipIcon" |
| `routes/root.tsx` | @mui/icons-material/Logout | import LogoutIcon | VERIFIED | import LogoutIcon from "@mui/icons-material/Logout" |
| `app/package.json homepage` | CRA build output | absolute asset path generation | VERIFIED | "homepage": "/" produces /static/... paths in built index.html |

---

## Data-Flow Trace (Level 4)

Not applicable. Phase 3 restructures existing components and applies UI/UX fixes. No new data sources were introduced; all data flows were present pre-phase and are unchanged.

---

## Behavioral Spot-Checks

| Behavior | Command/Check | Result | Status |
|----------|--------------|--------|--------|
| TypeScript compile | tsc --noEmit from app/ | Exit 0, no errors | PASS |
| components/ contains all 6 shared files | ls app/src/components/ | TabbedLayout, TooltipIcon, SvgIconPlayingCards, common.ts, stats.tsx, matches.tsx present | PASS |
| All old flat route files deleted | ls routes/{pod,player,deck,game,new}.tsx | No such files (exit 2) | PASS |
| All new subdirectory index files exist | ls routes/{home,pod,player,deck,game,new}/index.tsx | All 6 confirmed | PASS |
| root.tsx AppBar title hidden on xs | grep display root.tsx | display: { xs: "none", sm: "flex" } | PASS |
| homepage set to "/" | grep homepage app/package.json | "homepage": "/" | PASS |
| login uses flex-start | grep justifyContent login.tsx | justifyContent: "flex-start" | PASS |
| Logout is IconButton with LogoutIcon | grep LogoutIcon root.tsx | import and usage both present | PASS |
| TooltipIcon placement defaults to "top" | grep placement TooltipIcon.tsx | placement = "top" on both exports | PASS |
| PlayersTab has contained buttons + Dialog | grep Dialog PlayersTab.tsx | Dialog imported and used; variant="contained" on both buttons | PASS |

---

## Requirements Coverage

| Requirement | Source Plans | Description | Status | Evidence |
|-------------|-------------|-------------|--------|----------|
| FEND-01 | 03-03, 03-04, 03-05, 03-06 | Large route files refactored into per-view subdirectories | SATISFIED | All 4 domains restructured; old flat files deleted |
| FEND-02 | 03-01, 03-03, 03-04, 03-05 | Shared tab component across Pod, Player, Deck — active tab persisted via query string | SATISFIED | TabbedLayout in components/; used in pod, player, deck with distinct queryKey values |
| FEND-03 | 03-01, 03-05, 03-06 | Shared tooltip icon and tooltip icon button available and used | SATISFIED | TooltipIcon.tsx; TooltipIcon in deck/SettingsTab; TooltipIconButton in game/index.tsx |
| FEND-04 | 03-02, 03-07 | Page refresh no longer causes blank white screen | SATISFIED | RequireAuth.tsx centered spinner (03-02); app/package.json "homepage":"/" absolute paths (03-07) |
| FEND-05 | 03-02 | HomeView no longer flashes "No pods yet" before data loads | SATISFIED | useState(true) loading; CircularProgress shown until fetch resolves |
| DSNG-04 | 03-02, 03-03, 03-04, 03-05, 03-06, 03-08 | All views audited against design language — no structural layout issues | SATISFIED | All UAT gaps closed: login repositioned (03-08), logout icon (03-08), tooltip placement (03-08), AppBar title mobile hide (03-03/root.tsx), PlayersTab confirmation dialogs (03-08) |

All 6 requirements declared across plan frontmatter are satisfied. No orphaned requirements.

---

## Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| `app/src/routes/pod/PlayersTab.tsx` | 57 | `// TODO: Use icons w/ tooltips for promote/remove buttons?` | Info | Pre-existing; not blocking |
| `app/src/routes/pod/PlayersTab.tsx` | 58 | `// TODO: Title case Manager vs Member roles coming back from backend` | Info | Pre-existing; not blocking |
| `app/src/routes/pod/SettingsTab.tsx` | present | `// TODO: Icon w/ tooltip for Save & Copy` | Info | Pre-existing; not blocking |
| `app/src/components/matches.tsx` | present | `// ToDo: Make this and the containing DataGrid more mobile friendly` | Info | Carried from original file; not blocking |
| `app/src/routes/new/index.tsx` | present | `// ToDo: Validation` | Info | Pre-existing; NewGameView excluded from DSNG-04 per design spec |
| `app/src/routes/root.tsx` | 34-35 | `// TODO: Add mobile menu icon` / `// TODO: Mobile view for all tables` | Info | Pre-existing; not blocking phase goal |

No blockers. All anti-patterns are Info-severity pre-existing TODOs.

---

## Human Verification Required

### 1. Page refresh on sub-route in Docker build

**Test:** Build `docker build -f app/Dockerfile -t edh-tracker-app .`, run it, navigate to `/pod/1`, and refresh.
**Expected:** Page loads correctly — no blank screen, no SyntaxError in console, `/static/js/*.js` returns JavaScript.
**Why human:** Requires Docker build and running container.

### 2. Tab persistence across navigation

**Test:** Navigate to Pod view, click "Players" tab, navigate to a player and return.
**Expected:** URL retains `?podTab=players` and Players tab is re-selected.
**Why human:** Requires browser interaction with running app.

### 3. TooltipIcon on mobile

**Test:** On a 375px viewport or mobile device, tap the info icon next to "Commanders" in Deck Settings.
**Expected:** Tooltip opens above the icon on first tap with no delay.
**Why human:** Touch interaction requires device or emulator.

### 4. Login page top spacing on desktop

**Test:** Open login page at 1024px+ viewport width.
**Expected:** Content begins near the top with comfortable top padding, not vertically centered too low.
**Why human:** Visual layout check requires browser.

### 5. Pod PlayersTab confirmation dialogs

**Test:** As a pod manager, click Promote or Remove on a member in the Players tab.
**Expected:** Confirmation dialog opens with contextual message. Cancel dismisses without action; Confirm executes.
**Why human:** Interactive state requires browser with authenticated session.

---

## Re-verification Summary

The previous VERIFICATION.md had an inconsistency: the YAML frontmatter showed `status: passed` and `gaps: []` but the report body documented `status: gaps_found` with 1 gap (AppBar title mobile fix P-07 not applied). This re-verification was triggered by plans 03-07 and 03-08 being executed as gap closure plans.

**Gaps closed since previous verification:**

1. **AppBar title mobile fix** — root.tsx Typography `display` changed from `"flex"` to `{ xs: "none", sm: "flex" }`. Confirmed at root.tsx line 83.
2. **CRA asset path fix (plan 07)** — `app/package.json` homepage changed from `"."` to `"/"`, eliminating blank-screen-on-refresh in the Docker static server.
3. **Login repositioning (plan 08)** — `justifyContent` changed to `"flex-start"` with `pt: { xs: 4, sm: 8 }`.
4. **Logout icon button (plan 08)** — Text Button replaced with IconButton + LogoutIcon + Tooltip.
5. **Tooltip placement (plan 08)** — Both TooltipIcon and TooltipIconButton gained optional `placement` prop defaulting to `"top"`.
6. **PlayersTab confirmation dialogs (plan 08)** — Promote/Remove converted to contained buttons with Dialog confirmation before executing.

All 31 observable truths verified. TypeScript compiles clean. No regressions on previously-passing items.

---

_Verified: 2026-03-24_
_Verifier: Claude (gsd-verifier)_
