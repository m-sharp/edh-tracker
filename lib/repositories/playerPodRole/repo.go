package playerPodRole

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
)

type Repository struct {
	*base.Repo
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{Repo: base.NewRepo(client.GormDb)}
}

func NewRepositoryFromDB(db *gorm.DB) *Repository {
	return &Repository{Repo: base.NewRepo(db)}
}

func (r *Repository) GetRole(ctx context.Context, podID, playerID int) (*Model, error) {
	var m Model
	err := r.DB().WithContext(ctx).Where("pod_id = ? AND player_id = ?", podID, playerID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get role for player %d in pod %d: %w", playerID, podID, err)
	}
	return &m, nil
}

func (r *Repository) SetRole(ctx context.Context, podID, playerID int, role string) error {
	m := Model{PodID: podID, PlayerID: playerID, Role: role}
	err := r.DB().WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"role":       role,
			"deleted_at": nil,
		}),
	}).Create(&m).Error
	if err != nil {
		return fmt.Errorf("failed to set role %q for player %d in pod %d: %w", role, playerID, podID, err)
	}
	return nil
}

func (r *Repository) GetMembersWithRoles(ctx context.Context, podID int) ([]Model, error) {
	var rows []Model
	if err := r.DB().WithContext(ctx).Where("pod_id = ?", podID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to get members with roles for pod %d: %w", podID, err)
	}
	if rows == nil {
		return []Model{}, nil
	}
	return rows, nil
}

func (r *Repository) BulkAdd(ctx context.Context, podID int, playerIDs []int, role string) error {
	if len(playerIDs) == 0 {
		return nil
	}
	entries := make([]Model, len(playerIDs))
	for i, id := range playerIDs {
		entries[i] = Model{PodID: podID, PlayerID: id, Role: role}
	}
	if err := r.DB().WithContext(ctx).CreateInBatches(&entries, 100).Error; err != nil {
		return fmt.Errorf("failed to bulk insert player_pod_role records: %w", err)
	}
	return nil
}
