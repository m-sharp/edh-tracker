package game

import (
	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	gameResultRepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
)

type Model struct {
	base.GormModelBase
	Description string
	PodID       int
	FormatID    int
	Results     []gameResultRepo.Model `gorm:"foreignKey:GameID"`
}

func (Model) TableName() string { return "game" }
