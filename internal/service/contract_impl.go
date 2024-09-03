package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
)

type ContractImpl struct {
	repository repository.Contract
}

func NewContractImpl(repository repository.Contract) *ContractImpl {
	return &ContractImpl{
		repository: repository,
	}
}

func (s *ContractImpl) GetList(ctx context.Context) ([]domain.ContractMeta, error) {
	return s.repository.GetList(ctx)
}

func (s *ContractImpl) GetById(ctx context.Context, id domain.ContractId) (domain.Contract, error) {
	return s.repository.GetById(ctx, id)
}
