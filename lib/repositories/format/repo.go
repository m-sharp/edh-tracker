package format

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{db: client.GormDb}
}

func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
	var formats []Model
	if err := r.db.WithContext(ctx).Find(&formats).Error; err != nil {
		return nil, fmt.Errorf("failed to get Format records: %w", err)
	}
	if formats == nil {
		return []Model{}, nil
	}
	return formats, nil
}

func (r *Repository) GetById(ctx context.Context, id int) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Format record for id %d: %w", id, err)
	}
	return &m, nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Format record for name %q: %w", name, err)
	}
	return &m, nil
}
