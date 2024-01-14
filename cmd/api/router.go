package main

import "github.com/labstack/echo/v4"

func (app *Applicaton) SetupRouter(e *echo.Echo) {

	api := e.Group("/api")
	api.GET("/user", app.GetAllUsers)
	api.POST("/user", app.CreateUser)

	api.POST("/account/:id", app.CreateSavingAccount)
	api.POST("/account/deposit/:id", app.Deposit)
	api.POST("/account/withdraw/:id", app.Withdraw)
	api.GET("/account/transactions/:id", app.ViewTransactions)
	api.POST("/account/checktosave/:id", app.TransferCheckingToSaving)
	api.POST("/account/transferfunds/:id", app.TransferFunds)
}
