package internal

import (
	"github.com/labstack/echo/v4"

	"github.com/markmumba/chasebank/internal/handlers"
	"github.com/markmumba/chasebank/internal/middleware"
)

func SetupRouter(e *echo.Echo, app *handlers.Applicaton) {

	public := e.Group("/api")
	public.GET("/health", app.HealthCheck)
	public.POST("/user/create", app.CreateUser)
	public.POST("/login", app.Login)
	public.POST("/logout", app.Logout)

	protected := e.Group("/api")
	protected.Use(middleware.Authentication)

	protected.GET("/users", app.GetAllUsers)
	protected.GET("/user", app.GetUser)

	account := protected.Group("/account")
	account.GET("/", app.GetUserAccounts)
	account.GET("/transactions", app.ViewTransactions)

	account.POST("/create", app.CreateSavingAccount)
	account.POST("/deposit", app.Deposit)
	account.POST("/withdraw", app.Withdraw)
	account.POST("/checktosave", app.TransferCheckingToSaving)
	account.POST("/transferfunds", app.TransferFunds)
}
