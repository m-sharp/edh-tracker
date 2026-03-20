package deck

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.GormModelBase
	PlayerID int
	Name     string
	FormatID int
	Retired  bool
}

func (Model) TableName() string { return "deck" }

// UpdateFields holds the optional fields that may be updated on a deck.
// Only non-nil fields are applied.
type UpdateFields struct {
	Name     *string
	FormatID *int
	Retired  *bool
}
