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
		ImageUrl:    "http://localhost:9000/main/1.png",
		Name:        "Вклад Хороший",
		AnnualRate:  22,
		Description: "Очень выгодный вклад",
		Type:        "debit",
	},
	{
		Id:          2,
		ImageUrl:    "http://localhost:9000/main/2.png",
		Name:        "Вклад Долгосрочный",
		AnnualRate:  2,
		Description: "Невыгодный вклад",
		Type:        "credit",
	},
	{
		Id:          3,
		ImageUrl:    "http://localhost:9000/main/3.png",
		Name:        "Вклад Пополняй",
		AnnualRate:  17,
		Description: "Немного выгодный вклад",
		Type:        "credit",
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
		) && (filter.Type == nil || contract.Type == *filter.Type) {
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
