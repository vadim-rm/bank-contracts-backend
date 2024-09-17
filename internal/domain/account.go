package domain

type AccountId int

type Account struct {
	Id        AccountId
	Contracts []AccountContract
}
