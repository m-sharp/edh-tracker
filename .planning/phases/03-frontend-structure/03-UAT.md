---
status: complete
phase: 03-frontend-structure
source: 03-01-SUMMARY.md, 03-02-SUMMARY.md, 03-03-SUMMARY.md, 03-04-SUMMARY.md, 03-05-SUMMARY.md, 03-06-SUMMARY.md
started: 2026-03-24T05:00:00Z
updated: 2026-03-24T05:30:00Z
---

## Current Test

## Current Test

[testing complete]

## Tests

### 1. Tab state persists in URL
expected: Navigate to a Pod view. Click between the Decks, Players, Games, and Settings tabs. The URL query string should update with each click (e.g., ?podTab=players, ?podTab=games). Refresh the page while on a non-default tab — you should land back on the same tab, not reset to the first tab. Same behavior applies on Player view (?playerTab=...) and Deck view (?deckTab=...).
result: pass

### 2. No blank screen on refresh
expected: While logged in, navigate to any deep page (Pod, Player, Deck, or Game view). Hit browser refresh. You should NOT see a blank white screen. A centered loading spinner should appear briefly while auth is checked, then the page renders normally.
result: issue
reported: "Refreshing still has the broken white page behavior. Accompanied by Uncaught SyntaxError: Unexpected token '<' console error. Network requests show the server returning HTML index page for static asset requests (JS/CSS). The SPA handler in app/main.go falls back to index.html for any path not found on disk — including JS/CSS asset requests — so the app fails to boot."
severity: blocker

### 3. No "No pods yet" flash on home
expected: After logging in, navigate to the home page (or refresh it). Watch carefully — you should see a loading spinner first, then either navigate to your pod or (if you have no pods) see the empty state message. The message "No pods yet" or "Create your first pod" should NOT appear for a split second before your pods load.
result: pass

### 4. Login page polish
expected: Navigate to /login while logged out. The page should be vertically centered on screen (not pushed to the top). Above the "EDH Tracker" heading, a playing cards icon should be visible.
result: issue
reported: "Renders as expected, though the content with playing cards icon starts too low in the page. Needs to be brought up by changing to justify-content: flex-start and adding some top margin or padding on desktop."
severity: cosmetic

### 5. AppBar title hidden on mobile
expected: Open the app on a phone or resize the browser to ~375px wide. In the top AppBar, the "EDH Tracker" text label should NOT be visible — only the playing cards icon (if present), pod selector, avatar, and logout button should show. At desktop width (~900px+), the "EDH Tracker" text should reappear.
result: issue
reported: "AppBar title is correctly hidden. The LOGOUT text can clip or crowd very closely to right side on small screens — look to replace with an appropriate icon button."
severity: minor

### 6. Deck view: Commanders heading has tooltip
expected: Navigate to any deck's settings tab. Find the "Commanders" section heading. There should be a small info icon (ⓘ) next to or inline with the "Commanders" label. Hovering over it (desktop) or tapping it (mobile) should show a tooltip explaining how commanders work (disambiguation text).
result: issue
reported: "Tooltip shows, should shift to popping up ABOVE the icon, not BELOW."
severity: cosmetic

### 7. Game view: Icon buttons have tooltips
expected: Navigate to any game result view. Find the icon buttons for editing and removing results (pencil/edit and X/trash icons). Hovering over each button on desktop should show a tooltip label: "Edit result", "Remove result". The description edit button should show "Edit description".
result: pass

### 8. Mobile touch targets in Pod Players tab
expected: In the Pod Players tab, find the Promote and Remove buttons next to each player. On a phone or narrow viewport, these buttons should be comfortably tappable — they should be at least 44px tall with no need to tap precisely on the text. (You can verify by inspecting the element or simply tapping to confirm it feels natural.)
result: issue
reported: "Buttons look tall enough. I'd rather see 'contained' buttons — currently they appear as just colored text on black background. There should also be confirmation dialogs when clicking PROMOTE and REMOVE — the actions dangerously go straight through as is."
severity: major

## Summary

total: 8
passed: 3
issues: 5
pending: 0
skipped: 0

## Gaps

- truth: "Refreshing any page in the app does not produce a blank white screen — a centered loading spinner appears, then the page renders"
  status: failed
  reason: "User reported: SPA handler in app/main.go falls back to index.html for ALL non-existent paths, including JS/CSS static asset requests. Browser gets HTML when expecting JS, causing Uncaught SyntaxError: Unexpected token '<'. App does not boot at all after refresh."
  severity: blocker
  test: 2
  artifacts: [app/main.go]
  missing: []

- truth: "Login page content is positioned with adequate top spacing on desktop — not too low on screen"
  status: failed
  reason: "User reported: content starts too low. Fix: change justifyContent from center to flex-start and add top padding/margin on desktop viewports."
  severity: cosmetic
  test: 4
  artifacts: [app/src/routes/login.tsx]
  missing: []

- truth: "The AppBar logout control does not clip or crowd on small screens"
  status: failed
  reason: "User reported: LOGOUT text clips/crowds on right side of AppBar on small screens. Fix: replace text Button with an icon button (e.g., LogoutIcon)."
  severity: minor
  test: 5
  artifacts: [app/src/routes/root.tsx]
  missing: []

- truth: "TooltipIcon tooltips open above the icon, not below"
  status: failed
  reason: "User reported: tooltip pops up below the icon. Fix: add placement='top' to the Tooltip in TooltipIcon component (or at the call site in DeckView SettingsTab)."
  severity: cosmetic
  test: 6
  artifacts: [app/src/components/TooltipIcon.tsx, app/src/routes/deck/SettingsTab.tsx]
  missing: []

- truth: "Promote and Remove buttons in Pod Players tab are visually distinct (contained style) and require confirmation before executing"
  status: failed
  reason: "User reported: buttons appear as colored text on black background (not contained). Promote and Remove execute immediately with no confirmation dialog — dangerous for irreversible actions."
  severity: major
  test: 8
  artifacts: [app/src/routes/pod/PlayersTab.tsx]
  missing: []
