package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
)

type OrderImpl struct {
	repository repository.Order
}

func NewOrderImpl(repository repository.Order) *OrderImpl {
	return &OrderImpl{
		repository: repository,
	}
}

func (s *OrderImpl) GetById(ctx context.Context, id domain.OrderId) (domain.Order, error) {
	return s.repository.GetById(ctx, id)
}

func (s *OrderImpl) GetCurrentDraft(ctx context.Context) (domain.OrderMeta, error) {
	return s.repository.GetCurrentDraft(ctx)
}
