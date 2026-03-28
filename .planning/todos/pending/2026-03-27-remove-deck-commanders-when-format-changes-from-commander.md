---
created: 2026-03-27T02:26:11.462Z
title: Remove deck commanders when format changes from Commander
area: api
files:
  - lib/business/deck/functions.go
  - lib/repositories/deckCommander/repo.go
---

## Problem

When a deck's format is updated from "Commander" to another format (e.g., Standard, Modern, Legacy), the deck's associated commander entries remain in the database. Commanders are a Commander-format-specific concept — leaving them attached to non-Commander decks creates inconsistent data and could cause display/logic issues.

## Solution

In the deck update business logic, detect when the format is changing away from Commander format. If so, soft-delete (or hard-delete) all associated `deck_commander` entries for that deck as part of the same update operation. Likely in `lib/business/deck/functions.go` Update function, calling a `SoftDelete`/`DeleteByDeckID` method on the deckCommander repository.
