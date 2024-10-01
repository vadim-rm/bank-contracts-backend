package handler

import "github.com/gin-gonic/gin"

type Contract interface {
	GetList(ctx *gin.Context)
	Get(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	AddToAccount(ctx *gin.Context)
	UpdateImage(ctx *gin.Context)
}
