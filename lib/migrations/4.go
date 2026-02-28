package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	createGameTable = `CREATE TABLE game(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		description VARCHAR(256),
		ctime DATETIME DEFAULT NOW()
	);`
	destroyGameTable = `DROP TABLE game;`
)

// ToDo: Move this up in order
type Migration4 struct{}

func (m *Migration4) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createGameTable); err != nil {
		return lib.NewDBError(createGameTable, err)
	}
	return nil
}

func (m *Migration4) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, destroyGameTable); err != nil {
		return lib.NewDBError(destroyGameTable, err)
	}
	return nil
}

