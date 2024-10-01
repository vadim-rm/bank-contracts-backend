package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"io"
)

type Contract interface {
	GetList(ctx context.Context, filter dto.ContractsFilter) ([]domain.Contract, error)
	Get(ctx context.Context, id domain.ContractId) (domain.Contract, error)
	Create(ctx context.Context, input CreateContractInput) (domain.ContractId, error)
	Update(ctx context.Context, id domain.ContractId, input UpdateContractInput) error
	Delete(ctx context.Context, id domain.ContractId) error
	AddToCurrentDraft(ctx context.Context, id domain.ContractId) error
	UpdateImage(ctx context.Context, id domain.ContractId, input UpdateContractImageInput) error
}

type CreateContractInput struct {
	Name        string
	Fee         int32
	Description *string
	ImageUrl    *string
	Type        domain.ContractType
}

type UpdateContractInput struct {
	Name        *string
	Fee         *int32
	Description *string
	Type        *domain.ContractType
}

type UpdateContractImageInput struct {
	Image       io.Reader
	Size        int64
	ContentType string
}
