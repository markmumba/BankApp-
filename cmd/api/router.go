package main

import "github.com/labstack/echo/v4"

func (app *Applicaton) SetupRouter(e *echo.Echo) {

	api := e.Group("/api")
	api.GET("/user", app.GetAllUsers)
	api.POST("/user", app.CreateUser)
	api.GET("/user/:id", app.GetUser)

	account := api.Group("/account")
	account.GET("/:id", app.GetUserAccounts)
	account.POST("/:id", app.CreateSavingAccount)
	account.POST("/deposit/:id", app.Deposit)
	account.POST("/withdraw/:id", app.Withdraw)
	account.GET("/transactions/:id", app.ViewTransactions)
	account.POST("/checktosave/:id", app.TransferCheckingToSaving)
	account.POST("/transferfunds/:id", app.TransferFunds)
}
