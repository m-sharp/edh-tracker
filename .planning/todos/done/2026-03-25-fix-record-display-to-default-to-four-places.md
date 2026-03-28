---
created: 2026-03-25T01:37:54.447Z
title: Fix record display to default to four places
area: ui
files:
  - app/src/routes/pod/index.tsx
---

## Problem

The dynamic W/L/D-style record component is too aggressive in collapsing — if a player has only played games where they finished in the same place (e.g., always 1st), it renders just a single number instead of showing the full breakdown. This looks broken and loses context.

## Solution

Default to always showing four place columns (1st, 2nd, 3rd, 4th). Only add additional columns if data contains finishes beyond 4th place. This mirrors how a standard sports record is displayed — the shape is fixed, the numbers fill in.
