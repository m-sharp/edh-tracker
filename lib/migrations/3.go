package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	dropPlayers = `DELETE * FROM player;`
)

var (
	playerSeeds = []string{
		`INSERT INTO player (name) VALUES ("Mike");`,
		`INSERT INTO player (name) VALUES ("Tom");`,
		`INSERT INTO player (name) VALUES ("Dillon");`,
		`INSERT INTO player (name) VALUES ("Peter");`,
	}
)

type Migration3 struct{}

func (m *Migration3) Upgrade(ctx context.Context, client *lib.DBClient) error {
	for _, playerSeed := range playerSeeds {
		if _, err := client.Db.ExecContext(ctx, playerSeed); err != nil {
			return lib.NewDBError(playerSeed, err)
		}
	}
	return nil
}

func (m *Migration3) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropPlayers); err != nil {
		return lib.NewDBError(dropPlayers, err)
	}
	return nil
}
