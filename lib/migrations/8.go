package migrations

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	dropGames        = `DELETE * FROM game;`
	insertResultTmpl = `INSERT INTO game_result (game_id, deck_id, place, kill_count) VALUES (
		(SELECT id FROM game where description = "%s"),
		(SELECT id FROM deck where commander = "%s" AND player_id = (SELECT id FROM player where name = "%s")),
		%d,
		%d
	);`
)

type resultInfo struct {
	Player    string
	Commander string
	Place     int
	Kills     int
}

var (
	gameSeeds = []string{
		`INSERT INTO game (description, ctime) VALUES ('Game 1', '2023-04-06');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 2', '2023-04-06');`,
	}

	gameResultSeeds = map[string][]resultInfo{
		"Game 1": {
			{
				Player:    "Peter",
				Commander: "Thalisse, Reverent Medium",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Tom",
				Commander: "Zethi, Arcane Blademaster",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Go-Shintai of Life's Origin",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Jhoira, Weatherlight Captain",
				Place:     1,
				Kills:     3,
			},
		},
		"Game 2": {
			{
				Player:    "Peter",
				Commander: "Atraxa, Praetor's Voice",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Tom",
				Commander: "Old Stickfingers",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Rakdos, Lord of Riots",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Sakashima of a Thousand Faces",
				Place:     3,
				Kills:     1,
			},
		},
	}
)

type Migration8 struct{}

func (m *Migration8) Upgrade(ctx context.Context, client *lib.DBClient) error {
	for _, gameSeed := range gameSeeds {
		if _, err := client.Db.ExecContext(ctx, gameSeed); err != nil {
			return lib.NewDBError(gameSeed, err)
		}
	}

	for game, resultInfos := range gameResultSeeds {
		for _, result := range resultInfos {
			query := fmt.Sprintf(insertResultTmpl, game, result.Commander, result.Player, result.Place, result.Kills)
			if _, err := client.Db.ExecContext(ctx, query); err != nil {
				return lib.NewDBError(query, err)
			}
		}
	}
	return nil
}

func (m *Migration8) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropGames); err != nil {
		return lib.NewDBError(dropGames, err)
	}
	return nil
}
