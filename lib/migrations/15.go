package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	addGameFormatID  = `ALTER TABLE game ADD COLUMN format_id INT NOT NULL DEFAULT 1, ADD CONSTRAINT fk_game_format FOREIGN KEY (format_id) REFERENCES format(id);`
	dropGameFormatID = `ALTER TABLE game DROP FOREIGN KEY fk_game_format, DROP COLUMN format_id;`
)

type Migration15 struct{}

func (m *Migration15) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, addGameFormatID); err != nil {
		return lib.NewDBError(addGameFormatID, err)
	}
	return nil
}

func (m *Migration15) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropGameFormatID); err != nil {
		return lib.NewDBError(dropGameFormatID, err)
	}
	return nil
}
