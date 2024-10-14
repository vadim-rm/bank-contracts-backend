package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
)

type UserImpl struct {
	user  repository.User
	token repository.Token
}

func NewUserImpl(
	repository repository.User,
	token repository.Token,
) *UserImpl {
	return &UserImpl{
		user:  repository,
		token: token,
	}
}

func (s *UserImpl) Create(ctx context.Context, input CreateUserInput) (domain.UserId, error) {
	return s.user.Create(ctx, repository.CreateUserInput{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: generateHashString(input.Password),
	})
}

func (s *UserImpl) Update(ctx context.Context, id domain.UserId, input UpdateUserInput) error {
	var hash *string
	if input.Password != nil {
		hashedPassword := generateHashString(*input.Password)
		hash = &hashedPassword
	}

	return s.user.Update(ctx, id, repository.UpdateUserInput{
		Name:         input.Name,
		PasswordHash: hash,
	})
}

func (s *UserImpl) Authenticate(ctx context.Context, input AuthorizeInput) (dto.Token, error) {
	user, err := s.user.Get(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return dto.Token{}, domain.ErrInvalidCredentials
		}
		return dto.Token{}, fmt.Errorf("error getting user: %w", err)
	}

	if user.PasswordHash != generateHashString(input.Password) {
		return dto.Token{}, domain.ErrInvalidCredentials
	}

	token, err := s.token.Create(ctx, dto.TokenClaims{
		UserId:      user.ID,
		IsModerator: user.IsModerator,
	})
	if err != nil {
		return dto.Token{}, fmt.Errorf("error creating token: %w", err)
	}

	return token, err
}

func (s *UserImpl) Logout(ctx context.Context, token string) error {
	return s.token.Blacklist(ctx, token)
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
