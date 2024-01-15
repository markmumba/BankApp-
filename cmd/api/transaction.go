package main

//TODO get all transactions

import (
	"database/sql"
	"fmt"

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
