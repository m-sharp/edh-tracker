package models

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetAllFormats   = `SELECT id, name, created_at, updated_at, deleted_at FROM format WHERE deleted_at IS NULL;`
	GetFormatByID   = `SELECT id, name, created_at, updated_at, deleted_at FROM format WHERE id = ? AND deleted_at IS NULL LIMIT 1;`
	GetFormatByName = `SELECT id, name, created_at, updated_at, deleted_at FROM format WHERE name = ? AND deleted_at IS NULL LIMIT 1;`
)

type Format struct {
	Model
	Name string `json:"name" db:"name"`
}

type FormatRepository struct {
	client *lib.DBClient
}

func NewFormatRepository(client *lib.DBClient) *FormatRepository {
	return &FormatRepository{client: client}
}

func (f *FormatRepository) GetAll(ctx context.Context) ([]Format, error) {
	var formats []Format
	if err := f.client.Db.SelectContext(ctx, &formats, GetAllFormats); err != nil {
		return nil, fmt.Errorf("failed to get Format records: %w", err)
	}
	if formats == nil {
		return []Format{}, nil
	}
	return formats, nil
}

func (f *FormatRepository) GetById(ctx context.Context, id int) (*Format, error) {
	var formats []Format
	if err := f.client.Db.SelectContext(ctx, &formats, GetFormatByID, id); err != nil {
		return nil, fmt.Errorf("failed to get Format record for id %d: %w", id, err)
	}
	if len(formats) == 0 {
		return nil, nil
	}
	return &formats[0], nil
}

func (f *FormatRepository) GetByName(ctx context.Context, name string) (*Format, error) {
	var formats []Format
	if err := f.client.Db.SelectContext(ctx, &formats, GetFormatByName, name); err != nil {
		return nil, fmt.Errorf("failed to get Format record for name %q: %w", name, err)
	}
	if len(formats) == 0 {
		return nil, nil
	}
	return &formats[0], nil
}
