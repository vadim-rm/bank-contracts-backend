package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

type Contract interface {
	GetList(ctx context.Context, filter dto.ContractsFilter) ([]domain.Contract, error)
	Get(ctx context.Context, id domain.ContractId) (domain.Contract, error)
	Add(ctx context.Context, input AddContractInput) (domain.ContractId, error)
	Update(ctx context.Context, id domain.ContractId, input UpdateContractInput) error
	Delete(ctx context.Context, id domain.ContractId) error
	AddToAccount(ctx context.Context, input AddToAccountInput) error
}

type AddToAccountInput struct {
	AccountId  domain.AccountId
	ContractId domain.ContractId
	IsMain     bool
}

type AddContractInput struct {
	Name        string
	Fee         int32
	Description *string
	ImageUrl    *string
	Type        domain.ContractType
}

type UpdateContractInput struct {
	Name        *string
	Fee         *int32
	Description *string
	ImageUrl    *string
	Type        *domain.ContractType
}
