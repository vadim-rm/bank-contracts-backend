package auth

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

const tokenClaims = "tokenClaims"

func GetClaims(ctx context.Context) (dto.TokenClaims, error) {
	rawClaims := ctx.Value(tokenClaims)

	claims, ok := rawClaims.(dto.TokenClaims)
	if !ok {
		return dto.TokenClaims{}, domain.ErrActionNotPermitted
	}

	return claims, nil
}
