package gameResult

import (
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	deckRepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
)

type Model struct {
	base.GormModelBase
	GameID    int
	Place     int
	KillCount int

	DeckID int
	Deck   deckRepo.Model // BelongsTo via DeckID; populated only with GetByGameIDWithDeckInfo
}

func (Model) TableName() string { return "game_result" }
