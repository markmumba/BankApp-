package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/markmumba/chasebank/config"
	"github.com/markmumba/chasebank/internal/database"
)

type Applicaton struct {
	DB  *database.Queries
	Ctx context.Context
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())

	Applicaton := &Applicaton{
		DB:  database.New(config.Conn),
		Ctx: e.NewContext(nil, nil).Request().Context(),
	}

	defer config.Conn.Close()

	Applicaton.SetupRouter(e)
	e.Logger.Fatal(e.Start(":" + config.Config("PORT")))

}
