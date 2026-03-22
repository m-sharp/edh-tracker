package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createGameTable = `CREATE TABLE game(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		description VARCHAR(256),
		ctime DATETIME DEFAULT NOW()
	);`
	destroyGameTable = `DROP TABLE game;`
)

type Migration4 struct{}

func (m *Migration4) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(createGameTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", createGameTable, err)
	}
	return nil
}

func (m *Migration4) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(destroyGameTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", destroyGameTable, err)
	}
	return nil
}
