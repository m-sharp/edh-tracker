package user

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

func (r *Repository) GetByID(ctx context.Context, id int) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get User record for id %d: %w", id, err)
	}
	return &m, nil
}

func (r *Repository) GetByPlayerID(ctx context.Context, playerID int) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).Where("player_id = ?", playerID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get User record for player_id %d: %w", playerID, err)
	}
	return &m, nil
}

func (r *Repository) GetByOAuth(ctx context.Context, provider, subject string) (*Model, error) {
	var m Model
	err := r.db.WithContext(ctx).Where("oauth_provider = ? AND oauth_subject = ?", provider, subject).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get User record for oauth %s/%s: %w", provider, subject, err)
	}
	return &m, nil
}

func (r *Repository) GetRoleByName(ctx context.Context, name string) (*RoleModel, error) {
	var m RoleModel
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no UserRole found with name %q", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get UserRole record for name %q: %w", name, err)
	}
	return &m, nil
}

func (r *Repository) Add(ctx context.Context, playerID, roleID int) (int, error) {
	m := Model{PlayerID: playerID, RoleID: roleID}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert User record: %w", err)
	}
	return m.ID, nil
}

func (r *Repository) AddWithOAuth(
	ctx context.Context,
	playerID, roleID int,
	provider, subject, email, displayName, avatarURL string,
) (int, error) {
	m := Model{
		PlayerID:      playerID,
		RoleID:        roleID,
		OAuthProvider: &provider,
		OAuthSubject:  &subject,
		Email:         &email,
		DisplayName:   &displayName,
		AvatarURL:     &avatarURL,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to insert User record with OAuth: %w", err)
	}
	return m.ID, nil
}

type playerRow struct {
	ID   int    `gorm:"primaryKey"`
	Name string
}

func (r *Repository) CreatePlayerAndUser(
	ctx context.Context,
	playerName string,
	roleID int,
	provider, subject, email, displayName, avatarURL string,
) (*Model, error) {
	var result Model
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		pr := playerRow{Name: playerName}
		if err := tx.Table("player").Create(&pr).Error; err != nil {
			return fmt.Errorf("failed to insert player in CreatePlayerAndUser: %w", err)
		}

		m := Model{
			PlayerID:      pr.ID,
			RoleID:        roleID,
			OAuthProvider: &provider,
			OAuthSubject:  &subject,
			Email:         &email,
			DisplayName:   &displayName,
			AvatarURL:     &avatarURL,
		}
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("failed to insert user in CreatePlayerAndUser: %w", err)
		}
		result = m
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) BulkAdd(ctx context.Context, playerIDs []int, roleID int) error {
	if len(playerIDs) == 0 {
		return nil
	}
	entries := make([]Model, len(playerIDs))
	for i, id := range playerIDs {
		entries[i] = Model{PlayerID: id, RoleID: roleID}
	}
	if err := r.db.WithContext(ctx).CreateInBatches(&entries, 100).Error; err != nil {
		return fmt.Errorf("failed to bulk insert User records: %w", err)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&Model{}, id).Error; err != nil {
		return fmt.Errorf("failed to soft-delete User record: %w", err)
	}
	return nil
}
