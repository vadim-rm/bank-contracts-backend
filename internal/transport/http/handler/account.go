package handler

import "github.com/gin-gonic/gin"

type Account interface {
	Delete(ctx *gin.Context)
	GetById(ctx *gin.Context)
}
