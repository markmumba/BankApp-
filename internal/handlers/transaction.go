package handlers 

//TODO get all transactions

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

	var accountType AccountType
	var transactions []database.Transaction
	var jsonResp []Transaction

	id := c.Param("id")
	err := c.Bind(&accountType)
	parsedId := app.ConvertStringToUuid(id)
	userAccounts := app.FindAccountHelper(c, parsedId)

	found := false
	for _, acc := range userAccounts {
		if accountType.Type == acc.AccountType {
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
		account, err := app.DB.FindAccountById(app.Ctx, transaction.RecepientID.Int32)
		if err != nil {
			app.ServerError(c, "unable to find recipeint account")
		}
		newTransaction := Transaction{
			RecepientAccount: account.AccountNumber,
			Amount:           transaction.Amount,
			Type:             transaction.Type,
			Timestamp:        transaction.Timestamp.Time,
		}
		jsonResp = append(jsonResp, newTransaction)
	}
	return c.JSON(http.StatusOK, jsonResp)
}
