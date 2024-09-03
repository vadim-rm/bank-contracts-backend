package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
)

type ContractImpl struct {
}

func NewContractImpl() *ContractImpl {
	return &ContractImpl{}
}

func (r *ContractImpl) GetList(ctx context.Context) ([]domain.ContractMeta, error) {
	return []domain.ContractMeta{
		{
			Id:         1,
			Name:       "Вклад Это норм",
			AnnualRate: 22,
			ImageUrl:   "http://localhost:9000/main/norm.jpg",
		},
		{
			Id:         2,
			Name:       "Вклад Фиговый",
			AnnualRate: 2,
			ImageUrl:   "http://localhost:9000/main/bad.jpg",
		},
		{
			Id:         3,
			Name:       "Вклад Ништяк",
			AnnualRate: 17,
			ImageUrl:   "http://localhost:9000/main/nishtyak.jpg",
		},
	}, nil
}

func (r *ContractImpl) GetById(ctx context.Context, id domain.ContractId) (domain.Contract, error) {
	return domain.Contract{
		Id:          1,
		Name:        "Вклад Это норм",
		AnnualRate:  22,
		ImageUrl:    "http://localhost:9000/main/norm.jpg",
		Description: "Очень выгодный вклад",
	}, nil
}
