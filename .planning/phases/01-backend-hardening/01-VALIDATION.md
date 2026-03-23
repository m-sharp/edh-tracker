---
phase: 1
slug: backend-hardening
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-22
---

# Phase 1 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test (testify/assert + testify/require) |
| **Config file** | none — existing infrastructure |
| **Quick run command** | `go test ./lib/...` |
| **Full suite command** | `go test ./lib/...` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./lib/...`
- **After every plan wave:** Run `go test ./lib/...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 1-01-01 | 01 | 1 | AUTH-02 | unit | `go test ./lib/routers/... -run TestGame` | ✅ | ⬜ pending |
| 1-01-02 | 01 | 1 | SEC-01 | unit | `go test ./lib/business/deck/... -run TestCreate` | ✅ | ⬜ pending |
| 1-02-01 | 02 | 1 | SEC-02, SEC-03 | unit | `go test ./lib/business/game/... -run TestCreate` | ✅ | ⬜ pending |
| 1-03-01 | 03 | 2 | SEC-04, SEC-05 | unit | `go test ./lib/business/... -run TestValidate` | ✅ | ⬜ pending |
| 1-04-01 | 04 | 2 | PERF-01, PERF-02 | unit | `go test ./lib/repositories/... -run TestGetBatch` | ✅ | ⬜ pending |
| 1-05-01 | 05 | 3 | INFRA-02 | unit | `go test ./lib/...` | ✅ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

Existing infrastructure covers all phase requirements.

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Game creation by non-pod-member returns 403 | AUTH-02 | Integration with auth middleware + pod membership check | POST /api/game with JWT for player not in pod — expect 403 |
| Deck creation uses caller's player_id (not body) | SEC-01 | Requires live auth context | POST /api/deck with player_id in body ≠ JWT player — confirm deck created with JWT player_id |
| Game creation with failing result row leaves no orphaned game | SEC-02 | Requires DB fault injection | Simulate insert failure mid-transaction — confirm no game row in DB |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 30s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
