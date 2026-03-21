package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createFormatTable = `CREATE TABLE format (
		id         INT AUTO_INCREMENT PRIMARY KEY,
		name       VARCHAR(255) NOT NULL UNIQUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME NULL
	);`
	seedFormatTable = `INSERT INTO format (name) VALUES ('commander'), ('other');`
	dropFormatTable = `DROP TABLE format;`
)

type Migration11 struct{}

func (m *Migration11) Upgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{createFormatTable, seedFormatTable} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}

func (m *Migration11) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(dropFormatTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", dropFormatTable, err)
	}
	return nil
}
