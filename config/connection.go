package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/markmumba/chasebank/internal/database"
)

type DbInstance struct {
	DB *database.Queries
}

func ConnectDB() DbInstance {

	conn, err := sql.Open("postgres", Config("DATABASE_URL"))
	if err != nil {
		fmt.Printf(err.Error())
	}

	return DbInstance{
		DB: database.New(conn),
	}

}
