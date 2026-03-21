# Phase 4 — Preloading: pod.Model → []playerPodRole.Model

## Skill

Use the `/gorm` skill tool at the start of each implementation session for this phase.

## Goal

Add a HasMany `Members` association to `pod.Model` for convenience. There is no current
N+1 pattern for pod member loading — `GetMembersWithRoles` already issues a single query.
This phase adds an optional `GetByIDWithMembers` method for future scenarios where a pod
and its members are needed together (e.g., a pod detail page endpoint).

**Priority:** Low — additive, no existing behavior changes.

## Dependencies

None — standalone, no deps on other phases in this plan.

## Scope

- `lib/repositories/pod/model.go` — add `Members` HasMany field
- `lib/repositories/pod/repo.go` — add `GetByIDWithMembers` method
- `lib/repositories/interfaces.go` — add new `PodRepository` signature
- Business layer: optional `GetByIDWithMembers` + `EntityWithMembers` — only add if a
  concrete caller exists at implementation time

## Association Declaration

Add to `lib/repositories/pod/model.go`:

```go
import playerPodRoleRepo "github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"

type Model struct {
    base.GormModelBase
    Name string
    // Populated only when using GetByIDWithMembers.
    Members []playerPodRoleRepo.Model `gorm:"foreignKey:PodID"`
}

func (Model) TableName() string { return "pod" }
```

## New Repository Method

Add to `lib/repositories/pod/repo.go`:

```go
func (r *Repository) GetByIDWithMembers(ctx context.Context, podID int) (*Model, error) {
    var m Model
    err := r.db.WithContext(ctx).Preload("Members").First(&m, podID).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get Pod with members for id %d: %w", podID, err)
    }
    return &m, nil
}
```

## Interface Update

Add to `PodRepository` in `lib/repositories/interfaces.go`:

```go
GetByIDWithMembers(ctx context.Context, podID int) (*pod.Model, error)
```

## Business Layer Changes

This phase is primarily additive. No existing business functions need to change —
`GetMembersWithRoles` already works correctly via a direct repo call.

Optionally, if a combined pod+members fetch is ever needed, add:

```go
func GetByIDWithMembers(podRepo repos.PodRepository) GetByIDWithMembersFunc {
    return func(ctx context.Context, podID int) (*EntityWithMembers, error) {
        m, err := podRepo.GetByIDWithMembers(ctx, podID)
        if m == nil {
            return nil, err
        }
        members := make([]PlayerWithRole, 0, len(m.Members))
        for _, member := range m.Members {
            members = append(members, PlayerWithRole{PlayerID: member.PlayerID, Role: member.Role})
        }
        e := &EntityWithMembers{Entity: ToEntity(*m), Members: members}
        return e, nil
    }
}
```

Adding `GetByIDWithMembersFunc`, `EntityWithMembers`, and the `Functions` field update is
**optional** — only add if there is a concrete caller at implementation time.

## Tests

**`lib/repositories/pod/repo_test.go`:**

Add an integration test for `GetByIDWithMembers`. Seed a pod with at least two
`player_pod_role` rows and assert `m.Members` is populated with the expected entries
(correct `PlayerID` and `Role` values for each). Also seed a soft-deleted role row and
verify it is excluded from `m.Members`.

## Verification

1. `go vet ./lib/...` passes
2. `go test ./lib/repositories/pod/...` passes (or skips)
3. Existing pod endpoint smoke tests continue to pass
