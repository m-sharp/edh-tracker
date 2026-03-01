package gameResult

import "time"

type Model struct {
	ID        int        `db:"id"`
	GameID    int        `db:"game_id"`
	DeckID    int        `db:"deck_id"`
	Place     int        `db:"place"`
	KillCount int        `db:"kill_count"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
