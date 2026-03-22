package pod

import (
	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	playerPodRoleRepo "github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
)

type Model struct {
	base.GormModelBase
	Name    string
	Members []playerPodRoleRepo.Model `gorm:"foreignKey:PodID"`
}

func (*Model) TableName() string { return "pod" }

func (m *Model) AfterDelete(tx *gorm.DB) error {
	stmts := []string{
		`UPDATE player_pod SET deleted_at = NOW() WHERE pod_id = ? AND deleted_at IS NULL`,
		`UPDATE player_pod_role SET deleted_at = NOW() WHERE pod_id = ? AND deleted_at IS NULL`,
		`UPDATE pod_invite SET deleted_at = NOW() WHERE pod_id = ? AND deleted_at IS NULL`,
	}
	for _, s := range stmts {
		if err := tx.Exec(s, m.ID).Error; err != nil {
			return err
		}
	}
	return nil
}

type PlayerPodModel struct {
	base.GormModelBase
	PodID    int
	PlayerID int
}

func (*PlayerPodModel) TableName() string { return "player_pod" }

func (m *PlayerPodModel) AfterDelete(tx *gorm.DB) error {
	return tx.Exec(
		`UPDATE player_pod_role SET deleted_at = NOW()
         WHERE pod_id = ? AND player_id = ? AND deleted_at IS NULL`,
		m.PodID, m.PlayerID,
	).Error
}
