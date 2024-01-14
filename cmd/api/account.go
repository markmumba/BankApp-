package main

//FIXME  refuse withdrawal when money exceeds balance
//TODO view all accounts 
// TODO return balance as json after the transaction 
//FIXME understand where the empty curly braces are showing	
import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/internal/database"
)

func (app *Applicaton) CreateSavingAccount(c echo.Context) error {
	var jsonResp string

	id := c.Param("id")
	parsedId := app.ConvertStringToUuid(id)

	account, err := app.DB.CreateAccount(app.Ctx, database.CreateAccountParams{
		UserID:      uuid.NullUUID{UUID: parsedId, Valid: true},
		AccountType: Savings,
	})
	if err != nil {
		app.ServerError(c, err.Error())
	}

	jsonResp = account.AccountNumber

	return c.JSON(http.StatusOK, jsonResp)

}

func (app *Applicaton) Deposit(c echo.Context) error {
	var accountStruct Account
	var jsonResp Account

	id := c.Param("id")
	err := c.Bind(&accountStruct)
	parsedId := app.ConvertStringToUuid(id)

	userAccount, err := app.DB.FindAccount(app.Ctx, parsedId)
	if err != nil {
		app.ServerError(c, "Failed to retrieve all accounts")
	}

	for _, acc := range userAccount {
		if accountStruct.AccountType == acc.AccountType {
			newTotal := app.ConvertStringToDecimal(acc.Balance).Add(app.ConvertStringToDecimal(accountStruct.Amount))
			account, err := app.DB.Deposit(app.Ctx, database.DepositParams{
				Balance:   newTotal.String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				app.ServerError(c, err.Error())
			}
			fmt.Println(acc.AccountID)
			fmt.Println(accountStruct.Amount)
			fmt.Println(Deposit)
			err = app.SaveTransaction(c, acc.AccountID, accountStruct.Amount, Deposit)
			if err != nil {
				app.ServerError(c, err.Error())
			}
			jsonResp = Account{
				AccountType: account.AccountType,
				Amount:      account.Balance,
			}
		}

		break
	}

	return c.JSON(http.StatusOK, jsonResp)

}

func (app *Applicaton) Withdraw(c echo.Context) error {
	var accountStruct Account
	var jsonResp Account

	id := c.Param("id")
	err := c.Bind(&accountStruct)
	parsedId := app.ConvertStringToUuid(id)

	userAccount, err := app.DB.FindAccount(app.Ctx, parsedId)
	if err != nil {
		app.ServerError(c, "Failed to retrieve user accounts ")
	}

	for _, acc := range userAccount {
		if accountStruct.AccountType == acc.AccountType {
			newTotal := app.ConvertStringToDecimal(acc.Balance).Sub(app.ConvertStringToDecimal(accountStruct.Amount))
			account, err := app.DB.Withdraw(app.Ctx, database.WithdrawParams{
				Balance:   newTotal.String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				app.ServerError(c, "Failed to complete withdrawal")
			}
			err = app.SaveTransaction(c, acc.AccountID, accountStruct.Amount, Withdraw)
			if err != nil {
				app.ServerError(c, err.Error())
			}
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
	parsedId := app.ConvertStringToUuid(id)
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
	parsedId := app.ConvertStringToUuid(id)

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
			Balance:   app.ConvertStringToDecimal(account.Balance).Sub(app.ConvertStringToDecimal(fundInstance.Amount)).String(),
		})
		if err != nil {
			app.ServerError(c, err.Error())

		}
		err = qtx.CreditSaving(app.Ctx, database.CreditSavingParams{
			AccountID: account.AccountID,
			Balance:   app.ConvertStringToDecimal(account.Balance).Add(app.ConvertStringToDecimal(fundInstance.Amount)).String(),
		})
		if err != nil {
			app.ServerError(c, err.Error())
		}
		err = app.SaveTransactionFunds(c, account.AccountID, account.AccountID, fundInstance.Amount, TransferFunds)
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

func (app *Applicaton) TransferFunds(c echo.Context) error {

	var accountReceiving database.Account
	var accountSending database.Account

	type parameters struct {
		AccountType   string `json:"account_type"`
		AccountNumber string `json:"account_number"`
		Amount        string `json:"amount"`
	}

	params := parameters{}

	id := c.Param("id")
	parsedId := app.ConvertStringToUuid(id)
	err := c.Bind(&params)

	accounts, err := app.DB.FindAccount(app.Ctx, parsedId)
	if err != nil {
		app.ServerError(c, "Could not find the users accounts")
	}

	if err != nil {
		app.ServerError(c, err.Error())
	}

	tx, err := app.SDB.Begin()
	if err != nil {
		app.ServerError(c, err.Error())
	}
	defer tx.Rollback()

	qtx := app.DB.WithTx(tx)

	for _, account := range accounts {
		if account.AccountType == params.AccountType {
			accountSending, err = qtx.FindAccountById(app.Ctx, account.AccountID)
			if err != nil {
				app.ServerError(c, err.Error())
			}
		}
		break
	}

	accountReceiving, err = qtx.FindAccountByAccNo(app.Ctx, params.AccountNumber)

	if err != nil {
		app.ServerError(c, err.Error())
	}

	_, err = qtx.Withdraw(app.Ctx, database.WithdrawParams{
		Balance:   app.ConvertStringToDecimal(accountSending.Balance).Sub(app.ConvertStringToDecimal(params.Amount)).String(),
		AccountID: accountSending.AccountID,
	})

	_, err = qtx.Deposit(app.Ctx, database.DepositParams{
		Balance:   app.ConvertStringToDecimal(accountReceiving.Balance).Add(app.ConvertStringToDecimal(params.Amount)).String(),
		AccountID: accountReceiving.AccountID,
	})

	err = app.SaveTransactionFunds(c, accountSending.AccountID, accountReceiving.AccountID, params.Amount, TransferFunds)
	if err != nil {
		app.ServerError(c, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		app.ServerError(c, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"Success": "Transfer succesful"})
}
