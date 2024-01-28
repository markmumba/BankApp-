package handlers

//TODO get all transactions make viewtransaction to return all
//TODO get for specific accounts

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/internal/database"
)

func (app *Applicaton) SaveTransaction(c echo.Context, accoutId int32, amount string, typeTransaction string) error {

	err := app.DB.SaveTransaction(app.Ctx, database.SaveTransactionParams{

		AccountID: sql.NullInt32{Int32: accoutId, Valid: true},
		Amount:    amount,
		Type:      typeTransaction,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil

}

func (app *Applicaton) SaveTransactionFunds(c echo.Context, accoutId int32, recipientId int32, amount string, typeTransaction string) error {

	err := app.DB.SaveTransactionFunds(app.Ctx, database.SaveTransactionFundsParams{

		AccountID:   sql.NullInt32{Int32: accoutId, Valid: true},
		RecepientID: sql.NullInt32{Int32: recipientId, Valid: true},
		Amount:      amount,
		Type:        typeTransaction,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (app *Applicaton) ViewTransactions(c echo.Context) error {
	var allTransactions []database.ViewTransactionsRow
	var jsonResp []Transaction

	id := app.GetUserIdFromToken(c)
	parsedId := app.ConvertStringToUuid(id)
	userAccounts := app.FindAccountHelper(c, parsedId)

	for _, acc := range userAccounts {
		transactions, err := app.DB.ViewTransactions(app.Ctx, sql.NullInt32{Int32: acc.AccountID, Valid: true})
		if err != nil {
			app.ServerError(c, "Failed to retrieve the transactions")
			continue 
		}
		allTransactions = append(allTransactions, transactions...)
	}

	for _, transaction := range allTransactions {
		var recepientAccount string

		if transaction.Type == "deposit" || transaction.Type == "withdraw" {
			recepientAccount = ""
		} else {
			account, err := app.DB.FindAccountById(app.Ctx, transaction.RecepientID.Int32)
			if err != nil {
				app.ServerError(c, "Unable to find recipient account")
				continue 
			}
			recepientAccount = account.AccountNumber
		}

		newTransaction := Transaction{
			RecepientAccount: recepientAccount,
			Amount:           transaction.Amount,
			Type:             transaction.Type,
			Timestamp:        transaction.Timestamp.Time,
			AccountType:      transaction.AccountType,
		}
		jsonResp = append(jsonResp, newTransaction)
	}

	return c.JSON(http.StatusOK, jsonResp)
}
