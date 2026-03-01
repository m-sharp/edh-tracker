package deck

import "time"

// TODO: Make sure the generic "BaseModel" with ID, CreatedAt, UpdatedAt, and DeletedAt is embedded
type Model struct {
	ID        int        `db:"id"`
	PlayerID  int        `db:"player_id"`
	Name      string     `db:"name"`
	FormatID  int        `db:"format_id"`
	Retired   bool       `db:"retired"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
