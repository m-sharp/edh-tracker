package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createDeckTable = `CREATE TABLE deck(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		player_id INT,
		commander VARCHAR(256),
		retired BOOL DEFAULT FALSE,
		ctime DATETIME DEFAULT NOW(),
		FOREIGN KEY (player_id) REFERENCES player(id) ON DELETE CASCADE
	);`
	destroyDeckTable = `DROP TABLE deck;`
)

type Migration3 struct{}

func (m *Migration3) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(createDeckTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", createDeckTable, err)
	}
	return nil
}

func (m *Migration3) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(destroyDeckTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", destroyDeckTable, err)
	}
	return nil
}
