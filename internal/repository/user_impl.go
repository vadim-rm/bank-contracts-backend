package repository

import (
	"context"
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
		Name:         input.Name,
		PasswordHash: input.PasswordHash,
		Email:        input.Email,
	}

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return 0, err
	}

	return domain.UserId(user.ID), nil
}

func (r *UserImpl) Update(ctx context.Context, id domain.UserId, input UpdateUserInput) error {
	updateColumns := make([]string, 0)
	updateValues := make(map[string]any)

	if input.Name != nil {
		updateColumns = append(updateColumns, "Name")
		updateValues["Name"] = *input.Name
	}

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