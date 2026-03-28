package deck

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
		Preload("Commander.Commander").
		Preload("Commander.PartnerCommander").
		Preload("Player").
		Preload("Format")
}

func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
	var decks []Model
	if err := r.preloadAll(r.DB().WithContext(ctx)).Where("retired = ?", false).Find(&decks).Error; err != nil {
		return nil, fmt.Errorf("failed to get Deck records with associations: %w", err)
	}
	if decks == nil {
		return []Model{}, nil
	}
	return decks, nil
}

func (r *Repository) GetAllForPlayer(ctx context.Context, playerID int) ([]Model, error) {
	var decks []Model
	if err := r.preloadAll(r.DB().WithContext(ctx)).Where("player_id = ?", playerID).Find(&decks).Error; err != nil {
		return nil, fmt.Errorf("failed to get Deck records for player %d with associations: %w", playerID, err)
	}
	if decks == nil {
		return []Model{}, nil
	}
	return decks, nil
}

func (r *Repository) GetAllByPlayerPaginated(ctx context.Context, playerID, limit, offset int) ([]Model, int, error) {
	var total int64
	if err := r.DB().WithContext(ctx).Model(&Model{}).
		Where("player_id = ?", playerID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count Deck records for player %d: %w", playerID, err)
	}

	var decks []Model
	err := r.preloadAll(r.DB().WithContext(ctx)).
		Where("player_id = ?", playerID).
		Limit(limit).Offset(offset).
		Find(&decks).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get Deck records for player %d: %w", playerID, err)
	}
	if decks == nil {
		return []Model{}, int(total), nil
	}
	return decks, int(total), nil
}

func (r *Repository) GetAllByPodPaginated(ctx context.Context, podID, limit, offset int) ([]Model, int, error) {
	var total int64
	if err := r.DB().WithContext(ctx).Model(&Model{}).
		Joins("INNER JOIN player_pod_role ppr ON deck.player_id = ppr.player_id AND ppr.pod_id = ? AND ppr.deleted_at IS NULL", podID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count Deck records for pod %d: %w", podID, err)
	}

	var decks []Model
	err := r.preloadAll(r.DB().WithContext(ctx)).
		Joins("INNER JOIN player_pod_role ppr ON deck.player_id = ppr.player_id AND ppr.pod_id = ? AND ppr.deleted_at IS NULL", podID).
		Limit(limit).Offset(offset).
		Find(&decks).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get Deck records for pod %d: %w", podID, err)
	}
	if decks == nil {
		return []Model{}, int(total), nil
	}
	return decks, int(total), nil
}

func (r *Repository) GetById(ctx context.Context, deckID int) (*Model, error) {
	var m Model
	err := r.DB().WithContext(ctx).First(&m, deckID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get Deck record for id %d: %w", deckID, err)
	}
	return &m, nil
}

func (r *Repository) GetByIDHydrated(ctx context.Context, deckID int) (*Model, error) {
	var m Model
	err := r.preloadAll(r.DB().WithContext(ctx)).First(&m, deckID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get Deck record for id %d with associations: %w", deckID, err)
	}
	return &m, nil
}

func (r *Repository) Add(ctx context.Context, playerID int, name string, formatID int) (int, error) {
	m := Model{PlayerID: playerID, Name: name, FormatID: formatID}
	if err := r.DB().WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert Deck record: %w", err)
	}
	return m.ID, nil
}

func (r *Repository) BulkAdd(ctx context.Context, decks []Model) ([]Model, error) {
	if len(decks) == 0 {
		return []Model{}, nil
	}
	if err := r.DB().WithContext(ctx).CreateInBatches(&decks, 100).Error; err != nil {
		return nil, fmt.Errorf("failed to bulk insert Deck records: %w", err)
	}
	return decks, nil
}

func (r *Repository) GetAllByPlayerIDs(ctx context.Context, playerIDs []int) ([]Model, error) {
	if len(playerIDs) == 0 {
		return []Model{}, nil
	}
	var decks []Model
	if err := r.preloadAll(r.DB().WithContext(ctx)).Where("player_id IN ?", playerIDs).Find(&decks).Error; err != nil {
		return nil, fmt.Errorf("failed to get Deck records for player IDs with associations: %w", err)
	}
	if decks == nil {
		return []Model{}, nil
	}
	return decks, nil
}

func (r *Repository) Update(ctx context.Context, deckID int, fields UpdateFields) error {
	if !fields.HasChanges() {
		return nil
	}

	updated := Model{}
	if fields.Name != nil {
		updated.Name = *fields.Name
	}
	if fields.FormatID != nil {
		updated.FormatID = *fields.FormatID
	}
	if fields.Retired != nil {
		updated.Retired = *fields.Retired
	}

	result := r.DB().WithContext(ctx).Model(&Model{}).Where("id = ?", deckID).Updates(updated)
	if result.Error != nil {
		return fmt.Errorf("failed to update Deck record: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("unexpected rows affected by Deck update: got %d, expected 1", result.RowsAffected)
	}
	return nil
}

func (r *Repository) Retire(ctx context.Context, deckID int) error {
	result := r.DB().WithContext(ctx).Model(&Model{}).Where("id = ?", deckID).Update("retired", true)
	if result.Error != nil {
		return fmt.Errorf("failed to retire Deck: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("unexpected rows affected by Deck retirement: got %d, expected 1", result.RowsAffected)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	m := Model{}
	m.ID = id
	if err := r.DB().WithContext(ctx).Delete(&m).Error; err != nil {
		return fmt.Errorf("failed to soft-delete Deck record: %w", err)
	}
	return nil
}
