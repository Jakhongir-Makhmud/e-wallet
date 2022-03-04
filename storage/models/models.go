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
	Balance float64 `json:"balance"`
}

type WalletHistory struct {
	Id                     string `json:"id"`
	CurrentBalance         float64 `json:"current_balance"`
	TotalIncome            float64 `json:"total_income"`
	TotalExpense           float64 `json:"total_expense"`
	TotalIncomeOperations  int64  `json:"total_income_operations"`
	TotalExpenseOperations int64  `json:"tatal_expense_operations"`
	TotalOperations        int64  `json:"total_operations"`
}

type WalletFill struct {
	Id string `json:"id"`
	Amount float64 `json:"amount"`
}

type Err struct {
	Error string `json:"error"`
}