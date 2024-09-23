package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"time"
)

type Account interface {
	Create(ctx context.Context, input CreateAccountInput) (domain.AccountId, error)
	GetById(ctx context.Context, id domain.AccountId) (domain.Account, error)
	GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error)
}

type CreateAccountInput struct {
	CreatedAt time.Time
	Status    domain.AccountStatus
	Creator   domain.UserId
}
