package repository

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
)

type Order interface {
	GetById(ctx context.Context, id domain.OrderId) (domain.Order, error)
}
