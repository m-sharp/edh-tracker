package deckCommander

import (
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	commanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/commander"
)

type Model struct {
	base.GormModelBase
	DeckID int

	CommanderID int
	Commander   commanderRepo.Model `gorm:"foreignKey:CommanderID"`

	PartnerCommanderID *int
	PartnerCommander   *commanderRepo.Model `gorm:"foreignKey:PartnerCommanderID"`
}

func (Model) TableName() string { return "deck_commander" }
