package external_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/handler"
)

func Initialize(parent *gin.Engine, contractHandler handler.Contract, orderHandler handler.Account) {
	parent.GET("/contracts", contractHandler.GetList)
	parent.GET("/contracts/:id", contractHandler.GetById)

	parent.GET("/accounts/:id", orderHandler.GetById)
}
