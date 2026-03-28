package player

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
)

type Repository struct {
	*base.Repo
}

func NewRepository(client *lib.DBClient) *Repository {
	return &Repository{Repo: base.NewRepo(client.GormDb)}
}

func NewRepositoryFromDB(db *gorm.DB) *Repository {
	return &Repository{Repo: base.NewRepo(db)}
}

func (r *Repository) GetAll(ctx context.Context) ([]Model, error) {
	var players []Model
	if err := r.DB().WithContext(ctx).Find(&players).Error; err != nil {
		return nil, fmt.Errorf("failed to get Player records: %w", err)
	}
	if players == nil {
		return []Model{}, nil
	}
	return players, nil
}

func (r *Repository) GetById(ctx context.Context, playerID int) (*Model, error) {
	var m Model
	err := r.DB().WithContext(ctx).First(&m, playerID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Player record for id %d: %w", playerID, err)
	}
	return &m, nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*Model, error) {
	var m Model
	err := r.DB().WithContext(ctx).Where("name = ?", name).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Player record for name %q: %w", name, err)
	}
	return &m, nil
}

func (r *Repository) GetByNames(ctx context.Context, names []string) ([]Model, error) {
	if len(names) == 0 {
		return []Model{}, nil
	}
	var players []Model
	if err := r.DB().WithContext(ctx).Where("name IN ?", names).Find(&players).Error; err != nil {
		return nil, fmt.Errorf("failed to get Player records by names: %w", err)
	}
	return players, nil
}

func (r *Repository) Add(ctx context.Context, name string) (int, error) {
	m := Model{Name: name}
	if err := r.DB().WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert Player record: %w", err)
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
	if err := r.DB().WithContext(ctx).CreateInBatches(&models, 100).Error; err != nil {
		return nil, fmt.Errorf("failed to bulk insert Player records: %w", err)
	}
	return r.GetByNames(ctx, names)
}

func (r *Repository) Update(ctx context.Context, playerID int, name string) error {
	result := r.DB().WithContext(ctx).Model(&Model{}).Where("id = ?", playerID).Update("name", name)
	if result.Error != nil {
		return fmt.Errorf("failed to update Player record: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("unexpected number of rows affected by Player update: got 0, expected 1")
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	m := Model{}
	m.ID = id
	result := r.DB().WithContext(ctx).Delete(&m)
	if result.Error != nil {
		return fmt.Errorf("failed to soft-delete Player record: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("unexpected number of rows affected by Player soft-delete: got 0, expected 1")
	}
	return nil
}
