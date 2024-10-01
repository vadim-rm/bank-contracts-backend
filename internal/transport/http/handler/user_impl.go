package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/auth"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"net/http"
)

type UserImpl struct {
	service service.User
}

func NewUserImpl(service service.User) *UserImpl {
	return &UserImpl{service: service}
}

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserImpl) Create(ctx *gin.Context) {
	var request createUserRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	_, err = h.service.Create(ctx, service.CreateUserInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type updateUserRequest struct {
	Name         *string `json:"name,omitempty"`
	PasswordHash *string `json:"passwordHash,omitempty"`
}

func (h *UserImpl) Update(ctx *gin.Context) {
	var request updateUserRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	user := auth.GetUser()

	err = h.service.Update(ctx, user.ID, service.UpdateUserInput{
		Name:     request.Name,
		Password: request.PasswordHash,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
