package gameResult

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.GormModelBase
	GameID    int
	DeckID    int
	Place     int
	KillCount int
}

func (Model) TableName() string { return "game_result" }
