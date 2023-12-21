package models

import (
	"context"
	"fmt"
	"time"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetAllPlayers = `SELECT id, name, ctime FROM player;`
	GetPlayerByID = `SELECT id, name, ctime FROM player WHERE id = ?;`
	InsertPlayer  = `INSERT INTO player (name) VALUES (?);`

	playerValidationErr = "invalid Player: %s"
)

type Player struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"ctime" db:"ctime"`
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

func (p *PlayerProvider) GetAll(ctx context.Context) ([]Player, error) {
	var players []Player
	if err := p.client.Db.SelectContext(ctx, &players, GetAllPlayers); err != nil {
		return nil, fmt.Errorf("failed to get Player records: %w", err)
	}

	// Return an empty list instead of nil
	if players == nil {
		return []Player{}, nil
	}

	return players, nil
}

func (p *PlayerProvider) GetById(ctx context.Context, playerId int) (*Player, error) {
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

	return &players[0], nil
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
