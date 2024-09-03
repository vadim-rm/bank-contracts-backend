package domain

type ContractId int

type Contract struct {
	Id          ContractId
	Name        string
	Description string
	AnnualRate  uint8
	ImageUrl    string
}
