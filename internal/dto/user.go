package dto

import "github.com/vadim-rm/bmstu-web-backend/internal/domain"

type User struct {
	ID          domain.UserId
	Email       string
	IsModerator bool
}
