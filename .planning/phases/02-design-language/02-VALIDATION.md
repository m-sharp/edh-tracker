---
phase: 2
slug: design-language
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-23
---

# Phase 2 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | TypeScript compiler (tsc) — no automated test framework for frontend |
| **Config file** | `app/tsconfig.json` |
| **Quick run command** | `cd app && ./node_modules/.bin/tsc --noEmit` |
| **Full suite command** | `cd app && ./node_modules/.bin/tsc --noEmit` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `cd app && ./node_modules/.bin/tsc --noEmit`
- **After every plan wave:** Run `cd app && ./node_modules/.bin/tsc --noEmit`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 2-01-01 | 01 | 1 | DSNG-01 | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ | ⬜ pending |
| 2-01-02 | 01 | 1 | DSNG-01 | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ | ⬜ pending |
| 2-01-03 | 01 | 2 | DSNG-02 | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ | ⬜ pending |
| 2-02-01 | 02 | 2 | DSNG-02 | manual | n/a — visual inspection | n/a | ⬜ pending |
| 2-02-02 | 02 | 2 | DSNG-03 | manual | n/a — viewport resize to 375px | n/a | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

*Existing infrastructure covers all phase requirements.* No new test infrastructure needed — TypeScript compiler is pre-installed in `app/node_modules`.

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| MUI component overrides visually correct (AppBar dark navy, gold accents, no unstyled components) | DSNG-02 | Visual design can't be verified by compiler | Open app in browser, confirm AppBar is `#0d1b2a` (dark navy), primary buttons/accents are `#c9a227` (gold), background is `#1a2637` |
| Phone viewport usability | DSNG-03 | Viewport rendering is visual | Open PodView in DevTools with 375px viewport; confirm text readable without zoom, touch targets ≥ 44px, no horizontal overflow |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
