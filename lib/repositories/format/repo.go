package format

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	getAllFormats   = `SELECT id, name, created_at, updated_at, deleted_at FROM format WHERE deleted_at IS NULL;`
	getFormatByID   = `SELECT id, name, created_at, updated_at, deleted_at FROM format WHERE id = ? AND deleted_at IS NULL LIMIT 1;`
	getFormatByName = `SELECT id, name, created_at, updated_at, deleted_at FROM format WHERE name = ? AND deleted_at IS NULL LIMIT 1;`
)

type Repository struct {
	client *lib.DBClient
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
	var formats []Model
	if err := r.client.Db.SelectContext(ctx, &formats, getAllFormats); err != nil {
		return nil, fmt.Errorf("failed to get Format records: %w", err)
	}
	if formats == nil {
		return []Model{}, nil
	}
	return formats, nil
}

func (r *Repository) GetById(ctx context.Context, id int) (*Model, error) {
	var formats []Model
	if err := r.client.Db.SelectContext(ctx, &formats, getFormatByID, id); err != nil {
		return nil, fmt.Errorf("failed to get Format record for id %d: %w", id, err)
	}
	if len(formats) == 0 {
		return nil, nil
	}
	return &formats[0], nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*Model, error) {
	var formats []Model
	if err := r.client.Db.SelectContext(ctx, &formats, getFormatByName, name); err != nil {
		return nil, fmt.Errorf("failed to get Format record for name %q: %w", name, err)
	}
	if len(formats) == 0 {
		return nil, nil
	}
	return &formats[0], nil
}
