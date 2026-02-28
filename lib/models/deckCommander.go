package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetDeckCommanderByDeckId = `SELECT id, deck_id, commander_id, partner_commander_id, created_at, updated_at, deleted_at
								FROM deck_commander WHERE deck_id = ? AND deleted_at IS NULL LIMIT 1;`
	InsertDeckCommander = `INSERT INTO deck_commander (deck_id, commander_id, partner_commander_id) VALUES (?, ?, ?);`
)

type DeckCommander struct {
	Model
	DeckID             int  `json:"deck_id"              db:"deck_id"`
	CommanderID        int  `json:"commander_id"         db:"commander_id"`
	PartnerCommanderID *int `json:"partner_commander_id" db:"partner_commander_id"`
}

type DeckCommanderProvider struct {
	client *lib.DBClient
}

func NewDeckCommanderProvider(client *lib.DBClient) *DeckCommanderProvider {
	return &DeckCommanderProvider{client: client}
}

func (d *DeckCommanderProvider) GetByDeckId(ctx context.Context, deckID int) (*DeckCommander, error) {
	var rows []DeckCommander
	if err := d.client.Db.SelectContext(ctx, &rows, GetDeckCommanderByDeckId, deckID); err != nil {
		return nil, fmt.Errorf("failed to get DeckCommander record for deck %d: %w", deckID, err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func (d *DeckCommanderProvider) BulkAdd(ctx context.Context, entries []DeckCommander) error {
	if len(entries) == 0 {
		return nil
	}

	query := "INSERT INTO deck_commander (deck_id, commander_id, partner_commander_id) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?,?),", len(entries)), ",")
	args := make([]interface{}, 0, len(entries)*3)
	for _, e := range entries {
		args = append(args, e.DeckID, e.CommanderID, e.PartnerCommanderID)
	}
	if _, err := d.client.Db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk insert DeckCommander records: %w", err)
	}
	return nil
}

func (d *DeckCommanderProvider) Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
	result, err := d.client.Db.ExecContext(ctx, InsertDeckCommander, deckID, commanderID, partnerCommanderID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert DeckCommander record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return 0, fmt.Errorf("unexpected number of rows affected by DeckCommander insert: got %d, expected 1", numAffected)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new DeckCommander: %w", err)
	}

	return int(id), nil
}
