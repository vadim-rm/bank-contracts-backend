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

// Delete
// @Summary Удаление договора из заявки на счёт
// @Description Удаляет указанный договор из заявки на счёт по их ID.
// @Tags account-contracts
// @Accept  json
// @Produce  json
// @Param accountId path int true "ID заявки на счёт"
// @Param contractId path int true "ID договора"
// @Success 200 "Договор успешно удалён из заявки на счёт"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Заявка на счёт или договор не найдены"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /accounts/{accountId}/contract/{contractId} [delete]
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

// SetMain
// @Summary Установить договор основным в заявке на счёт
// @Description Устанавливает указанный договор в качестве основного для заявки на счёт по их ID.
// @Tags account-contracts
// @Accept  json
// @Produce  json
// @Param accountId path int true "ID заявки на счёт"
// @Param contractId path int true "ID договора"
// @Success 200 "Договор успешно установлен как основной в заявке на счёт"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Заявка на счёт или договор не найдены"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /accounts/{accountId}/contract/{contractId}/main [put]
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
