package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetGameResultsByGameID = `SELECT game_result.id, game_result.game_id, game_result.deck_id,
								deck.name    AS deck_name,
								c.name       AS commander_name,
								pc.name      AS partner_commander_name,
								game_result.place, game_result.kill_count,
								game_result.created_at, game_result.updated_at, game_result.deleted_at
							  FROM game_result
							  INNER JOIN deck ON game_result.deck_id = deck.id
							  LEFT JOIN deck_commander dc ON dc.deck_id = deck.id AND dc.deleted_at IS NULL
							  LEFT JOIN commander c  ON dc.commander_id         = c.id  AND c.deleted_at  IS NULL
							  LEFT JOIN commander pc ON dc.partner_commander_id = pc.id AND pc.deleted_at IS NULL
							  WHERE game_result.game_id = ? AND game_result.deleted_at IS NULL;`

	InsertGameResult     = `INSERT INTO game_result (game_id, deck_id, place, kill_count) VALUES (?, ?, ?, ?);`
	SoftDeleteGameResult = `UPDATE game_result SET deleted_at = NOW() WHERE id = ?;`

	gameResultValidationErr = "invalid Game Result: %s"
)

type GameResult struct {
	Model
	GameId               int     `json:"game_id"    db:"game_id"`
	DeckId               int     `json:"deck_id"    db:"deck_id"`
	DeckName             string  `json:"deck_name"  db:"deck_name"`
	CommanderName        *string `json:"commander_name,omitempty"         db:"commander_name"`
	PartnerCommanderName *string `json:"partner_commander_name,omitempty" db:"partner_commander_name"`
	Place                int     `json:"place"      db:"place"`
	Kills                int     `json:"kill_count" db:"kill_count"`
	Points               int     `json:"points"`
}

func (g *GameResult) Validate() error {
	if g.DeckId == 0 {
		return fmt.Errorf(gameResultValidationErr, "missing DeckId")
	}
	if g.Place == 0 {
		return fmt.Errorf(gameResultValidationErr, "missing Place")
	}
	if g.Place < 1 {
		return fmt.Errorf(gameResultValidationErr, "Place cannot be less than 1")
	}
	if g.Kills < 0 {
		return fmt.Errorf(gameResultValidationErr, "Kills cannot be less than 0")
	}

	return nil
}

type GameResultRepository struct {
	client *lib.DBClient
}

func NewGameResultRepository(client *lib.DBClient) *GameResultRepository {
	return &GameResultRepository{client: client}
}

func (gr *GameResultRepository) GetByGameId(ctx context.Context, gameId int) ([]GameResult, error) {
	var results []GameResult
	if err := gr.client.Db.SelectContext(ctx, &results, GetGameResultsByGameID, gameId); err != nil {
		return nil, fmt.Errorf("failed to get Game Results for Game %d: %w", gameId, err)
	}

	if results == nil {
		return nil, fmt.Errorf("failed to get Game Results for Game %d: no results found", gameId)
	}

	for i := range results {
		results[i].Points = getPointsForPlace(results[i].Kills, results[i].Place)
	}

	return results, nil
}

func (gr *GameResultRepository) BulkAdd(ctx context.Context, results []GameResult) error {
	if len(results) == 0 {
		return nil
	}

	query := "INSERT INTO game_result (game_id, deck_id, place, kill_count) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?,?,?),", len(results)), ",")
	args := make([]interface{}, 0, len(results)*4)
	for _, result := range results {
		args = append(args, result.GameId, result.DeckId, result.Place, result.Kills)
	}
	if _, err := gr.client.Db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk insert GameResult records: %w", err)
	}

	return nil
}

func (gr *GameResultRepository) SoftDelete(ctx context.Context, id int) error {
	result, err := gr.client.Db.ExecContext(ctx, SoftDeleteGameResult, id)
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
