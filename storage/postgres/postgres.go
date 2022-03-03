package postgres

import (
	"e-wallet/storage/models"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		db: db,
	}
}
// This method returs balance of the wallet
func (d Database) GetBalance(w models.Wallet) (*models.Wallet, error) {

	query := `select balance where wallet_id = $1 and delted_at is null`

	err := d.db.QueryRow(query, w.Id).Scan(&w.Balance)
	if err != nil {
		return nil, err
	}
	return &models.Wallet{w.Id, w.Balance}, nil
}
// This method checks wheather wallet exists or not, if not it returns error and Wallet with null fields
func (d Database) CheckWalletExists(w models.Wallet) (*models.Wallet, error) {
	var exists int
	query := `select count(*) where wallet_id = $1 and delted_at is null`
	err := d.db.QueryRow(query,w.Id).Scan(&exists)
	if err != nil {
		return nil,err
	}

	if exists != 1 {
		return &models.Wallet{Id: "",Balance: 0},err
	}

	return d.GetBalance(w)
}

// This method returns total operations and total expense and total income for current month
func (d Database) GetTotals(w models.Wallet) (*models.WalletHistory,error) {
	var walletHistory = models.WalletHistory{}
	query := `SELECT balance, SUM(wi.amount),SUM(we.amount),COUNT(wi.amount),COUNT(we.amount) FROM wallets w JOIN wallets_income wi USING(wallet_id) JOIN wallets_expenses we USING(wallet_id) GROUP BY wallet_id having wallet_id = $1`

	err := d.db.DB.QueryRow(query,w.Id).Scan(
		&walletHistory.CurrentBalance,
		&walletHistory.TotalIncome,
		&walletHistory.TotalExpense,
		&walletHistory.TotalIncomeOperations,
		&walletHistory.TotalExpenseOperations,
	)

	if err != nil {
		return nil,err
	}

	walletHistory.Id = w.Id
	walletHistory.TotalOperations = walletHistory.TotalExpenseOperations + walletHistory.TotalIncomeOperations

	return &walletHistory,nil

}

