package dto

import "github.com/vadim-rm/bmstu-web-backend/internal/domain"

type TokenClaims struct {
	UserId      domain.UserId
	IsModerator bool
}
