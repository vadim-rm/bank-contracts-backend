package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
)

type AccountContractsImpl struct {
	repository repository.AccountContracts
}

func NewAccountContractsImpl(repository repository.AccountContracts) *AccountContractsImpl {
	return &AccountContractsImpl{repository: repository}
}

func (s *AccountContractsImpl) RemoveContractFromAccount(ctx context.Context, contractId domain.ContractId, accountId domain.AccountId) error {
	return s.repository.RemoveContractFromAccount(ctx, contractId, accountId)
}

func (s *AccountContractsImpl) SetMain(ctx context.Context, contractId domain.ContractId, accountId domain.AccountId) error {
	return s.repository.SetMain(ctx, contractId, accountId)
}
