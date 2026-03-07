package deckCommander

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getDeckCommanderByDeckId = `SELECT id, deck_id, commander_id, partner_commander_id, created_at, updated_at, deleted_at
								  FROM deck_commander WHERE deck_id = ? AND deleted_at IS NULL LIMIT 1;`
	insertDeckCommander         = `INSERT INTO deck_commander (deck_id, commander_id, partner_commander_id) VALUES (?, ?, ?);`
	deleteDeckCommanderByDeckID = `UPDATE deck_commander SET deleted_at = NOW() WHERE deck_id = ? AND deleted_at IS NULL;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetByDeckId(ctx context.Context, deckID int) (*Model, error) {
	var rows []Model
	if err := r.client.Db.SelectContext(ctx, &rows, getDeckCommanderByDeckId, deckID); err != nil {
		return nil, fmt.Errorf("failed to get DeckCommander record for deck %d: %w", deckID, err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func (r *Repository) Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertDeckCommander, deckID, commanderID, partnerCommanderID)
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

func (r *Repository) DeleteByDeckID(ctx context.Context, deckID int) error {
	if _, err := r.client.Db.ExecContext(ctx, deleteDeckCommanderByDeckID, deckID); err != nil {
		return fmt.Errorf("failed to soft-delete DeckCommander records for deck %d: %w", deckID, err)
	}
	return nil
}

func (r *Repository) BulkAdd(ctx context.Context, entries []Model) error {
	if len(entries) == 0 {
		return nil
	}

	query := "INSERT INTO deck_commander (deck_id, commander_id, partner_commander_id) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?,?),", len(entries)), ",")
	args := make([]interface{}, 0, len(entries)*3)
	for _, e := range entries {
		args = append(args, e.DeckID, e.CommanderID, e.PartnerCommanderID)
	}
	if _, err := r.client.Db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk insert DeckCommander records: %w", err)
	}
	return nil
}
