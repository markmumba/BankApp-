package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/chasebank/config"
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

	userList := []User{}

	users, err := app.DB.GetAllUsers(app.Ctx)
	if err != nil {
		app.ServerError(c, err.Error())
	}
	for _, user := range users {
		userList = append(userList, app.ConvertDBUser(user))
	}

	return c.JSON(http.StatusOK, userList)
}

func (app *Applicaton) GetUser(c echo.Context) error {
	accountDetailList := []map[string]string{}
	id := app.GetUserIdFromToken(c)

	parseId := app.ConvertStringToUuid(id)

	user, err := app.DB.FindUser(app.Ctx, parseId)
	if err != nil {
		app.ServerError(c, "Failed to locate user")
	}
	accounts, err := app.DB.FindAccount(app.Ctx, parseId)
	if err != nil {
		app.ServerError(c, "could not find accounts ")
	}
	for _, account := range accounts {

		accountDetails := map[string]string{
			"accountNumber": account.AccountNumber,
			"accountType":   account.AccountType,
			"balance":       account.Balance,
			"dateJoined":    account.DateOpened.Time.String(),
		}

		accountDetailList = append(accountDetailList, accountDetails)

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":    map[string]string{"username": user.Username, "fullname": user.FullName},
		"account": accountDetailList,
	})
}

func (app *Applicaton) Login(c echo.Context) error {

	var user User
	err := c.Bind(&user)
	if err != nil {
		app.ServerError(c, "unable to get details")
	}
	userDetalis, err := app.DB.FindUserByEmail(app.Ctx, user.Email)
	if err != nil {
		app.ServerError(c, "Failed to retrive user details")

	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetalis.PasswordHash), []byte(user.Password))
	if err != nil {
		app.ServerError(c, "password does not match ")
	}

	claims := &CustomClaims{
		userDetalis.UserID.String(),
		jwt.RegisteredClaims{
			Issuer:    userDetalis.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config("SECRET_KEY")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 72),
		HttpOnly: true,
	})

	user.UserId = userDetalis.UserID.String()
	user.Token = t
	user.Password = ""
	return c.JSON(http.StatusOK, user)

}
func (app *Applicaton) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Logout successful",
	})
}
