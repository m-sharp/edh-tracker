package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	createGameResultTable = `CREATE TABLE game_result(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		game_id INT,
		deck_id INT,
		place INT,
		kill_count INT DEFAULT 0,
		FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE,
		FOREIGN KEY (deck_id) REFERENCES deck(id) ON DELETE CASCADE
	);`
	destroyGameResultTable = `DROP TABLE game_result;`
)

type Migration5 struct{}

func (m *Migration5) Upgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(createGameResultTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", createGameResultTable, err)
	}
	return nil
}

func (m *Migration5) Downgrade(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(destroyGameResultTable).Error; err != nil {
		return fmt.Errorf("query %q: %w", destroyGameResultTable, err)
	}
	return nil
}
