package external_routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vadim-rm/bmstu-web-backend/docs"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/handler"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/middleware"
)

func Initialize(
	parent *gin.Engine,
	authMiddleware *middleware.AuthMiddleware,
	contractHandler handler.Contract,
	accountHandler handler.Account,
	accountContractsHandler handler.AccountContracts,
	usersHandler handler.User,
) {
	parent.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	contracts := parent.Group("contracts", authMiddleware.WithOptionalAuth)
	{
		initializeContracts(contracts, authMiddleware, contractHandler)
	}

	accounts := parent.Group("accounts")
	{
		initializeAccounts(accounts, authMiddleware, accountHandler)
	}

	users := parent.Group("users")
	{
		initializeUsers(users, authMiddleware, usersHandler)
	}
	initializeAccountContracts(parent, authMiddleware, accountContractsHandler)
}

func initializeContracts(parent *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, contractHandler handler.Contract) {
	parent.GET("", contractHandler.GetList)
	parent.GET(":id", contractHandler.Get)

	parent.Group("", authMiddleware.WithAuth).POST(":id/draft", contractHandler.AddToAccount)

	moderator := parent.Group("", authMiddleware.WithAuth, authMiddleware.WithModerator)
	{
		moderator.POST("", contractHandler.Create)
		moderator.PUT(":id", contractHandler.Update)
		moderator.DELETE(":id", contractHandler.Delete)
		moderator.PUT(":id/image", contractHandler.UpdateImage)
	}
}

func initializeAccounts(parent *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, accountsHandler handler.Account) {
	authorized := parent.Group("", authMiddleware.WithAuth)
	{
		authorized.GET("", accountsHandler.GetList)
		authorized.GET(":accountId", accountsHandler.Get)
		authorized.PUT(":accountId", accountsHandler.Update)
		authorized.PUT(":accountId/submit", accountsHandler.Submit)
		authorized.DELETE(":accountId", accountsHandler.Delete)

		authorized.Group("", authMiddleware.WithModerator).PUT(":accountId/complete", accountsHandler.Complete)
	}
}

func initializeAccountContracts(parent *gin.Engine, authMiddleware *middleware.AuthMiddleware, accountContractsHandler handler.AccountContracts) {
	authorized := parent.Group("", authMiddleware.WithAuth)
	{
		authorized.DELETE("/accounts/:accountId/contract/:contractId", accountContractsHandler.Delete)
		authorized.PUT("/accounts/:accountId/contract/:contractId/main", accountContractsHandler.SetMain)
	}
}

func initializeUsers(parent *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, usersHandler handler.User) {
	parent.POST("", usersHandler.Create)
	parent.POST("login", usersHandler.Authenticate)

	authorized := parent.Group("", authMiddleware.WithAuth)
	{
		authorized.PUT("", usersHandler.Update)
		authorized.POST("logout", usersHandler.Logout)
	}
}
