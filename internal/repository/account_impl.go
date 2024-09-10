package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
)

type AccountImpl struct {
}

func NewAccountImpl() *AccountImpl {
	return &AccountImpl{}
}

func (r *AccountImpl) GetById(ctx context.Context, id domain.AccountId) (domain.Account, error) {
	return domain.Account{
		Id: "40817810099910004312",
		Contracts: []domain.Contract{
			{
				Id:          1,
				ImageUrl:    "http://localhost:9000/main/1.png",
				Name:        "Вклад Хороший",
				AnnualRate:  22,
				Description: "Очень выгодный вклад",
			},
			{
				Id:          2,
				ImageUrl:    "http://localhost:9000/main/2.png",
				Name:        "Вклад Долгосрочный",
				AnnualRate:  2,
				Description: "Невыгодный вклад",
			},
		},
	}, nil
}

func (r *AccountImpl) GetCurrentDraft(ctx context.Context) (dto.Account, error) {
	return dto.Account{
		Id:    "40817810099910004312",
		Count: 2,
	}, nil
}
