package external_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/handler"
)

func Initialize(
	parent *gin.Engine,
	contractHandler handler.Contract,
	accountHandler handler.Account,
	accountContractsHandler handler.AccountContracts,
	usersHandler handler.User,
) {
	contracts := parent.Group("contracts")
	{
		initializeContracts(contracts, contractHandler)
	}

	accounts := parent.Group("accounts")
	{
		initializeAccounts(accounts, accountHandler)
	}

	users := parent.Group("users")
	{
		initializeUsers(users, usersHandler)
	}
	initializeAccountContracts(parent, accountContractsHandler)
}

func initializeContracts(parent *gin.RouterGroup, contractHandler handler.Contract) {
	parent.GET("", contractHandler.GetList)
	parent.GET(":id", contractHandler.Get)
	parent.POST("", contractHandler.Create)
	parent.PUT(":id", contractHandler.Update)
	parent.DELETE(":id", contractHandler.Delete)
	parent.POST(":id/draft", contractHandler.AddToAccount)
	parent.PUT(":id/image", contractHandler.UpdateImage)
}

func initializeAccounts(parent *gin.RouterGroup, accountsHandler handler.Account) {
	parent.GET("", accountsHandler.GetList)
	parent.GET(":accountId", accountsHandler.Get)
	parent.PUT(":accountId", accountsHandler.Update)
	parent.PUT(":accountId/submit", accountsHandler.Submit)
	parent.PUT(":accountId/complete", accountsHandler.Complete)
	parent.DELETE(":accountId", accountsHandler.Delete)
}

func initializeAccountContracts(parent *gin.Engine, accountContractsHandler handler.AccountContracts) {
	parent.DELETE("/accounts/:accountId/contract/:contractId", accountContractsHandler.Delete)
	parent.PUT("/accounts/:accountId/contract/:contractId/main", accountContractsHandler.SetMain)
}

func initializeUsers(parent *gin.RouterGroup, usersHandler handler.User) {
	parent.POST("", usersHandler.Create)
	parent.PUT("", usersHandler.Update)
}
