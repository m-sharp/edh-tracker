package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	createDeckCommanderTable = `CREATE TABLE deck_commander (
		id                   INT AUTO_INCREMENT PRIMARY KEY,
		deck_id              INT NOT NULL,
		commander_id         INT NOT NULL,
		partner_commander_id INT NULL,
		created_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at           DATETIME NULL,
		CONSTRAINT fk_dc_deck    FOREIGN KEY (deck_id)              REFERENCES deck(id),
		CONSTRAINT fk_dc_cmd     FOREIGN KEY (commander_id)         REFERENCES commander(id),
		CONSTRAINT fk_dc_partner FOREIGN KEY (partner_commander_id) REFERENCES commander(id)
	);`
	dropDeckCommanderTable = `DROP TABLE deck_commander;`
)

type Migration13 struct{}

func (m *Migration13) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createDeckCommanderTable); err != nil {
		return lib.NewDBError(createDeckCommanderTable, err)
	}
	return nil
}

func (m *Migration13) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropDeckCommanderTable); err != nil {
		return lib.NewDBError(dropDeckCommanderTable, err)
	}
	return nil
}
