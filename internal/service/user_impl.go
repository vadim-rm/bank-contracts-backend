package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
)

type UserImpl struct {
	repository repository.User
}

func NewUserImpl(repository repository.User) *UserImpl {
	return &UserImpl{repository: repository}
}

func (s *UserImpl) Create(ctx context.Context, input CreateUserInput) (domain.UserId, error) {
	return s.repository.Create(ctx, repository.CreateUserInput{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: input.Password,
	})
}

func (s *UserImpl) Update(ctx context.Context, id domain.UserId, input UpdateUserInput) error {
	return s.repository.Update(ctx, id, repository.UpdateUserInput{
		Name:         input.Name,
		PasswordHash: input.Password,
	})
}
