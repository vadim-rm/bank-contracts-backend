package handler

import "github.com/gin-gonic/gin"

type Account interface {
	GetById(ctx *gin.Context)
}
