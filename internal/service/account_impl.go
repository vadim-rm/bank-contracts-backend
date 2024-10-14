package service

import (
	"context"
	"fmt"
	"github.com/vadim-rm/bmstu-web-backend/internal/auth"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
	"time"
)

type AccountImpl struct {
	repository repository.Account
}

func NewAccountImpl(repository repository.Account) *AccountImpl {
	return &AccountImpl{
		repository: repository,
	}
}

func (s *AccountImpl) GetList(ctx context.Context, filter dto.AccountsFilter) ([]domain.Account, error) {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting claims: %w", err)
	}

	accountsFilter := repository.GetListInput{
		From:   filter.From,
		Status: filter.Status,
	}

	if !claims.IsModerator {
		accountsFilter.CreatorId = &claims.UserId
	}

	return s.repository.GetList(ctx, accountsFilter)
}

func (s *AccountImpl) Get(ctx context.Context, id domain.AccountId) (domain.Account, error) {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return domain.Account{}, fmt.Errorf("error getting claims: %w", err)
	}

	account, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Account{}, err
	}

	if account.Creator != claims.UserId && !claims.IsModerator {
		return domain.Account{}, domain.ErrActionNotPermitted
	}

	return account, nil
}

func (s *AccountImpl) Update(ctx context.Context, id domain.AccountId, input UpdateAccountInput) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return fmt.Errorf("error getting claims: %w", err)
	}

	account, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	}

	if account.Creator != claims.UserId {
		return domain.ErrActionNotPermitted
	}

	return s.repository.Update(ctx, id, repository.UpdateAccountInput{
		Number: &input.Number,
	})
}

func (s *AccountImpl) Submit(ctx context.Context, id domain.AccountId) error {
	account, err := s.repository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error loading account: %w", err)
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return fmt.Errorf("error getting claims: %w", err)
	}

	if account.Creator != claims.UserId {
		return domain.ErrActionNotPermitted
	}

	if account.Number == nil {
		return domain.ErrAccountNumberEmpty
	}

	if account.Status != domain.AccountStatusDraft {
		return domain.ErrInvalidTargetStatus
	}

	var totalFee int32
	for _, contract := range account.Contracts {
		totalFee += contract.Fee
	}

	status := domain.AccountStatusApplied

	now := time.Now()
	return s.repository.Update(ctx, id, repository.UpdateAccountInput{
		RequestedAt: &now,
		Status:      &status,
		TotalFee:    &totalFee,
	})
}

func (s *AccountImpl) Complete(ctx context.Context, id domain.AccountId, status domain.AccountStatus) error {
	account, err := s.repository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting account: %w", err)
	}

	if !(account.Status == domain.AccountStatusApplied &&
		(status == domain.AccountStatusRejected || status == domain.AccountStatusFinalized)) {
		return domain.ErrInvalidTargetStatus
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return fmt.Errorf("error getting claims: %w", err)
	}

	now := time.Now()
	return s.repository.Update(ctx, id, repository.UpdateAccountInput{
		Moderator:  &claims.UserId,
		FinishedAt: &now,
		Status:     &status,
	})
}

func (s *AccountImpl) GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error) {
	return s.repository.GetCurrentDraft(ctx, userId)
}

func (s *AccountImpl) Delete(ctx context.Context, id domain.AccountId) error {
	return s.repository.Delete(ctx, id)
}
