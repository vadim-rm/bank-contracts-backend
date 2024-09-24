package entity

import (
	"time"
)

type Account struct {
	ID          uint
	CreatedAt   time.Time `gorm:"not null"`
	RequestedAt *time.Time
	FinishedAt  *time.Time
	Status      string  `gorm:"size:20;not null"`
	Number      *string `gorm:"size:20"`

	Creator   uint `gorm:"not null"`
	Moderator *uint

	Contracts []Contract `gorm:"many2many:account_contracts"`

	Deleted bool `gorm:"not null"`
}
