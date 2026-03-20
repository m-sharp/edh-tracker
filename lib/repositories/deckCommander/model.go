package deckCommander

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.GormModelBase
	DeckID             int
	CommanderID        int
	PartnerCommanderID *int
}

func (Model) TableName() string { return "deck_commander" }
