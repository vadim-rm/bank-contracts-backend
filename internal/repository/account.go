package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

type Account interface {
	GetById(ctx context.Context, id domain.AccountId) (domain.Account, error)
	GetCurrentDraft(ctx context.Context) (dto.Account, error)
}
