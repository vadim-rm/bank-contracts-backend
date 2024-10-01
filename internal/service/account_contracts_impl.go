package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
)

type AccountContractsImpl struct {
	accountContractsRepository repository.AccountContracts
	accountsRepository         repository.Account
}

func NewAccountContractsImpl(
	accountContractsRepository repository.AccountContracts,
	accountsRepository repository.Account,
) *AccountContractsImpl {
	return &AccountContractsImpl{
		accountContractsRepository: accountContractsRepository,
		accountsRepository:         accountsRepository,
	}
}

func (s *AccountContractsImpl) RemoveContractFromAccount(ctx context.Context, contractId domain.ContractId, accountId domain.AccountId) error {
	account, err := s.accountsRepository.Get(ctx, accountId)
	if err != nil {
		return err
	}

	if account.Status != domain.AccountStatusDraft {
		return domain.ErrWrongAccountStatus
	}

	if len(account.Contracts) == 1 {
		return s.accountsRepository.Delete(ctx, accountId)
	}

	return s.accountContractsRepository.RemoveContractFromAccount(ctx, contractId, accountId)
}

func (s *AccountContractsImpl) SetMain(ctx context.Context, contractId domain.ContractId, accountId domain.AccountId) error {
	return s.accountContractsRepository.SetMain(ctx, contractId, accountId)
}
