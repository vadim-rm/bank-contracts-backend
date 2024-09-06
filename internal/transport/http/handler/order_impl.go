package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"net/http"
)

type OrderImpl struct {
	service service.Order
}

func NewOrderImpl(service service.Order) *OrderImpl {
	return &OrderImpl{
		service: service,
	}
}

type getOrderByIdRequest struct {
	Id int `uri:"id"`
}

func (h *OrderImpl) GetById(ctx *gin.Context) {
	var request getOrderByIdRequest
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	order, err := h.service.GetById(ctx, domain.OrderId(request.Id))
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

	ctx.HTML(http.StatusOK, "order.gohtml", gin.H{
		"Id":        order.Id,
		"Contracts": order.Contracts,
	})
}
