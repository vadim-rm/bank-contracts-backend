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
) {
	contracts := parent.Group("contracts")
	{
		initializeContracts(contracts, contractHandler)
	}

	accounts := parent.Group("accounts")
	{
		initializeAccounts(accounts, accountHandler)
		initializeAccountContracts(accounts, accountContractsHandler)
	}
}

func initializeContracts(parent *gin.RouterGroup, contractHandler handler.Contract) {
	parent.GET("", contractHandler.GetList)
	parent.GET(":id", contractHandler.Get)
	parent.POST("", contractHandler.Create)
	parent.PUT(":id", contractHandler.Update)
	parent.DELETE(":id", contractHandler.Delete)
	parent.POST(":id/draft", contractHandler.AddToAccount)
	// todo. add image upload
}

func initializeAccounts(parent *gin.RouterGroup, accountsHandler handler.Account) {
	parent.GET("", accountsHandler.GetList)
	parent.GET(":id", accountsHandler.Get)
	parent.PUT(":id", accountsHandler.Update)
	parent.PUT(":id/submit", accountsHandler.Submit)
	parent.PUT(":id/complete", accountsHandler.Complete)
	parent.POST(":id", accountsHandler.Delete)
}

func initializeAccountContracts(parent *gin.RouterGroup, accountContractsHandler handler.AccountContracts) {
	parent.DELETE("/accounts/:accountId/contract/:contractId", accountContractsHandler.Delete)
	parent.PUT("/accounts/:accountId/contract/:contractId/main", accountContractsHandler.SetMain)
}

// todo. crud users
