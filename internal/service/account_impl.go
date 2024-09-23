package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
)

type AccountImpl struct {
	repository repository.Account
}

func NewAccountImpl(repository repository.Account) *AccountImpl {
	return &AccountImpl{
		repository: repository,
	}
}

func (s *AccountImpl) GetById(ctx context.Context, id domain.AccountId) (domain.Account, error) {
	return s.repository.GetById(ctx, id)
}

func (s *AccountImpl) GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error) {
	return s.repository.GetCurrentDraft(ctx, userId)
}

func (s *AccountImpl) Delete(ctx context.Context, id domain.AccountId) error {
	return s.repository.Delete(ctx, id)
}
