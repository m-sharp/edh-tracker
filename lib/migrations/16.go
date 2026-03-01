package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	addFormatIndexes        = `ALTER TABLE format ADD INDEX idx_format_name (name), ADD INDEX idx_format_deleted_at (deleted_at);`
	addCommanderIndexes     = `ALTER TABLE commander ADD INDEX idx_commander_name (name), ADD INDEX idx_commander_deleted_at (deleted_at);`
	addDeckCommanderIndexes = `ALTER TABLE deck_commander ADD INDEX idx_dc_deck_id (deck_id), ADD INDEX idx_dc_deleted_at (deleted_at);`
	addDeckNewIndexes       = `ALTER TABLE deck ADD INDEX idx_deck_format_id (format_id), ADD INDEX idx_deck_name (name);`
	addGameFormatIndex      = `ALTER TABLE game ADD INDEX idx_game_format_id (format_id);`

	dropFormatIndexes        = `ALTER TABLE format DROP INDEX idx_format_name, DROP INDEX idx_format_deleted_at;`
	dropCommanderIndexes     = `ALTER TABLE commander DROP INDEX idx_commander_name, DROP INDEX idx_commander_deleted_at;`
	dropDeckCommanderIndexes = `ALTER TABLE deck_commander DROP INDEX idx_dc_deck_id, DROP INDEX idx_dc_deleted_at;`
	dropDeckNewIndexes       = `ALTER TABLE deck DROP INDEX idx_deck_format_id, DROP INDEX idx_deck_name;`
	dropGameFormatIndex      = `ALTER TABLE game DROP INDEX idx_game_format_id;`
)

type Migration16 struct{}

func (m *Migration16) Upgrade(ctx context.Context, client *lib.DBClient) error {
	for _, stmt := range []string{
		addFormatIndexes,
		addCommanderIndexes,
		addDeckCommanderIndexes,
		addDeckNewIndexes,
		addGameFormatIndex,
	} {
		if _, err := client.Db.ExecContext(ctx, stmt); err != nil {
			return lib.NewDBError(stmt, err)
		}
	}
	return nil
}

func (m *Migration16) Downgrade(ctx context.Context, client *lib.DBClient) error {
	for _, stmt := range []string{
		dropFormatIndexes,
		dropCommanderIndexes,
		dropDeckCommanderIndexes,
		dropDeckNewIndexes,
		dropGameFormatIndex,
	} {
		if _, err := client.Db.ExecContext(ctx, stmt); err != nil {
			return lib.NewDBError(stmt, err)
		}
	}
	return nil
}
