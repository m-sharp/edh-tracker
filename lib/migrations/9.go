package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	addPodIDToGame    = `ALTER TABLE game ADD COLUMN pod_id INT NOT NULL, ADD CONSTRAINT fk_game_pod FOREIGN KEY (pod_id) REFERENCES pod (id);`
	dropPodIDFromGame = `ALTER TABLE game DROP FOREIGN KEY fk_game_pod, DROP COLUMN pod_id;`
)

type Migration9 struct{}

func (m *Migration9) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, addPodIDToGame); err != nil {
		return lib.NewDBError(addPodIDToGame, err)
	}
	return nil
}

func (m *Migration9) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropPodIDFromGame); err != nil {
		return lib.NewDBError(dropPodIDFromGame, err)
	}
	return nil
}

