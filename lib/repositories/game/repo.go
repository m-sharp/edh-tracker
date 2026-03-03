package game

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getGameByID      = `SELECT id, description, pod_id, format_id, created_at, updated_at, deleted_at FROM game WHERE id = ? AND deleted_at IS NULL;`
	getGamesByDeckId = `SELECT game.id, game.description, game.pod_id, game.format_id, game.created_at, game.updated_at, game.deleted_at
						  FROM (game INNER JOIN game_result on game.id = game_result.game_id)
						 WHERE game_result.deck_id = ?
						   AND game.deleted_at IS NULL
						   AND game_result.deleted_at IS NULL;`
	getGamesByPodId = `SELECT id, description, pod_id, format_id, created_at, updated_at, deleted_at FROM game WHERE pod_id = ? AND deleted_at IS NULL;`
	insertGame      = `INSERT INTO game (description, pod_id, format_id) VALUES (?, ?, ?);`
	softDeleteGame  = `UPDATE game SET deleted_at = NOW() WHERE id = ?;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetAllByPod(ctx context.Context, podID int) ([]Model, error) {
	var games []Model
	if err := r.client.Db.SelectContext(ctx, &games, getGamesByPodId, podID); err != nil {
		return nil, fmt.Errorf("failed to get Game records for pod %d: %w", podID, err)
	}
	if games == nil {
		return []Model{}, nil
	}
	return games, nil
}

func (r *Repository) GetAllByDeck(ctx context.Context, deckID int) ([]Model, error) {
	var games []Model
	if err := r.client.Db.SelectContext(ctx, &games, getGamesByDeckId, deckID); err != nil {
		return nil, fmt.Errorf("failed to get Game records for deck %d: %w", deckID, err)
	}
	if games == nil {
		return []Model{}, nil
	}
	return games, nil
}

func (r *Repository) GetById(ctx context.Context, gameID int) (*Model, error) {
	var games []Model
	if err := r.client.Db.SelectContext(ctx, &games, getGameByID, gameID); err != nil {
		return nil, fmt.Errorf("failed to get Game record for id %d: %w", gameID, err)
	}
	if len(games) == 0 {
		return nil, nil
	}
	return &games[0], nil
}

func (r *Repository) Add(ctx context.Context, description string, podID, formatID int) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertGame, description, podID, formatID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert Game record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return 0, fmt.Errorf("unexpected number of rows affected by Game insert: got %d, expected 1", numAffected)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new Game: %w", err)
	}

	return int(id), nil
}

func (r *Repository) BulkAdd(ctx context.Context, games []Model) ([]int, error) {
	if len(games) == 0 {
		return []int{}, nil
	}

	query := "INSERT INTO game (description, pod_id, format_id) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?,?),", len(games)), ",")
	args := make([]interface{}, 0, len(games)*3)
	for _, g := range games {
		args = append(args, g.Description, g.PodID, g.FormatID)
	}
	result, err := r.client.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk insert Game records: %w", err)
	}

	firstID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID for bulk Game insert: %w", err)
	}

	ids := make([]int, len(games))
	for i := range games {
		ids[i] = int(firstID) + i
	}
	return ids, nil
}

// TODO: Soft deleting a game should also delete all associated GameResult records
// TODO: Will need to look for other cascading deletes
func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	result, err := r.client.Db.ExecContext(ctx, softDeleteGame, id)
	if err != nil {
		return fmt.Errorf("failed to soft-delete Game record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by soft-delete: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Game soft-delete: got %d, expected 1", numAffected)
	}

	return nil
}
