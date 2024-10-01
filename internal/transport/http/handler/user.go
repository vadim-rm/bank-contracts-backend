package handler

import "github.com/gin-gonic/gin"

type User interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
}
