---
created: 2026-03-25T01:34:30.053Z
title: Make playing cards icon clickable to go home on mobile
area: ui
files:
  - app/src/index.tsx
---

## Problem

On mobile, there's no easy way to navigate back to the home screen. The playing cards icon in the app header/nav area is decorative but not interactive — users on mobile lack a tap target to return home.

## Solution

Wrap the playing cards icon in a `Link` (React Router) or `IconButton` component that navigates to `/`. Should be visible in the top-level layout/nav so it's available on all pages.
