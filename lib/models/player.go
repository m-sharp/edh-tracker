package models

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetAllPlayers  = `SELECT id, name, created_at, updated_at, deleted_at FROM player WHERE deleted_at IS NULL;`
	GetPlayerByID  = `SELECT id, name, created_at, updated_at, deleted_at FROM player WHERE id = ? AND deleted_at IS NULL;`
	GetPlayerStats = `SELECT DISTINCT game_result.game_id, game_result.place, game_result.kill_count
						FROM (game_result INNER JOIN deck on game_result.deck_id = deck.id)
					  WHERE deck.player_id = ?
					    AND deck.deleted_at IS NULL
					    AND game_result.deleted_at IS NULL;`
	InsertPlayer     = `INSERT INTO player (name) VALUES (?);`
	SoftDeletePlayer = `UPDATE player SET deleted_at = NOW() WHERE id = ?;`

	playerValidationErr = "invalid Player: %s"
)

type Player struct {
	Model
	Name string `json:"name" db:"name"`
}

type PlayerWithStats struct {
	Player
	Stats
}

func (p *Player) Validate() error {
	if p.Name == "" {
		return fmt.Errorf(playerValidationErr, "missing Name")
	}

	return nil
}

type PlayerProvider struct {
	client *lib.DBClient
}

func NewPlayerProvider(client *lib.DBClient) *PlayerProvider {
	return &PlayerProvider{
		client: client,
	}
}

func (p *PlayerProvider) GetAll(ctx context.Context) ([]PlayerWithStats, error) {
	var players []Player
	if err := p.client.Db.SelectContext(ctx, &players, GetAllPlayers); err != nil {
		return nil, fmt.Errorf("failed to get Player records: %w", err)
	}

	// Return an empty list instead of nil
	if players == nil {
		return []PlayerWithStats{}, nil
	}

	var withStats []PlayerWithStats
	for _, player := range players {
		var gameStats GameStats
		if err := p.client.Db.SelectContext(ctx, &gameStats, GetPlayerStats, player.ID); err != nil {
			return nil, fmt.Errorf("failed to get Player statistics: %w", err)
		}

		withStats = append(withStats, PlayerWithStats{Player: player, Stats: gameStats.ToStats()})
	}

	return withStats, nil
}

func (p *PlayerProvider) GetById(ctx context.Context, playerId int) (*PlayerWithStats, error) {
	var players []Player
	if err := p.client.Db.SelectContext(ctx, &players, GetPlayerByID, playerId); err != nil {
		return nil, fmt.Errorf("failed to get Player record for id %d: %w", playerId, err)
	}

	if len(players) == 0 || len(players) > 1 {
		return nil, fmt.Errorf(
			"unexpected number of players returned for ID %d: got %d, expected 1",
			playerId, len(players),
		)
	}

	var gameStats GameStats
	if err := p.client.Db.SelectContext(ctx, &gameStats, GetPlayerStats, playerId); err != nil {
		return nil, fmt.Errorf("failed to get Player statistics: %w", err)
	}

	return &PlayerWithStats{
		Player: players[0],
		Stats:  gameStats.ToStats(),
	}, nil
}

func (p *PlayerProvider) Add(ctx context.Context, name string) error {
	result, err := p.client.Db.ExecContext(ctx, InsertPlayer, name)
	if err != nil {
		return fmt.Errorf("failed to insert Player record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Player insert: got %d, expected 1", numAffected)
	}

	return nil
}

func (p *PlayerProvider) SoftDelete(ctx context.Context, id int) error {
	result, err := p.client.Db.ExecContext(ctx, SoftDeletePlayer, id)
	if err != nil {
		return fmt.Errorf("failed to soft-delete Player record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by soft-delete: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Player soft-delete: got %d, expected 1", numAffected)
	}

	return nil
}
