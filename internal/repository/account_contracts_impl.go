package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository/entity"
	"gorm.io/gorm"
)

type AccountContractsImpl struct {
	db *gorm.DB
}

func NewAccountContractsImpl(db *gorm.DB) *AccountContractsImpl {
	return &AccountContractsImpl{db: db}
}

func (r *AccountContractsImpl) RemoveContractFromAccount(ctx context.Context, contractId domain.ContractId, accountId domain.AccountId) error {
	return r.db.WithContext(ctx).
		Where("account_id = ? AND contract_id = ?", accountId, contractId).
		Delete(entity.AccountContracts{}).Error
}

func (r *AccountContractsImpl) SetMain(ctx context.Context, contractId domain.ContractId, accountId domain.AccountId) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := r.db.WithContext(ctx).Model(entity.AccountContracts{}).
			Where("account_id = ? AND contract_id = ?", accountId, contractId).
			Updates(entity.AccountContracts{IsMain: true}).Error
		if err != nil {
			return err
		}

		err = r.db.WithContext(ctx).Model(entity.AccountContracts{}).
			Where("account_id = ? AND contract_id != ?", accountId, contractId).
			Select("IsMain").
			Updates(entity.AccountContracts{IsMain: false}).Error
		if err != nil {
			return err
		}

		return nil
	}, nil)
}
