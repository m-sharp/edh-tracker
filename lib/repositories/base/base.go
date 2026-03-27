package base

import (
	"time"

	"gorm.io/gorm"
)

type GormModelBase struct {
	ID        int            `gorm:"primaryKey;column:id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// Repo is the base repository struct that each Repository should embed.
type Repo struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// StartTX sets the underlying gorm.DB to a passed transaction wrapped gorm.DB. Always defer EndTX after calling.
func (r *Repo) StartTX(tx *gorm.DB) {
	r.tx = tx
}

// EndTX should be deferred after calling StartTX. Unsets the underlying gorm.DB so that the regular connection is used.
func (r *Repo) EndTX() {
	r.tx = nil
}

// DB returns the underlying gorm.DB.
func (r *Repo) DB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}

	return r.db
}
