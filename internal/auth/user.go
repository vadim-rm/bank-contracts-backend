package auth

import (
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

var user = dto.User{
	ID:          0,
	Email:       "mockemail@example.com",
	IsModerator: true,
}

func GetUser() dto.User {
	return user
}

var moderator = dto.User{
	ID:          1,
	Email:       "mockemailmoderator@example.com",
	IsModerator: true,
}

func GetModerator() dto.User {
	return moderator
}
