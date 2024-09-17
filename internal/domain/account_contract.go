package domain

type AccountContract struct {
	Id       ContractId
	Name     string
	Fee      int
	ImageUrl string
	Type     ContractType
	IsMain   bool
}
