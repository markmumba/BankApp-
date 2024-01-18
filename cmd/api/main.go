package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/markmumba/chasebank/config"
	"github.com/markmumba/chasebank/internal/database"
)

type Applicaton struct {
	SDB *sql.DB
	DB  *database.Queries
	Ctx context.Context
}

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173/"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	conn, err := sql.Open("postgres", config.Config("DATABASE_URL"))
	if err != nil {
		log.Println(err.Error())
	}
	appContex := context.Background()

	Applicaton := &Applicaton{
		SDB: conn,
		DB:  database.New(conn),
		Ctx: appContex,
	}

	defer conn.Close()

	Applicaton.SetupRouter(e)
	e.Logger.Fatal(e.Start(":" + config.Config("PORT")))

}
