package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/auth"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"net/http"
	"strings"
	"time"
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
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
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
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	err = h.service.Update(ctx, claims.UserId, service.UpdateUserInput{
		Name:     request.Name,
		Password: request.PasswordHash,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type authenticateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authenticateUserResponse struct {
	ExpiresAt   time.Time `json:"expiresAt"`
	AccessToken string    `json:"accessToken"`
}

func (h *UserImpl) Authenticate(ctx *gin.Context) {
	var request authenticateUserRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	token, err := h.service.Authenticate(ctx, service.AuthorizeInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, authenticateUserResponse{
		ExpiresAt:   token.ExpiresAt,
		AccessToken: token.Token,
	})
}

func (h *UserImpl) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token, _ = strings.CutPrefix(token, "Bearer ")

	err := h.service.Logout(ctx, token)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
