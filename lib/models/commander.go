package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetCommanderByID     = `SELECT id, name, created_at, updated_at, deleted_at FROM commander WHERE id = ? AND deleted_at IS NULL LIMIT 1;`
	GetCommanderByName   = `SELECT id, name, created_at, updated_at, deleted_at FROM commander WHERE name = ? AND deleted_at IS NULL LIMIT 1;`
	GetCommandersByNames = `SELECT id, name FROM commander WHERE name IN (%s) AND deleted_at IS NULL`
	InsertCommander      = `INSERT INTO commander (name) VALUES (?);`
)

type Commander struct {
	Model
	Name string `json:"name" db:"name"`
}

type CommanderRepository struct {
	client *lib.DBClient
}

func NewCommanderRepository(client *lib.DBClient) *CommanderRepository {
	return &CommanderRepository{client: client}
}

func (c *CommanderRepository) GetById(ctx context.Context, id int) (*Commander, error) {
	var commanders []Commander
	if err := c.client.Db.SelectContext(ctx, &commanders, GetCommanderByID, id); err != nil {
		return nil, fmt.Errorf("failed to get Commander record for id %d: %w", id, err)
	}
	if len(commanders) == 0 {
		return nil, nil
	}
	return &commanders[0], nil
}

func (c *CommanderRepository) GetByName(ctx context.Context, name string) (*Commander, error) {
	var commanders []Commander
	if err := c.client.Db.SelectContext(ctx, &commanders, GetCommanderByName, name); err != nil {
		return nil, fmt.Errorf("failed to get Commander record for name %q: %w", name, err)
	}
	if len(commanders) == 0 {
		return nil, nil
	}
	return &commanders[0], nil
}

func (c *CommanderRepository) BulkAdd(ctx context.Context, names []string) ([]Commander, error) {
	if len(names) == 0 {
		return []Commander{}, nil
	}

	insertQuery := "INSERT INTO commander (name) VALUES " + strings.TrimSuffix(strings.Repeat("(?),", len(names)), ",")
	args := make([]interface{}, len(names))
	for i, name := range names {
		args[i] = name
	}
	if _, err := c.client.Db.ExecContext(ctx, insertQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to bulk insert Commander records: %w", err)
	}

	inPlaceholders := strings.TrimSuffix(strings.Repeat("?,", len(names)), ",")
	selectQuery := fmt.Sprintf(GetCommandersByNames, inPlaceholders)
	var commanders []Commander
	if err := c.client.Db.SelectContext(ctx, &commanders, selectQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to select inserted Commanders: %w", err)
	}

	return commanders, nil
}

func (c *CommanderRepository) Add(ctx context.Context, name string) (int, error) {
	result, err := c.client.Db.ExecContext(ctx, InsertCommander, name)
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
