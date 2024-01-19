package handlers 

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/internal/database"
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

func (app *Applicaton) CheckBalance(balance string, amount string) bool {

	if app.ConvertStringToDecimal(balance).GreaterThan(app.ConvertStringToDecimal(amount)) {
		return true
	} else {
		return false
	}

}

func (app *Applicaton) DepositHelper(balance string, amount string) decimal.Decimal {
	return app.ConvertStringToDecimal(balance).Add(app.ConvertStringToDecimal(amount))
}

func (app *Applicaton) WithdrawHelper(balance string, amount string) decimal.Decimal {
	return app.ConvertStringToDecimal(balance).Sub(app.ConvertStringToDecimal(amount))
}

func (app *Applicaton) FindAccountHelper(c echo.Context, parseId uuid.UUID) []database.FindAccountRow {

	userAccounts, err := app.DB.FindAccount(app.Ctx, parseId)
	if err != nil {
		app.ServerError(c, "Failed to retrieve user accounts")

	}
	return userAccounts
}

