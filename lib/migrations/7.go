package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	createGameResultTable = `CREATE TABLE game_result(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		game_id INT,
		deck_id INT,
		place INT,
		kill_count INT DEFAULT 0,
		FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE,
		FOREIGN KEY (deck_id) REFERENCES deck(id) ON DELETE CASCADE
	);`
	destroyGameResultTable = `DROP TABLE game_result;`
)

type Migration7 struct{}

func (m *Migration7) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createGameResultTable); err != nil {
		return lib.NewDBError(createGameResultTable, err)
	}
	return nil
}

func (m *Migration7) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, destroyGameResultTable); err != nil {
		return lib.NewDBError(destroyGameResultTable, err)
	}
	return nil
}
