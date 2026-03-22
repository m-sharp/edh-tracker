package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createMigrationTable = `CREATE TABLE migration(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		success BIT,
		ctime DATETIME
	);`
	destroyMigrationTable = `DROP TABLE migration;`
)

type Migration1 struct{}

func (m *Migration1) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(createMigrationTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", createMigrationTable, err)
	}
	return nil
}

func (m *Migration1) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(destroyMigrationTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", destroyMigrationTable, err)
	}
	return nil
}
