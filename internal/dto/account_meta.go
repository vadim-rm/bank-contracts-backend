package dto

import "github.com/vadim-rm/bmstu-web-backend/internal/domain"

type Account struct {
	Id    domain.AccountId
	Count int
}
