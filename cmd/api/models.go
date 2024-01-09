package main

const (
	Checking = "checking"
	Savings  = "savings"
)

type Response struct {
	Message string `json:"message"`
}
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
