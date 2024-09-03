package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
)

type Contract interface {
	GetList(ctx context.Context) ([]domain.ContractMeta, error)
	GetById(ctx context.Context, id domain.ContractId) (domain.Contract, error)
}
