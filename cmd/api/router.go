package main

import (

	"github.com/labstack/echo/v4"
)

func (app *Applicaton) SetupRouter(e *echo.Echo) {
	
	public := e.Group("/api")
	public.POST("/user/create",app.CreateUser)
	public.POST("/login", app.Login)

	api := e.Group("/api")
	api.Use(Authentication)

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
