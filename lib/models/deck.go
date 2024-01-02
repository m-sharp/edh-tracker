package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetAllDecks       = `SELECT id, player_id, commander, retired, ctime FROM deck;`
	GetDecksForPlayer = `SELECT id, player_id, commander, retired, ctime FROM deck WHERE player_id = ?;`
	GetDeckByID       = `SELECT deck.id, deck.player_id, deck.commander, deck.retired, deck.ctime, player.name
							FROM (deck INNER JOIN player on deck.player_id = player.id)
						 WHERE deck.id = ?;`
	GetDeckStats = `SELECT DISTINCT game_result.game_id, game_result.place, game_result.kill_count
						FROM (game_result INNER JOIN deck on game_result.deck_id = deck.id)
					WHERE deck.id = ?;`

	DeckExists = `SELECT COUNT(*) FROM deck where player_id = ? AND commander = ?;`

	InsertDeck = `INSERT INTO deck (player_id, commander) VALUES (?, ?);`

	RetireDeck = `UPDATE deck SET retired = TRUE WHERE id = ?;`

	deckValidationErr = "invalid Deck: %s"
)

var (
	ErrDeckExists = errors.New("a deck for the specified commander already exists")
)

type Deck struct {
	Id         int       `json:"id" db:"id"`
	PlayerId   int       `json:"player_id" db:"player_id"`
	PlayerName string    `json:"player_name,omitempty" db:"name"`
	Commander  string    `json:"commander" db:"commander"`
	Retired    bool      `json:"retired" db:"retired"`
	CreatedAt  time.Time `json:"ctime" db:"ctime"`
}

type DeckWithStats struct {
	Deck
	Stats
}

func (d *Deck) Validate() error {
	if d.PlayerId == 0 {
		return fmt.Errorf(deckValidationErr, "missing PlayerId")
	}
	if d.Commander == "" {
		return fmt.Errorf(deckValidationErr, "missing Commander")
	}

	return nil
}

type DeckProvider struct {
	client *lib.DBClient
}

func NewDeckProvider(client *lib.DBClient) *DeckProvider {
	return &DeckProvider{
		client: client,
	}
}

func (d *DeckProvider) GetAll(ctx context.Context) ([]Deck, error) {
	var decks []Deck
	if err := d.client.Db.SelectContext(ctx, &decks, GetAllDecks); err != nil {
		return nil, fmt.Errorf("failed to get Deck records: %w", err)
	}

	// Return an empty list instead of nil
	if decks == nil {
		return []Deck{}, nil
	}

	return decks, nil
}

func (d *DeckProvider) GetAllForPlayer(ctx context.Context, playerId int) ([]Deck, error) {
	var decks []Deck
	if err := d.client.Db.SelectContext(ctx, &decks, GetDecksForPlayer, playerId); err != nil {
		return nil, fmt.Errorf("failed to get Deck records for player %d: %w", playerId, err)
	}

	// Return an empty list instead of nil
	if decks == nil {
		return []Deck{}, nil
	}

	return decks, nil
}

func (d *DeckProvider) GetById(ctx context.Context, deckId int) (*DeckWithStats, error) {
	var decks []Deck
	if err := d.client.Db.SelectContext(ctx, &decks, GetDeckByID, deckId); err != nil {
		return nil, fmt.Errorf("failed to get Deck record for id %d: %w", deckId, err)
	}

	if len(decks) == 0 || len(decks) > 1 {
		return nil, fmt.Errorf(
			"unexpected number of decks returned for ID %d: got %d, expected 1",
			deckId, len(decks),
		)
	}

	var gameStats GameStats
	if err := d.client.Db.SelectContext(ctx, &gameStats, GetDeckStats, deckId); err != nil {
		return nil, fmt.Errorf("failed to get Deck statistics: %w", err)
	}

	return &DeckWithStats{
		Deck:  decks[0],
		Stats: gameStats.ToStats(),
	}, nil
}

func (d *DeckProvider) Add(ctx context.Context, playerId int, commander string) error {
	var preexisting int
	if err := d.client.Db.QueryRowContext(ctx, DeckExists, playerId, commander).Scan(&preexisting); err != nil {
		return fmt.Errorf("failed to check if player %d has preexisting deck for %s: %w", playerId, commander, err)
	}

	if preexisting >= 1 {
		return ErrDeckExists
	}

	result, err := d.client.Db.ExecContext(ctx, InsertDeck, playerId, commander)
	if err != nil {
		return fmt.Errorf("failed to insert Deck record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Deck insert: got %d, expected 1", numAffected)
	}

	return nil
}

func (d *DeckProvider) Retire(ctx context.Context, deckId int) error {
	result, err := d.client.Db.ExecContext(ctx, RetireDeck, deckId)
	if err != nil {
		return fmt.Errorf("failed to retire Deck: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by retirement: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Deck retirement: got %d, expected 1", numAffected)
	}

	return nil
}
