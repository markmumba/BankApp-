package main

import (
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/config"
)

type DbModel struct {
    
}

func main() {
	e := echo.New()
	SetupRouter(e)
	e.Logger.Fatal(e.Start(":" + config.Config("PORT")))

}
