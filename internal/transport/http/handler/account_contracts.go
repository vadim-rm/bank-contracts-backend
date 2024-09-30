package handler

import "github.com/gin-gonic/gin"

type AccountContracts interface {
	Delete(ctx *gin.Context)
	SetMain(ctx *gin.Context)
}
