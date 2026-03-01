package player

import "time"

// TODO: Make sure the generic "BaseModel" with ID, CreatedAt, UpdatedAt, and DeletedAt is embedded
type Model struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
