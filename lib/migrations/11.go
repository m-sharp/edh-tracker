package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	createFormatTable = `CREATE TABLE format (
		id         INT AUTO_INCREMENT PRIMARY KEY,
		name       VARCHAR(255) NOT NULL UNIQUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME NULL
	);`
	seedFormatTable = `INSERT INTO format (name) VALUES ('commander'), ('other');`
	dropFormatTable = `DROP TABLE format;`
)

type Migration11 struct{}

func (m *Migration11) Upgrade(ctx context.Context, client *lib.DBClient) error {
	for _, stmt := range []string{createFormatTable, seedFormatTable} {
		if _, err := client.Db.ExecContext(ctx, stmt); err != nil {
			return lib.NewDBError(stmt, err)
		}
	}
	return nil
}

func (m *Migration11) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropFormatTable); err != nil {
		return lib.NewDBError(dropFormatTable, err)
	}
	return nil
}
