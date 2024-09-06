package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
)

type OrderImpl struct {
}

func NewOrderImpl() *OrderImpl {
	return &OrderImpl{}
}

var orderContracts = contracts[:2]

func (r *OrderImpl) GetById(ctx context.Context, id domain.OrderId) (domain.Order, error) {
	return domain.Order{
		Id:        id,
		Contracts: orderContracts,
	}, nil
}

func (r *OrderImpl) GetCurrentDraft(ctx context.Context) (domain.OrderMeta, error) {
	return domain.OrderMeta{
		Id:    1,
		Count: len(orderContracts),
	}, nil
}
