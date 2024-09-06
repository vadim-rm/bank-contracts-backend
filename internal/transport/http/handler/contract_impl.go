package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
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

type getListOfContractsRequest struct {
	Query string `form:"query"`
}

func (h *ContractImpl) GetList(ctx *gin.Context) {
	var request getListOfContractsRequest
	if err := ctx.BindQuery(&request); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	contracts, err := h.service.GetList(ctx, dto.ContractsFilter{
		Name: request.Query,
	})
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
		"Contracts": contracts,
		"Cart": gin.H{
			"Id":    1,
			"Count": 5,
		},
		"Query": request.Query,
	})
}

type getContractByIdRequest struct {
	Id   int    `uri:"id"`
	From string `form:"from"`
}

func (h *ContractImpl) GetById(ctx *gin.Context) {
	var request getContractByIdRequest
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	if err := ctx.BindQuery(&request); err != nil {
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
		"Contract": contract,
		"Cart": gin.H{
			"Id":    1,
			"Count": 5,
		},
		"From": request.From,
	})
}
