package models

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetCommanderByID   = `SELECT id, name, created_at, updated_at, deleted_at FROM commander WHERE id = ? AND deleted_at IS NULL LIMIT 1;`
	GetCommanderByName = `SELECT id, name, created_at, updated_at, deleted_at FROM commander WHERE name = ? AND deleted_at IS NULL LIMIT 1;`
	InsertCommander    = `INSERT INTO commander (name) VALUES (?);`
)

type Commander struct {
	Model
	Name string `json:"name" db:"name"`
}

type CommanderProvider struct {
	client *lib.DBClient
}

func NewCommanderProvider(client *lib.DBClient) *CommanderProvider {
	return &CommanderProvider{client: client}
}

func (c *CommanderProvider) GetById(ctx context.Context, id int) (*Commander, error) {
	var commanders []Commander
	if err := c.client.Db.SelectContext(ctx, &commanders, GetCommanderByID, id); err != nil {
		return nil, fmt.Errorf("failed to get Commander record for id %d: %w", id, err)
	}
	if len(commanders) == 0 {
		return nil, nil
	}
	return &commanders[0], nil
}

func (c *CommanderProvider) GetByName(ctx context.Context, name string) (*Commander, error) {
	var commanders []Commander
	if err := c.client.Db.SelectContext(ctx, &commanders, GetCommanderByName, name); err != nil {
		return nil, fmt.Errorf("failed to get Commander record for name %q: %w", name, err)
	}
	if len(commanders) == 0 {
		return nil, nil
	}
	return &commanders[0], nil
}

func (c *CommanderProvider) Add(ctx context.Context, name string) (int, error) {
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
