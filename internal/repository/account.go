package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"time"
)

type Account interface {
	GetList(ctx context.Context, id domain.UserId, filter dto.AccountsFilter) ([]domain.Account, error)
	Create(ctx context.Context, input CreateAccountInput) (domain.AccountId, error)
	Get(ctx context.Context, id domain.AccountId) (domain.Account, error)
	Update(ctx context.Context, id domain.AccountId, input UpdateAccountInput) error
	GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error)
	Delete(ctx context.Context, id domain.AccountId) error
}

type CreateAccountInput struct {
	CreatedAt time.Time
	Status    domain.AccountStatus
	Creator   domain.UserId
}

type UpdateAccountInput struct {
	RequestedAt *time.Time
	FinishedAt  *time.Time
	Status      *domain.AccountStatus
	Number      *domain.AccountNumber
	Moderator   *domain.UserId
}
