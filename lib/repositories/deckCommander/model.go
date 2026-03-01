package deckCommander

import "time"

type Model struct {
	ID                 int        `db:"id"`
	DeckID             int        `db:"deck_id"`
	CommanderID        int        `db:"commander_id"`
	PartnerCommanderID *int       `db:"partner_commander_id"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at"`
}
