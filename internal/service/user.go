package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

type User interface {
	Create(ctx context.Context, input CreateUserInput) (domain.UserId, error)
	Update(ctx context.Context, id domain.UserId, input UpdateUserInput) error
	Logout(ctx context.Context, token string) error
	Authenticate(ctx context.Context, input AuthorizeInput) (dto.Token, error)
	Get(ctx context.Context, email string) (domain.User, error)
}

type CreateUserInput struct {
	Login    string
	Password string
}

type UpdateUserInput struct {
	Password *string
}

type AuthorizeInput struct {
	Login    string
	Password string
}
