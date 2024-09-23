package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
	"time"
)

type ContractImpl struct {
	contractRepository repository.Contract
	accountRepository  repository.Account
}

func NewContractImpl(
	contractRepository repository.Contract,
	accountRepository repository.Account,
) *ContractImpl {
	return &ContractImpl{
		contractRepository: contractRepository,
		accountRepository:  accountRepository,
	}
}

func (s *ContractImpl) GetList(ctx context.Context, filter dto.ContractsFilter) ([]domain.Contract, error) {
	return s.contractRepository.GetList(ctx, filter)
}

func (s *ContractImpl) GetById(ctx context.Context, id domain.ContractId) (domain.Contract, error) {
	return s.contractRepository.GetById(ctx, id)
}

func (s *ContractImpl) AddToCurrentDraft(ctx context.Context, id domain.ContractId) error {
	account, err := s.getOrCreateAccount(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving account: %w", err)
	}

	err = s.contractRepository.AddToAccount(ctx, repository.AddToAccountInput{
		AccountId:  account.Id,
		ContractId: id,
		IsMain:     account.Count == 0,
	})
	if err != nil {
		return fmt.Errorf("error adding contract to account: %w", err)
	}
	return nil
}

func (s *ContractImpl) getOrCreateAccount(ctx context.Context) (dto.Account, error) {
	draft, err := s.accountRepository.GetCurrentDraft(ctx, 0)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			accountId, err := s.accountRepository.Create(ctx, repository.CreateAccountInput{
				CreatedAt: time.Now(),
				Status:    domain.AccountStatusDraft,
				Creator:   0,
			})
			if err != nil {
				return dto.Account{}, fmt.Errorf("error creating new draft account: %w", err)
			}

			return dto.Account{Id: accountId}, nil
		}
		return dto.Account{}, fmt.Errorf("error getting draft account: %w", err)
	}

	return draft, nil
}
