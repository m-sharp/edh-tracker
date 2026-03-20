package commander

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.GormModelBase
	Name string `gorm:"column:name"`
}

func (Model) TableName() string { return "commander" }
