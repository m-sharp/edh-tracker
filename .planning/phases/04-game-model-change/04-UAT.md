---
status: complete
phase: 04-game-model-change
source: 04-01-SUMMARY.md, 04-02-SUMMARY.md, 04-03-SUMMARY.md
started: 2026-03-24T20:00:00Z
updated: 2026-03-24T20:10:00Z
---

## Current Test

[testing complete]

## Tests

### 1. Add Result — No Player Picker
expected: Open a game page and click the "Add Result" button. The modal should show only a deck selector, place field, and kill count field — no player dropdown or player picker.
result: issue
reported: "The Deck selector in the Add Result modal should show '<DeckName> - (<Player Name>)' like the new game form. Place and Kills should have min & max bounds like the new game form."
severity: major

### 2. Add Result — Submit requires deck only
expected: In the Add Result modal, leave the deck unselected — the Submit button should be disabled. Select any deck — Submit should become enabled. Kill count and place alone do not enable Submit.
result: pass

### 3. New Game — Stacked card layout, no player picker
expected: Navigate to a pod's New Game page. The form should display stacked cards (one per result entry), each with a deck selector, place field, and kills field. No player picker or player dropdown should appear anywhere on this page.
result: issue
reported: "Cards appear correctly but the X to remove a card is left in a top row all by itself, leaving a lot of dead space. A better view would have the X still right aligned, but on the same row as the input fields. Also: the New Game link should be front and center on the Pod page, not hidden in the Game tab that is not loaded first — it looks awkward all on a line by itself at the top of the Game tab."
severity: major

### 4. New Game — Dynamic place and kills bounds
expected: On the New Game form with 3 cards, the Place field should accept values 1–3 and the Kills field should accept 0–3. Add a 4th card — both bounds should update to 1–4 / 0–4. The bounds follow the number of cards, not a hardcoded 4.
result: pass

### 5. New Game — Redirect to game page after submit
expected: Fill out a valid new game (select a format, fill in deck/place/kills for each card) and submit. You should be redirected to the newly created game's detail page (e.g., /pod/:podId/game/:id) — not left on the New Game form.
result: pass

### 6. New Game — Format label renders correctly
expected: On the New Game form, the Format selector should display its label correctly: the label text is visible and the outlined input border has a proper notch around it (label not overlapping the border, not hidden).
result: pass

### 7. Record displays correct number of places
expected: View a player's record on their profile (or deck stats). A record from a 3-player game should show 3 slots (e.g., "2 / 0 / 1"), not always 4. A record from a 5-player game should show 5 slots. The display is dynamic based on actual pod size, not hardcoded.
result: pass

## Summary

total: 7
passed: 5
issues: 2
pending: 0
skipped: 0
blocked: 0

## Gaps

- truth: "Add Result modal: deck selector shows '<DeckName> - (<Player Name>)' and Place/Kills have min/max bounds matching the New Game form"
  status: failed
  reason: "User reported: The Deck selector in the Add Result modal should show '<DeckName> - (<Player Name>)' like the new game form. Place and Kills should have min & max bounds like the new game form."
  severity: major
  test: 1
  root_cause: "Two issues: (1) getOptionLabel in AddResultModal returns only d.name (game/index.tsx line 214); d.player_name is available on the Deck type and just needs to be included. Same issue exists in EditResultModal line 144. (2) Place/Kills TextFields in AddResultModal have no inputProps min/max. The bound equivalent is game.results.length + 1 (total players in the game), which is available at the GameResultsGrid call site but not passed into the modal — needs a new playerCount prop."
  artifacts:
    - path: "app/src/routes/game/index.tsx"
      issue: "getOptionLabel returns d.name only (line 214); no inputProps on Place/Kills fields (lines 217-230); EditResultModal same label issue (line 144)"
  missing:
    - "Change getOptionLabel to `${d.name} (${d.player_name})` in AddResultModal and EditResultModal"
    - "Add playerCount prop to AddResultModal; pass game.results.length + 1 from GameResultsGrid; add inputProps min/max to Place and Kills fields"
  debug_session: ""

- truth: "New Game card layout: remove button (X) is inline with the input fields on the same row, not isolated above them. New Game entry point is prominently accessible from the Pod page, not buried in the Games tab."
  status: failed
  reason: "User reported: The X to remove a card is in a top row all by itself, leaving dead space — should be right-aligned but on the same row as the input fields. Also: the New Game link should be front and center on the Pod page, not hidden in the Game tab that is not loaded first."
  severity: major
  test: 3
  root_cause: "Two issues: (1) In new/index.tsx, the remove button is wrapped in its own full-width Box (justifyContent: flex-end, mb: 1) as a sibling to the Autocomplete and Place/Kills fields — it occupies an entire card row alone. Fix: put fields and button in a single row-flex container. (2) New Game button lives in GamesTab.tsx (lines 54-60) inside the Games tab (tab index 2, not the default). The button should be lifted to pod/index.tsx level — visible on the pod page regardless of active tab, e.g. in the page header beside the pod name."
  artifacts:
    - path: "app/src/routes/new/index.tsx"
      issue: "Remove button in isolated Box row above fields (lines 146-189); needs inline row layout with fields"
    - path: "app/src/routes/pod/GamesTab.tsx"
      issue: "New Game button at lines 54-60 inside Games tab (tab index 2)"
    - path: "app/src/routes/pod/index.tsx"
      issue: "New Game button needs to move here — into page header (line 44 area), above tab strip"
  missing:
    - "Restructure card interior in new/index.tsx: row-flex container with fields column on left, X button on right"
    - "Move New Game button from GamesTab.tsx to pod/index.tsx header area; remove it from GamesTab"
  debug_session: ""
