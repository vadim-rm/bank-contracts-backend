package domain

import "time"

type AccountId int

type AccountStatus string

const (
	AccountStatusDraft AccountStatus = "draft"
)

type AccountNumber string

type Account struct {
	Id AccountId

	CreatedAt   time.Time
	RequestedAt *time.Time
	FinishedAt  *time.Time
	Status      AccountStatus
	Number      *AccountNumber

	Creator   UserId
	Moderator *UserId

	Contracts []Contract
}
