package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

func (app *Applicaton) ServerError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": message})
}

func (app *Applicaton) ConvertStringToUuid(id string) uuid.UUID {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Println(err.Error())
	}
	return parsedId
}

func (app *Applicaton) ConvertStringToDecimal(amount string) decimal.Decimal {
	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		log.Println(err.Error())
		return decimal.Zero
	}
	return decimalAmount
}

func (app *Applicaton) ConvertStringToInt32(stringInt string) int32 {
	integer, err := strconv.Atoi(stringInt)
	if err != nil {
		log.Println(err.Error())
	}
	return int32(integer)
}
