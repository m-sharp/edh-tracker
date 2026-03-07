package deck

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getAllDecks = `SELECT id, player_id, name, format_id, retired, created_at, updated_at, deleted_at
					FROM deck WHERE retired = 0 AND deleted_at IS NULL;`
	getDecksForPlayer = `SELECT id, player_id, name, format_id, retired, created_at, updated_at, deleted_at
						   FROM deck WHERE player_id = ? AND deleted_at IS NULL;`
	getDeckByID = `SELECT id, player_id, name, format_id, retired, created_at, updated_at, deleted_at
					 FROM deck WHERE id = ? AND deleted_at IS NULL;`
	getDecksForPlayerIDs = `SELECT id, player_id, name, format_id, retired, created_at, updated_at, deleted_at
							  FROM deck WHERE player_id IN (?) AND deleted_at IS NULL;`
	getBulkAddedDecks = `SELECT id, player_id, name, format_id, retired, created_at, updated_at, deleted_at
                           FROM deck WHERE player_id IN (%s) AND name IN (%s) AND deleted_at IS NULL`
	insertDeck     = `INSERT INTO deck (player_id, name, format_id) VALUES (?, ?, ?);`
	retireDeck     = `UPDATE deck SET retired = TRUE WHERE id = ?;`
	softDeleteDeck = `UPDATE deck SET deleted_at = NOW() WHERE id = ?;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
	var decks []Model
	if err := r.client.Db.SelectContext(ctx, &decks, getAllDecks); err != nil {
		return nil, fmt.Errorf("failed to get Deck records: %w", err)
	}
	if decks == nil {
		return []Model{}, nil
	}
	return decks, nil
}

func (r *Repository) GetAllForPlayer(ctx context.Context, playerID int) ([]Model, error) {
	var decks []Model
	if err := r.client.Db.SelectContext(ctx, &decks, getDecksForPlayer, playerID); err != nil {
		return nil, fmt.Errorf("failed to get Deck records for player %d: %w", playerID, err)
	}
	if decks == nil {
		return []Model{}, nil
	}
	return decks, nil
}

func (r *Repository) GetById(ctx context.Context, deckID int) (*Model, error) {
	var decks []Model
	if err := r.client.Db.SelectContext(ctx, &decks, getDeckByID, deckID); err != nil {
		return nil, fmt.Errorf("failed to get Deck record for id %d: %w", deckID, err)
	}
	if len(decks) == 0 {
		return nil, nil
	}
	return &decks[0], nil
}

func (r *Repository) Add(ctx context.Context, playerID int, name string, formatID int) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertDeck, playerID, name, formatID)
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

func (r *Repository) BulkAdd(ctx context.Context, decks []Model) ([]Model, error) {
	if len(decks) == 0 {
		return []Model{}, nil
	}

	insertQuery := "INSERT INTO deck (player_id, name, format_id) VALUES " + strings.TrimSuffix(strings.Repeat("(?,?,?),", len(decks)), ",")
	insertArgs := make([]interface{}, 0, len(decks)*3)
	for _, d := range decks {
		insertArgs = append(insertArgs, d.PlayerID, d.Name, d.FormatID)
	}
	if _, err := r.client.Db.ExecContext(ctx, insertQuery, insertArgs...); err != nil {
		return nil, fmt.Errorf("failed to bulk insert Deck records: %w", err)
	}

	playerIDArgs := make([]interface{}, len(decks))
	nameArgs := make([]interface{}, len(decks))
	for i, d := range decks {
		playerIDArgs[i] = d.PlayerID
		nameArgs[i] = d.Name
	}
	inPlayerIDs := strings.TrimSuffix(strings.Repeat("?,", len(decks)), ",")
	inNames := strings.TrimSuffix(strings.Repeat("?,", len(decks)), ",")
	selectQuery := fmt.Sprintf(getBulkAddedDecks, inPlayerIDs, inNames)
	selectArgs := append(playerIDArgs, nameArgs...)

	var result []Model
	if err := r.client.Db.SelectContext(ctx, &result, selectQuery, selectArgs...); err != nil {
		return nil, fmt.Errorf("failed to select inserted Decks: %w", err)
	}

	return result, nil
}

func (r *Repository) GetAllByPlayerIDs(ctx context.Context, playerIDs []int) ([]Model, error) {
	if len(playerIDs) == 0 {
		return []Model{}, nil
	}

	query, args, err := sqlx.In(getDecksForPlayerIDs, playerIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to build GetAllByPlayerIDs query: %w", err)
	}
	query = r.client.Db.Rebind(query)

	var decks []Model
	if err = r.client.Db.SelectContext(ctx, &decks, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get Deck records for player IDs: %w", err)
	}
	if decks == nil {
		return []Model{}, nil
	}
	return decks, nil
}

func (r *Repository) Update(ctx context.Context, deckID int, fields UpdateFields) error {
	setClauses := []string{}
	args := []interface{}{}

	if fields.Name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *fields.Name)
	}
	if fields.FormatID != nil {
		setClauses = append(setClauses, "format_id = ?")
		args = append(args, *fields.FormatID)
	}
	if fields.Retired != nil {
		setClauses = append(setClauses, "retired = ?")
		args = append(args, *fields.Retired)
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, deckID)
	query := "UPDATE deck SET " + strings.Join(setClauses, ", ") + " WHERE id = ? AND deleted_at IS NULL;"

	result, err := r.client.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update Deck record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by update: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Deck update: got %d, expected 1", numAffected)
	}

	return nil
}

func (r *Repository) Retire(ctx context.Context, deckID int) error {
	result, err := r.client.Db.ExecContext(ctx, retireDeck, deckID)
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

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	result, err := r.client.Db.ExecContext(ctx, softDeleteDeck, id)
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
