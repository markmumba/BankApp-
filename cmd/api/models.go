package main

const (
	Checking = "checking"
	Savings = "savings"
)
type User struct {
	UserName string `json:"user_name"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}


type Account struct {

}