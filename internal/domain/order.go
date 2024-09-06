package domain

type OrderId int

type Order struct {
	Id        OrderId
	Contracts []Contract
}
