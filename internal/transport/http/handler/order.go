package handler

import "github.com/gin-gonic/gin"

type Order interface {
	GetById(ctx *gin.Context)
}
