package podInvite

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *Repository) GetByCode(ctx context.Context, code string) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).Where("invite_code = ?", code).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pod invite by code: %w", err)
	}
	return &m, nil
}

func (r *Repository) Add(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error {
	m := Model{
		PodID:             podID,
		InviteCode:        code,
		CreatedByPlayerID: createdByPlayerID,
		ExpiresAt:         expiresAt,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("failed to insert pod invite: %w", err)
	}
	return nil
}

func (r *Repository) IncrementUsedCount(ctx context.Context, code string) error {
	result := r.db.WithContext(ctx).Model(&Model{}).
		Where("invite_code = ?", code).
		Update("used_count", gorm.Expr("used_count + 1"))
	if result.Error != nil {
		return fmt.Errorf("failed to increment used_count for invite %q: %w", code, result.Error)
	}
	return nil
}
