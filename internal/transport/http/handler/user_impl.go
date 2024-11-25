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
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Create
// @Summary Создание нового пользователя
// @Description Создает нового пользователя с указанными именем, электронной почтой и паролем.
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body createUserRequest true "Данные для создания пользователя"
// @Success 200 "Пользователь успешно создан"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /users [post]
func (h *UserImpl) Create(ctx *gin.Context) {
	var request createUserRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	_, err = h.service.Create(ctx, service.CreateUserInput{
		Login:    request.Login,
		Password: request.Password,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type updateUserRequest struct {
	Password *string `json:"password,omitempty"`
}

// Update
// @Summary Обновление данных текущего пользователя
// @Description Обновляет имя и/или пароль текущего пользователя.
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body updateUserRequest true "Данные для обновления пользователя"
// @Success 200 "Данные пользователя успешно обновлены"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 401 {object} errorResponse "Неавторизован"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /users [put]
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
		Password: request.Password,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type authenticateUserRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authenticateUserResponse struct {
	ExpiresAt   time.Time `json:"expiresAt"`
	AccessToken string    `json:"accessToken"`
	Login       string    `json:"login"`
	IsModerator bool      `json:"isModerator"`
}

// Authenticate
// @Summary Аутентификация пользователя
// @Description Аутентифицирует пользователя по электронной почте и паролю и возвращает токен доступа.
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body authenticateUserRequest true "Данные для аутентификации пользователя"
// @Success 200 {object} authenticateUserResponse "Успешная аутентификация, токен доступа"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 401 {object} errorResponse "Неавторизован"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /users/login [post]
func (h *UserImpl) Authenticate(ctx *gin.Context) {
	var request authenticateUserRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		newErrorResponse(ctx, errors.Join(domain.ErrBadRequest, err))
		return
	}

	token, err := h.service.Authenticate(ctx, service.AuthorizeInput{
		Login:    request.Login,
		Password: request.Password,
	})
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	user, err := h.service.Get(ctx, request.Login)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, authenticateUserResponse{
		ExpiresAt:   token.ExpiresAt,
		AccessToken: token.Token,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	})
}

// Logout
// @Summary Выход пользователя
// @Description Завершает сеанс пользователя, аннулируя токен доступа.
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 "Пользователь успешно вышел"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 401 {object} errorResponse "Неавторизован"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Security Bearer
// @Router /users/logout [post]
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
