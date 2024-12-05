package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
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
	Status *string    `form:"status,omitempty"`
	From   *time.Time `form:"from,omitempty"`
	To     *time.Time `form:"to,omitempty"`
}

type getAccountListResponse struct {
	Accounts []getAccountsListAccount `json:"accounts"`
}

type getAccountsListAccount struct {
	Id int `json:"id"`

	CreatedAt   time.Time  `json:"createdAt"`
	RequestedAt *time.Time `json:"requestedAt"`
	FinishedAt  *time.Time `json:"finishedAt"`
	Status      string     `json:"status"`
	Number      *string    `json:"number"`

	Creator   string  `json:"creator"`
	Moderator *string `json:"moderator"`

	TotalFee *int32 `json:"totalFee"`
}

// GetList
// @Summary Получение списка заявок на счёт
// @Description Возвращает список всех заявок на счёт с возможностью фильтрации по статусу и дате.
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param status query string false "Фильтр по статуса"
// @Param from query string false "Фильтр по дате"
// @Param to query string false "Фильтр по дате"
// @Success 200 {object} getAccountListResponse "Список заявок на счёт"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /accounts [get]
func (h *AccountImpl) GetList(ctx *gin.Context) {
	var request getAccountListRequest
	if err := ctx.BindQuery(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	accounts, err := h.service.GetList(ctx, dto.AccountsFilter{
		From:   request.From,
		To:     request.To,
		Status: (*domain.AccountStatus)(request.Status),
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	response := make([]getAccountsListAccount, 0, len(accounts))
	for _, account := range accounts {
		accountResponse := getAccountsListAccount{
			Id:          int(account.Id),
			CreatedAt:   account.CreatedAt,
			RequestedAt: account.RequestedAt,
			FinishedAt:  account.FinishedAt,
			Status:      string(account.Status),
			Number:      (*string)(account.Number),
			Creator:     account.CreatorUser.Login,

			TotalFee: account.TotalFee,
		}

		if account.ModeratorUser != nil {
			accountResponse.Moderator = &account.ModeratorUser.Login
		}
		response = append(response, accountResponse)
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

	CreatedAt   time.Time  `json:"createdAt"`
	RequestedAt *time.Time `json:"requestedAt"`
	FinishedAt  *time.Time `json:"finishedAt"`
	Status      string     `json:"status"`
	Number      *string    `json:"number"`

	Creator   string  `json:"creator"`
	Moderator *string `json:"moderator"`

	TotalFee *int32 `json:"totalFee"`

	Contracts []accountContractResponse `json:"contracts"`
}

type accountContractResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Fee         int32   `json:"fee"`
	Description *string `json:"description"`
	ImageUrl    *string `json:"imageUrl"`
	Type        string  `json:"type"`

	IsMain bool `json:"isMain"`
}

// Get
// @Summary Получение информации о заявке на счёт
// @Description Возвращает подробную информацию о конкретной заявке на счёт по её ID, включая связанные договоры.
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param accountId path int true "ID заявки на счёт"
// @Success 200 {object} getAccountResponse "Детали заявки на счёт"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Заявка на счёт не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /accounts/{accountId} [get]
func (h *AccountImpl) Get(ctx *gin.Context) {
	var request getAccountRequest
	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	account, err := h.service.Get(ctx, domain.AccountId(request.Id))
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	contracts := make([]accountContractResponse, 0, len(account.Contracts))
	for _, contract := range account.Contracts {
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

	accountResponse := getAccountResponse{
		Id:          int(account.Id),
		CreatedAt:   account.CreatedAt,
		RequestedAt: account.RequestedAt,
		FinishedAt:  account.FinishedAt,
		Status:      string(account.Status),
		Number:      (*string)(account.Number),
		Creator:     account.CreatorUser.Login,
		Contracts:   contracts,

		TotalFee: account.TotalFee,
	}

	if account.ModeratorUser != nil {
		accountResponse.Moderator = &account.ModeratorUser.Login
	}

	ctx.JSON(http.StatusOK, accountResponse)
}

type updateAccountRequest struct {
	Id     int    `uri:"accountId"`
	Number string `json:"number"`
}

// Update
// @Summary Обновление информации о заявке на счёт
// @Description Обновляет номер заявки на счёт по её ID.
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param accountId path int true "ID заявки на счёт"
// @Param request body updateAccountRequest true "Данные для обновления заявки"
// @Success 200 {object} getAccountResponse
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Заявка на счёт не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /accounts/{accountId} [put]
func (h *AccountImpl) Update(ctx *gin.Context) {
	var request updateAccountRequest

	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	err := h.service.Update(ctx, domain.AccountId(request.Id), service.UpdateAccountInput{
		Number: domain.AccountNumber(request.Number),
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	account, err := h.service.Get(ctx, domain.AccountId(request.Id))
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	contracts := make([]accountContractResponse, 0, len(account.Contracts))
	for _, contract := range account.Contracts {
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

	accountResponse := getAccountsListAccount{
		Id:          int(account.Id),
		CreatedAt:   account.CreatedAt,
		RequestedAt: account.RequestedAt,
		FinishedAt:  account.FinishedAt,
		Status:      string(account.Status),
		Number:      (*string)(account.Number),
		Creator:     account.CreatorUser.Login,

		TotalFee: account.TotalFee,
	}

	if account.ModeratorUser != nil {
		accountResponse.Moderator = &account.ModeratorUser.Login
	}

	ctx.JSON(http.StatusOK, accountResponse)
}

type submitAccountRequest struct {
	Id int `uri:"accountId"`
}

// Submit
// @Summary Отправка заявки на счёт
// @Description Отправляет заявку на счёт по её ID для дальнейшей обработки.
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param accountId path int true "ID заявки на счёт"
// @Success 200 {object} getAccountResponse
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Заявка на счёт не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /accounts/{accountId}/submit [put]
func (h *AccountImpl) Submit(ctx *gin.Context) {
	var request submitAccountRequest

	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	err := h.service.Submit(ctx, domain.AccountId(request.Id))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	account, err := h.service.Get(ctx, domain.AccountId(request.Id))
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	contracts := make([]accountContractResponse, 0, len(account.Contracts))
	for _, contract := range account.Contracts {
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

	accountResponse := getAccountsListAccount{
		Id:          int(account.Id),
		CreatedAt:   account.CreatedAt,
		RequestedAt: account.RequestedAt,
		FinishedAt:  account.FinishedAt,
		Status:      string(account.Status),
		Number:      (*string)(account.Number),
		Creator:     account.CreatorUser.Login,

		TotalFee: account.TotalFee,
	}

	if account.ModeratorUser != nil {
		accountResponse.Moderator = &account.ModeratorUser.Login
	}

	ctx.JSON(http.StatusOK, accountResponse)
}

type completeAccountRequest struct {
	Id     int    `uri:"accountId"`
	Status string `json:"status"`
}

// Complete
// @Summary Завершение заявки на счёт
// @Description Завершает заявку на счёт, обновляя её статус по ID заявки.
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param accountId path int true "ID заявки на счёт"
// @Param request body completeAccountRequest true "Данные для завершения заявки"
// @Success 200 "Заявка на счёт успешно завершена"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Заявка на счёт не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /accounts/{accountId}/complete [put]
func (h *AccountImpl) Complete(ctx *gin.Context) {
	var request completeAccountRequest

	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
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

// Delete
// @Summary Удаление заявки на счёт
// @Description Удаляет заявку на счёт по её ID.
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param accountId path int true "ID заявки на счёт"
// @Success 200 "Заявка на счёт успешно удалена"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Заявка на счёт не найдена"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /accounts/{accountId} [delete]
func (h *AccountImpl) Delete(ctx *gin.Context) {
	var request deleteRequest
	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	err := h.service.Delete(ctx, domain.AccountId(request.Id))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
