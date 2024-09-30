package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

type Account interface {
	GetList(ctx context.Context, id domain.UserId, filter dto.AccountsFilter) ([]domain.Account, error)
	Get(ctx context.Context, id domain.AccountId) (domain.Account, error)
	Update(ctx context.Context, id domain.AccountId, input UpdateAccountInput) error
	Submit(ctx context.Context, id domain.AccountId) error
	Complete(ctx context.Context, id domain.AccountId, status domain.AccountStatus) error
	Delete(ctx context.Context, id domain.AccountId) error

	GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error)
}

type UpdateAccountInput struct {
	Number domain.AccountNumber
}
