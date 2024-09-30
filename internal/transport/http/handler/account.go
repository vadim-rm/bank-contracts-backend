package handler

import "github.com/gin-gonic/gin"

type Account interface {
	GetList(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Submit(ctx *gin.Context)
	Complete(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetById(ctx *gin.Context)
}
