package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	alterPlayer = `ALTER TABLE player
		RENAME COLUMN ctime TO created_at,
		ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ADD COLUMN deleted_at DATETIME NULL;`
	revertPlayer = `ALTER TABLE player
		RENAME COLUMN created_at TO ctime,
		DROP COLUMN updated_at,
		DROP COLUMN deleted_at;`

	alterDeck = `ALTER TABLE deck
		RENAME COLUMN ctime TO created_at,
		ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ADD COLUMN deleted_at DATETIME NULL;`
	revertDeck = `ALTER TABLE deck
		RENAME COLUMN created_at TO ctime,
		DROP COLUMN updated_at,
		DROP COLUMN deleted_at;`

	alterGame = `ALTER TABLE game
		RENAME COLUMN ctime TO created_at,
		ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ADD COLUMN deleted_at DATETIME NULL;`
	revertGame = `ALTER TABLE game
		RENAME COLUMN created_at TO ctime,
		DROP COLUMN updated_at,
		DROP COLUMN deleted_at;`

	alterGameResult = `ALTER TABLE game_result
		ADD COLUMN created_at DATETIME NOT NULL DEFAULT NOW(),
		ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ADD COLUMN deleted_at DATETIME NULL;`
	revertGameResult = `ALTER TABLE game_result
		DROP COLUMN created_at,
		DROP COLUMN updated_at,
		DROP COLUMN deleted_at;`
)

type Migration6 struct{}

func (m *Migration6) Upgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{alterPlayer, alterDeck, alterGame, alterGameResult} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}

func (m *Migration6) Downgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{revertGameResult, revertGame, revertDeck, revertPlayer} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}
