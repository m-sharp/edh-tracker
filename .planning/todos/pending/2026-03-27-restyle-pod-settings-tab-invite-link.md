---
created: 2026-03-27T02:23:13.195Z
title: Restyle pod settings tab invite link
area: ui
files:
  - app/src/routes/pod.tsx
---

## Problem

The pod settings tab displays the generated invite link as plain text, which looks ugly and out of place. It needs to be restyled to look polished — ideally as a copyable link component or a styled card/chip similar to how other settings are displayed.

## Solution

Restyle the invite link section in the pod settings tab. Consider using a MUI `TextField` with a copy-to-clipboard button, or a styled box that makes the link visually distinct and easy to share. Align with the restyle work being done on deck and player settings tabs.
