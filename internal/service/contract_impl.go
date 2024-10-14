package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/vadim-rm/bmstu-web-backend/internal/auth"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
	"time"
)

type ContractImpl struct {
	contractRepository repository.Contract
	accountRepository  repository.Account
	imageRepository    repository.Image
}

func NewContractImpl(
	contractRepository repository.Contract,
	accountRepository repository.Account,
	imageRepository repository.Image,
) *ContractImpl {
	return &ContractImpl{
		contractRepository: contractRepository,
		accountRepository:  accountRepository,
		imageRepository:    imageRepository,
	}
}

func (s *ContractImpl) GetList(ctx context.Context, filter dto.ContractsFilter) ([]domain.Contract, error) {
	return s.contractRepository.GetList(ctx, filter)
}

func (s *ContractImpl) Get(ctx context.Context, id domain.ContractId) (domain.Contract, error) {
	return s.contractRepository.Get(ctx, id)
}

func (s *ContractImpl) Create(ctx context.Context, input CreateContractInput) (domain.ContractId, error) {
	return s.contractRepository.Add(ctx, repository.AddContractInput{
		Name:        input.Name,
		Fee:         input.Fee,
		Description: input.Description,
		ImageUrl:    input.ImageUrl,
		Type:        input.Type,
	})
}

func (s *ContractImpl) Update(ctx context.Context, id domain.ContractId, input UpdateContractInput) error {
	return s.contractRepository.Update(ctx, id, repository.UpdateContractInput{
		Name:        input.Name,
		Fee:         input.Fee,
		Description: input.Description,
		Type:        input.Type,
	})
}

func (s *ContractImpl) Delete(ctx context.Context, id domain.ContractId) error {
	return s.contractRepository.Delete(ctx, id)
}

func (s *ContractImpl) AddToCurrentDraft(ctx context.Context, id domain.ContractId) error {
	account, err := s.getOrCreateAccount(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving account: %w", err)
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return fmt.Errorf("error getting claims: %w", err)
	}
	if account.Creator != claims.UserId {
		return domain.ErrActionNotPermitted
	}

	err = s.contractRepository.AddToAccount(ctx, repository.AddToAccountInput{
		AccountId:  account.Id,
		ContractId: id,
		IsMain:     account.Count == 0,
	})
	if err != nil {
		return fmt.Errorf("error adding contract to account: %w", err)
	}
	return nil
}

func (s *ContractImpl) getOrCreateAccount(ctx context.Context) (dto.Account, error) {
	draft, err := s.accountRepository.GetCurrentDraft(ctx, 0)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			claims, err := auth.GetClaims(ctx)
			if err != nil {
				return dto.Account{}, fmt.Errorf("error getting claims: %w", err)
			}

			accountId, err := s.accountRepository.Create(ctx, repository.CreateAccountInput{
				CreatedAt: time.Now(),
				Status:    domain.AccountStatusDraft,
				Creator:   claims.UserId,
			})
			if err != nil {
				return dto.Account{}, fmt.Errorf("error creating new draft account: %w", err)
			}

			return dto.Account{Id: accountId}, nil
		}
		return dto.Account{}, fmt.Errorf("error getting draft account: %w", err)
	}

	return draft, nil
}

func (s *ContractImpl) UpdateImage(ctx context.Context, id domain.ContractId, input UpdateContractImageInput) error {
	imageUrl, err := s.imageRepository.Upload(ctx, repository.UploadImageInput{
		Image:       input.Image,
		Size:        input.Size,
		Name:        fmt.Sprintf("%d", id),
		ContentType: input.ContentType,
	})
	if err != nil {
		return err
	}

	return s.contractRepository.Update(ctx, id, repository.UpdateContractInput{
		ImageUrl: &imageUrl,
	})
}
