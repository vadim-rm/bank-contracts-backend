package domain

type ContractId int

type ContractType string

type Contract struct {
	Id          ContractId
	Name        string
	Fee         int
	Description string
	ImageUrl    string
	Type        ContractType
}
