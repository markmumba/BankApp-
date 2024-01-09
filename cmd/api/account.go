package main

import (
	"database/sql"
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to Create account"})
	}

	fmt.Println(account)
	jsonResp = Response{
		Message: "Savings account has been successfully been created",
	}
	return c.JSON(http.StatusOK, jsonResp)

}

func (app *Applicaton) Deposit(c echo.Context) error {
	var accountStruct Account
	var jsonResp Account

	id := c.Param("id")
	err := c.Bind(&accountStruct)
	parsedId := ConvertStringToUuid(id)

	userAccount, err := app.DB.FindAccount(app.ctx, parsedId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user accounts"})
	}

	for _, acc := range userAccount {
		if accountStruct.AccountType == acc.AccountType {
			newTotal := ConvertStringToDecimal(acc.Balance).Add(ConvertStringToDecimal(accountStruct.Amount))
			account, err := app.DB.Deposit(app.ctx, database.DepositParams{
				Balance:   newTotal.String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to make Deposit "})
			}
			app.SaveTransaction(acc.AccountID, 0, accountStruct.Amount, Deposit)
			jsonResp = Account{
				AccountType: account.AccountType,
				Amount:      account.Balance,
			}
		}
	}

	return c.JSON(http.StatusOK, jsonResp)

}

func (app *Applicaton) Withdraw(c echo.Context) error {
	var accountStruct Account
	var jsonResp Account

	id := c.Param("id")
	err := c.Bind(&accountStruct)
	parsedId := ConvertStringToUuid(id)

	userAccount, err := app.DB.FindAccount(app.ctx, parsedId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user accounts"})
	}

	for _, acc := range userAccount {
		if accountStruct.AccountType == acc.AccountType {
			newTotal := ConvertStringToDecimal(acc.Balance).Sub(ConvertStringToDecimal(accountStruct.Amount))
			account, err := app.DB.Withdraw(app.ctx, database.WithdrawParams{
				Balance:   newTotal.String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to complete withdrawal"})
			}
			app.SaveTransaction(acc.AccountID, 0, accountStruct.Amount, Withdraw)
			jsonResp = Account{
				AccountType: account.AccountType,
				Amount:      account.Balance,
			}
		}
	}

	return c.JSON(http.StatusOK, jsonResp)

}

func (app *Applicaton) ViewTransactions(c echo.Context) error {
	var account Account
	var transactions []database.Transaction
	var jsonResp []Transaction

	id := c.Param("id")
	err := c.Bind(&account)
	parsedId := ConvertStringToUuid(id)
	userAccount, err := app.DB.FindAccount(app.ctx, parsedId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user accounts"})

	}

	found := false
	for _, acc := range userAccount {
		if account.AccountType == acc.AccountType {
			transactions, err = app.DB.ViewTransactions(app.ctx, sql.NullInt32{Int32: acc.AccountID, Valid: true})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve transactions"})
			}
			found = true
			break
		}
	}
	if !found {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
	}

	for _, transaction := range transactions {
		newTransaction := Transaction{
			RecepientID: transaction.RecepientID.Int32,
			Amount:      transaction.Amount,
			Type:        transaction.Type,
			Timestamp:   transaction.Timestamp.Time,
		}
		jsonResp = append(jsonResp, newTransaction)
	}
	return c.JSON(http.StatusOK, jsonResp)
}
