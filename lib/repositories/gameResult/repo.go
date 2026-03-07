package gameResult

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getGameResultsByGameID = `SELECT id, game_id, deck_id, place, kill_count, created_at, updated_at, deleted_at
								FROM game_result WHERE game_id = ? AND deleted_at IS NULL;`
	getGameResultByID = `SELECT id, game_id, deck_id, place, kill_count, created_at, updated_at, deleted_at
						   FROM game_result WHERE id = ? AND deleted_at IS NULL;`
	insertGameResult = `INSERT INTO game_result (game_id, deck_id, place, kill_count) VALUES (?, ?, ?, ?);`
	updateGameResult = `UPDATE game_result SET place = ?, kill_count = ?, deck_id = ? WHERE id = ? AND deleted_at IS NULL;`

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

	softDeleteGameResult = `UPDATE game_result SET deleted_at = NOW() WHERE id = ?;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetByID(ctx context.Context, resultID int) (*Model, error) {
	var results []Model
	if err := r.client.Db.SelectContext(ctx, &results, getGameResultByID, resultID); err != nil {
		return nil, fmt.Errorf("failed to get GameResult record for id %d: %w", resultID, err)
	}
	if len(results) == 0 {
		return nil, nil
	}
	return &results[0], nil
}

func (r *Repository) Add(ctx context.Context, m Model) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertGameResult, m.GameID, m.DeckID, m.Place, m.KillCount)
	if err != nil {
		return 0, fmt.Errorf("failed to insert GameResult record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return 0, fmt.Errorf("unexpected number of rows affected by GameResult insert: got %d, expected 1", numAffected)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new GameResult: %w", err)
	}

	return int(id), nil
}

func (r *Repository) Update(ctx context.Context, resultID, place, killCount, deckID int) error {
	result, err := r.client.Db.ExecContext(ctx, updateGameResult, place, killCount, deckID, resultID)
	if err != nil {
		return fmt.Errorf("failed to update GameResult record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by update: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by GameResult update: got %d, expected 1", numAffected)
	}

	return nil
}

func (r *Repository) GetByGameId(ctx context.Context, gameID int) ([]Model, error) {
	var results []Model
	if err := r.client.Db.SelectContext(ctx, &results, getGameResultsByGameID, gameID); err != nil {
		return nil, fmt.Errorf("failed to get GameResult records for game %d: %w", gameID, err)
	}
	if results == nil {
		return []Model{}, nil
	}
	return results, nil
}

func (r *Repository) GetStatsForPlayer(ctx context.Context, playerID int) (*Aggregate, error) {
	var stats gameStats
	if err := r.client.Db.SelectContext(ctx, &stats, getStatsForPlayer, playerID); err != nil {
		return nil, fmt.Errorf("failed to get stats for player %d: %w", playerID, err)
	}
	agg := stats.toAggregate()
	return &agg, nil
}

func (r *Repository) GetStatsForDeck(ctx context.Context, deckID int) (*Aggregate, error) {
	var stats gameStats
	if err := r.client.Db.SelectContext(ctx, &stats, getStatsForDeck, deckID); err != nil {
		return nil, fmt.Errorf("failed to get stats for deck %d: %w", deckID, err)
	}
	agg := stats.toAggregate()
	return &agg, nil
}

func (r *Repository) BulkAdd(ctx context.Context, results []Model) error {
	if len(results) == 0 {
		return nil
	}

	query := "INSERT INTO game_result (game_id, deck_id, place, kill_count) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?,?,?),", len(results)), ",")
	args := make([]interface{}, 0, len(results)*4)
	for _, result := range results {
		args = append(args, result.GameID, result.DeckID, result.Place, result.KillCount)
	}
	if _, err := r.client.Db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk insert GameResult records: %w", err)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	result, err := r.client.Db.ExecContext(ctx, softDeleteGameResult, id)
	if err != nil {
		return fmt.Errorf("failed to soft-delete GameResult record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by soft-delete: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by GameResult soft-delete: got %d, expected 1", numAffected)
	}

	return nil
}
