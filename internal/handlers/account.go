package handlers 

//FIXME  refuse withdrawal when money exceeds balance
//TODO view all accounts
// TODO return balance as json after the transaction
import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/internal/database"
)

func (app *Applicaton) CreateSavingAccount(c echo.Context) error {
	var jsonResp string

	id := app.GetUserIdFromToken(c)
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

func (app *Applicaton) GetUserAccounts(c echo.Context) error {

	accountDetails := map[string]map[string]string{}


	id := app.GetUserIdFromToken(c)
	parseId := app.ConvertStringToUuid(id)

	userAccounts := app.FindAccountHelper(c, parseId)

	for _, account := range userAccounts {
		innermap := map[string]string{}
		innermap[account.AccountNumber] = account.Balance
		accountDetails[account.AccountType] = innermap
	}

	return c.JSON(http.StatusOK, accountDetails)

}

func (app *Applicaton) Deposit(c echo.Context) error {
	var accountStruct Account
	var jsonResp Account

	id := app.GetUserIdFromToken(c)
	err := c.Bind(&accountStruct)
	if err != nil {
		app.ServerError(c, "Failed to get account details")
	}
	parsedId := app.ConvertStringToUuid(id)

	userAccounts := app.FindAccountHelper(c, parsedId)

	for _, acc := range userAccounts {
		if accountStruct.AccountType == acc.AccountType {
			account, err := app.DB.Deposit(app.Ctx, database.DepositParams{
				Balance:   app.DepositHelper(acc.Balance, accountStruct.Amount).String(),
				AccountID: acc.AccountID,
			})
			if err != nil {
				app.ServerError(c, err.Error())
			}
			err = app.SaveTransaction(c, acc.AccountID, accountStruct.Amount, Deposit)
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

func (app *Applicaton) Withdraw(c echo.Context) error {
	var accountStruct Account
	var jsonResp Account

	err := c.Bind(&accountStruct)
	if err != nil {
		app.ServerError(c, "failed to get account details")
	}
	id := app.GetUserIdFromToken(c)
	parsedId := app.ConvertStringToUuid(id)
	userAccounts := app.FindAccountHelper(c, parsedId)

	for _, acc := range userAccounts {
		if accountStruct.AccountType == acc.AccountType {
			if app.CheckBalance(acc.Balance, accountStruct.Amount) {
				account, err := app.DB.Withdraw(app.Ctx, database.WithdrawParams{
					Balance:   app.WithdrawHelper(acc.Balance, accountStruct.Amount).String(),
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
			} else {
				return c.JSON(http.StatusOK, map[string]string{"invalid": "Balance is not enough to make transaction"})
			}
		}
	}

	return c.JSON(http.StatusOK, jsonResp)

}

func (app *Applicaton) TransferCheckingToSaving(c echo.Context) error {

	var accountInstance Account

	id := app.GetUserIdFromToken(c)
	parseId := app.ConvertStringToUuid(id)
	err := c.Bind(&accountInstance)

	userAccounts := app.FindAccountHelper(c, parseId)

	tx, err := app.SDB.Begin()
	if err != nil {
		app.ServerError(c, err.Error())
	}
	defer tx.Rollback()

	qtx := app.DB.WithTx(tx)

	for _, account := range userAccounts {

		err := qtx.CheckToSave(app.Ctx, database.CheckToSaveParams{
			Balance:     app.WithdrawHelper(account.Balance, accountInstance.Amount).String(),
			AccountID:   account.AccountID,
			AccountType: Checking,
		})
		if err != nil {
			app.ServerError(c, "failed to  withdraw from checking account")
		}

		err = qtx.CheckToSave(app.Ctx, database.CheckToSaveParams{
			Balance:     app.DepositHelper(account.Balance, accountInstance.Amount).String(),
			AccountID:   account.AccountID,
			AccountType: Savings,
		})
		if err != nil {
			app.ServerError(c, "failed to deposit funds in savings account ")
		}
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

	id := app.GetUserIdFromToken(c)
	parsedId := app.ConvertStringToUuid(id)
	err := c.Bind(&params)
	if err != nil {
		app.ServerError(c, err.Error())
	}

	userAccounts := app.FindAccountHelper(c, parsedId)

	tx, err := app.SDB.Begin()
	if err != nil {
		app.ServerError(c, err.Error())
	}
	defer tx.Rollback()

	qtx := app.DB.WithTx(tx)

	for _, account := range userAccounts {
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
	if app.CheckBalance(accountSending.Balance, params.Amount) {

		_, err = qtx.Withdraw(app.Ctx, database.WithdrawParams{
			Balance:   app.WithdrawHelper(accountSending.Balance, params.Amount).String(),
			AccountID: accountSending.AccountID,
		})

		_, err = qtx.Deposit(app.Ctx, database.DepositParams{
			Balance:   app.DepositHelper(accountReceiving.Balance, params.Amount).String(),
			AccountID: accountReceiving.AccountID,
		})

		err = app.SaveTransactionFunds(c, accountSending.AccountID, accountReceiving.AccountID, params.Amount, TransferFunds)
		if err != nil {
			app.ServerError(c, err.Error())
		}
	} else {
		app.ServerError(c, "No enough Funds to transfer")
	}

	err = tx.Commit()
	if err != nil {
		app.ServerError(c, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"Success": "Transfer succesful"})
}