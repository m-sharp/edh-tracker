package deckCommander

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.ModelBase
	DeckID             int  `db:"deck_id"`
	CommanderID        int  `db:"commander_id"`
	PartnerCommanderID *int `db:"partner_commander_id"`
}
