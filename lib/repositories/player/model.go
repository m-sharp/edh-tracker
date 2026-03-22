package player

import (
	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
)

type Model struct {
	base.GormModelBase
	Name string
}

func (*Model) TableName() string { return "player" }

func (m *Model) AfterDelete(tx *gorm.DB) error {
	stmts := []string{
		`UPDATE deck_commander dc INNER JOIN deck d ON dc.deck_id = d.id
         SET dc.deleted_at = NOW() WHERE d.player_id = ? AND dc.deleted_at IS NULL`,
		`UPDATE deck SET deleted_at = NOW() WHERE player_id = ? AND deleted_at IS NULL`,
		`UPDATE user SET deleted_at = NOW() WHERE player_id = ? AND deleted_at IS NULL`,
		`UPDATE player_pod SET deleted_at = NOW() WHERE player_id = ? AND deleted_at IS NULL`,
		`UPDATE player_pod_role SET deleted_at = NOW() WHERE player_id = ? AND deleted_at IS NULL`,
	}
	for _, s := range stmts {
		if err := tx.Exec(s, m.ID).Error; err != nil {
			return err
		}
	}
	return nil
}
