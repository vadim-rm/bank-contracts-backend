package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"net/http"
)

type AccountImpl struct {
	service service.Account
}

func NewAccountImpl(service service.Account) *AccountImpl {
	return &AccountImpl{
		service: service,
	}
}

type getAccountByIdRequest struct {
	Id int `uri:"id"`
}

func (h *AccountImpl) GetById(ctx *gin.Context) {
	var request getAccountByIdRequest
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	order, err := h.service.Get(ctx, domain.AccountId(request.Id))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{
					"error": err.Error(),
				},
			)
			return
		}
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.HTML(http.StatusOK, "account.gohtml", gin.H{
		"Id":        order.Id,
		"Contracts": order.Contracts,
	})
}

type deleteRequest struct {
	Id int `uri:"id"`
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

	ctx.Redirect(http.StatusSeeOther, "/contracts")
}
