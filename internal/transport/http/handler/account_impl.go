package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/auth"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"net/http"
	"time"
)

type AccountImpl struct {
	service service.Account
}

func NewAccountImpl(service service.Account) *AccountImpl {
	return &AccountImpl{
		service: service,
	}
}

type getAccountListRequest struct {
	Status *string    `json:"type,omitempty"`
	From   *time.Time `json:"from,omitempty"`
}

type getAccountListResponse struct {
	Accounts []getAccountsListAccount `json:"accounts"`
}

type getAccountsListAccount struct {
	Id int `json:"id"`

	CreatedAt   time.Time  `json:"createdAt"`
	RequestedAt *time.Time `json:"requestedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`
	Status      string     `json:"status"`
	Number      *string    `json:"number,omitempty"`

	Creator   domain.UserId  `json:"creator"`
	Moderator *domain.UserId `json:"moderator,omitempty"`
}

func (h *AccountImpl) GetList(ctx *gin.Context) {
	var request getAccountListRequest
	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	user := auth.GetUser()

	accounts, err := h.service.GetList(ctx, user.ID, dto.AccountsFilter{
		From:   request.From,
		Status: (*domain.AccountStatus)(request.Status),
	})
	if err != nil {
		newErrorResponse(ctx, err)
	}

	response := make([]getAccountsListAccount, 0, len(accounts))
	for _, account := range accounts {
		response = append(response, getAccountsListAccount{
			Id:          int(account.Id),
			CreatedAt:   account.CreatedAt,
			RequestedAt: account.RequestedAt,
			FinishedAt:  account.FinishedAt,
			Status:      string(account.Status),
			Number:      (*string)(account.Number),
			Creator:     account.Creator,
			Moderator:   account.Moderator,
		})
	}

	ctx.JSON(http.StatusOK, getAccountListResponse{
		Accounts: response,
	})
}

type getAccountRequest struct {
	Id int `uri:"accountId"`
}

type getAccountResponse struct {
	Id int `json:"id"`

	CreatedAt   time.Time             `json:"createdAt"`
	RequestedAt *time.Time            `json:"requestedAt,omitempty"`
	FinishedAt  *time.Time            `json:"finishedAt,omitempty"`
	Status      string                `json:"status"`
	Number      *domain.AccountNumber `json:"number,omitempty"`

	Creator   domain.UserId  `json:"creator"`
	Moderator *domain.UserId `json:"moderator"`

	Contracts []accountContractResponse `json:"contracts"`
}

type accountContractResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Fee         int32   `json:"fee"`
	Description *string `json:"description,omitempty"`
	ImageUrl    *string `json:"imageUrl,omitempty"`
	Type        string  `json:"type"`

	IsMain bool `json:"isMain"`
}

func (h *AccountImpl) Get(ctx *gin.Context) {
	var request getAccountRequest
	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	order, err := h.service.Get(ctx, domain.AccountId(request.Id))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	contracts := make([]accountContractResponse, 0, len(order.Contracts))
	for _, contract := range order.Contracts {
		contracts = append(contracts, accountContractResponse{
			Id:          int(contract.Id),
			Name:        contract.Name,
			Fee:         contract.Fee,
			Description: contract.Description,
			ImageUrl:    contract.ImageUrl,
			Type:        string(contract.Type),
			IsMain:      contract.IsMain,
		})
	}

	ctx.JSON(http.StatusOK, getAccountResponse{
		Id:          int(order.Id),
		CreatedAt:   order.CreatedAt,
		RequestedAt: order.RequestedAt,
		FinishedAt:  order.FinishedAt,
		Status:      string(order.Status),
		Number:      order.Number,
		Creator:     order.Creator,
		Moderator:   order.Moderator,
		Contracts:   contracts,
	})
}

type updateAccountRequest struct {
	Id     int    `uri:"accountId"`
	Number string `json:"number"`
}

func (h *AccountImpl) Update(ctx *gin.Context) {
	var request updateAccountRequest

	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	err := h.service.Update(ctx, domain.AccountId(request.Id), service.UpdateAccountInput{
		Number: domain.AccountNumber(request.Number),
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type submitAccountRequest struct {
	Id int `uri:"accountId"`
}

func (h *AccountImpl) Submit(ctx *gin.Context) {
	var request submitAccountRequest

	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	err := h.service.Submit(ctx, domain.AccountId(request.Id))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type completeAccountRequest struct {
	Id     int    `uri:"accountId"`
	Status string `json:"status"`
}

func (h *AccountImpl) Complete(ctx *gin.Context) {
	var request completeAccountRequest

	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	err := h.service.Complete(ctx, domain.AccountId(request.Id), domain.AccountStatus(request.Status))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type deleteRequest struct {
	Id int `uri:"accountId"`
}

func (h *AccountImpl) Delete(ctx *gin.Context) {
	var request deleteRequest
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	err := h.service.Delete(ctx, domain.AccountId(request.Id))
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.Status(http.StatusOK)
}
