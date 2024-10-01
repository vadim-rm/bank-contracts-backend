package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
)

type User interface {
	Create(ctx context.Context, input CreateUserInput) (domain.UserId, error)
	Update(ctx context.Context, id domain.UserId, input UpdateUserInput) error
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserInput struct {
	Name     *string
	Password *string
}
