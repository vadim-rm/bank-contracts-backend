package domain

type ContractId int

type ContractType string

type Contract struct {
	Id          ContractId
	Name        string
	Fee         *int32
	Description *string
	ImageUrl    *string
	Type        *ContractType

	IsMain bool // todo. move to dto
}
