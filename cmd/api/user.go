package main


//TODO get user accounts 
//TODO 

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
		app.ServerError(c, "Failed to bind to user struct")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		fmt.Println(err.Error())
	}

	result, err := app.DB.CreateUser(c.Request().Context(), database.CreateUserParams{
		UserID:       uuid.New(),
		Username:     user.UserName,
		Email:        user.Email,
		FullName:     user.FullName,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		app.ServerError(c, err.Error())
	}

	account, err := app.DB.CreateAccount(c.Request().Context(), database.CreateAccountParams{
		UserID:      uuid.NullUUID{UUID: result.UserID, Valid: true},
		AccountType: Checking,
	})

	if err != nil {
		app.ServerError(c, err.Error())
	}

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



func (app *Applicaton) GetAllUsers(c echo.Context) error {

	userList := []database.User{}

	users, err := app.DB.GetAllUsers(app.Ctx)
	if err != nil {
		app.ServerError(c, err.Error())
	}
	for _, user := range users {
		userList = append(userList, user)
	}

	return c.JSON(http.StatusOK, userList)
}
