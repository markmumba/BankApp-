package handlers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Checking      = "checking"
	Savings       = "savings"
	Deposit       = "deposit"
	Withdraw      = "withdraw"
	TransferFunds = "transferFunds"
)

type CustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type User struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Account struct {
	AccountType string `json:"account_type"`
	Amount      string `json:"amount"`
}

type Transaction struct {
	RecepientAccount string    `json:"recepient_account"`
	Amount           string    `json:"amount"`
	Type             string    `json:"type"`
	Timestamp        time.Time `json:"timestamp"`
	AccountType      string    `json:"account_type"`
}
