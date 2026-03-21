package game

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.GormModelBase
	Description string
	PodID       int
	FormatID    int
}

func (Model) TableName() string { return "game" }
