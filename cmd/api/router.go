package main

import "github.com/labstack/echo/v4"

func (app *Applicaton) SetupRouter(e *echo.Echo) {

	api := e.Group("/api")
	api.POST("/user", app.CreateUser)
}
