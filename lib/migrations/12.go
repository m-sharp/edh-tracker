package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	createCommanderTable = `CREATE TABLE commander (
		id         INT AUTO_INCREMENT PRIMARY KEY,
		name       VARCHAR(255) NOT NULL UNIQUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME NULL
	);`
	dropCommanderTable = `DROP TABLE commander;`
)

type Migration12 struct{}

func (m *Migration12) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createCommanderTable); err != nil {
		return lib.NewDBError(createCommanderTable, err)
	}
	return nil
}

func (m *Migration12) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropCommanderTable); err != nil {
		return lib.NewDBError(dropCommanderTable, err)
	}
	return nil
}
