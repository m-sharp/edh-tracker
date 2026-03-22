package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	addGameFormatID  = `ALTER TABLE game ADD COLUMN format_id INT NOT NULL DEFAULT 1, ADD CONSTRAINT fk_game_format FOREIGN KEY (format_id) REFERENCES format(id);`
	dropGameFormatID = `ALTER TABLE game DROP FOREIGN KEY fk_game_format, DROP COLUMN format_id;`
)

type Migration15 struct{}

func (m *Migration15) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(addGameFormatID).Error; err != nil {
		return fmt.Errorf("query %q: %w", addGameFormatID, err)
	}
	return nil
}

func (m *Migration15) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(dropGameFormatID).Error; err != nil {
		return fmt.Errorf("query %q: %w", dropGameFormatID, err)
	}
	return nil
}
