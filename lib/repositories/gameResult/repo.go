package gameResult

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
)

const (
	// TODO: Better way to do deck and player stats? Is this too burdensome here? Should we have views?
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

	getStatsForDecks = `SELECT deck.id AS deck_id, game_result.game_id, game_result.place, game_result.kill_count,
						       (SELECT COUNT(*) FROM game_result gr2
						         WHERE gr2.game_id = game_result.game_id
						           AND gr2.deleted_at IS NULL) AS player_count
						  FROM game_result
						  INNER JOIN deck ON game_result.deck_id = deck.id
						 WHERE deck.id IN ?
						   AND game_result.deleted_at IS NULL`

	getStatsForPlayersInPod = `SELECT deck.player_id, game_result.game_id, game_result.place, game_result.kill_count,
					            (SELECT COUNT(*) FROM game_result gr2
					              WHERE gr2.game_id = game_result.game_id
					                AND gr2.deleted_at IS NULL) AS player_count
					       FROM game_result
					       INNER JOIN deck ON game_result.deck_id = deck.id
					       INNER JOIN game ON game_result.game_id = game.id
					      WHERE game.pod_id = ?
					        AND deck.player_id IN ?
					        AND game.deleted_at IS NULL
					        AND deck.deleted_at IS NULL
					        AND game_result.deleted_at IS NULL`
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

func (r *Repository) GetByGameID(ctx context.Context, gameID int) ([]Model, error) {
	var results []Model
	err := r.DB().WithContext(ctx).
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

func (r *Repository) GetByID(ctx context.Context, resultID int) (*Model, error) {
	var m Model
	err := r.DB().WithContext(ctx).First(&m, resultID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get GameResult record for id %d: %w", resultID, err)
	}
	return &m, nil
}

func (r *Repository) Add(ctx context.Context, m Model) (int, error) {
	if err := r.DB().WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert GameResult record: %w", err)
	}
	return m.ID, nil
}

func (r *Repository) BulkAdd(ctx context.Context, results []Model) error {
	if len(results) == 0 {
		return nil
	}
	if err := r.DB().WithContext(ctx).CreateInBatches(&results, 100).Error; err != nil {
		return fmt.Errorf("failed to bulk insert GameResult records: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, resultID, place, killCount, deckID int) error {
	result := r.DB().WithContext(ctx).Model(&Model{}).Where("id = ?", resultID).Updates(Model{
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
	if err := r.DB().WithContext(ctx).Delete(&Model{}, id).Error; err != nil {
		return fmt.Errorf("failed to soft-delete GameResult record: %w", err)
	}
	return nil
}

func (r *Repository) GetStatsForPlayer(ctx context.Context, playerID int) (*Aggregate, error) {
	var stats gameStats
	if err := r.DB().WithContext(ctx).Raw(getStatsForPlayer, playerID).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get stats for player %d: %w", playerID, err)
	}
	agg := stats.toAggregate()
	return &agg, nil
}

func (r *Repository) GetStatsForDeck(ctx context.Context, deckID int) (*Aggregate, error) {
	var stats gameStats
	if err := r.DB().WithContext(ctx).Raw(getStatsForDeck, deckID).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get stats for deck %d: %w", deckID, err)
	}
	agg := stats.toAggregate()
	return &agg, nil
}

func (r *Repository) GetStatsForDecks(ctx context.Context, deckIDs []int) (map[int]*Aggregate, error) {
	result := make(map[int]*Aggregate, len(deckIDs))
	if len(deckIDs) == 0 {
		return result, nil
	}

	var rows []gameStatWithDeck
	if err := r.DB().WithContext(ctx).Raw(getStatsForDecks, deckIDs).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to get batch stats for decks: %w", err)
	}

	// Group by deck ID
	grouped := make(map[int]gameStats)
	for _, row := range rows {
		grouped[row.DeckID] = append(grouped[row.DeckID], gameStat{
			GameID:      row.GameID,
			Place:       row.Place,
			KillCount:   row.KillCount,
			PlayerCount: row.PlayerCount,
		})
	}

	// Convert each group to Aggregate
	for _, deckID := range deckIDs {
		if stats, ok := grouped[deckID]; ok {
			agg := stats.toAggregate()
			result[deckID] = &agg
		} else {
			// No games for this deck — return zero aggregate
			result[deckID] = &Aggregate{Record: map[int]int{}}
		}
	}

	return result, nil
}

func (r *Repository) GetStatsForPlayersInPod(ctx context.Context, podID int, playerIDs []int) (map[int]*Aggregate, error) {
	result := make(map[int]*Aggregate, len(playerIDs))
	if len(playerIDs) == 0 {
		return result, nil
	}

	var rows []gameStatWithPlayer
	if err := r.DB().WithContext(ctx).Raw(getStatsForPlayersInPod, podID, playerIDs).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to get batch stats for players in pod %d: %w", podID, err)
	}

	// Group by player ID
	grouped := make(map[int]gameStats)
	for _, row := range rows {
		grouped[row.PlayerID] = append(grouped[row.PlayerID], gameStat{
			GameID:      row.GameID,
			Place:       row.Place,
			KillCount:   row.KillCount,
			PlayerCount: row.PlayerCount,
		})
	}

	// Convert each group to Aggregate
	for _, playerID := range playerIDs {
		if stats, ok := grouped[playerID]; ok {
			agg := stats.toAggregate()
			result[playerID] = &agg
		} else {
			result[playerID] = &Aggregate{Record: map[int]int{}}
		}
	}

	return result, nil
}
