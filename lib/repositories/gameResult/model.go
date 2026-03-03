package gameResult

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.ModelBase
	GameID    int `db:"game_id"`
	DeckID    int `db:"deck_id"`
	Place     int `db:"place"`
	KillCount int `db:"kill_count"`
}
