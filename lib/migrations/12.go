package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createCommanderTable = `CREATE TABLE commander (
		id         INT AUTO_INCREMENT PRIMARY KEY,
		name       VARCHAR(255) NOT NULL UNIQUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME NULL
	);`
	dropCommanderTable = `DROP TABLE commander;`
)

type Migration12 struct{}

func (m *Migration12) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(createCommanderTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", createCommanderTable, err)
	}
	return nil
}

func (m *Migration12) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(dropCommanderTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", dropCommanderTable, err)
	}
	return nil
}
