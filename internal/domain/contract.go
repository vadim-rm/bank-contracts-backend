package domain

type ContractId int

type Contract struct {
	Id          ContractId
	Name        string
	AnnualRate  uint8
	Description string
	ImageId     int
}
