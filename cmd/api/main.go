package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/config"
	"github.com/markmumba/chasebank/internal/database"
)
type DbInstance struct {
	DB *database.Queries
}

func main() {
	e := echo.New()

	conn,err := sql.Open("postgres",config.Config("DATABASE_URL"))
	if err != nil {
		fmt.Print(err.Error())
	}
	DbInstance := DbInstance{
		DB: database.New(conn),
	}
	defer conn.Close() 

	fmt.Print(DbInstance)
	SetupRouter(e)
	e.Logger.Fatal(e.Start(":" + config.Config("PORT")))

}
