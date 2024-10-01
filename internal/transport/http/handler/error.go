package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func newErrorResponse(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(
		http.StatusInternalServerError,
		gin.H{
			"error": err.Error(),
		},
	)
	// todo change status code (errors is bind error gin)
}
