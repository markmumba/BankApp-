package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/internal/database"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

type newuserDetails struct {
	UserName      string          `json:"user_name"`
	Email         string          `json:"email"`
	FullName      string          `json:"full_name"`
	AccountNumber string          `json:"account_number"`
	AccountType   string          `json:"account_type"`
	Balance       decimal.Decimal `json:"balance"`
}

func (app *Applicaton) CreateUser(c echo.Context) error {

	var user User
	err := c.Bind(&user)
	if err != nil {
		fmt.Println("was unable to bind to the struct")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		fmt.Println(err.Error())
	}

	result, err := app.DB.CreateUser(app.ctx, database.CreateUserParams{
		Username:     user.UserName,
		Email:        user.Email,
		FullName:     user.FullName,
		PasswordHash: string(hashedPassword),
	})

	account, err := app.DB.CreateAccount(app.ctx, database.CreateAccountParams{
		UserID:      uuid.NullUUID{UUID: result.UserID, Valid: true},
		AccountType: Checking,
	})

	balance, err := decimal.NewFromString(account.Balance)
	if err != nil {
		fmt.Println("Unable to convert to decimal")
	}

	newUserDetails := newuserDetails{
		UserName:      result.Username,
		Email:         result.Email,
		FullName:      result.FullName,
		AccountNumber: account.AccountNumber,
		AccountType:   account.AccountType,
		Balance:       balance,
	}

	return c.JSON(http.StatusCreated, newUserDetails)

}
