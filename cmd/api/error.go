package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *Applicaton) ServerError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": message})
}
