package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
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

func (m *Migration3) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createDeckTable); err != nil {
		return lib.NewDBError(createDeckTable, err)
	}
	return nil
}

func (m *Migration3) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, destroyDeckTable); err != nil {
		return lib.NewDBError(destroyDeckTable, err)
	}
	return nil
}

func (m *Migration3) RecordMigration() bool { return true }
