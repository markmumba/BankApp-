package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/internal/database"
	"github.com/shopspring/decimal"
)

func ConvertStringToUuid(id string) uuid.UUID {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		fmt.Println(err)
	}
	return parsedId
}

func ConvertStringToDecimal(amount string) decimal.Decimal {
	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		fmt.Println(err.Error())
		return decimal.Zero
	}
	return decimalAmount
}

func (app *Applicaton) CreateSavingAccount(c echo.Context) error {
	var jsonResp Response

	id := c.Param("id")
	parsedId := ConvertStringToUuid(id)

	account, err := app.DB.CreateAccount(app.ctx, database.CreateAccountParams{
		UserID:      uuid.NullUUID{UUID: parsedId, Valid: true},
		AccountType: Savings,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(account)
	jsonResp = Response{
		Message: "Savings account has been successfully been created",
	}
	return c.JSON(http.StatusOK, jsonResp)

}

func (app *Applicaton) Deposit(c echo.Context) error {
	var account Account
	var jsonResp Account

	id := c.Param("id")
	err := c.Bind(&account)
	parsedId := ConvertStringToUuid(id)

	userAccount, err := app.DB.FindAccount(app.ctx, parsedId)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, acc := range userAccount {
		if account.AccountType == acc.AccountType {
			newTotal := ConvertStringToDecimal(acc.Balance).Add(ConvertStringToDecimal(account.Amount))
			account, err := app.DB.Deposit(app.ctx, database.DepositParams{
				Balance:   newTotal.String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				fmt.Println(err.Error())
			}
			jsonResp = Account{
				AccountType: account.AccountType,
				Amount:      account.Balance,
			}
		}
	}

	return c.JSON(http.StatusOK, jsonResp)

}

func (app *Applicaton) Withdraw(c echo.Context) error {
	var account Account
	var jsonResp Account

	id := c.Param("id")
	err := c.Bind(&account)
	parsedId := ConvertStringToUuid(id)

	userAccount, err := app.DB.FindAccount(app.ctx, parsedId)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, acc := range userAccount {
		if account.AccountType == acc.AccountType {
			newTotal := ConvertStringToDecimal(acc.Balance).Sub(ConvertStringToDecimal(account.Amount))
			account, err := app.DB.Withdraw(app.ctx, database.WithdrawParams{
				Balance:   newTotal.String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				fmt.Println(err.Error())
			}
			jsonResp = Account{
				AccountType: account.AccountType,
				Amount:      account.Balance,
			}
		}
	}

	return c.JSON(http.StatusOK, jsonResp)

}
