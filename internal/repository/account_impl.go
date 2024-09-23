package repository

import (
	"context"
	"errors"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository/entity"
	"gorm.io/gorm"
)

type AccountImpl struct {
	db *gorm.DB
}

func NewAccountImpl(db *gorm.DB) *AccountImpl {
	return &AccountImpl{
		db: db,
	}
}

func (r *AccountImpl) Create(ctx context.Context, input CreateAccountInput) (domain.AccountId, error) {
	account := entity.Account{
		CreatedAt: input.CreatedAt,
		Status:    string(input.Status),
		Creator:   uint(input.Creator),
	}
	err := r.db.WithContext(ctx).Select("CreatedAt", "Status", "Creator").Create(&account).Error
	if err != nil {
		return 0, err
	}

	return domain.AccountId(account.ID), nil
}

func (r *AccountImpl) GetById(ctx context.Context, id domain.AccountId) (domain.Account, error) {
	var account entity.Account

	err := r.db.WithContext(ctx).Preload("Contracts").
		Where(entity.Account{
			ID: uint(id),
		}).
		Where("deleted = ?", false).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Account{}, domain.ErrNotFound
		}
	}

	return account.ToDomain(), nil
}

func (r *AccountImpl) GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error) {
	var account entity.Account

	err := r.db.WithContext(ctx).
		Where(entity.Account{
			Creator: uint(userId),
		}).
		Where("deleted = ?", false).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.Account{}, domain.ErrNotFound
		}
		return dto.Account{}, err
	}

	var count int64
	err = r.db.WithContext(ctx).Table("account_contracts").
		Where(entity.AccountContracts{AccountID: account.ID}).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.Account{}, domain.ErrNotFound
		}
		return dto.Account{}, err
	}

	return dto.Account{
		Id:    domain.AccountId(account.ID),
		Count: int(count),
	}, nil
}

func (r *AccountImpl) Delete(ctx context.Context, id domain.AccountId) error {
	return r.db.WithContext(ctx).Exec("UPDATE accounts SET deleted = true WHERE id = ?", id).Error
}
