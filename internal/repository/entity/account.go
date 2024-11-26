package entity

import (
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
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

	CreatorUser   User  `gorm:"foreignKey:Creator"`
	ModeratorUser *User `gorm:"foreignKey:Moderator"`

	TotalFee *int32

	Contracts []Contract `gorm:"many2many:account_contracts"`
}

func (a Account) ToDomain() domain.Account {
	account := domain.Account{
		Id:          domain.AccountId(a.ID),
		CreatedAt:   a.CreatedAt,
		RequestedAt: a.RequestedAt,
		FinishedAt:  a.FinishedAt,
		Status:      domain.AccountStatus(a.Status),
		Number:      (*domain.AccountNumber)(a.Number),
		CreatorUser: a.CreatorUser.ToDomain(),
		TotalFee:    a.TotalFee,
	}

	if a.ModeratorUser != nil {
		u := a.ModeratorUser.ToDomain()
		account.ModeratorUser = &u
	}

	return account
}
