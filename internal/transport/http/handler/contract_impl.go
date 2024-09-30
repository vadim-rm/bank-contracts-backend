package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/auth"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"net/http"
)

type ContractImpl struct {
	contractService service.Contract
	accountService  service.Account
}

func NewContractImpl(
	contractService service.Contract,
	orderService service.Account,
) *ContractImpl {
	return &ContractImpl{
		contractService: contractService,
		accountService:  orderService,
	}
}

type getListOfContractsRequest struct {
	ContractNameFilter string               `form:"contractName"`
	ContractTypeFilter *domain.ContractType `form:"contractType,omitempty"`
}

type getListOfContractsResponse struct {
	Contracts []contractResponse `json:"contracts"`
	Account   *accountResponse   `json:"account,omitempty"`
}

type contractResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Fee         int32   `json:"fee"`
	Description *string `json:"description,omitempty"`
	ImageUrl    *string `json:"imageUrl,omitempty"`
	Type        string  `json:"type"`
}

type accountResponse struct {
	Id    int `json:"id"`
	Count int `json:"count"`
}

func (h *ContractImpl) GetList(ctx *gin.Context) {
	var request getListOfContractsRequest
	if err := ctx.BindQuery(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	contracts, err := h.contractService.GetList(ctx, dto.ContractsFilter{
		Name: request.ContractNameFilter,
		Type: request.ContractTypeFilter,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}
	contractsResponse := make([]contractResponse, 0, len(contracts))
	for _, contract := range contracts {
		contractsResponse = append(contractsResponse, contractResponse{
			Id:          int(contract.Id),
			Name:        contract.Name,
			Fee:         contract.Fee,
			Description: contract.Description,
			ImageUrl:    contract.ImageUrl,
			Type:        string(contract.Type),
		})
	}

	response := getListOfContractsResponse{Contracts: contractsResponse}
	user := auth.GetUser()
	account, err := h.accountService.GetCurrentDraft(ctx, user.ID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		newErrorResponse(ctx, err)
		return
	}

	if !errors.Is(err, domain.ErrNotFound) {
		response.Account = &accountResponse{
			Id:    int(account.Id),
			Count: account.Count,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

type getContractByIdRequest struct {
	Id int `uri:"id"`
}

func (h *ContractImpl) Get(ctx *gin.Context) {
	var request getContractByIdRequest
	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	if err := ctx.BindQuery(&request); err != nil {
		newErrorResponse(ctx, err)
		return
	}

	contract, err := h.contractService.Get(ctx, domain.ContractId(request.Id))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, contractResponse{
		Id:          int(contract.Id),
		Name:        contract.Name,
		Fee:         contract.Fee,
		Description: contract.Description,
		ImageUrl:    contract.ImageUrl,
		Type:        string(contract.Type),
	})
}

type addToAccountRequest struct {
	Id int `uri:"id"`
}

func (h *ContractImpl) AddToAccount(ctx *gin.Context) {
	var request addToAccountRequest
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	err := h.contractService.AddToCurrentDraft(ctx, domain.ContractId(request.Id))
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

	ctx.Redirect(http.StatusSeeOther, "/contracts")
}
