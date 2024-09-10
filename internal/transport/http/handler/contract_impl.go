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
	contractService service.Contract
	orderService    service.Account
}

func NewContractImpl(
	contractService service.Contract,
	orderService service.Account,
) *ContractImpl {
	return &ContractImpl{
		contractService: contractService,
		orderService:    orderService,
	}
}

type getListOfContractsRequest struct {
	ContractNameFilter string               `form:"contractName"`
	ContractTypeFilter *domain.ContractType `form:"contractType,omitempty"`
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

	contracts, err := h.contractService.GetList(ctx, dto.ContractsFilter{
		Name: request.ContractNameFilter,
		Type: request.ContractTypeFilter,
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

	account, err := h.orderService.GetCurrentDraft(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	filter := ""
	if request.ContractTypeFilter != nil {
		filter = string(*request.ContractTypeFilter)
	}

	ctx.HTML(http.StatusOK, "contracts.gohtml", gin.H{
		"Contracts": contracts,
		"Account": gin.H{
			"Id":    account.Id,
			"Count": account.Count,
		},
		"ContractNameFilter": request.ContractNameFilter,
		"ContractTypeFilter": filter,
	})
}

type getContractByIdRequest struct {
	Id int `uri:"id"`
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

	contract, err := h.contractService.GetById(ctx, domain.ContractId(request.Id))
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.HTML(http.StatusOK, "contract.gohtml", contract)
}
