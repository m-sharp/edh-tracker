package player

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getAllPlayers     = `SELECT id, name, created_at, updated_at, deleted_at FROM player WHERE deleted_at IS NULL;`
	getPlayerByID     = `SELECT id, name, created_at, updated_at, deleted_at FROM player WHERE id = ? AND deleted_at IS NULL;`
	getPlayerByName   = `SELECT id, name, created_at, updated_at, deleted_at FROM player WHERE name = ? AND deleted_at IS NULL LIMIT 1;`
	getPlayersByNames = `SELECT id, name, created_at, updated_at, deleted_at FROM player WHERE name IN (?) AND deleted_at IS NULL`
	insertPlayer      = `INSERT INTO player (name) VALUES (?);`
	softDeletePlayer  = `UPDATE player SET deleted_at = NOW() WHERE id = ?;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
	var players []Model
	if err := r.client.Db.SelectContext(ctx, &players, getAllPlayers); err != nil {
		return nil, fmt.Errorf("failed to get Player records: %w", err)
	}
	if players == nil {
		return []Model{}, nil
	}
	return players, nil
}

func (r *Repository) GetById(ctx context.Context, playerID int) (*Model, error) {
	var players []Model
	if err := r.client.Db.SelectContext(ctx, &players, getPlayerByID, playerID); err != nil {
		return nil, fmt.Errorf("failed to get Player record for id %d: %w", playerID, err)
	}
	if len(players) == 0 {
		return nil, nil
	}
	return &players[0], nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*Model, error) {
	var players []Model
	if err := r.client.Db.SelectContext(ctx, &players, getPlayerByName, name); err != nil {
		return nil, fmt.Errorf("failed to get Player record for name %q: %w", name, err)
	}
	if len(players) == 0 {
		return nil, nil
	}
	return &players[0], nil
}

func (r *Repository) GetByNames(ctx context.Context, names []string) ([]Model, error) {
	if len(names) == 0 {
		return []Model{}, nil
	}
	query, args, err := sqlx.In(getPlayersByNames, names)
	if err != nil {
		return nil, fmt.Errorf("failed to build GetByNames query: %w", err)
	}
	query = r.client.Db.Rebind(query)

	var players []Model
	if err = r.client.Db.SelectContext(ctx, &players, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get Player records by names: %w", err)
	}

	return players, nil
}

func (r *Repository) Add(ctx context.Context, name string) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertPlayer, name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert Player record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return 0, fmt.Errorf("unexpected number of rows affected by Player insert: got %d, expected 1", numAffected)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new Player: %w", err)
	}

	return int(id), nil
}

func (r *Repository) BulkAdd(ctx context.Context, names []string) ([]Model, error) {
	if len(names) == 0 {
		return []Model{}, nil
	}

	insertQuery := "INSERT INTO player (name) VALUES " + strings.TrimSuffix(strings.Repeat("(?),", len(names)), ",")
	args := make([]interface{}, len(names))
	for i, name := range names {
		args[i] = name
	}
	if _, err := r.client.Db.ExecContext(ctx, insertQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to bulk insert Player records: %w", err)
	}

	return r.GetByNames(ctx, names)
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	result, err := r.client.Db.ExecContext(ctx, softDeletePlayer, id)
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
