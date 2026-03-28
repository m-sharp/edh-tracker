---
phase: 5
slug: pod-deck-ux
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-25
---

# Phase 5 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test (backend) / tsc --noEmit (frontend) |
| **Config file** | none — existing infrastructure |
| **Quick run command** | `cd app && ./node_modules/.bin/tsc --noEmit` |
| **Full suite command** | `go vet ./lib/... && cd app && ./node_modules/.bin/tsc --noEmit` |
| **Estimated runtime** | ~15 seconds |

---

## Sampling Rate

- **After every task commit:** Run `cd app && ./node_modules/.bin/tsc --noEmit`
- **After every plan wave:** Run `go vet ./lib/... && cd app && ./node_modules/.bin/tsc --noEmit`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 15 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 5-01-01 | 01 | 1 | POD-01 | compile | `go vet ./lib/...` | ✅ | ⬜ pending |
| 5-01-02 | 01 | 1 | POD-01 | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ | ⬜ pending |
| 5-02-01 | 02 | 1 | POD-02 | compile | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ | ⬜ pending |
| 5-03-01 | 03 | 1 | POD-03/POD-04 | compile | `go vet ./lib/... && cd app && ./node_modules/.bin/tsc --noEmit` | ✅ | ⬜ pending |
| 5-04-01 | 04 | 2 | DECK-01/DECK-03 | compile | `go vet ./lib/... && cd app && ./node_modules/.bin/tsc --noEmit` | ✅ | ⬜ pending |
| 5-05-01 | 05 | 2 | DECK-02 | manual | n/a | ✅ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

*Existing infrastructure covers all phase requirements.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Commander disambiguation tooltip visible on deck update | DECK-02 | Visual UI element, already implemented per research | Open deck settings, hover commander field, verify tooltip appears |
| Empty state prompt visible for new user with no pods | POD-02 | Requires test account with no pod membership | Log in as new user, visit /, verify "Create or join a pod" prompt is shown |
| Create new deck dialog opens from UI | DECK-01 | Requires live browser session | Navigate to pod view, open deck tab, click "New Deck", verify dialog opens |
| Create new pod flow in PodSelector | POD-01 | Requires live browser session | From home/pod page, select "Create new pod" from selector, verify dialog and navigation |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 15s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
