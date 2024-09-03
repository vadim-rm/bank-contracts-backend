package external_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/handler"
)

func Initialize(parent *gin.Engine, contractHandler handler.Contract) {
	parent.GET("/contracts", contractHandler.GetList)
	parent.GET("/contracts/:id", contractHandler.GetById)
}
