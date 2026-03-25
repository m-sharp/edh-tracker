---
phase: 4
slug: game-model-change
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-24
---

# Phase 4 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go: testify v1.11.1 + `net/http/httptest`; Frontend: TypeScript compiler |
| **Config file** | Go: none (standard `go test`); Frontend: `app/tsconfig.json` |
| **Quick run command** | `go vet ./lib/...` (backend) or `cd app && ./node_modules/.bin/tsc --noEmit` (frontend) |
| **Full suite command** | `go vet ./lib/... && go test ./lib/... && cd app && ./node_modules/.bin/tsc --noEmit` |
| **Estimated runtime** | ~15 seconds (Go) + ~5 seconds (TypeScript) |

---

## Sampling Rate

- **After every task commit:** Run `go vet ./lib/...` (after backend changes) or `./node_modules/.bin/tsc --noEmit` from `app/` (after frontend changes)
- **After every plan wave:** Run full suite: `go test ./lib/...` + `./node_modules/.bin/tsc --noEmit`
- **Before `/gsd:verify-work`:** Both suites green; manual mobile visual check at 375px
- **Max feedback latency:** ~20 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 4-01-01 | 01 | 1 | GAME-01 | unit (Go compile + router test) | `go vet ./lib/... && go test ./lib/routers/ -run TestGameRouter_AddGameResult` | ✅ game_test.go | ⬜ pending |
| 4-01-02 | 01 | 1 | GAME-01 | type-check | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ types.ts | ⬜ pending |
| 4-02-01 | 02 | 1 | GAME-04 | type-check | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ stats.tsx | ⬜ pending |
| 4-03-01 | 03 | 2 | GAME-01, GAME-02 | type-check | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ new/index.tsx | ⬜ pending |
| 4-03-02 | 03 | 2 | GAME-03 | manual visual | open browser at 375px, verify no horizontal scroll | — | ⬜ pending |
| 4-04-01 | 04 | 2 | GAME-01, GAME-02 | type-check | `cd app && ./node_modules/.bin/tsc --noEmit` | ✅ game/index.tsx | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

Existing infrastructure covers all phase requirements. No Wave 0 test file creation needed.

- `lib/routers/game_test.go` already exists and tests `AddGameResult` — will be **updated** (not created) as part of GAME-01 backend cleanup
- `app/tsconfig.json` already exists — TypeScript type checking covers all frontend changes

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Game form is visually clean at 375px with no horizontal scroll | GAME-03 | CSS/layout behavior requires visual verification | Open NewGameView in Chrome DevTools at 375px × 667px; add 3 deck cards; verify no horizontal scrollbar; confirm ❌ buttons are tappable |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 20s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
