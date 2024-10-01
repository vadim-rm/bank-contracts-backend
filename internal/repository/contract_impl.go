package repository

import (
	"context"
	"errors"
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

func (r *ContractImpl) Get(ctx context.Context, id domain.ContractId) (domain.Contract, error) {
	var contract entity.Contract

	err := r.db.WithContext(ctx).Where(entity.Contract{ID: uint(id)}).First(&contract).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Contract{}, domain.ErrNotFound
		}
		return domain.Contract{}, err
	}

	return contract.ToDomain(), nil
}

func (r *ContractImpl) Add(ctx context.Context, input AddContractInput) (domain.ContractId, error) {
	contract := entity.Contract{
		Name:        input.Name,
		Fee:         input.Fee,
		Description: input.Description,
		ImageUrl:    input.ImageUrl,
		Type:        string(input.Type),
	}

	err := r.db.WithContext(ctx).Create(&contract).Error
	if err != nil {
		return 0, err
	}

	return domain.ContractId(contract.ID), nil
}

func (r *ContractImpl) Update(ctx context.Context, id domain.ContractId, input UpdateContractInput) error {
	contract := entity.Contract{
		ID: uint(id),
	}

	updateColumns := make([]string, 0)
	updateValues := make(map[string]any)

	if input.Name != nil {
		updateColumns = append(updateColumns, "Name")
		updateValues["Name"] = *input.Name
	}

	if input.Fee != nil {
		updateColumns = append(updateColumns, "Fee")
		updateValues["Fee"] = *input.Fee
	}

	if input.Description != nil {
		updateColumns = append(updateColumns, "Description")
		updateValues["Description"] = *input.Description
	}

	if input.ImageUrl != nil {
		updateColumns = append(updateColumns, "ImageUrl")
		updateValues["ImageUrl"] = *input.ImageUrl
	}

	if input.Type != nil {
		updateColumns = append(updateColumns, "Type")
		updateValues["Type"] = *input.Type
	}

	err := r.db.WithContext(ctx).Model(&contract).Select(updateColumns).Updates(updateValues).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ContractImpl) Delete(ctx context.Context, id domain.ContractId) error {
	return r.db.WithContext(ctx).Delete(&entity.Contract{}, id).Error
}

func (r *ContractImpl) AddToAccount(ctx context.Context, input AddToAccountInput) error {
	accountContract := entity.AccountContracts{
		AccountID:  uint(input.AccountId),
		ContractID: uint(input.ContractId),
		IsMain:     input.IsMain,
	}

	return r.db.WithContext(ctx).Create(&accountContract).Error
}
