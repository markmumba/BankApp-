package main

import (
	"database/sql"
	"fmt"

	"github.com/markmumba/chasebank/internal/database"
)


func (app *Applicaton) SaveTransaction (accoutId int32,recipientId int32,amount string , typeTransaction string ){

	err := app.DB.SaveTransaction(app.Ctx,database.SaveTransactionParams{
		AccountID: sql.NullInt32{Int32:accoutId ,Valid: true},
		RecepientID: sql.NullInt32{Int32: recipientId,Valid: true},
		Amount: amount,
		Type: typeTransaction,

	}) 
	if err != nil {
		fmt.Println(err.Error())
	}

}
