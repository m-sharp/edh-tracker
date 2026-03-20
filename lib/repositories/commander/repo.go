package commander

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

func NewRepositoryFromDB(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetById(ctx context.Context, id int) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Commander record for id %d: %w", id, err)
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
		return nil, fmt.Errorf("failed to get Commander record for name %q: %w", name, err)
	}
	return &m, nil
}

func (r *Repository) GetByNames(ctx context.Context, names []string) ([]Model, error) {
	if len(names) == 0 {
		return []Model{}, nil
	}
	var commanders []Model
	if err := r.db.WithContext(ctx).Where("name IN ?", names).Find(&commanders).Error; err != nil {
		return nil, fmt.Errorf("failed to get Commander records by names: %w", err)
	}
	return commanders, nil
}

func (r *Repository) Add(ctx context.Context, name string) (int, error) {
	m := Model{Name: name}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert Commander record: %w", err)
	}
	return m.ID, nil
}

func (r *Repository) BulkAdd(ctx context.Context, names []string) ([]Model, error) {
	if len(names) == 0 {
		return []Model{}, nil
	}
	models := make([]Model, len(names))
	for i, n := range names {
		models[i] = Model{Name: n}
	}
	if err := r.db.WithContext(ctx).CreateInBatches(&models, 100).Error; err != nil {
		return nil, fmt.Errorf("failed to bulk insert Commander records: %w", err)
	}
	return r.GetByNames(ctx, names)
}
