package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/config"
)


func main() {
	e := echo.New()

	DbInstance := config.ConnectDB()
	fmt.Print(DbInstance)
	SetupRouter(e)
	e.Logger.Fatal(e.Start(":" + config.Config("PORT")))

}
