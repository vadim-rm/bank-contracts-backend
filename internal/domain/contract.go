package domain

type ContractId int

type ContractType string

type Contract struct {
	Id          ContractId
	Name        string
	AnnualRate  uint8
	Description string
	ImageUrl    string
	Type        ContractType
}
