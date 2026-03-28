---
phase: quick
plan: 260326-x0y
subsystem: frontend
tags: [ux, pod-selector, react, state-management]
dependency_graph:
  requires: []
  provides: [pod-selector-refresh]
  affects: [app/src/routes/root.tsx]
tech_stack:
  added: []
  patterns: [re-fetch-after-create, navigation-driven-refresh]
key_files:
  modified:
    - app/src/routes/root.tsx
decisions:
  - "Re-fetch pods inside handleCreatePod (await GetPodsForPlayer after PostPod) so PodSelector state updates before navigation"
  - "Second useEffect with [podId, playerId] deps + pods.find guard handles HomeView creation path where PodSelector is never notified of the new pod"
metrics:
  duration: "5min"
  completed: "2026-03-26"
  tasks: 1
  files: 1
---

# Quick Task 260326-x0y: Fix AppBar PodSelector Not Updating After Pod Creation

## One-liner

Re-fetch pods list in PodSelector after creation (dialog path) and on navigation to unknown pod (HomeView path) so dropdown stays current without page refresh.

## What Was Done

### Task 1: Re-fetch pods after PodSelector dialog creation and on navigation to unknown pod

Two changes to `PodSelector` in `app/src/routes/root.tsx`:

**Change 1 — handleCreatePod re-fetch**

Added `await GetPodsForPlayer(playerId).then(setPods)` immediately after `PostPod` succeeds, before closing the dialog and navigating. This ensures the new pod is in `pods` state before the component re-renders with the new route.

**Change 2 — Navigation-triggered useEffect**

Added a second `useEffect` watching `[podId, playerId]`. When `podId` is set and the current pod is not found in `pods`, it re-fetches. This handles the case where the user creates a pod from HomeView (which uses its own `PostPod` call and navigates directly), meaning the PodSelector was never notified. The `pods.find` guard inside the effect prevents spurious fetches when the pod is already loaded.

The `eslint-disable-line react-hooks/exhaustive-deps` comment suppresses the linter warning about the intentional omission of `pods` from deps — including `pods` would cause the effect to re-trigger after the fetch resolves (loop risk). Reading `pods` inside the effect body at execution time is correct behavior here.

## Commits

| Task | Commit | Description |
|------|--------|-------------|
| 1    | 1316e85 | feat(260326-x0y): re-fetch pods after create and on navigation to unknown pod |

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None.

## Self-Check

- [x] `app/src/routes/root.tsx` modified with both changes
- [x] Commit `1316e85` exists
- [x] TypeScript compile passed (no output, exit 0)
