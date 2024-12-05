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

func (r *AccountImpl) GetList(ctx context.Context, filter GetListInput) ([]domain.Account, error) {
	var dbAccounts []entity.Account

	query := r.db.Preload("CreatorUser").Preload("ModeratorUser").WithContext(ctx).Where(
		"status != ? AND status != ?",
		domain.AccountStatusDeleted, domain.AccountStatusDraft,
	).Order("id DESC")

	if filter.CreatorId != nil {
		query = query.Where("creator = ? OR moderator = ?", *filter.CreatorId, *filter.CreatorId)
	}

	if filter.Status != nil {
		query = query.Where(map[string]any{
			"status": *filter.Status,
		})
	}

	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}

	if filter.To != nil {
		query = query.Where("created_at <= ?", *filter.To)
	}

	if err := query.Find(&dbAccounts).Error; err != nil {
		return nil, err
	}

	accounts := make([]domain.Account, 0)
	for _, account := range dbAccounts {
		accounts = append(accounts, account.ToDomain())
	}

	return accounts, nil
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
	Fee         int32
	Description *string
	ImageUrl    *string
	Type        string

	Deleted bool
	IsMain  bool
}

func (r *AccountImpl) Get(ctx context.Context, id domain.AccountId) (domain.Account, error) {
	var account entity.Account

	err := r.db.WithContext(ctx).Preload("CreatorUser").Preload("ModeratorUser").Where(entity.Account{ID: uint(id)}).
		Where("status != ?", domain.AccountStatusDeleted).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Account{}, domain.ErrNotFound
		}
		return domain.Account{}, err
	}
	var contracts []contractWithIsMain
	err = r.db.WithContext(ctx).Table("contracts").
		Select("contracts.id as id, contracts.name, contracts.fee, contracts.description, contracts.image_url, contracts.type, account_contracts.is_main").
		Joins("JOIN account_contracts ON account_contracts.contract_id = contracts.id").
		Joins("JOIN accounts ON account_contracts.account_id = accounts.id").
		Where("account_contracts.account_id = ? AND accounts.status <> ?", id, domain.AccountStatusDeleted).
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
			Type:        domain.ContractType(c.Type),
			IsMain:      c.IsMain,
		}
		accountContracts = append(accountContracts, accountContract)
	}

	domainAccount := domain.Account{
		Id:          domain.AccountId(account.ID),
		CreatedAt:   account.CreatedAt,
		RequestedAt: account.RequestedAt,
		FinishedAt:  account.FinishedAt,
		Number:      (*domain.AccountNumber)(account.Number),
		Status:      domain.AccountStatus(account.Status),
		Creator:     domain.UserId(account.Creator),
		CreatorUser: account.CreatorUser.ToDomain(),
		Contracts:   accountContracts,
		TotalFee:    account.TotalFee,
	}

	if account.ModeratorUser != nil {
		u := account.ModeratorUser.ToDomain()
		domainAccount.ModeratorUser = &u
	}

	return domainAccount, nil
}

func (r *AccountImpl) GetCurrentDraft(ctx context.Context, userId domain.UserId) (dto.Account, error) {
	var account entity.Account

	err := r.db.WithContext(ctx).
		Where(entity.Account{
			Creator: uint(userId),
		}).
		Where("status = ?", domain.AccountStatusDraft).
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
		Id:      domain.AccountId(account.ID),
		Count:   int(count),
		Creator: domain.UserId(account.Creator),
	}, nil
}

func (r *AccountImpl) Update(ctx context.Context, id domain.AccountId, input UpdateAccountInput) error {
	account := entity.Account{
		ID: uint(id),
	}

	updateColumns := make([]string, 0)
	updateValues := make(map[string]any)

	if input.RequestedAt != nil {
		updateColumns = append(updateColumns, "RequestedAt")
		updateValues["RequestedAt"] = *input.RequestedAt
	}

	if input.FinishedAt != nil {
		updateColumns = append(updateColumns, "FinishedAt")
		updateValues["FinishedAt"] = *input.FinishedAt
	}

	if input.Status != nil {
		updateColumns = append(updateColumns, "Status")
		updateValues["Status"] = *input.Status
	}

	if input.Number != nil {
		updateColumns = append(updateColumns, "Number")
		updateValues["Number"] = *input.Number
	}

	if input.Moderator != nil {
		updateColumns = append(updateColumns, "Moderator")
		updateValues["Moderator"] = *input.Moderator
	}

	if input.TotalFee != nil {
		updateColumns = append(updateColumns, "TotalFee")
		updateValues["TotalFee"] = *input.TotalFee
	}

	err := r.db.WithContext(ctx).Model(&account).Select(updateColumns).Updates(updateValues).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountImpl) Delete(ctx context.Context, id domain.AccountId) error {
	return r.db.WithContext(ctx).Exec("UPDATE accounts SET status = ? WHERE id = ?", domain.AccountStatusDeleted, id).Error
}
