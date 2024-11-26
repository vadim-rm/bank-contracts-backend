package domain

import "time"

type AccountId int

type AccountStatus string

const (
	AccountStatusDraft     AccountStatus = "draft"
	AccountStatusDeleted   AccountStatus = "deleted"
	AccountStatusApplied   AccountStatus = "applied"
	AccountStatusFinalized AccountStatus = "finalized"
	AccountStatusRejected  AccountStatus = "rejected"
)

type AccountNumber string

type Account struct {
	Id AccountId

	CreatedAt   time.Time
	RequestedAt *time.Time
	FinishedAt  *time.Time
	Status      AccountStatus
	Number      *AccountNumber

	Creator       UserId
	Moderator     *UserId
	CreatorUser   User
	ModeratorUser *User

	TotalFee *int32

	Contracts []AccountContract
}
