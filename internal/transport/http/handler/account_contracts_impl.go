package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"net/http"
)

type AccountContractsImpl struct {
	service service.AccountContracts
}

func NewAccountContractsImpl(service service.AccountContracts) *AccountContractsImpl {
	return &AccountContractsImpl{service: service}
}

type deleteAccountContractRequest struct {
	AccountId  int `uri:"accountId"`
	ContractId int `uri:"contractId"`
}

func (h *AccountContractsImpl) Delete(ctx *gin.Context) {
	var request deleteAccountContractRequest

	err := ctx.BindUri(&request)
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	err = h.service.RemoveContractFromAccount(ctx, domain.ContractId(request.ContractId), domain.AccountId(request.AccountId))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type setMainAccountContractRequest struct {
	AccountId  int `uri:"accountId"`
	ContractId int `uri:"contractId"`
}

func (h *AccountContractsImpl) SetMain(ctx *gin.Context) {
	var request setMainAccountContractRequest

	err := ctx.BindUri(&request)
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	err = h.service.SetMain(ctx, domain.ContractId(request.ContractId), domain.AccountId(request.AccountId))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
