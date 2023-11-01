package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	createPlayerTable = `CREATE TABLE player(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(256),
		ctime DATETIME
	);`
	destroyPlayerTable = `DROP TABLE player;`
)

type Migration2 struct{}

func (m *Migration2) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createPlayerTable); err != nil {
		return lib.NewDBError(createPlayerTable, err)
	}
	return nil
}

func (m *Migration2) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, destroyPlayerTable); err != nil {
		return lib.NewDBError(destroyPlayerTable, err)
	}
	return nil
}
