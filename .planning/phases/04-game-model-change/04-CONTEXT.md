# Phase 4: Game Model Change - Context

**Gathered:** 2026-03-24
**Status:** Ready for planning

<domain>
## Phase Boundary

Remove the player field from game entry (decks are the unit, player is implicit via deck ownership). Redesign the new game form to be visually clean on mobile. Fix the record stat display to work correctly for games with any number of players.

**In scope:**
- GAME-01: Remove player field from game creation — both NewGameView and AddResultModal in GameView
- GAME-02: Deck picker in game form shows `DeckName (PlayerName)` format
- GAME-03: New game form is a stacked card layout, clean and usable at 375px
- GAME-04: Record stat component dynamically shows places with data (not hardcoded to 4)

**Backend note:** The backend is already functionally deck-only. `gameResult.Model` has no `PlayerID` field; `gameResult.InputEntity` (game creation) has no `player_id`; `AddResult` business function accepts `playerID` but silently drops it. Phase 4 includes API contract cleanup to match the existing behaviour.

**Out of scope:** Pod/deck functional gaps (Phase 5), auth interceptor (Phase 6), match display mobile polish (`MatchUpDisplay` TODO — deferred until DSNG-02 sweep).

</domain>

<decisions>
## Implementation Decisions

### Deck Picker Label (GAME-02)

- **D-01:** Deck picker label format: `DeckName (PlayerName)` — e.g. `"My Rakdos Deck (Mike)"`. Consistent across all cases (with or without commander set, partner commanders).
- **D-02:** Applies to the Autocomplete in NewGameView. Does NOT change the deck label in other views (GameView result display still uses `deck_name`).

### New Game Form Layout (GAME-03)

- **D-03:** Layout: stacked card per player. Each player entry is a Box/Card:
  ```
  ┌────────────────────────────┐
  │  [ Deck Autocomplete    ❌ ] │
  │  [ Place ]  [ Kills ]       │
  └────────────────────────────┘
  ```
  The ❌ removes the card. Deck picker spans full width. Place and Kills side-by-side below.

- **D-04:** Place field always starts blank (no auto-populate by card position). User sets every place explicitly.

- **D-05:** Format field is required and always visible at the top of the form (no collapse).

- **D-06:** Description field is hidden behind an expand link ("+ Add description"). Not shown by default.

- **D-07:** "Add Player" button renamed to "Add Deck" — reinforces that decks are the unit of entry.

- **D-08:** A Remove button (❌) appears on each card. Minimum cards enforced at 2 (no submitting with fewer than 2 entries) — Claude's discretion on exact validation behavior.

### Variable-Player Record Display (GAME-04)

- **D-09:** `Record` component becomes dynamic: iterate the keys of the `RecordDict` from `1` to `max(keys)`, show counts separated by ` / `. Places with no entries show `0`.
  - 4-player meta: `5 / 3 / 2 / 1`
  - 3-player meta: `5 / 3 / 2`
  - Mixed (up to place 6): `5 / 3 / 2 / 1 / 0 / 1`

- **D-10:** `RecordComparator` updated to compare dynamically: iterate from place 1 up to the max place in either record dict, return the first non-zero difference.

### Player Field Removal (GAME-01)

- **D-11:** NewGameView: remove player picker from `GameInput`. Remove `GetPlayersForPod` call from `newGameLoader`. Remove `players` from `NewGameData` type and `NewGameResult` type (`player_id` field removed).

- **D-12:** GameView AddResultModal: remove the Player Autocomplete. Player is implicit from deck selection.

- **D-13:** Backend cleanup: remove `player_id` from `addGameResultRequest` struct in `lib/routers/game.go`. Remove `playerID int` parameter from `AddResultFunc` type and `AddResult` constructor in the business layer.

### Claude's Discretion

- Initial number of cards shown when NewGameView loads (2 is the natural default for a 2-player minimum)
- Whether the Remove button is a small icon button (❌) in the card header or a full-width button below each card
- Exact MUI spacing, elevation, and border style for player entry cards
- Whether Place and Kills use `TextField` with `type="number"` or `Select` with numeric options
- Whether the `MatchUpDisplay` in `matches.tsx` is updated for mobile in this phase (TODO comment present — researcher should assess)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase Requirements
- `.planning/REQUIREMENTS.md` §Games — GAME-01 through GAME-04
- `.planning/ROADMAP.md` §Phase 4 — goal, success criteria, `UI hint: yes`

### Design System (required for any UI work)
- `.planning/phases/02-design-language/02-CONTEXT.md` — color palette, typography, spacing decisions
- `.planning/phases/02-design-language/02-UI-SPEC.md` — established MUI theme and component patterns

### Files to modify
- `app/src/routes/new/index.tsx` — NewGameView + newGameLoader + createGame action (primary redesign target)
- `app/src/routes/game/index.tsx` — GameView with AddResultModal (player picker removal)
- `app/src/components/stats.tsx` — `Record` component + `RecordComparator` (GAME-04 fix)
- `app/src/types.ts` — `NewGameData`, `NewGameResult` types (remove player fields)
- `lib/routers/game.go` — `addGameResultRequest` struct (remove player_id)
- `lib/business/game/functions.go` — `AddResult` constructor + `AddResultFunc` type (remove playerID param)
- `lib/business/game/types.go` — `AddResultFunc` type alias (remove playerID param)

### No external specs
- Requirements fully captured in decisions above and REQUIREMENTS.md §Games.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `TabbedLayout` (`app/src/components/TabbedLayout.tsx`): not directly needed for this phase but available
- `TooltipIconButton` (`app/src/components/TooltipIcon.tsx`): use for the ❌ remove button on each player card
- MUI `Autocomplete` already used in NewGameView for deck picker — keep, just change `getOptionLabel`
- MUI `Card` or `Box` with `sx={{ border: 1, borderRadius: 1, p: 2 }}` for player entry cards

### Established Patterns
- Form submission uses `useSubmit` + React Router `action` (not useEffect/fetch) — preserve in redesign
- Data fetching via React Router `loader` (`newGameLoader`) — stays, just removes the players fetch
- `window.location.reload()` used in GameView after mutations — existing pattern, keep as-is
- `LoaderFunctionArgs` imported from `@remix-run/router/utils` — note import source

### Integration Points
- `app/src/http.ts`: `GetAllDecksForPod` (used in newGameLoader) stays; `GetPlayersForPod` call removed from newGameLoader only (still used in GameView for the players list in AddResultModal header display)
- `app/src/types.ts`: `NewGameData.players` and `NewGameResult.player_id` removed; `NewGameResultWithGame.player_id` also likely dead after cleanup
- The `RecordDict` type (`{ [key: number]: number }`) already supports dynamic keys — no type change needed, only display logic changes

### Backend (already deck-only)
- `gameResult.InputEntity` has no `player_id` — create game endpoint is already correct
- `gameResult.Model` has no `PlayerID` field — DB doesn't store player per result (player derived from deck FK)
- `AddResult` business function: `playerID int` parameter is accepted but silently dropped when building the Model — cleanup target

</code_context>

<specifics>
## Specific Ideas

- Card layout is the primary mobile fix — the existing 4-column row at 375px was the main problem; the card approach mirrors how MUI form patterns naturally scale
- "Add Deck" (not "Add Player") — chosen to reinforce that decks are the entry unit
- Format is required — it always was conceptually required, and making it explicit in the UI prevents incomplete entries
- Description is genuinely optional ("Add a game description!") — hide it by default to reduce form weight

</specifics>

<deferred>
## Deferred Ideas

- `MatchUpDisplay` mobile polish — has a TODO comment; not in GAME-03 scope, deferred to DSNG-02 sweep (Phase 5 or later)
- Sorting/reordering player cards by drag — not in scope; table ordering is separate UX
- Pending todo: "Security review all API and frontend route authorization" — reviewed, belongs in Phase 6

</deferred>

---

*Phase: 04-game-model-change*
*Context gathered: 2026-03-24*
