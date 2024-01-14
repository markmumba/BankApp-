package main

import "time"

const (
	Checking      = "checking"
	Savings       = "savings"
	Deposit       = "deposit"
	Withdraw      = "withdraw"
	TransferFunds = "transfer_funds"
)


type User struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type Account struct {
	AccountType string `json:"account_type"`
	Amount      string `json:"amount"`
}

type Transaction struct {
	RecepientID int32     `json:"recepient_id"`
	Amount      string    `json:"amount"`
	Type        string    `json:"type"`
	Timestamp   time.Time `json:"timestamp"`
}
