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

// GetList
// @Summary Получение списка договоров
// @Description Возвращает список договоров с возможностью фильтрации по названию и типу договора
// @Tags contracts
// @Accept  json
// @Produce  json
// @Param contractName query string false "Фильтр по названию договора"
// @Param contractType query string false "Фильтр по типу договора"
// @Success 200 {object} getListOfContractsResponse "Список договоров и учетная запись"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /contracts [get]
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
	claims, err := auth.GetClaims(ctx)
	if err == nil {
		account, err := h.accountService.GetCurrentDraft(ctx, claims.UserId)
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
	}

	ctx.JSON(http.StatusOK, response)
}

type getContractByIdRequest struct {
	Id int `uri:"id"`
}

// Get
// @Summary Получение информации о договоре
// @Description Возвращает информацию о договоре по его ID.
// @Tags contracts
// @Accept  json
// @Produce  json
// @Param id path int true "ID договора"
// @Success 200 {object} contractResponse "Информация о договоре"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Договор не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /contracts/{id} [get]
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

// Create
// @Summary Создание нового договора
// @Description Создает новый договор с указанными данными.
// @Tags contracts
// @Accept  json
// @Produce  json
// @Param request body createRequest true "Данные для создания договора"
// @Success 201 {object} createResponse "Успешное создание договора"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /contracts [post]
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

// Update
// @Summary Обновление данных договора
// @Description Обновляет данные существующего договора по его ID.
// @Tags contracts
// @Accept  json
// @Produce  json
// @Param id path int true "ID договора"
// @Param request body updateRequest true "Данные для обновления договора"
// @Success 200 "Договор успешно обновлен"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Договор не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /contracts/{id} [put]
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

// Delete godoc
// @Summary Удаление договора
// @Description Удаляет существующий договор по его ID.
// @Tags contracts
// @Accept  json
// @Produce  json
// @Param id path int true "ID договора"
// @Success 200 "Договор успешно удален"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Договор не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /contracts/{id} [delete]
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

// AddToAccount
// @Summary Добавление договора в заявку на счёт
// @Description Добавляет существующий договор в текущую заявку на счёт по его ID.
// @Tags contracts
// @Accept  json
// @Produce  json
// @Param id path int true "ID договора"
// @Success 200 "Договор успешно добавлен в заявку на счёт"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Договор не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /contracts/{id}/draft [post]
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

// UpdateImage
// @Summary Обновление изображения договора
// @Description Обновляет изображение для договора по его ID
// @Tags contracts
// @Accept  multipart/form-data
// @Produce  json
// @Param id path int true "ID договора"
// @Param image formData file true "Файл изображения для загрузки"
// @Success 200 "Изображение успешно обновлено"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 404 {object} errorResponse "Договор не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /contracts/{id}/image [put]
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
