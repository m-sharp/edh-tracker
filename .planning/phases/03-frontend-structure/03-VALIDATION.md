---
phase: 3
slug: frontend-structure
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-23
---

# Phase 3 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | jest 29.x (react-scripts / Create React App) |
| **Config file** | `app/package.json` (jest config inline) |
| **Quick run command** | `cd app && ./node_modules/.bin/tsc --noEmit` |
| **Full suite command** | `cd app && npm test -- --watchAll=false` |
| **Estimated runtime** | ~30 seconds (tsc), ~60 seconds (jest) |

---

## Sampling Rate

- **After every task commit:** Run `cd app && ./node_modules/.bin/tsc --noEmit`
- **After every plan wave:** Run `cd app && npm test -- --watchAll=false`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 60 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 3-01-01 | 01 | 1 | FEND-01 | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ❌ W0 | ⬜ pending |
| 3-01-02 | 01 | 1 | FEND-02 | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ❌ W0 | ⬜ pending |
| 3-01-03 | 01 | 1 | FEND-03 | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ❌ W0 | ⬜ pending |
| 3-02-01 | 02 | 2 | FEND-04 | manual | Refresh `/pod/:podId` in browser — no blank screen | N/A | ⬜ pending |
| 3-02-02 | 02 | 2 | FEND-05 | manual | Load HomeView with no pods — verify no flash of "No pods yet" | N/A | ⬜ pending |
| 3-03-01 | 03 | 3 | DSNG-04 | manual | Visually audit Login, Home, Player, Deck, Game views | N/A | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- Existing TypeScript infrastructure covers compile-time verification.
- No new test files required — all automated checks use `tsc --noEmit`.

*Existing infrastructure covers all phase requirements for automated checks.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| No blank screen on page refresh | FEND-04 | Browser navigation state not testable with jest | Load app, navigate to `/pod/:podId`, refresh browser, confirm page renders |
| No "No pods yet" flash on HomeView | FEND-05 | Loading state timing not reproducible in unit tests | Clear cache, load HomeView, confirm no flash of empty state before data arrives |
| Visual audit of all views | DSNG-04 | Layout/spacing correctness is visual | Open each view (Login, Home, Player, Deck, Game) in browser, verify consistent typography and spacing per Phase 2 design system |
| Tab active state persists in query string | FEND-02 | URL state requires real browser navigation | On Pod/Player/Deck views, click tabs and verify URL updates; refresh and verify correct tab is active |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 60s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
