package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	addDeckNameColumn   = `ALTER TABLE deck ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';`
	copyCommanderToName = `UPDATE deck SET name = commander;`
	addDeckFormatID     = `ALTER TABLE deck ADD COLUMN format_id INT NOT NULL DEFAULT 1;`
	addDeckFormatFKey   = `ALTER TABLE deck ADD CONSTRAINT fk_deck_format FOREIGN KEY (format_id) REFERENCES format(id);`
	dropDeckCommander   = `ALTER TABLE deck DROP COLUMN commander;`

	addDeckCommanderBack        = `ALTER TABLE deck ADD COLUMN commander VARCHAR(255) NOT NULL DEFAULT '';`
	copyNameToCommander         = `UPDATE deck SET commander = name;`
	dropDeckFormatFKey          = `ALTER TABLE deck DROP FOREIGN KEY fk_deck_format;`
	dropDeckFormatID            = `ALTER TABLE deck DROP COLUMN format_id;`
	dropDeckName                = `ALTER TABLE deck DROP COLUMN name;`
	addBackDeckCommanderIndexes = `ALTER TABLE deck ADD INDEX idx_deck_commander (commander), ADD INDEX idx_deck_player_commander (player_id, commander);`
)

type Migration14 struct{}

func (m *Migration14) Upgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{
		addDeckNameColumn,
		copyCommanderToName,
		addDeckFormatID,
		addDeckFormatFKey,
		dropDeckCommander,
	} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}

func (m *Migration14) Downgrade(ctx context.Context, db *gorm.DB) error {
	for _, stmt := range []string{
		addDeckCommanderBack,
		copyNameToCommander,
		dropDeckFormatFKey,
		dropDeckFormatID,
		dropDeckName,
		addBackDeckCommanderIndexes,
	} {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			return fmt.Errorf("query %q: %w", stmt, err)
		}
	}
	return nil
}
