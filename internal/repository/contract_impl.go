package repository

import (
	"context"
	"fmt"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository/entity"
	"gorm.io/gorm"
)

type ContractImpl struct {
	db *gorm.DB
}

func NewContractImpl(db *gorm.DB) *ContractImpl {
	return &ContractImpl{
		db: db,
	}
}

func (r *ContractImpl) GetList(ctx context.Context, filter dto.ContractsFilter) ([]domain.Contract, error) {
	var dbContracts []entity.Contract

	query := r.db.WithContext(ctx).Where(
		"name ILIKE ?",
		fmt.Sprintf("%%%s%%", filter.Name),
	)

	if filter.Type != nil {
		query = query.Where(map[string]any{
			"type": *filter.Type,
		})
	}

	if err := query.Find(&dbContracts).Error; err != nil {
		return nil, err
	}

	contracts := make([]domain.Contract, 0)
	for _, contract := range dbContracts {
		contracts = append(contracts, contract.ToDomain())
	}

	return contracts, nil
}

func (r *ContractImpl) GetById(ctx context.Context, id domain.ContractId) (domain.Contract, error) {
	var contract entity.Contract

	err := r.db.WithContext(ctx).Where(entity.Contract{ID: uint(id)}).First(&contract).Error
	if err != nil {
		return domain.Contract{}, err
	}

	return contract.ToDomain(), nil
}
