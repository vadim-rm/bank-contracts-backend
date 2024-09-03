package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"net/http"
)

type Config struct {
	DebugCors     bool
	TemplatesPath string
}

func New(config Config) *gin.Engine {
	router := gin.Default()

	if config.DebugCors {
		router.Use(cors.Default())
	}

	router.LoadHTMLGlob(config.TemplatesPath)

	router.Static("/static", "./static") //todo add to config

	router.NoRoute(func(c *gin.Context) {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": domain.ErrNotFound.Error(),
			},
		)
	})

	return router
}
