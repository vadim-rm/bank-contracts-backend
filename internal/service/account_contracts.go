package service

import (
	"context"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
)

type AccountContracts interface {
	RemoveContractFromAccount(ctx context.Context, contractId domain.ContractId, accountId domain.AccountId) error
	SetMain(ctx context.Context, contractId domain.ContractId, accountId domain.AccountId) error
}
