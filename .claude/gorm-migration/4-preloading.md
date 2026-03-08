# Phase 4 — GORM Preloading & Associations

## Status
Superseded

This plan has been broken into targeted sub-phases. See:
- 4a-gameResults.md
- 4b-deckCommander.md
- 4c-deckBelongsTo.md
- 4d-gameResultDeck.md
- 4e-podMembers.md

## Goal

Leverage GORM's `Preload` and association declarations to eliminate N+1 query patterns in
the repository layer. Currently, the business layer orchestrates separate repository calls
to assemble enriched Entities — some of these can be pushed down into the repository layer
as preloaded associations, reducing round-trips and simplifying business logic.

## Background

- No inter-repository imports currently exist at the repository layer — all packages only
  import `lib`, `base`, and the standard library. There is therefore **no circular import risk**
  when adding association fields that reference other repository model types.
- The primary N+1 pattern today: fetching games for a pod, then fetching game_results
  per-game in a loop inside the business layer.
- GORM supports HasMany, BelongsTo, HasOne associations on model structs, with
  `Preload("AssociationField")` to eagerly load them in a second batched query.

## Candidate Associations (To Be Reviewed)

| Parent Model        | Association Field                    | Type     | FK              | Benefit |
|---------------------|--------------------------------------|----------|-----------------|---------|
| `game.Model`        | `Results []gameResult.Model`         | HasMany  | `game_id`       | Eliminates per-game gameResult loop |
| `deck.Model`        | `Commander *deckCommander.Model`     | HasOne   | `deck_id`       | Single fetch for deck + commander |
| `pod.Model`         | `Players []playerPodRole.Model`      | HasMany  | `pod_id`        | Members + roles in one call |

Additional candidates should be evaluated during review:
- `deck.Model` → `Player player.Model` (BelongsTo via `player_id`)
- `deck.Model` → `Format format.Model` (BelongsTo via `format_id`)

## Approach

1. **Additive — new hydrated methods alongside flat ones.** Do not replace existing flat
   repository methods. Instead, add new methods (e.g., `GetAllByPodWithResults`) that use
   `Preload` and return models with association fields populated.
2. **Association fields on Model structs.** Add the association fields with appropriate
   `gorm:"foreignKey:..."` tags. Existing code using the flat model is unaffected since
   the new fields default to nil/empty.
3. **Business layer simplification.** Where a hydrated repository method exists, the
   business layer can use the preloaded data instead of making a separate repo call. This
   may allow removing some dependency closures from the business `Functions` structs.
4. **Interface updates.** New repository interface methods added to `lib/repositories/interfaces.go`.

## Constraints / Open Questions

- Which associations are worth adding? Evaluate by query frequency and complexity of
  current business layer orchestration.
- Should association fields live on the existing `Model` struct, or on a separate
  `HydratedModel` / `ModelWithResults` type to keep the flat model free of GORM association
  complexity?
- `GetMembersWithRoles` in the pod business layer currently joins `player_pod_role` manually.
  A preloaded association could replace this, but the current implementation also enriches with
  player name — review how far down the preloading should go.
- Phase 4 depends on all Phase 1–3 repos being GORM-migrated first (association types must
  use GORM models, not sqlx models).

## Out of Scope

- Business layer Entity construction changes beyond what is necessary to use preloaded data.
- AutoMigrate or schema changes — all schemas are already correct from existing migrations.
- Cross-layer associations (e.g., preloading from the business layer downward).
