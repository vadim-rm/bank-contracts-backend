package repository

import (
	"context"
	"errors"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository/entity"
	"gorm.io/gorm"
)

type UserImpl struct {
	db *gorm.DB
}

func NewUserImpl(db *gorm.DB) *UserImpl {
	return &UserImpl{db: db}
}

func (r *UserImpl) Create(ctx context.Context, input CreateUserInput) (domain.UserId, error) {
	user := entity.User{
		Login:        input.Login,
		PasswordHash: input.PasswordHash,
	}

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return 0, err
	}

	return domain.UserId(user.ID), nil
}

func (r *UserImpl) Get(ctx context.Context, login string) (domain.User, error) {
	var user entity.User

	err := r.db.WithContext(ctx).Where(entity.User{Login: login}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	return user.ToDomain(), nil
}

func (r *UserImpl) Update(ctx context.Context, id domain.UserId, input UpdateUserInput) error {
	updateColumns := make([]string, 0)
	updateValues := make(map[string]any)

	if input.PasswordHash != nil {
		updateColumns = append(updateColumns, "PasswordHash")
		updateValues["PasswordHash"] = *input.PasswordHash
	}

	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Select(updateColumns).Updates(updateValues).Error
	if err != nil {
		return err
	}

	return nil
}
