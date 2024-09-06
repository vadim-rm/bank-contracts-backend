package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"strings"
)

var contracts = []domain.Contract{
	{
		Id:          1,
		ImageId:     1,
		Name:        "Вклад Это норм",
		AnnualRate:  22,
		Description: "Очень выгодный вклад",
	},
	{
		Id:          2,
		ImageId:     2,
		Name:        "Вклад Фиговый",
		AnnualRate:  2,
		Description: "Невыгодный вклад",
	},
	{
		Id:          3,
		ImageId:     3,
		Name:        "Вклад Ништяк",
		AnnualRate:  17,
		Description: "Немного выгодный вклад",
	},
}

type ContractImpl struct {
}

func NewContractImpl() *ContractImpl {
	return &ContractImpl{}
}

func (r *ContractImpl) GetList(ctx context.Context, filter dto.ContractsFilter) ([]domain.Contract, error) {
	filteredContracts := make([]domain.Contract, 0)
	for _, contract := range contracts {
		if strings.Contains(
			strings.ToLower(contract.Name),
			strings.ToLower(filter.Name),
		) {
			filteredContracts = append(filteredContracts, contract)
		}
	}

	return filteredContracts, nil
}

func (r *ContractImpl) GetById(ctx context.Context, id domain.ContractId) (domain.Contract, error) {
	for _, contract := range contracts {
		if contract.Id == id {
			return contract, nil
		}
	}
	return domain.Contract{}, domain.ErrNotFound
}
