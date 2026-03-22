package deckCommander

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

func (r *Repository) GetByDeckId(ctx context.Context, deckID int) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).Where("deck_id = ?", deckID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get DeckCommander record for deck %d: %w", deckID, err)
	}
	return &m, nil
}

func (r *Repository) Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
	m := Model{
		DeckID:             deckID,
		CommanderID:        commanderID,
		PartnerCommanderID: partnerCommanderID,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert DeckCommander record: %w", err)
	}
	return m.ID, nil
}

func (r *Repository) DeleteByDeckID(ctx context.Context, deckID int) error {
	if err := r.db.WithContext(ctx).Where("deck_id = ?", deckID).Delete(&Model{}).Error; err != nil {
		return fmt.Errorf("failed to soft-delete DeckCommander records for deck %d: %w", deckID, err)
	}
	return nil
}

func (r *Repository) BulkAdd(ctx context.Context, entries []Model) error {
	if len(entries) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(&entries, 100).Error; err != nil {
		return fmt.Errorf("failed to bulk insert DeckCommander records: %w", err)
	}
	return nil
}
