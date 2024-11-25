package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
)

type User interface {
	Create(ctx context.Context, input CreateUserInput) (domain.UserId, error)
	Get(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, id domain.UserId, input UpdateUserInput) error
}

type CreateUserInput struct {
	Login        string
	PasswordHash string
}

type UpdateUserInput struct {
	PasswordHash *string
}
