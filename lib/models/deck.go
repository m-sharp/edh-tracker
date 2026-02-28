package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetAllDecks = `SELECT deck.id, deck.player_id, player.name AS player_name, deck.name, deck.format_id,
						format.name AS format_name, deck.retired,
						deck.created_at, deck.updated_at, deck.deleted_at,
						dc.commander_id, c.name AS commander_name,
						dc.partner_commander_id, pc.name AS partner_commander_name
					FROM deck
					INNER JOIN player ON deck.player_id = player.id
					INNER JOIN format ON deck.format_id = format.id
					LEFT JOIN deck_commander dc ON dc.deck_id = deck.id AND dc.deleted_at IS NULL
					LEFT JOIN commander c  ON dc.commander_id         = c.id  AND c.deleted_at  IS NULL
					LEFT JOIN commander pc ON dc.partner_commander_id = pc.id AND pc.deleted_at IS NULL
					WHERE deck.retired = 0
					  AND deck.deleted_at IS NULL
					  AND player.deleted_at IS NULL;`

	GetDecksForPlayer = `SELECT deck.id, deck.player_id, player.name AS player_name, deck.name, deck.format_id,
						format.name AS format_name, deck.retired,
						deck.created_at, deck.updated_at, deck.deleted_at,
						dc.commander_id, c.name AS commander_name,
						dc.partner_commander_id, pc.name AS partner_commander_name
					FROM deck
					INNER JOIN player ON deck.player_id = player.id
					INNER JOIN format ON deck.format_id = format.id
					LEFT JOIN deck_commander dc ON dc.deck_id = deck.id AND dc.deleted_at IS NULL
					LEFT JOIN commander c  ON dc.commander_id         = c.id  AND c.deleted_at  IS NULL
					LEFT JOIN commander pc ON dc.partner_commander_id = pc.id AND pc.deleted_at IS NULL
					WHERE deck.player_id = ? AND deck.deleted_at IS NULL;`

	GetDeckByID = `SELECT deck.id, deck.player_id, player.name AS player_name, deck.name, deck.format_id,
						format.name AS format_name, deck.retired,
						deck.created_at, deck.updated_at, deck.deleted_at,
						dc.commander_id, c.name AS commander_name,
						dc.partner_commander_id, pc.name AS partner_commander_name
					FROM deck
					INNER JOIN player ON deck.player_id = player.id
					INNER JOIN format ON deck.format_id = format.id
					LEFT JOIN deck_commander dc ON dc.deck_id = deck.id AND dc.deleted_at IS NULL
					LEFT JOIN commander c  ON dc.commander_id         = c.id  AND c.deleted_at  IS NULL
					LEFT JOIN commander pc ON dc.partner_commander_id = pc.id AND pc.deleted_at IS NULL
					WHERE deck.id = ? AND deck.deleted_at IS NULL AND player.deleted_at IS NULL;`

	GetDeckStats = `SELECT DISTINCT game_result.game_id, game_result.place, game_result.kill_count
						FROM (game_result INNER JOIN deck on game_result.deck_id = deck.id)
					WHERE deck.id = ? AND game_result.deleted_at IS NULL;`

	InsertDeck     = `INSERT INTO deck (player_id, name, format_id) VALUES (?, ?, ?);`
	RetireDeck     = `UPDATE deck SET retired = TRUE WHERE id = ?;`
	SoftDeleteDeck = `UPDATE deck SET deleted_at = NOW() WHERE id = ?;`

	deckValidationErr = "invalid Deck: %s"
)

type DeckCommanderEntry struct {
	CommanderID          int     `json:"commander_id"`
	CommanderName        string  `json:"commander_name"`
	PartnerCommanderID   *int    `json:"partner_commander_id,omitempty"`
	PartnerCommanderName *string `json:"partner_commander_name,omitempty"`
}

type Deck struct {
	Model
	PlayerID   int                 `json:"player_id"   db:"player_id"`
	PlayerName string              `json:"player_name,omitempty" db:"player_name"`
	Name       string              `json:"name"        db:"name"`
	FormatID   int                 `json:"format_id"   db:"format_id"`
	FormatName string              `json:"format_name" db:"format_name"`
	Retired    bool                `json:"retired"     db:"retired"`
	Commanders *DeckCommanderEntry `json:"commanders,omitempty"`
}

// deckRow is a flat scan struct for queries that LEFT JOIN commander tables.
type deckRow struct {
	Model
	PlayerID             int            `db:"player_id"`
	PlayerName           string         `db:"player_name"`
	Name                 string         `db:"name"`
	FormatID             int            `db:"format_id"`
	FormatName           string         `db:"format_name"`
	Retired              bool           `db:"retired"`
	CommanderID          sql.NullInt64  `db:"commander_id"`
	CommanderName        sql.NullString `db:"commander_name"`
	PartnerCommanderID   sql.NullInt64  `db:"partner_commander_id"`
	PartnerCommanderName sql.NullString `db:"partner_commander_name"`
}

func (r *deckRow) toDeck() Deck {
	deck := Deck{
		Model:      r.Model,
		PlayerID:   r.PlayerID,
		PlayerName: r.PlayerName,
		Name:       r.Name,
		FormatID:   r.FormatID,
		FormatName: r.FormatName,
		Retired:    r.Retired,
	}
	if r.CommanderID.Valid {
		cmdID := int(r.CommanderID.Int64)
		entry := &DeckCommanderEntry{
			CommanderID:   cmdID,
			CommanderName: r.CommanderName.String,
		}
		if r.PartnerCommanderID.Valid {
			partnerID := int(r.PartnerCommanderID.Int64)
			entry.PartnerCommanderID = &partnerID
			entry.PartnerCommanderName = &r.PartnerCommanderName.String
		}
		deck.Commanders = entry
	}
	return deck
}

type DeckWithStats struct {
	Deck
	Stats
}

func (d *Deck) Validate() error {
	if d.PlayerID == 0 {
		return fmt.Errorf(deckValidationErr, "missing PlayerID")
	}
	if d.Name == "" {
		return fmt.Errorf(deckValidationErr, "missing Name")
	}
	if d.FormatID == 0 {
		return fmt.Errorf(deckValidationErr, "missing FormatID")
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

func (d *DeckProvider) GetAll(ctx context.Context) ([]DeckWithStats, error) {
	var rows []deckRow
	if err := d.client.Db.SelectContext(ctx, &rows, GetAllDecks); err != nil {
		return nil, fmt.Errorf("failed to get Deck records: %w", err)
	}

	if rows == nil {
		return []DeckWithStats{}, nil
	}

	var result []DeckWithStats
	for _, row := range rows {
		deck := row.toDeck()
		var gameStats GameStats
		if err := d.client.Db.SelectContext(ctx, &gameStats, GetDeckStats, deck.ID); err != nil {
			return nil, fmt.Errorf("failed to get Deck statistics: %w", err)
		}
		result = append(result, DeckWithStats{Deck: deck, Stats: gameStats.ToStats()})
	}

	return result, nil
}

func (d *DeckProvider) GetAllForPlayer(ctx context.Context, playerID int) ([]DeckWithStats, error) {
	var rows []deckRow
	if err := d.client.Db.SelectContext(ctx, &rows, GetDecksForPlayer, playerID); err != nil {
		return nil, fmt.Errorf("failed to get Deck records for player %d: %w", playerID, err)
	}

	if rows == nil {
		return []DeckWithStats{}, nil
	}

	var result []DeckWithStats
	for _, row := range rows {
		deck := row.toDeck()
		var gameStats GameStats
		if err := d.client.Db.SelectContext(ctx, &gameStats, GetDeckStats, deck.ID); err != nil {
			return nil, fmt.Errorf("failed to get Deck statistics: %w", err)
		}
		result = append(result, DeckWithStats{Deck: deck, Stats: gameStats.ToStats()})
	}

	return result, nil
}

func (d *DeckProvider) GetById(ctx context.Context, deckID int) (*DeckWithStats, error) {
	var rows []deckRow
	if err := d.client.Db.SelectContext(ctx, &rows, GetDeckByID, deckID); err != nil {
		return nil, fmt.Errorf("failed to get Deck record for id %d: %w", deckID, err)
	}

	if len(rows) == 0 || len(rows) > 1 {
		return nil, fmt.Errorf(
			"unexpected number of decks returned for ID %d: got %d, expected 1",
			deckID, len(rows),
		)
	}

	deck := rows[0].toDeck()
	var gameStats GameStats
	if err := d.client.Db.SelectContext(ctx, &gameStats, GetDeckStats, deckID); err != nil {
		return nil, fmt.Errorf("failed to get Deck statistics: %w", err)
	}

	return &DeckWithStats{
		Deck:  deck,
		Stats: gameStats.ToStats(),
	}, nil
}

func (d *DeckProvider) Add(ctx context.Context, playerID int, name string, formatID int) (int, error) {
	result, err := d.client.Db.ExecContext(ctx, InsertDeck, playerID, name, formatID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert Deck record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return 0, fmt.Errorf("unexpected number of rows affected by Deck insert: got %d, expected 1", numAffected)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new Deck: %w", err)
	}

	return int(id), nil
}

func (d *DeckProvider) Retire(ctx context.Context, deckID int) error {
	result, err := d.client.Db.ExecContext(ctx, RetireDeck, deckID)
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

func (d *DeckProvider) SoftDelete(ctx context.Context, id int) error {
	result, err := d.client.Db.ExecContext(ctx, SoftDeleteDeck, id)
	if err != nil {
		return fmt.Errorf("failed to soft-delete Deck record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by soft-delete: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Deck soft-delete: got %d, expected 1", numAffected)
	}

	return nil
}
