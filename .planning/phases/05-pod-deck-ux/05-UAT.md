---
status: complete
phase: 05-pod-deck-ux
source: [05-01-SUMMARY.md, 05-02-SUMMARY.md, 05-03-SUMMARY.md, 05-04-SUMMARY.md, 05-05-SUMMARY.md]
started: 2026-03-26T00:00:00Z
updated: 2026-03-26T00:00:00Z
---

## Current Test
<!-- OVERWRITE each test - shows where we are -->

## Current Test

[testing complete]

## Tests

### 1. Record shows 4 place columns
expected: Open any deck stats or player stats that has fewer than 4 distinct finish positions. The record display still shows at least 4 columns (e.g. a deck with only 1 win shows "1 / 0 / 0 / 0" — not just "1").
result: pass

### 2. Pod Decks tab default sort
expected: Open a pod's Decks tab. The deck list is sorted by record descending by default, without needing to click any column header.
result: issue
reported: "Decks tab is sorted by record correctly. The issue around pagination limiting results remains - the truely best/greatest record deck does not appear at the top of the list as it is on a different page of results."
severity: major

### 3. Playing cards icon links home
expected: Click the playing cards icon in the AppBar. You are taken to the home page (/).
result: pass

### 4. Player Decks tab hides retired decks
expected: On a player's Decks tab, any retired decks are not shown in the list. There is also no "Is Retired" boolean column in the grid.
result: issue
reported: "This is the behavior, but it was not desired. The desired behavior was to have retired decks filtered out of the data grid by default - they should still be accessible via the data grid filtering options. As is, there is no way to ever find one's retired decks."
severity: major

### 5. Create Pod removed from Player Settings
expected: Open your player's Settings tab. There is no "Create New Pod" form, text field, or button anywhere on the page.
result: pass

### 6. New game Discard button
expected: Open a pod's New Game page. A "Discard" button appears alongside the Submit button. Clicking Discard takes you back to the pod page without creating a game.
result: pass

### 7. Home screen onboarding (no pods)
expected: Log in as a user with no pods (or view as a new user). The home screen shows a "Welcome to EDH Tracker" heading and a "Create a Pod" button instead of the pod dashboard.
result: pass

### 8. Create Pod from home screen
expected: On the onboarding home screen, click "Create a Pod". A dialog opens with a Pod Name text field, a "Discard" button, and a "Create Pod" submit button. After creating a pod, you are automatically navigated to the new pod's page (/pod/{id}).
result: issue
reported: "Creating a new pod works and redirects to the new pod's page correctly. On redirect, the pod selector state in the nav bar still shows as empty with 'No pods' and 'Create new pod' as the options. Even on hard refresh, the pod does not show up for the player. A PodPlayer record was NOT created for the player - leaving this pod effectively orphaned."
severity: blocker

### 9. AppBar "Create new pod" option
expected: Open the pod selector dropdown in the AppBar. At the bottom, separated by a divider and styled in primary color, there is a "Create new pod" option. Selecting it opens the create pod dialog, and after creation you land on the new pod's page.
result: pass

### 10. Players tab card layout
expected: Open a pod's Players tab. Each player is displayed in a card (not a flat list). Each card shows the player's name (as a clickable link to their profile) and pod-scoped stats: record (W-L-D), points, and kills.
result: pass

### 11. Manager chip on Players tab
expected: On the Players tab, players who are pod managers have a "Manager" chip on their card. If you are the pod manager, you see promote (PersonAdd) and remove (PersonOff) icon buttons on other players' cards. Non-managers do not see these action buttons.
result: pass

### 12. Promote/remove dialog copy
expected: Click promote (PersonAdd) on a player card. The dialog title includes the player's name, the cancel button says "Never mind", and the confirm button says "Make Manager". Similarly for remove: title includes name, cancel is "Never mind", confirm is "Remove".
result: issue
reported: "Copy is right, would rather have 'Cancel' over 'Never mind'"
severity: cosmetic

### 13. Add Deck button (owner only)
expected: On your own player profile's Decks tab, an "Add Deck" button appears above the deck grid (and in the empty state). View another player's Decks tab — the "Add Deck" button is not visible.
result: pass

### 14. Deck creation form and conditional commander fields
expected: Click "Add Deck" to open /deck/new. The form shows Name and Format fields. Select the "Commander" format — Commander and Partner Commander autocomplete fields appear. Select a different format — those fields are hidden.
result: pass

### 15. FreeSolo commander creation in new deck form
expected: On the /deck/new form with Commander format selected, type a new commander name that doesn't exist yet. A "Create "{name}"" option appears in the dropdown. Selecting it creates the commander inline and fills the field — no separate page or form needed.
result: issue
reported: "On submitting the new commander name, the 'Failed to create commander. Try again.' error text is shown. On refresh, the new commander is in the autocomplete. The network tab shows that a POST to /api/commander succeeds with a 201 status response. No body is included with the response though - handling of the commander post probably expects an ID"
severity: major

### 16. After deck creation navigation
expected: Complete the /deck/new form and submit. You are automatically navigated to the new deck's detail page (/player/{playerId}/deck/{newDeckId}).
result: pass

### 17. DeckSettingsTab freeSolo commander
expected: Open an existing deck's Settings tab. The Commander autocomplete field supports typing a new commander name (freeSolo). A "Create "{name}"" option appears, and selecting it creates the commander inline without leaving the page.
result: pass
notes: "Same POST /api/commander no-body bug applies (already logged under test 15)"

## Summary

total: 17
passed: 12
issues: 5
pending: 0
skipped: 0
blocked: 0

## Gaps

- truth: "Pod Decks tab opens sorted by record descending so the best deck appears at the top"
  status: failed
  reason: "User reported: Decks tab is sorted by record correctly. The issue around pagination limiting results remains - the truely best/greatest record deck does not appear at the top of the list as it is on a different page of results."
  severity: major
  test: 2
  root_cause: ""
  artifacts: []
  missing: []
  debug_session: ""

- truth: "Creating a pod via the home screen onboarding dialog adds the creator as a pod member, so the pod appears in their pod selector immediately"
  status: failed
  reason: "User reported: pod creation navigates correctly but no PodPlayer record was created for the creator. Pod selector shows 'No pods' even after hard refresh. Pod is orphaned — creator cannot access it."
  severity: blocker
  test: 8
  root_cause: "Backend POST /api/pod creates the pod record but does not add the calling player as a member. Fix: backend pod Create business function should create the pod AND add the caller as a pod member (Manager role) in a single transaction."
  artifacts: []
  missing:
    - "lib/business/pod/functions.go: Create function must wrap pod insert + player_pod_role insert in a transaction"
  debug_session: ""

- truth: "FreeSolo commander creation fills the autocomplete field inline after the POST succeeds"
  status: failed
  reason: "User reported: POST /api/commander returns 201 with no body. Frontend shows 'Failed to create commander' error. Commander is actually created (visible after refresh). Frontend's PostCommander handler expects an ID in the response body to wire the new option."
  severity: major
  test: 15
  root_cause: "POST /api/commander returns 201 with empty body — same pattern fixed for pod/deck in plan 05-01 but missed for commander. Backend lib/routers/commander.go needs to return {\"id\": N}. Frontend PostCommander in http.ts expects Promise<{id: number}> to set the selected commander."
  artifacts:
    - path: "lib/routers/commander.go"
      issue: "CommanderCreate handler returns 201 with no body; needs to return {\"id\": N}"
    - path: "app/src/http.ts"
      issue: "PostCommander return type and error handling expects id from response body"
  missing:
    - "lib/routers/commander.go: capture LastInsertId from business layer and write JSON {\"id\": N} like pod/deck handlers"
    - "app/src/http.ts: PostCommander return type should be Promise<{id: number}>"
  debug_session: ""

- truth: "Cancel buttons in promote/remove dialogs say 'Cancel'"
  status: failed
  reason: "User reported: copy is right per spec but prefers 'Cancel' over 'Never mind'"
  severity: cosmetic
  test: 12
  root_cause: "UI-SPEC specified 'Never mind'; user preference is 'Cancel' — simple string change in PlayersTab.tsx dialog buttons"
  artifacts:
    - path: "app/src/routes/pod/PlayersTab.tsx"
      issue: "Cancel button label 'Never mind' should be 'Cancel'"
  missing:
    - "Change 'Never mind' to 'Cancel' in promote and remove dialog cancel buttons"
  debug_session: ""

- truth: "Retired decks are hidden by default via DataGrid filter state but remain accessible when user clears/removes the filter — the 'Is Retired' column must be present for filtering to work"
  status: failed
  reason: "User reported: retired decks are hard-filtered out in JS (visibleRows), removing them entirely from the grid with no way to access them. The Is Retired column was also removed. Desired behavior: use DataGrid initialState filterModel to hide retired by default, keeping the column available for the user to remove the filter."
  severity: major
  test: 4
  root_cause: ""
  artifacts: []
  missing: []
  debug_session: ""
