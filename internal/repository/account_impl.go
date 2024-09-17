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
		Id: id,
		Contracts: []domain.AccountContract{
			{

				Id:       2,
				ImageUrl: "http://localhost:9000/main/2.png",
				Name:     "Интернет-эквайринг для крупного бизнеса",
				Fee:      2000,
				Type:     "acquiring",

				IsMain: true,
			},
			{

				Id:       3,
				ImageUrl: "http://localhost:9000/main/3.png",
				Name:     "Расчётный счёт",
				Fee:      350,
				Type:     "account",
			},
		},
	}, nil
}

func (r *AccountImpl) GetCurrentDraft(ctx context.Context) (dto.Account, error) {
	return dto.Account{
		Id:    1,
		Count: 2,
	}, nil
}
