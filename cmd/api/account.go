package main

//TODO  refuse withdrawal when money exceeds balance

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/internal/database"
	"github.com/shopspring/decimal"
)

func ConvertStringToUuid(id string) uuid.UUID {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Println(err.Error())
	}
	return parsedId
}

func ConvertStringToDecimal(amount string) decimal.Decimal {
	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		log.Println(err.Error())
		return decimal.Zero
	}
	return decimalAmount
}

func (app *Applicaton) CreateSavingAccount(c echo.Context) error {
	var jsonResp Response

	id := c.Param("id")
	parsedId := ConvertStringToUuid(id)

	account, err := app.DB.CreateAccount(app.Ctx, database.CreateAccountParams{
		UserID:      uuid.NullUUID{UUID: parsedId, Valid: true},
		AccountType: Savings,
	})
	if err != nil {
		app.ServerError(c, "Failed to create account ")
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

	userAccount, err := app.DB.FindAccount(app.Ctx, parsedId)
	if err != nil {
		app.ServerError(c, "Failed to retrieve all accounts")
	}

	for _, acc := range userAccount {
		if accountStruct.AccountType == acc.AccountType {
			newTotal := ConvertStringToDecimal(acc.Balance).Add(ConvertStringToDecimal(accountStruct.Amount))
			account, err := app.DB.Deposit(app.Ctx, database.DepositParams{
				Balance:   newTotal.String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				app.ServerError(c, "Failed to make Deposit")
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

	userAccount, err := app.DB.FindAccount(app.Ctx, parsedId)
	if err != nil {
		app.ServerError(c, "Failed to retrieve user accounts ")
	}

	for _, acc := range userAccount {
		if accountStruct.AccountType == acc.AccountType {
			newTotal := ConvertStringToDecimal(acc.Balance).Sub(ConvertStringToDecimal(accountStruct.Amount))
			account, err := app.DB.Withdraw(app.Ctx, database.WithdrawParams{
				Balance:   newTotal.String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				app.ServerError(c, "Failed to complete withdrawal")
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
	userAccount, err := app.DB.FindAccount(app.Ctx, parsedId)
	if err != nil {
		app.ServerError(c, "Failed to retrieve user accounts")

	}

	found := false
	for _, acc := range userAccount {
		if account.AccountType == acc.AccountType {
			transactions, err = app.DB.ViewTransactions(app.Ctx, sql.NullInt32{Int32: acc.AccountID, Valid: true})
			if err != nil {
				app.ServerError(c, "Failed to retrieve the transactions ")
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

func (app *Applicaton) TransferCheckingToSaving(c echo.Context) error {
	type Funds struct {
		Amount string `json:"amount"`
	}
	var fundInstance Funds

	id := c.Param("id")
	err := c.Bind(&fundInstance)
	if err != nil {
		app.ServerError(c, err.Error())
	}
	parsedId := ConvertStringToUuid(id)

	accounts, err := app.DB.FindAccount(app.Ctx, parsedId)
	if err != nil {
		app.ServerError(c, "Could not find the users accounts")
	}

	tx, err := app.SDB.Begin()
	if err != nil {
		app.ServerError(c, err.Error())
	}
	defer tx.Rollback()

	qtx := app.DB.WithTx(tx)

	for _, account := range accounts {
		err = qtx.DebitChecking(app.Ctx, database.DebitCheckingParams{
			AccountID: account.AccountID,
			Balance:   ConvertStringToDecimal(account.Balance).Sub(ConvertStringToDecimal(fundInstance.Amount)).String(),
		})
		if err != nil {
			app.ServerError(c, err.Error())

		}
		err = qtx.CreditSaving(app.Ctx, database.CreditSavingParams{
			AccountID: account.AccountID,
			Balance:   ConvertStringToDecimal(account.Balance).Add(ConvertStringToDecimal(fundInstance.Amount)).String(),
		})
		if err != nil {
			app.ServerError(c, err.Error())
		}

		break
	}

		err = tx.Commit()
	if err != nil {
		app.ServerError(c, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"Success": "Transfer succesful"})

}

func (app *Applicaton) TransferFunds (c *echo.Context) error {
	return nil 
}