---
plan: 05-05
phase: 05-pod-deck-ux
status: complete
completed: 2026-03-26
commits:
  - 8b4d7ab
  - f0a18d5
---

# Plan 05-05: Deck Creation Route — Summary

## What Was Built

1. **Types** (`app/src/types.ts`) — Added `NewDeckRequest` and `NewDeckData` interfaces.

2. **PostDeck** (`app/src/http.ts`) — New `PostDeck(body: NewDeckRequest): Promise<{ id: number }>` function.

3. **/deck/new route** (`app/src/routes/deck/new/index.tsx`) — New deck creation form with:
   - Name (required) and Format (required) fields
   - Commander + Partner Commander fields (conditional, Commander format only)
   - freeSolo Autocomplete on commander fields — typing a new name creates it inline via PostCommander
   - Create option displayed as `Create "{input}"` via filterOptions
   - After deck creation, auto-navigates to `/player/{callerId}/deck/{newDeckId}`
   - "Create Deck" submit and "Discard" cancel buttons

4. **Route registration** (`app/src/index.tsx`) — `/deck/new` path registered with `newDeckLoader`.

5. **Add Deck button** (`app/src/routes/player/DecksTab.tsx`) — Owner-only "Add Deck" button linking to `/deck/new` appears above the DataGrid (and in empty state with owner-aware copy).

6. **DeckSettingsTab freeSolo** (`app/src/routes/deck/SettingsTab.tsx`) — Both Commander and Partner Commander Autocomplete fields updated to support freeSolo with inline commander creation via PostCommander. Removed `getOptionKey` for freeSolo compatibility.

## Key Files

- `app/src/routes/deck/new/index.tsx` — New route component + loader
- `app/src/http.ts` — PostDeck function
- `app/src/types.ts` — NewDeckRequest, NewDeckData
- `app/src/index.tsx` — Route registration
- `app/src/routes/player/DecksTab.tsx` — Add Deck button (owner-only)
- `app/src/routes/deck/SettingsTab.tsx` — freeSolo commander Autocomplete

## Self-Check: PASSED

- /deck/new route renders deck creation form ✓
- Commander fields conditional on Commander format ✓
- freeSolo Autocomplete creates commanders inline ✓
- After deck creation, navigates to player/deck view ✓
- Add Deck button owner-only on PlayerDecksTab ✓
- DeckSettingsTab Autocomplete supports freeSolo ✓
- TypeScript: exit 0 ✓
