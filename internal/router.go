package internal

import (
	"github.com/labstack/echo/v4"

	"github.com/markmumba/chasebank/internal/handlers"
	//"github.com/markmumba/chasebank/internal/middleware"
)

func SetupRouter(e *echo.Echo, app *handlers.Applicaton) {

	public := e.Group("/api")
	public.POST("/user/create", app.CreateUser)
	public.POST("/login", app.Login)
	public.POST("/logout", app.Logout)

	protected := e.Group("/api")
	//protected.Use(middleware.Authentication)

	protected.GET("/users", app.GetAllUsers)
	protected.GET("/user", app.GetUser)

	account := protected.Group("/account")
	account.GET("/:id", app.GetUserAccounts)
	account.GET("/transactions/:id", app.ViewTransactions)

	account.POST("/:id", app.CreateSavingAccount)
	account.POST("/deposit/:id", app.Deposit)
	account.POST("/withdraw/:id", app.Withdraw)
	account.POST("/checktosave/:id", app.TransferCheckingToSaving)
	account.POST("/transferfunds/:id", app.TransferFunds)
}
