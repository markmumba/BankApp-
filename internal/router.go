package internal

import (
	"github.com/labstack/echo/v4"

	"github.com/markmumba/chasebank/internal/handlers"
	"github.com/markmumba/chasebank/internal/middleware"
)

func SetupRouter(e *echo.Echo, app *handlers.Applicaton) {

	public := e.Group("/apip")
	public.POST("/user/create", app.CreateUser)
	public.POST("/login", app.Login)

	api := e.Group("/api")
	api.Use(middleware.Authentication)

	api.GET("/users", app.GetAllUsers)
	api.GET("/user/:id", app.GetUser)

	account := api.Group("/account")
	account.GET("/:id", app.GetUserAccounts)
	account.GET("/transactions/:id", app.ViewTransactions)

	account.POST("/:id", app.CreateSavingAccount)
	account.POST("/deposit/:id", app.Deposit)
	account.POST("/withdraw/:id", app.Withdraw)
	account.POST("/checktosave/:id", app.TransferCheckingToSaving)
	account.POST("/transferfunds/:id", app.TransferFunds)
}
