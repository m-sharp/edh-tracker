package deck

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.ModelBase
	PlayerID int    `db:"player_id"`
	Name     string `db:"name"`
	FormatID int    `db:"format_id"`
	Retired  bool   `db:"retired"`
}
