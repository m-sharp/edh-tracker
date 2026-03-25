# Phase 4: Game Model Change - Research

**Researched:** 2026-03-24
**Domain:** React/TypeScript form redesign, Go API cleanup, MUI v5 component patterns
**Confidence:** HIGH

## Summary

Phase 4 is a well-bounded, largely mechanical change: remove player selection from game entry, redesign the new game form for mobile, fix the record display to be dynamic, and clean up the backend to match what it already does in practice. The backend is already functionally deck-only — the only backend work is removing dead parameters. All design decisions are locked in 04-CONTEXT.md and the approved 04-UI-SPEC.md.

The primary complexity is in `app/src/routes/new/index.tsx`, which is being fully replaced. The current implementation uses a flat 4-column row layout that breaks at 375px. The replacement uses a stacked MUI `Card` per deck entry with a full-width `Autocomplete`, side-by-side Place/Kills `TextField` fields, and a `TooltipIconButton` remove button.

There are no new dependencies, no schema changes, and no data migrations required. All MUI components needed (`Card`, `Autocomplete`, `TextField`, `Select`, `Button`, `Typography`) are already installed and in use in the project. The `TooltipIconButton` reusable component is already available at `app/src/components/TooltipIcon.tsx`.

**Primary recommendation:** Address backend cleanup first (Go compile-checked with `go vet ./lib/...`), then frontend changes in this order: types.ts → stats.tsx → game/index.tsx → new/index.tsx. Run `./node_modules/.bin/tsc --noEmit` from `app/` after each frontend file to catch breakage early.

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Deck Picker Label (GAME-02)**
- D-01: Deck picker label format: `DeckName (PlayerName)` — e.g. `"My Rakdos Deck (Mike)"`. Consistent across all cases (with or without commander set, partner commanders).
- D-02: Applies to the Autocomplete in NewGameView only. Does NOT change the deck label in other views (GameView result display still uses `deck_name`).

**New Game Form Layout (GAME-03)**
- D-03: Layout: stacked card per player. Each player entry is a Box/Card with ❌ in card header, full-width Deck Autocomplete, Place and Kills side-by-side below.
- D-04: Place field always starts blank (no auto-populate by card position).
- D-05: Format field is required and always visible at the top of the form (no collapse).
- D-06: Description field is hidden behind an expand link ("+ Add description"). Not shown by default.
- D-07: "Add Player" button renamed to "Add Deck".
- D-08: Remove button (❌) on each card. Minimum 2 cards enforced.

**Variable-Player Record Display (GAME-04)**
- D-09: `Record` component iterates keys of `RecordDict` from 1 to max(keys), shows counts separated by ` / `. Places with no entries show `0`.
- D-10: `RecordComparator` updated to compare dynamically: iterate from place 1 up to max place in either record dict, return first non-zero difference.

**Player Field Removal (GAME-01)**
- D-11: NewGameView: remove player picker from `GameInput`. Remove `GetPlayersForPod` call from `newGameLoader`. Remove `players` from `NewGameData` type and `NewGameResult` type (`player_id` field removed).
- D-12: GameView AddResultModal: remove the Player Autocomplete. Player is implicit from deck selection.
- D-13: Backend cleanup: remove `player_id` from `addGameResultRequest` struct in `lib/routers/game.go`. Remove `playerID int` parameter from `AddResultFunc` type and `AddResult` constructor in the business layer.

### Claude's Discretion

- Initial number of cards shown when NewGameView loads (2 is the natural default for a 2-player minimum)
- Whether the Remove button is a small icon button (❌) in the card header or a full-width button below each card
- Exact MUI spacing, elevation, and border style for player entry cards
- Whether Place and Kills use `TextField` with `type="number"` or `Select` with numeric options
- Whether the `MatchUpDisplay` in `matches.tsx` is updated for mobile in this phase (TODO comment present — researcher should assess)

### Deferred Ideas (OUT OF SCOPE)

- `MatchUpDisplay` mobile polish — has a TODO comment; not in GAME-03 scope, deferred to DSNG-02 sweep (Phase 5 or later)
- Sorting/reordering player cards by drag — not in scope
- Pending todo: "Security review all API and frontend route authorization" — belongs in Phase 6
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| GAME-01 | Games do not require a player field — decks are the unit of game entry; player is implicit via deck ownership | Backend `AddResult` already drops `playerID`; remove from type, struct, and UI |
| GAME-02 | Deck picker in game form displays owner name alongside commander name (e.g., "Rakdos, Lord of Riots (Mike)") | `Deck.player_name` field confirmed present in `types.ts`; `getOptionLabel` change only |
| GAME-03 | New game form is visually clean and easy to use on mobile | Full redesign of `new/index.tsx` using stacked `Card` layout per 04-UI-SPEC.md |
| GAME-04 | Record renderer supports any number of places (not hardcoded to 4-player games) | `RecordDict` type already dynamic; only display/comparator logic changes |
</phase_requirements>

---

## Standard Stack

### Core (already installed — no new installs required)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| @mui/material | v5.15.2 | Card, Box, Button, TextField, Select, Autocomplete, Typography | Project design system; all components used in this phase are already imported elsewhere |
| @mui/icons-material | v5.15.3 | CloseIcon, AddIcon, PublishIcon | Already used in new/index.tsx and other routes |
| react-router-dom | 6.21.1 | useSubmit, useLoaderData, useParams, LoaderFunctionArgs | Form submission pattern already established; preserve exactly |

### No New Dependencies

All components needed for this phase (`Card`, `Autocomplete`, `TextField`, `Select`, `MenuItem`, `Button`, `Typography`) are already installed. `Card` is the one component not yet used in `new/index.tsx` but it is used elsewhere in the project (`app/src/routes/pod/`, etc.).

**No installation step required.**

---

## Architecture Patterns

### File Modification Map

This phase is entirely modifications to existing files — no new files created.

```
lib/
├── business/game/types.go         # Remove playerID from AddResultFunc signature
├── business/game/functions.go     # Remove playerID param from AddResult constructor
└── routers/game.go                # Remove PlayerID field from addGameResultRequest; update call site

app/src/
├── types.ts                       # Remove player_id from NewGameResult; remove players from NewGameData
├── http.ts                        # Remove GetPlayersForPod from newGameLoader (function itself stays)
├── components/stats.tsx           # Dynamic Record + RecordComparator
└── routes/
    ├── new/index.tsx              # Full redesign (primary complexity)
    └── game/index.tsx             # Remove Player Autocomplete from AddResultModal
```

### Pattern 1: State Management for Stacked Card List

The existing `new/index.tsx` uses a `ResultsMap` keyed by index (`{[key: number]: NewGameResult}`). The redesign must support removing cards by arbitrary index (not just removing from the end). The pattern that fits is an array of card states with stable keys.

**Recommended:** Replace `ResultsMap` + `numPlayers` with `useState<CardState[]>` where each entry carries a `key` (for React reconciliation) plus the deck/place/kills values.

```typescript
// Source: inferred from existing pattern + D-08 (minimum 2 enforcement)
interface CardState {
  key: number;          // stable React key, never reused after removal
  deckId: number | null;
  place: number | null;
  kills: number | null;
}

const [cards, setCards] = useState<CardState[]>([
  { key: 0, deckId: null, place: null, kills: null },
  { key: 1, deckId: null, place: null, kills: null },
]);
const [nextKey, setNextKey] = useState(2);
```

Adding a card appends with `nextKey`, incrementing `nextKey`. Removing filters by key. Updating patches by key.

**Why not re-index?** When cards are removed mid-list, re-indexing shifts all keys and React remounts unchanged components. Stable keys avoid this.

### Pattern 2: Deck Autocomplete with `getOptionLabel`

The current `getOptionLabel` in `new/index.tsx` already references `deck.player_name` (rendered as `— ${deck.player_name}`). The change is format only:

```typescript
// Current (line 206 of new/index.tsx):
getOptionLabel={(deck: Deck) => `${deck.name}${deck.commanders ? ` (${deck.commanders.commander_name})` : ""} — ${deck.player_name}`}

// New (per D-01):
getOptionLabel={(deck: Deck) => `${deck.name} (${deck.player_name})`}
```

The `Deck` interface in `types.ts` already has `player_name: string` — no type change needed for this.

Note: D-01 says the format is `DeckName (PlayerName)` consistently regardless of whether a commander is set. The existing logic that conditionally appends commander name is removed entirely in the new form.

### Pattern 3: Dynamic Record Component

```typescript
// Source: 04-UI-SPEC.md implementation pattern + existing stats.tsx
export function Record({ record }: RecordProps): ReactElement {
  const maxPlace = Math.max(...Object.keys(record).map(Number), 1);
  const parts = Array.from({ length: maxPlace }, (_, i) => record[i + 1] ?? 0);
  return <span className="record">{parts.join(" / ")}</span>;
}

export function RecordComparator(record1: RecordDict, record2: RecordDict): number {
  const maxPlace = Math.max(
    ...Object.keys(record1).map(Number),
    ...Object.keys(record2).map(Number),
    1
  );
  for (let place = 1; place <= maxPlace; place++) {
    const diff = (record1[place] ?? 0) - (record2[place] ?? 0);
    if (diff !== 0) return diff;
  }
  return 0;
}
```

Edge case: `Math.max()` with an empty spread returns `-Infinity`. The `, 1` guard in `maxPlace` prevents this when `record` is empty. Both `Record` and `RecordComparator` need this guard.

### Pattern 4: Backend `AddResult` Cleanup

`AddResult` in `functions.go` currently accepts `playerID` but ignores it (the `gameResultRepository.Model` has no `PlayerID` field — confirmed at line 146-151 of functions.go). The cleanup is:

1. `lib/business/game/types.go` line 19: Change `AddResultFunc func(ctx context.Context, gameID, deckID, playerID, place, killCount int) (int, error)` to remove `playerID int`.
2. `lib/business/game/functions.go` line 207: Remove `playerID int` from the closure signature.
3. `lib/routers/game.go` line 355: Remove `PlayerID int` field from `addGameResultRequest`.
4. `lib/routers/game.go` line 390: Update the `g.games.AddResult(...)` call to drop `req.PlayerID`.
5. `lib/routers/game_test.go` lines 331, 337: Update the `AddResult` mock closure signature and the `addGameResultRequest` literal.

The test at line 337 uses `addGameResultRequest{GameID: 1, DeckID: 10, PlayerID: 42, Place: 2, KillCount: 1}`. After removing `PlayerID` from the struct, the `PlayerID: 42` field reference must be removed from the test body too.

### Pattern 5: `AddResultModal` Player Removal

`AddResultModal` in `game/index.tsx` (line 196-247) uses `playerId` state and guards `handleAdd` with `if (!playerId || !deckId) return`. After removing the Player Autocomplete:
- Remove `playerId` useState
- Remove Player Autocomplete JSX (lines 213-217)
- Remove `player_id: playerId` from the `PostGameResult(...)` call (line 204)
- Update the submit disabled condition to `disabled={!deckId}` only

The `players: PlayerWithRole[]` prop on `AddResultModal` is still needed because the component receives it (GameResultsGrid passes it through). However, after this change, `players` is not used inside `AddResultModal` itself. The cleanest approach per D-12 is to remove the `players` prop from `AddResultModalProps` entirely (since it's no longer consumed), but the `GameResultsGrid` still receives `players` from the loader for other purposes (the DataGrid shows `player_id` in the `deck_name` cell link). Verify whether `players: PlayerWithRole[]` on `AddResultModalProps` and the corresponding prop pass at line 358 should simply be removed.

The `GetPlayersForPod` call stays in `gameLoader` (line 44-50 of game/index.tsx) — it is still used to build `players: PlayerWithRole[]` for `isManager` check and for the DataGrid display. D-12 only removes the player picker UI from AddResultModal; the data fetch is not removed from GameView.

### Pattern 6: Submit Validation in NewGameView

The existing submit uses `useSubmit` with `encType: "application/json"`. The redesigned form must continue using this exact pattern. The `NewGame` interface shape does not change — only `NewGameResult` loses `player_id`.

Submit guard (per 04-UI-SPEC.md):
```typescript
const isSubmittable =
  formatID !== 0 &&
  cards.every((c) => c.deckId !== null && c.place !== null && c.kills !== null);
```

On submit, map `cards` to `results: NewGameResult[]`:
```typescript
const result: NewGame = {
  description: desc,
  format_id: formatID,
  pod_id: Number(podId),
  results: cards.map((c) => ({
    deck_id: c.deckId!,
    place: c.place!,
    kill_count: c.kills!,
  })),
};
submit(result as { [key: string]: any }, { method: "post", encType: "application/json" });
```

### Anti-Patterns to Avoid

- **Re-indexing card array keys on removal:** When a card is removed, do not reset keys to 0..N. Use stable keys to prevent React remounting unchanged cards.
- **Using `index` as React key for a list that can be reordered/removed:** MUI `Autocomplete` internal state (selected value) will reset on other cards if keys shift.
- **Leaving `player_id` in `NewGameResult` in types.ts:** The type must be cleaned up before the NewGameView rewrite or TypeScript will flag missing field on the new `results` map construction.
- **Applying `DeckName (PlayerName)` format to GameView:** D-02 explicitly limits this to NewGameView only.
- **Calling `go build ./...`:** Crawls `app/node_modules/`; use `go vet ./lib/...` for compile checks.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Stable list keys for removable items | Custom key generator | `useState<number>` counter incrementing with each add | Simple, correct, no external dep |
| Mobile touch targets | Custom CSS min-height | `sx={{ minHeight: 44, minWidth: 44 }}` on MUI components | Consistent with project spacing contract |
| Form submission with JSON body | Direct fetch in onClick | `useSubmit` with `encType: "application/json"` | Already established pattern; React Router handles action lifecycle |
| TypeScript compile check | `npm run build` | `./node_modules/.bin/tsc --noEmit` from `app/` | 3-5 seconds vs 60-90 seconds |

---

## Common Pitfalls

### Pitfall 1: `Math.max()` on empty record dict

**What goes wrong:** `Math.max(...Object.keys({}).map(Number))` returns `-Infinity`, causing `Array.from({ length: -Infinity })` to throw or return an empty array silently.

**Why it happens:** `RecordDict` can theoretically be empty for a player with no games. The `stats.tsx` component is called from DataGrid columns, so it must be defensive.

**How to avoid:** Use `Math.max(...Object.keys(record).map(Number), 1)` — the `, 1` floor ensures maxPlace is at least 1. Same pattern in `RecordComparator` for both records.

**Warning signs:** TypeScript won't catch this; it's a runtime issue. Verify with an empty `{}` record in testing.

### Pitfall 2: `NewGameResultWithGame` in types.ts still has `player_id`

**What goes wrong:** D-11 removes `player_id` from `NewGameResult`. But `NewGameResultWithGame` (line 109-115 of types.ts) also has `player_id: number`. This interface is used by `PostGameResult` in `http.ts` — the GameView AddResultModal path. After D-12 removes the player picker from AddResultModal, the `player_id` field in `PostGameResult`'s body should also be removed.

**Why it happens:** Two separate types share the player_id field for different flows (new game creation vs. adding a result to an existing game).

**How to avoid:** Check `NewGameResultWithGame` explicitly and decide whether `player_id` should be removed from it too. The backend `addGameResultRequest` struct is losing `PlayerID` (D-13), so passing it from the frontend would have no effect anyway. Remove it for consistency.

### Pitfall 3: `game_test.go` references `addGameResultRequest.PlayerID` and `AddResultFunc` signature

**What goes wrong:** After removing `PlayerID` from `addGameResultRequest` and `playerID` from `AddResultFunc`, the existing test at line 337 of `game_test.go` references both. Forgetting to update the test file causes `go vet ./lib/...` to fail.

**Why it happens:** The test file directly uses the struct literal and the func signature.

**How to avoid:** Include `lib/routers/game_test.go` in the backend cleanup wave. Specifically:
- Line 331: `AddResult: func(ctx context.Context, gameID, deckID, playerID, place, killCount int) (int, error)` → remove `playerID int`
- Line 337: `addGameResultRequest{GameID: 1, DeckID: 10, PlayerID: 42, Place: 2, KillCount: 1}` → remove `PlayerID: 42`

### Pitfall 4: `GetPlayersForPod` removal scope

**What goes wrong:** D-11 says remove `GetPlayersForPod` from `newGameLoader`. The same function is also called in `gameLoader` in `game/index.tsx` (line 44). If the executor conflates the two and removes it from `gameLoader` as well, the GameView breaks (manager check and DataGrid player links depend on it).

**Why it happens:** Both loaders call `GetPlayersForPod`; the instructions reference only one.

**How to avoid:** The CONTEXT.md integration points note explicitly: "keep the function itself [in http.ts], still used in GameView for the players list." Only remove the `GetPlayersForPod` call from `newGameLoader` in `new/index.tsx`.

### Pitfall 5: `description` field collapse state vs. form submission

**What goes wrong:** The description `TextField` exists in the DOM (just hidden via `display: none` or conditional render). If it's conditionally rendered (`showDescription && <TextField ...>`), the value is lost when the user collapses it. If it persists in state even when hidden, the value is correctly submitted.

**Why it happens:** Toggling visibility via conditional render unmounts the component and resets its state.

**How to avoid:** Keep `desc` in `useState` at the parent component level (already the pattern in the existing code). The expand/collapse toggle only controls visibility (`showDescription` boolean), not whether the value is tracked. The `TextField` value is always driven by `desc` state, not internal component state.

---

## Code Examples

### Dynamic Record Component (from 04-UI-SPEC.md)

```typescript
// Source: .planning/phases/04-game-model-change/04-UI-SPEC.md
export function Record({ record }: RecordProps): ReactElement {
  const maxPlace = Math.max(...Object.keys(record).map(Number), 1);
  const parts = Array.from({ length: maxPlace }, (_, i) => record[i + 1] ?? 0);
  return <span className="record">{parts.join(" / ")}</span>;
}
```

### Card Header Row with TooltipIconButton (from 04-UI-SPEC.md)

```typescript
// Source: .planning/phases/04-game-model-change/04-UI-SPEC.md
// TooltipIconButton already used in game/index.tsx — import from ../../components/TooltipIcon
<Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 1 }}>
  <Typography variant="body2">Deck {index + 1}</Typography>
  <TooltipIconButton
    title={cards.length <= 2 ? "Minimum 2 entries required" : "Remove"}
    icon={<CloseIcon />}
    onClick={() => removeCard(card.key)}
    color="error"
    size="small"
    sx={{ minHeight: 44, minWidth: 44 }}
    disabled={cards.length <= 2}
  />
</Box>
```

### Backend AddResult After Cleanup

```go
// Source: lib/business/game/types.go (after D-13 cleanup)
type AddResultFunc func(ctx context.Context, gameID, deckID, place, killCount int) (int, error)

// Source: lib/business/game/functions.go (after D-13 cleanup)
func AddResult(gameResultRepo repos.GameResultRepository) AddResultFunc {
    return func(ctx context.Context, gameID, deckID, place, killCount int) (int, error) {
        return gameResultRepo.Add(ctx, gameResultRepository.Model{
            GameID:    gameID,
            DeckID:    deckID,
            Place:     place,
            KillCount: killCount,
        })
    }
}

// Source: lib/routers/game.go (after D-13 cleanup)
type addGameResultRequest struct {
    GameID    int `json:"game_id"`
    DeckID    int `json:"deck_id"`
    Place     int `json:"place"`
    KillCount int `json:"kill_count"`
}
// Call site: g.games.AddResult(ctx, req.GameID, req.DeckID, req.Place, req.KillCount)
```

### Deck Autocomplete Label (from 04-CONTEXT.md D-01)

```typescript
// Source: .planning/phases/04-game-model-change/04-CONTEXT.md D-01
getOptionLabel={(deck: Deck) => `${deck.name} (${deck.player_name})`}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Flat 4-column row layout in new/index.tsx | Stacked Card per deck entry | Phase 4 | Eliminates horizontal overflow at 375px |
| Hardcoded 4-place Record display | Dynamic Record iterating RecordDict keys | Phase 4 | Supports 2-player, 3-player, 5-player, etc. |
| `AddResult` accepts and silently drops `playerID` | `AddResult` removes `playerID` from signature | Phase 4 | API contract matches actual behaviour |
| Player picker in NewGameView and AddResultModal | Deck-only entry; player implicit from deck ownership | Phase 4 | Reduces form friction; consistent with data model |

---

## Open Questions

1. **`NewGameResultWithGame.player_id` in types.ts**
   - What we know: `NewGameResultWithGame` (used by `PostGameResult` in http.ts → AddResultModal) has `player_id: number`. The backend's `addGameResultRequest` is losing `PlayerID` (D-13).
   - What's unclear: D-11 only explicitly names `NewGameResult.player_id` for removal. `NewGameResultWithGame.player_id` is a separate type.
   - Recommendation: Remove `player_id` from `NewGameResultWithGame` too. Passing it would be a no-op after D-13, and leaving it creates inconsistency. Mark this as a task in the plan.

2. **`players` prop on `AddResultModalProps`**
   - What we know: After removing the Player Autocomplete from `AddResultModal`, the `players: PlayerWithRole[]` prop is no longer consumed inside that component.
   - What's unclear: Whether removing the prop is safe without a broader refactor of `GameResultsGrid` which constructs the modal.
   - Recommendation: Remove the `players` prop from `AddResultModalProps` and the prop pass at the `AddResultModal` JSX site in `GameResultsGrid`. The `players` array stays in `GameResultsGrid`'s props (for the `isManager` check path), just stop threading it into `AddResultModal`.

---

## Environment Availability

Step 2.6: SKIPPED — this phase is code/configuration changes only. No external tools, services, or runtimes are introduced beyond what is already running (Go toolchain, Node.js already confirmed in use).

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go: testify v1.11.1 + `net/http/httptest`; Frontend: TypeScript compiler (type checking) |
| Config file | Go: none (standard `go test`); Frontend: `app/tsconfig.json` |
| Quick run command | `go vet ./lib/... && go test ./lib/...` (backend); `cd /mnt/d/msharp/Documents/projects/edh-tracker/app && ./node_modules/.bin/tsc --noEmit` (frontend) |
| Full suite command | Same — no separate full suite |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|--------------|
| GAME-01 | Backend: `addGameResultRequest` has no `player_id`; `AddResultFunc` has no `playerID` param | unit (Go compile + existing router test) | `go vet ./lib/... && go test ./lib/routers/ -run TestGameRouter_AddGameResult` | ✅ game_test.go |
| GAME-01 | Frontend: `NewGameResult` has no `player_id`; form submits without player | type-check | `./node_modules/.bin/tsc --noEmit` | ✅ types.ts |
| GAME-02 | Deck Autocomplete shows `DeckName (PlayerName)` format | type-check (no runtime test) | `./node_modules/.bin/tsc --noEmit` | ✅ new/index.tsx |
| GAME-03 | Form renders without horizontal scroll at 375px | manual visual (no automated) | manual-only: open browser at 375px | — |
| GAME-04 | Record with 3 keys shows 3 parts; Record with empty key shows `0` | unit (TypeScript) | `./node_modules/.bin/tsc --noEmit` (type safety only; logic must be manual or tested via smoke test) | ✅ stats.tsx |

### Sampling Rate

- **Per task commit:** `go vet ./lib/...` (after backend changes) or `./node_modules/.bin/tsc --noEmit` (after frontend changes)
- **Per wave merge:** `go test ./lib/...` + `./node_modules/.bin/tsc --noEmit`
- **Phase gate:** Both suites green before `/gsd:verify-work`; manual mobile visual check at 375px

### Wave 0 Gaps

None — existing test infrastructure covers all phase requirements. The Go router tests for `AddGameResult` already exist in `game_test.go` and will need updating (not creation) as part of the GAME-01 backend task.

---

## Project Constraints (from CLAUDE.md)

These directives from CLAUDE.md apply to all implementation work in this phase:

- **Compile check:** Use `go vet ./lib/...` (not `go build ./...` or `go build ./lib/...`)
- **Frontend type check:** Use `./node_modules/.bin/tsc --noEmit` from `app/` (not `npm run build`, not `npx tsc`)
- **Frontend verify skill:** Run after any edit to files under `app/src/` — the `frontend-verify` skill must be triggered
- **No framework changes:** Stay on Go + Gorilla Mux + GORM + MySQL backend; React + MUI + React Router v6 frontend
- **No breaking DB changes:** No schema migrations in this phase (none needed)
- **JSON tags:** Go structs use `snake_case` JSON tags; TypeScript interfaces mirror `snake_case`
- **Architecture layers:** HTTP handlers in `lib/routers/`, business logic in `lib/business/`, no leakage between layers
- **Import source:** `LoaderFunctionArgs` imported from `@remix-run/router/utils` (not `react-router-dom`)
- **HTTP calls centralized:** All fetch calls in `app/src/http.ts`; never call `fetch` directly from components
- **Credentials:** Every `fetch` call must include `credentials: "include"`
- **GSD workflow:** All file changes go through GSD entry points (`/gsd:execute-phase`)

---

## Sources

### Primary (HIGH confidence)

- Direct source read: `app/src/routes/new/index.tsx` — current state of NewGameView, GameInput, loader, action
- Direct source read: `app/src/routes/game/index.tsx` — AddResultModal player picker confirmed at lines 213-217
- Direct source read: `app/src/components/stats.tsx` — current hardcoded 4-place Record implementation
- Direct source read: `app/src/types.ts` — confirmed `NewGameResult.player_id`, `NewGameData.players`, `NewGameResultWithGame.player_id`
- Direct source read: `lib/routers/game.go` — confirmed `addGameResultRequest.PlayerID` at line 355, call site at line 390
- Direct source read: `lib/business/game/functions.go` — confirmed `AddResult` accepts but silently drops `playerID` (Model has no PlayerID field)
- Direct source read: `lib/business/game/types.go` — confirmed `AddResultFunc` signature includes `playerID int`
- Direct source read: `lib/routers/game_test.go` — confirmed test references to `AddResult` mock (line 331) and `addGameResultRequest{..., PlayerID: 42, ...}` (line 337) that need updating
- Direct source read: `lib/business/gameResult/entity.go` — confirmed `InputEntity` has no `player_id`; already deck-only
- Direct source read: `.planning/phases/04-game-model-change/04-CONTEXT.md` — all implementation decisions
- Direct source read: `.planning/phases/04-game-model-change/04-UI-SPEC.md` — visual and interaction contract (approved)
- Direct source read: `.planning/phases/02-design-language/02-UI-SPEC.md` — inherited design system

### Secondary (MEDIUM confidence)

- CLAUDE.md project instructions — coding conventions, compile commands, forbidden patterns
- MEMORY.md project memory — architecture decisions, Go testing patterns

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all libraries confirmed installed and in use
- Architecture: HIGH — all patterns derived from direct source reads of files to be modified
- Pitfalls: HIGH — identified from direct inspection of existing code and cross-referencing decisions
- Backend cleanup: HIGH — `gameResult.InputEntity` already has no `player_id`; the param is truly dead

**Research date:** 2026-03-24
**Valid until:** 2026-04-24 (stable stack; no external dependencies)
