package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

type Contract interface {
	GetList(ctx context.Context, filter dto.ContractsFilter) ([]domain.Contract, error)
	GetById(ctx context.Context, id domain.ContractId) (domain.Contract, error)
	AddToAccount(ctx context.Context, input AddToAccountInput) error
}

type AddToAccountInput struct {
	AccountId  domain.AccountId
	ContractId domain.ContractId
	IsMain     bool
}
