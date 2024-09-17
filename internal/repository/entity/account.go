package entity

import "time"

type Account struct {
	ID          uint
	CreatedAt   time.Time  `gorm:"not null"`
	RequestedAt *time.Time `gorm:"not null"`
	FinishedAt  *time.Time `gorm:"not null"`
	Status      string     `gorm:"size:20;not null"`
	Number      *string    `gorm:"size:20"`

	Creator   uint
	Moderator *uint

	Contracts []Contract `gorm:"many2many:account_contracts"`
}
