package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/markmumba/chasebank/config"
	"github.com/markmumba/chasebank/internal"
	"github.com/markmumba/chasebank/internal/handlers"
)

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	appContex := context.Background()

	Applicaton := &handlers.Applicaton{
		SDB: config.ConnectDB(config.Config("DATABASE_URL")).SDB,
		DB:  config.ConnectDB("DATABASE_URL").DB,
		Ctx: appContex,
	}

	defer config.ConnectDB(config.Config("DATABASE_URL")).SDB.Close()

	internal.SetupRouter(e, Applicaton)
	e.Logger.Fatal(e.Start(":" + config.Config("PORT")))

}
