package models

import "time"

type Model struct {
	ID        int        `json:"id"         db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-"          db:"deleted_at"`
}
