package domain

type AccountId string

type Account struct {
	Id        AccountId
	Contracts []Contract
}
