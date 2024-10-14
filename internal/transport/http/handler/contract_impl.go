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
	Account   accountResponse    `json:"account"`
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
	Id    *int `json:"id"`
	Count int  `json:"count"`
}

func (h *ContractImpl) GetList(ctx *gin.Context) {
	var request getListOfContractsRequest
	if err := ctx.BindQuery(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
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
		response.Account = accountResponse{
			Id:    (*int)(&account.Id),
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
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
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

type createRequest struct {
	Name        string  `json:"name"`
	Fee         int32   `json:"fee"`
	Description *string `json:"description,omitempty"`
	Type        string  `json:"type"`
}

type createResponse struct {
	Id int `json:"id"`
}

func (h *ContractImpl) Create(ctx *gin.Context) {
	var request createRequest
	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	id, err := h.contractService.Create(ctx, service.CreateContractInput{
		Name:        request.Name,
		Fee:         request.Fee,
		Description: request.Description,
		Type:        domain.ContractType(request.Type),
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, createResponse{
		Id: int(id),
	})
}

type updateRequest struct {
	Id int `uri:"id"`

	Name        *string              `json:"name,omitempty"`
	Fee         *int32               `json:"fee,omitempty"`
	Description *string              `json:"description,omitempty"`
	Type        *domain.ContractType `json:"type,omitempty"`
}

func (h *ContractImpl) Update(ctx *gin.Context) {
	var request updateRequest
	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	if err := ctx.BindJSON(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	err := h.contractService.Update(ctx, domain.ContractId(request.Id), service.UpdateContractInput{
		Name:        request.Name,
		Fee:         request.Fee,
		Description: request.Description,
		Type:        request.Type,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type deleteContractRequest struct {
	Id int `uri:"id"`
}

func (h *ContractImpl) Delete(ctx *gin.Context) {
	var request deleteContractRequest
	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	err := h.contractService.Delete(ctx, domain.ContractId(request.Id))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type addToAccountRequest struct {
	Id int `uri:"id"`
}

func (h *ContractImpl) AddToAccount(ctx *gin.Context) {
	var request addToAccountRequest
	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	err := h.contractService.AddToCurrentDraft(ctx, domain.ContractId(request.Id))
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type updateImageRequest struct {
	Id int `uri:"id"`
}

func (h *ContractImpl) UpdateImage(ctx *gin.Context) {
	var request updateImageRequest
	if err := ctx.BindUri(&request); err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}
	defer file.Close()

	err = h.contractService.UpdateImage(ctx, domain.ContractId(request.Id), service.UpdateContractImageInput{
		Image:       file,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
