package game

import (
	"gorm.io/gorm"

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

func (*Model) TableName() string { return "game" }

func (m *Model) AfterDelete(tx *gorm.DB) error {
	return tx.Exec(
		`UPDATE game_result SET deleted_at = NOW() WHERE game_id = ? AND deleted_at IS NULL`,
		m.ID,
	).Error
}
