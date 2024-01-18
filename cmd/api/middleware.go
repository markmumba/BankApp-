package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/config"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Config("SECRET_KEY")), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "unauthorized",
			})
		}

		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "unauthorized",
			})

		
		}
		return next(c)
	}
}
