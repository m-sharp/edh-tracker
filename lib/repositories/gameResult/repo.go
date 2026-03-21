package gameResult

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getStatsForPlayer = `SELECT game_result.game_id, game_result.place, game_result.kill_count,
						        (SELECT COUNT(*) FROM game_result gr2
						          WHERE gr2.game_id = game_result.game_id
						            AND gr2.deleted_at IS NULL) AS player_count
						   FROM game_result INNER JOIN deck ON game_result.deck_id = deck.id
						  WHERE deck.player_id = ?
						    AND deck.deleted_at IS NULL
						    AND game_result.deleted_at IS NULL;`

	getStatsForDeck = `SELECT game_result.game_id, game_result.place, game_result.kill_count,
						      (SELECT COUNT(*) FROM game_result gr2
						        WHERE gr2.game_id = game_result.game_id
						          AND gr2.deleted_at IS NULL) AS player_count
						 FROM game_result INNER JOIN deck ON game_result.deck_id = deck.id
						WHERE deck.id = ? AND game_result.deleted_at IS NULL;`
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

func (r *Repository) GetByGameIDWithDeckInfo(ctx context.Context, gameID int) ([]Model, error) {
	var results []Model
	err := r.db.WithContext(ctx).
		Preload("Deck.Commander.Commander").
		Preload("Deck.Commander.PartnerCommander").
		Preload("Deck.Player").
		Where("game_id = ?", gameID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get GameResults with deck info for game %d: %w", gameID, err)
	}
	return results, nil
}

func (r *Repository) GetByGameId(ctx context.Context, gameID int) ([]Model, error) {
	var results []Model
	if err := r.db.WithContext(ctx).Where("game_id = ?", gameID).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get GameResult records for game %d: %w", gameID, err)
	}
	if results == nil {
		return []Model{}, nil
	}
	return results, nil
}

func (r *Repository) GetByID(ctx context.Context, resultID int) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).First(&m, resultID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get GameResult record for id %d: %w", resultID, err)
	}
	return &m, nil
}

func (r *Repository) Add(ctx context.Context, m Model) (int, error) {
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert GameResult record: %w", err)
	}
	return m.ID, nil
}

func (r *Repository) BulkAdd(ctx context.Context, results []Model) error {
	if len(results) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(&results, 100).Error; err != nil {
		return fmt.Errorf("failed to bulk insert GameResult records: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, resultID, place, killCount, deckID int) error {
	result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", resultID).Updates(Model{
		Place:     place,
		KillCount: killCount,
		DeckID:    deckID,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update GameResult record: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by GameResult update: got %d, expected 1", result.RowsAffected)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&Model{}, id).Error; err != nil {
		return fmt.Errorf("failed to soft-delete GameResult record: %w", err)
	}
	return nil
}

func (r *Repository) GetStatsForPlayer(ctx context.Context, playerID int) (*Aggregate, error) {
	var stats gameStats
	if err := r.db.WithContext(ctx).Raw(getStatsForPlayer, playerID).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get stats for player %d: %w", playerID, err)
	}
	agg := stats.toAggregate()
	return &agg, nil
}

func (r *Repository) GetStatsForDeck(ctx context.Context, deckID int) (*Aggregate, error) {
	var stats gameStats
	if err := r.db.WithContext(ctx).Raw(getStatsForDeck, deckID).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get stats for deck %d: %w", deckID, err)
	}
	agg := stats.toAggregate()
	return &agg, nil
}
