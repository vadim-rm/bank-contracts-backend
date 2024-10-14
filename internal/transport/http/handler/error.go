package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/pkg/logger"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError
	if errors.Is(err, domain.ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, domain.ErrAccountNumberEmpty) {
		code = http.StatusBadRequest
	} else if errors.Is(err, domain.ErrInvalidTargetStatus) {
		code = http.StatusBadRequest
	} else if errors.Is(err, domain.ErrActionNotPermitted) {
		code = http.StatusForbidden
	} else if errors.Is(err, domain.ErrWrongAccountStatus) {
		code = http.StatusBadRequest
	} else if errors.Is(err, domain.ErrBadRequest) {
		code = http.StatusBadRequest
	} else if errors.Is(err, domain.ErrInvalidCredentials) {
		code = http.StatusUnauthorized
	} else if errors.Is(err, domain.ErrUnauthenticated) {
		code = http.StatusUnauthorized
	}

	logger.Error(err.Error())

	ctx.AbortWithStatusJSON(
		code,
		errorResponse{
			Error: err.Error(),
		},
	)
}
