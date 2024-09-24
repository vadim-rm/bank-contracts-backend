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

type contractWithIsMain struct {
	ID          uint
	Name        string
	Fee         *int32
	Description *string
	ImageUrl    *string
	Type        *string

	Deleted bool
	IsMain  bool
}

func (r *AccountImpl) GetById(ctx context.Context, id domain.AccountId) (domain.Account, error) {
	var account entity.Account

	err := r.db.WithContext(ctx).Where(entity.Account{ID: uint(id)}).
		Where("deleted = ?", false).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Account{}, domain.ErrNotFound
		}
		return domain.Account{}, err
	}
	var contracts []contractWithIsMain
	err = r.db.WithContext(ctx).Table("contracts").
		Select("contracts.id as contract_id, contracts.name, contracts.fee, contracts.description, contracts.image_url, contracts.type, account_contracts.is_main").
		Joins("JOIN account_contracts ON account_contracts.contract_id = contracts.id").
		Where("account_contracts.account_id = ?", id).
		Order("account_contracts.is_main DESC").
		Scan(&contracts).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Account{}, domain.ErrNotFound
		}
		return domain.Account{}, err
	}

	var accountContracts []domain.AccountContract
	for _, c := range contracts {
		accountContract := domain.AccountContract{
			Id:          domain.ContractId(c.ID),
			Name:        c.Name,
			Fee:         c.Fee,
			Description: c.Description,
			ImageUrl:    c.ImageUrl,
			Type:        (*domain.ContractType)(c.Type),
			IsMain:      c.IsMain,
		}
		accountContracts = append(accountContracts, accountContract)
	}

	return domain.Account{
		Id:          domain.AccountId(account.ID),
		CreatedAt:   account.CreatedAt,
		RequestedAt: account.RequestedAt,
		FinishedAt:  account.FinishedAt,
		Status:      domain.AccountStatus(account.Status),
		Creator:     domain.UserId(account.Creator),
		Moderator:   (*domain.UserId)(account.Moderator),
		Contracts:   accountContracts,
	}, nil
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
