package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
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

func (m *Migration14) Upgrade(ctx context.Context, client *lib.DBClient) error {
	for _, stmt := range []string{
		addDeckNameColumn,
		copyCommanderToName,
		addDeckFormatID,
		addDeckFormatFKey,
		dropDeckCommander,
	} {
		if _, err := client.Db.ExecContext(ctx, stmt); err != nil {
			return lib.NewDBError(stmt, err)
		}
	}
	return nil
}

func (m *Migration14) Downgrade(ctx context.Context, client *lib.DBClient) error {
	for _, stmt := range []string{
		addDeckCommanderBack,
		copyNameToCommander,
		dropDeckFormatFKey,
		dropDeckFormatID,
		dropDeckName,
		addBackDeckCommanderIndexes,
	} {
		if _, err := client.Db.ExecContext(ctx, stmt); err != nil {
			return lib.NewDBError(stmt, err)
		}
	}
	return nil
}
