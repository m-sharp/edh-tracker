package deck

import (
	"gorm.io/gorm"

	"github.com/m-sharp/edh-tracker/lib/repositories/base"
	deckCommanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	formatRepo "github.com/m-sharp/edh-tracker/lib/repositories/format"
	playerRepo "github.com/m-sharp/edh-tracker/lib/repositories/player"
)

type Model struct {
	base.GormModelBase
	Name    string
	Retired bool

	PlayerID int
	Player   playerRepo.Model

	FormatID int
	Format   formatRepo.Model

	Commander *deckCommanderRepo.Model `gorm:"foreignKey:DeckID"`
}

func (*Model) TableName() string { return "deck" }

func (m *Model) AfterDelete(tx *gorm.DB) error {
	return tx.Exec(
		`UPDATE deck_commander SET deleted_at = NOW() WHERE deck_id = ? AND deleted_at IS NULL`,
		m.ID,
	).Error
}

// UpdateFields holds the optional fields that may be updated on a deck.
// Only non-nil fields are applied.
type UpdateFields struct {
	Name     *string
	FormatID *int
	Retired  *bool
}

func (u UpdateFields) HasChanges() bool {
	return u.Name != nil || u.FormatID != nil || u.Retired != nil
}
