package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/markmumba/chasebank/config"
	"github.com/markmumba/chasebank/internal/database"
)

type Applicaton struct {
	DB  *database.Queries
	ctx context.Context
}

type DbInstance struct {
	DB *database.Queries
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())

	conn, err := sql.Open("postgres", config.Config("DATABASE_URL"))

	if err != nil {
		fmt.Print(err.Error())
	}
	Applicaton := &Applicaton{
		DB:  database.New(conn),
		ctx: e.NewContext(nil,nil).Request().Context(),
	}

	defer conn.Close()

	Applicaton.SetupRouter(e)
	e.Logger.Fatal(e.Start(":" + config.Config("PORT")))

}
