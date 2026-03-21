package pod

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{db: client.GormDb}
}

func NewRepositoryFromDB(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
	var pods []Model
	if err := r.db.WithContext(ctx).Find(&pods).Error; err != nil {
		return nil, fmt.Errorf("failed to get Pod records: %w", err)
	}
	if pods == nil {
		return []Model{}, nil
	}
	return pods, nil
}

func (r *Repository) GetByID(ctx context.Context, podID int) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).First(&m, podID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Pod record for id %d: %w", podID, err)
	}
	return &m, nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Pod record for name %q: %w", name, err)
	}
	return &m, nil
}

// player_pod.deleted_at must be explicit — GORM only auto-filters the primary model (pod).
func (r *Repository) GetByPlayerID(ctx context.Context, playerID int) ([]Model, error) {
	var pods []Model
	err := r.db.WithContext(ctx).
		Joins("INNER JOIN player_pod ON pod.id = player_pod.pod_id").
		Where("player_pod.player_id = ? AND player_pod.deleted_at IS NULL", playerID).
		Find(&pods).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get Pod records for player %d: %w", playerID, err)
	}
	if pods == nil {
		return []Model{}, nil
	}
	return pods, nil
}

func (r *Repository) GetIDsByPlayerID(ctx context.Context, playerID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).
		Model(&PlayerPodModel{}).
		Where("player_id = ?", playerID).
		Pluck("pod_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get Pod IDs for player %d: %w", playerID, err)
	}
	if ids == nil {
		return []int{}, nil
	}
	return ids, nil
}

func (r *Repository) GetPlayerIDs(ctx context.Context, podID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).
		Model(&PlayerPodModel{}).
		Where("pod_id = ?", podID).
		Pluck("player_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get player IDs for pod %d: %w", podID, err)
	}
	if ids == nil {
		return []int{}, nil
	}
	return ids, nil
}

func (r *Repository) Add(ctx context.Context, name string) (int, error) {
	m := Model{Name: name}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert Pod record: %w", err)
	}
	return m.ID, nil
}

func (r *Repository) BulkAddPlayers(ctx context.Context, podID int, playerIDs []int) error {
	if len(playerIDs) == 0 {
		return nil
	}
	entries := make([]PlayerPodModel, len(playerIDs))
	for i, id := range playerIDs {
		entries[i] = PlayerPodModel{PodID: podID, PlayerID: id}
	}
	if err := r.db.WithContext(ctx).CreateInBatches(&entries, 100).Error; err != nil {
		return fmt.Errorf("failed to bulk insert PlayerPod records: %w", err)
	}
	return nil
}

func (r *Repository) AddPlayerToPod(ctx context.Context, podID, playerID int) error {
	m := PlayerPodModel{PodID: podID, PlayerID: playerID}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("failed to insert PlayerPod record: %w", err)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, podID int) error {
	if err := r.db.WithContext(ctx).Delete(&Model{}, podID).Error; err != nil {
		return fmt.Errorf("failed to soft delete pod %d: %w", podID, err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, podID int, name string) error {
	result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", podID).Update("name", name)
	if result.Error != nil {
		return fmt.Errorf("failed to update pod %d name: %w", podID, result.Error)
	}
	return nil
}

func (r *Repository) RemovePlayer(ctx context.Context, podID, playerID int) error {
	err := r.db.WithContext(ctx).
		Where("pod_id = ? AND player_id = ?", podID, playerID).
		Delete(&PlayerPodModel{}).Error
	if err != nil {
		return fmt.Errorf("failed to remove player %d from pod %d: %w", playerID, podID, err)
	}
	return nil
}
