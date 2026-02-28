package models

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetGameByID      = `SELECT id, description, pod_id, format_id, created_at, updated_at, deleted_at FROM game WHERE id = ? AND deleted_at IS NULL;`
	GetGamesByDeckId = `SELECT game.id, game.description, game.pod_id, game.format_id, game.created_at, game.updated_at, game.deleted_at
							FROM (game INNER JOIN game_result on game.id = game_result.game_id)
						  WHERE game_result.deck_id = ?
						    AND game.deleted_at IS NULL
						    AND game_result.deleted_at IS NULL;`
	GetGamesByPodId = `SELECT id, description, pod_id, format_id, created_at, updated_at, deleted_at FROM game WHERE pod_id = ? AND deleted_at IS NULL;`

	InsertGame     = `INSERT INTO game (description, pod_id, format_id) VALUES (?, ?, ?);`
	SoftDeleteGame = `UPDATE game SET deleted_at = NOW() WHERE id = ?;`
)

type Game struct {
	Model
	Description string `json:"description" db:"description"`
	PodID       int    `json:"pod_id"      db:"pod_id"`
	FormatID    int    `json:"format_id"   db:"format_id"`
}

type GameDetails struct {
	Game
	Results []GameResult `json:"results"`
}

type GameProvider struct {
	log         *zap.Logger
	client      *lib.DBClient
	gameResults *GameResultProvider
}

func NewGameProvider(log *zap.Logger, client *lib.DBClient, gameResults *GameResultProvider) *GameProvider {
	return &GameProvider{
		log:         log.Named("GameProvider"),
		client:      client,
		gameResults: gameResults,
	}
}

func (g *GameProvider) GetAllByPod(ctx context.Context, podId int) ([]GameDetails, error) {
	var games []Game
	if err := g.client.Db.SelectContext(ctx, &games, GetGamesByPodId, podId); err != nil {
		return nil, fmt.Errorf("failed to get Game records: %w", err)
	}

	if games == nil {
		return []GameDetails{}, nil
	}

	var details []GameDetails
	for _, game := range games {
		results, err := g.gameResults.GetByGameId(ctx, game.ID)
		if err != nil {
			g.log.Warn("Failed to get game results for game, dropping from results", zap.Any("Game", game))
			continue
		}

		details = append(details, GameDetails{Game: game, Results: results})
	}

	return details, nil
}

func (g *GameProvider) GetAllByDeck(ctx context.Context, deckId int) ([]GameDetails, error) {
	var games []Game
	if err := g.client.Db.SelectContext(ctx, &games, GetGamesByDeckId, deckId); err != nil {
		return nil, fmt.Errorf("failed to get Game records: %w", err)
	}

	if games == nil {
		return []GameDetails{}, nil
	}

	var details []GameDetails
	for _, game := range games {
		results, err := g.gameResults.GetByGameId(ctx, game.ID)
		if err != nil {
			g.log.Warn("Failed to get game results for game, dropping from results", zap.Any("Game", game))
			continue
		}

		details = append(details, GameDetails{Game: game, Results: results})
	}

	return details, nil
}

func (g *GameProvider) GetGameById(ctx context.Context, gameId int) (*GameDetails, error) {
	var games []Game
	if err := g.client.Db.SelectContext(ctx, &games, GetGameByID, gameId); err != nil {
		return nil, fmt.Errorf("failed to get Game record for id %d: %w", gameId, err)
	}

	if len(games) == 0 || len(games) > 1 {
		return nil, fmt.Errorf(
			"unexpected number of games returned for ID %d: got %d, expected 1",
			gameId, len(games),
		)
	}

	game := games[0]
	results, err := g.gameResults.GetByGameId(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	return &GameDetails{
		Game:    game,
		Results: results,
	}, nil
}

func (g *GameProvider) Add(ctx context.Context, description string, podID int, formatID int, results ...GameResult) error {
	r, err := g.client.Db.ExecContext(ctx, InsertGame, description, podID, formatID)
	if err != nil {
		return fmt.Errorf("failed to insert Game record: %w", err)
	}

	numAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Game insert: got %d, expected 1", numAffected)
	}

	newId, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	for i := range results {
		results[i].GameId = int(newId)
	}
	return g.gameResults.BulkAdd(ctx, results)
}

func (g *GameProvider) BulkAdd(ctx context.Context, games []GameDetails) error {
	if len(games) == 0 {
		return nil
	}

	// Phase A: bulk insert all games
	gameQuery := "INSERT INTO game (description, pod_id, format_id) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?,?),", len(games)), ",")
	gameArgs := make([]interface{}, 0, len(games)*3)
	for _, game := range games {
		gameArgs = append(gameArgs, game.Description, game.PodID, game.FormatID)
	}
	r, err := g.client.Db.ExecContext(ctx, gameQuery, gameArgs...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert Game records: %w", err)
	}
	firstGameID, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID for bulk Game insert: %w", err)
	}

	// Phase B: flatten all results with sequential game IDs
	var allResults []GameResult
	for i, game := range games {
		gameID := int(firstGameID) + i
		for _, result := range game.Results {
			result.GameId = gameID
			allResults = append(allResults, result)
		}
	}

	// Phase C: bulk insert all results via GameResultProvider
	return g.gameResults.BulkAdd(ctx, allResults)
}

// TODO: Soft deleting a game should also delete all associated GameResult records
// TODO: Will need to look for other cascading deletes
func (g *GameProvider) SoftDelete(ctx context.Context, id int) error {
	result, err := g.client.Db.ExecContext(ctx, SoftDeleteGame, id)
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
