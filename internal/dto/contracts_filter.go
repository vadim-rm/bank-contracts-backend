package dto

import "github.com/vadim-rm/bmstu-web-backend/internal/domain"

type ContractsFilter struct {
	Name string
	Type *domain.ContractType
}