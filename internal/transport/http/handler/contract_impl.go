package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"net/http"
)

type ContractImpl struct {
	service service.Contract
}

func NewContractImpl(service service.Contract) *ContractImpl {
	return &ContractImpl{
		service: service,
	}
}

func (h *ContractImpl) GetList(ctx *gin.Context) {
	contracts, err := h.service.GetList(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.HTML(http.StatusOK, "contracts.gohtml", gin.H{
		"contracts": contracts,
	})
}

type getByIdRequest struct {
	Id int `uri:"id"`
}

func (h *ContractImpl) GetById(ctx *gin.Context) {
	var request getByIdRequest
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	contract, err := h.service.GetById(ctx, domain.ContractId(request.Id))
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.HTML(http.StatusOK, "contract.gohtml", gin.H{
		"contract": contract,
	})
}
