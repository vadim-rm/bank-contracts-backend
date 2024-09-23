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

	Contracts []Contract `gorm:"many2many:account_contracts"`

	Deleted bool `gorm:"not null"`
}

func (a Account) ToDomain() domain.Account {
	contracts := make([]domain.Contract, 0, len(a.Contracts))
	for _, contract := range a.Contracts {
		contracts = append(contracts, contract.ToDomain())
	}

	return domain.Account{
		Id: domain.AccountId(a.ID),

		CreatedAt:   a.CreatedAt,
		RequestedAt: a.RequestedAt,
		FinishedAt:  a.FinishedAt,
		Status:      domain.AccountStatus(a.Status),
		Number:      (*domain.AccountNumber)(a.Number),

		Creator:   domain.UserId(a.Creator),
		Moderator: (*domain.UserId)(a.Moderator),

		Contracts: contracts,
	}
}
