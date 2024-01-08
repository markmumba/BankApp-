package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/config"
)

func Home(c echo.Context) error {
	return c.String(http.StatusOK, " We are running the server ")
}


func ( *config.DbInstance) CreateUser ()