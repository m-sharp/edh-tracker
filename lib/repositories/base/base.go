package base

import (
	"time"

	"gorm.io/gorm"
)

// ModelBase is used by sqlx-based repositories (pre-GORM).
type ModelBase struct {
	ID        int        `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// GormModelBase is used by GORM-based repositories.
type GormModelBase struct {
	ID        int            `gorm:"primaryKey;column:id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
