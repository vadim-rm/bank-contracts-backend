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

func (s *AccountImpl) GetList(ctx context.Context, id domain.UserId, filter dto.AccountsFilter) ([]domain.Account, error) {
	return s.repository.GetList(ctx, id, filter)
}

func (s *AccountImpl) Get(ctx context.Context, id domain.AccountId) (domain.Account, error) {
	return s.repository.Get(ctx, id)
}

func (s *AccountImpl) Update(ctx context.Context, id domain.AccountId, input UpdateAccountInput) error {
	return s.repository.Update(ctx, id, repository.UpdateAccountInput{
		Number: &input.Number,
	})
}

func (s *AccountImpl) Submit(ctx context.Context, id domain.AccountId) error {
	account, err := s.repository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error loading account: %w", err)
	}

	if account.Number == nil {
		return domain.ErrAccountNumberEmpty
	}

	if account.Status != domain.AccountStatusDraft {
		return domain.ErrInvalidTargetStatus
	}

	status := domain.AccountStatusApplied

	now := time.Now()
	return s.repository.Update(ctx, id, repository.UpdateAccountInput{
		RequestedAt: &now,
		Status:      &status,
	})
}

func (s *AccountImpl) Complete(ctx context.Context, id domain.AccountId, status domain.AccountStatus) error {
	account, err := s.repository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting account: %w", err)
	}

	user := auth.GetModerator()
	if !user.IsModerator {
		return domain.ErrActionNotPermitted
	}

	if !(account.Status == domain.AccountStatusApplied &&
		(status == domain.AccountStatusRejected || status == domain.AccountStatusFinalized)) {
		return domain.ErrInvalidTargetStatus
	}

	now := time.Now()
	return s.repository.Update(ctx, id, repository.UpdateAccountInput{
		Moderator:  &user.ID,
		FinishedAt: &now,
		// todo рассчитывается при завершении заявки (вычисление стоимости заказа, даты доставки в течении месяца, вычисления в м-м).
	})
}

func (s *AccountImpl) GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error) {
	return s.repository.GetCurrentDraft(ctx, userId)
}

func (s *AccountImpl) Delete(ctx context.Context, id domain.AccountId) error {
	return s.repository.Delete(ctx, id)
}
