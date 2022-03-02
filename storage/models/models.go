package models

type User struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type NotIdentifiedUser struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Wallet struct {
	Id      string `json:"id"`
	Balance string `json:"balance"`
}

type WalletHistory struct {
	Id                     string `json:"id"`
	CurrentBalance         string `json:"current_balance"`
	TotalIncome            string `json:"total_income"`
	TotalExpense           string `json:"total_expense"`
	TotalIncomeOperations  int64  `json:"total_income_operations"`
	TotalExpenseOperations int64  `json:"tatal_expense_operations"`
	TotalOperations        int64  `json:"total_operations"`
}
