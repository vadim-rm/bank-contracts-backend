package external_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/handler"
)

func Initialize(parent *gin.Engine, contractHandler handler.Contract, accountHandler handler.Account) {
	parent.GET("/contracts", contractHandler.GetList)
	parent.GET("/contracts/:id", contractHandler.GetById)
	parent.POST("/contracts/:id/add-to-account", contractHandler.AddToAccount)

	parent.GET("/accounts/:id", accountHandler.GetById)
	parent.POST("/accounts/:id", accountHandler.Delete)
}
