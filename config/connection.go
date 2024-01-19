package config

import (
	"database/sql"
	"log"

	"github.com/markmumba/chasebank/internal/database"
)

type DBInstance struct {
	SDB *sql.DB
	DB  *database.Queries
}

func ConnectDB(dbUrl string) *DBInstance {

	conn, err := sql.Open("postgres", Config(dbUrl))
	if err != nil {
		log.Println(err.Error())
	}

	qtx := database.New(conn)

	return &DBInstance{
		conn,
		qtx,
	}
}
