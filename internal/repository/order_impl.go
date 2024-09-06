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

func (r *OrderImpl) GetById(ctx context.Context, id domain.OrderId) (domain.Order, error) {
	return domain.Order{
		Id:        id,
		Contracts: contracts[:2],
	}, nil
}
