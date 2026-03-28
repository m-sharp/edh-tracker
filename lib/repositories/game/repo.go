package game

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

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

func (r *Repository) preloadAll(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Results.Deck.Commander.Commander").
		Preload("Results.Deck.Commander.PartnerCommander").
		Preload("Results.Deck.Player")
}

func (r *Repository) GetAllByPod(ctx context.Context, podID int) ([]Model, error) {
	var games []Model
	err := r.preloadAll(r.DB().WithContext(ctx)).
		Where("pod_id = ?", podID).
		Find(&games).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get Game records with results for pod %d: %w", podID, err)
	}
	if games == nil {
		return []Model{}, nil
	}
	return games, nil
}

func (r *Repository) GetAllByPodPaginated(ctx context.Context, podID, limit, offset int) ([]Model, int, error) {
	var total int64
	if err := r.DB().WithContext(ctx).Model(&Model{}).Where("pod_id = ?", podID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count Game records for pod %d: %w", podID, err)
	}

	var games []Model
	err := r.preloadAll(r.DB().WithContext(ctx)).
		Where("pod_id = ?", podID).
		Limit(limit).Offset(offset).
		Find(&games).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get Game records with results for pod %d: %w", podID, err)
	}
	if games == nil {
		return []Model{}, int(total), nil
	}
	return games, int(total), nil
}

func (r *Repository) GetAllByDeck(ctx context.Context, deckID int) ([]Model, error) {
	var games []Model
	err := r.preloadAll(r.DB().WithContext(ctx)).
		Joins("INNER JOIN game_result ON game.id = game_result.game_id").
		Where("game_result.deck_id = ? AND game_result.deleted_at IS NULL", deckID).
		Find(&games).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get Game records with results for deck %d: %w", deckID, err)
	}
	if games == nil {
		return []Model{}, nil
	}
	return games, nil
}

func (r *Repository) GetAllByDeckPaginated(ctx context.Context, deckID, limit, offset int) ([]Model, int, error) {
	var total int64
	if err := r.DB().WithContext(ctx).Model(&Model{}).
		Joins("INNER JOIN game_result ON game.id = game_result.game_id").
		Where("game_result.deck_id = ? AND game_result.deleted_at IS NULL", deckID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count Game records for deck %d: %w", deckID, err)
	}

	var games []Model
	err := r.preloadAll(r.DB().WithContext(ctx)).
		Joins("INNER JOIN game_result ON game.id = game_result.game_id").
		Where("game_result.deck_id = ? AND game_result.deleted_at IS NULL", deckID).
		Limit(limit).Offset(offset).
		Find(&games).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get Game records with results for deck %d: %w", deckID, err)
	}
	if games == nil {
		return []Model{}, int(total), nil
	}
	return games, int(total), nil
}

func (r *Repository) GetAllByPlayerID(ctx context.Context, playerID int) ([]Model, error) {
	var games []Model
	err := r.preloadAll(r.DB().WithContext(ctx)).
		Select("game.*").
		Joins("INNER JOIN game_result ON game.id = game_result.game_id").
		Joins("INNER JOIN deck ON game_result.deck_id = deck.id").
		Where("deck.player_id = ? AND game_result.deleted_at IS NULL AND deck.deleted_at IS NULL", playerID).
		Distinct().
		Find(&games).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get Game records with results for player %d: %w", playerID, err)
	}
	if games == nil {
		return []Model{}, nil
	}
	return games, nil
}

func (r *Repository) GetAllByPlayerIDPaginated(ctx context.Context, playerID, limit, offset int) ([]Model, int, error) {
	var total int64
	if err := r.DB().WithContext(ctx).Model(&Model{}).
		Select("game.*").
		Joins("INNER JOIN game_result ON game.id = game_result.game_id").
		Joins("INNER JOIN deck ON game_result.deck_id = deck.id").
		Where("deck.player_id = ? AND game_result.deleted_at IS NULL AND deck.deleted_at IS NULL", playerID).
		Distinct().
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count Game records for player %d: %w", playerID, err)
	}

	var games []Model
	err := r.preloadAll(r.DB().WithContext(ctx)).
		Select("game.*").
		Joins("INNER JOIN game_result ON game.id = game_result.game_id").
		Joins("INNER JOIN deck ON game_result.deck_id = deck.id").
		Where("deck.player_id = ? AND game_result.deleted_at IS NULL AND deck.deleted_at IS NULL", playerID).
		Distinct().
		Limit(limit).Offset(offset).
		Find(&games).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get Game records with results for player %d: %w", playerID, err)
	}
	if games == nil {
		return []Model{}, int(total), nil
	}
	return games, int(total), nil
}

func (r *Repository) GetByID(ctx context.Context, gameID int) (*Model, error) {
	var m Model
	err := r.preloadAll(r.DB().WithContext(ctx)).First(&m, gameID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Game record with results for id %d: %w", gameID, err)
	}
	return &m, nil
}

func (r *Repository) Add(ctx context.Context, description string, podID, formatID int) (int, error) {
	m := Model{Description: description, PodID: podID, FormatID: formatID}
	if err := r.DB().WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert Game record: %w", err)
	}
	return m.ID, nil
}

func (r *Repository) BulkAdd(ctx context.Context, games []Model) ([]int, error) {
	if len(games) == 0 {
		return []int{}, nil
	}
	if err := r.DB().WithContext(ctx).CreateInBatches(&games, 100).Error; err != nil {
		return nil, fmt.Errorf("failed to bulk insert Game records: %w", err)
	}
	ids := make([]int, len(games))
	for i, g := range games {
		ids[i] = g.ID
	}
	return ids, nil
}

func (r *Repository) Update(ctx context.Context, gameID int, description string) error {
	result := r.DB().WithContext(ctx).Model(&Model{}).Where("id = ?", gameID).Update("description", description)
	if result.Error != nil {
		return fmt.Errorf("failed to update Game record: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Game update: got %d, expected 1", result.RowsAffected)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	m := Model{}
	m.ID = id
	if err := r.DB().WithContext(ctx).Delete(&m).Error; err != nil {
		return fmt.Errorf("failed to soft-delete Game record: %w", err)
	}
	return nil
}
