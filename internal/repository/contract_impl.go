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
		Name:        "Эквайринг с терминалом для малого бизнеса",
		Fee:         1000,
		Description: "Договор для принятия платежей оффлайн",
		Type:        "acquiring",
	},
	{
		Id:          2,
		ImageUrl:    "http://localhost:9000/main/2.png",
		Name:        "Интернет-эквайринг для крупного бизнеса",
		Fee:         2000,
		Description: "Договор для принятия платежей на Ваш счёт",
		Type:        "acquiring",
	},
	{
		Id:          3,
		ImageUrl:    "http://localhost:9000/main/3.png",
		Name:        "Расчётный счёт",
		Fee:         350,
		Description: "Договор для выполнения расчётов с другими организациями",
		Type:        "account",
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
