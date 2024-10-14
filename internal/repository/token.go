package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

type Token interface {
	Create(ctx context.Context, input dto.TokenClaims) (dto.Token, error)
	GetClaims(ctx context.Context, token string) (dto.TokenClaims, error)
	Blacklist(ctx context.Context, token string) error
}
