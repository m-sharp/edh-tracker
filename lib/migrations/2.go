package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createPlayerTable = `CREATE TABLE player(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(256),
		ctime DATETIME DEFAULT NOW()
	);`
	destroyPlayerTable = `DROP TABLE player;`
)

type Migration2 struct{}

func (m *Migration2) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(createPlayerTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", createPlayerTable, err)
	}
	return nil
}

func (m *Migration2) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(destroyPlayerTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", destroyPlayerTable, err)
	}
	return nil
}
