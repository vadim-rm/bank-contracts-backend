package handler

import "github.com/gin-gonic/gin"

type Contract interface {
	GetList(ctx *gin.Context)
	GetById(ctx *gin.Context)
}
