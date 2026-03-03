package commander

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getCommanderByID     = `SELECT id, name, created_at, updated_at, deleted_at FROM commander WHERE id = ? AND deleted_at IS NULL LIMIT 1;`
	getCommanderByName   = `SELECT id, name, created_at, updated_at, deleted_at FROM commander WHERE name = ? AND deleted_at IS NULL LIMIT 1;`
	getCommandersByNames = `SELECT id, name, created_at, updated_at, deleted_at FROM commander WHERE name IN (?) AND deleted_at IS NULL;`
	insertCommander      = `INSERT INTO commander (name) VALUES (?);`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetById(ctx context.Context, id int) (*Model, error) {
	var commanders []Model
	if err := r.client.Db.SelectContext(ctx, &commanders, getCommanderByID, id); err != nil {
		return nil, fmt.Errorf("failed to get Commander record for id %d: %w", id, err)
	}
	if len(commanders) == 0 {
		return nil, nil
	}
	return &commanders[0], nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*Model, error) {
	var commanders []Model
	if err := r.client.Db.SelectContext(ctx, &commanders, getCommanderByName, name); err != nil {
		return nil, fmt.Errorf("failed to get Commander record for name %q: %w", name, err)
	}
	if len(commanders) == 0 {
		return nil, nil
	}
	return &commanders[0], nil
}

func (r *Repository) GetByNames(ctx context.Context, names []string) ([]Model, error) {
	if len(names) == 0 {
		return []Model{}, nil
	}
	query, args, err := sqlx.In(getCommandersByNames, names)
	if err != nil {
		return nil, fmt.Errorf("failed to build GetByNames query: %w", err)
	}
	query = r.client.Db.Rebind(query)

	var commanders []Model
	if err = r.client.Db.SelectContext(ctx, &commanders, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get Commander records by names: %w", err)
	}

	return commanders, nil
}

func (r *Repository) Add(ctx context.Context, name string) (int, error) {
	result, err := r.client.Db.ExecContext(ctx, insertCommander, name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert Commander record: %w", err)
	}

	numAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return 0, fmt.Errorf("unexpected number of rows affected by Commander insert: got %d, expected 1", numAffected)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID for new Commander: %w", err)
	}

	return int(id), nil
}

func (r *Repository) BulkAdd(ctx context.Context, names []string) ([]Model, error) {
	if len(names) == 0 {
		return []Model{}, nil
	}

	insertQuery := "INSERT INTO commander (name) VALUES " + strings.TrimSuffix(strings.Repeat("(?),", len(names)), ",")
	args := make([]interface{}, len(names))
	for i, name := range names {
		args[i] = name
	}
	if _, err := r.client.Db.ExecContext(ctx, insertQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to bulk insert Commander records: %w", err)
	}

	return r.GetByNames(ctx, names)
}
