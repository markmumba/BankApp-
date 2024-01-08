package main

import "github.com/labstack/echo/v4"



func SetupRouter(e *echo.Echo) {
e.GET("/",Home)
}