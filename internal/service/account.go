package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

type Account interface {
	GetById(ctx context.Context, id domain.AccountId) (domain.Account, error)
	GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error)
}
