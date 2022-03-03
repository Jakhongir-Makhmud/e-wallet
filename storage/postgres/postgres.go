package postgres

import (
	"e-wallet/storage/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	Limit        float64 = 100000
	LimitNotIden float64 = 10000
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
	return &models.Wallet{Id: w.Id, Balance: w.Balance}, nil
}

// This method checks wheather wallet exists or not, if not it returns error and Wallet with null fields
func (d Database) CheckWalletExists(w models.Wallet) (*models.Wallet, error) {
	var exists int
	query := `select count(*) where wallet_id = $1 and delted_at is null`
	err := d.db.QueryRow(query, w.Id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if exists != 1 {
		return &models.Wallet{Id: "", Balance: 0}, err
	}

	return d.GetBalance(w)
}

// This method returns total operations and total expense and total income for current month
func (d Database) GetTotals(w models.Wallet) (*models.WalletHistory, error) {
	var walletHistory = models.WalletHistory{}
	query := `SELECT balance, SUM(wi.amount),SUM(we.amount),COUNT(wi.amount),COUNT(we.amount) FROM wallets w JOIN wallets_income wi USING(wallet_id) JOIN wallets_expenses we USING(wallet_id) GROUP BY wallet_id having wallet_id = $1 and month(we.created_at) = month(now()) and month(wi.created_at) = month(now())`

	err := d.db.DB.QueryRow(query, w.Id).Scan(
		&walletHistory.CurrentBalance,
		&walletHistory.TotalIncome,
		&walletHistory.TotalExpense,
		&walletHistory.TotalIncomeOperations,
		&walletHistory.TotalExpenseOperations,
	)

	if err != nil {
		return nil, err
	}

	walletHistory.Id = w.Id
	walletHistory.TotalOperations = walletHistory.TotalExpenseOperations + walletHistory.TotalIncomeOperations

	return &walletHistory, nil

}

// This method is used to FillWallet
func (d Database) FillWallet(w models.WalletFill) (*models.Wallet, error) {
	isIdentified, err := d.isIdentified(w.Id)
	if err != nil {
		return nil, err
	}

	wallet, err := d.GetBalance(models.Wallet{Id: w.Id})

	currentBalance := wallet.Balance + w.Amount

	if !isIdentified {
		if currentBalance > LimitNotIden {
			return nil, fmt.Errorf("You are not identified user, so your limit is %v", LimitNotIden)
		}
	}

	if currentBalance > Limit {
		return nil, fmt.Errorf("Your limit is %v", Limit)
	}

	if isIdentified == false {
		return nil, nil
	}
	queryWallet := `UPDATE wallets SET balance = balance + $2, updated_at = $3 WHERE wallet_id = $1 AND deleted_at deleted_at IS NULL`
	queryIncome := `INSERT INTO wallets_income (income_id,wallet_id,amount,created_at) values ($1,$2,$3,$4)`

	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	now := time.Now().Format(time.RFC3339)

	_, err = tx.Exec(queryWallet, w.Id, w.Amount, now)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(queryIncome)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return d.GetBalance(models.Wallet{Id: w.Id})

}

func (d Database) isIdentified(id string) (bool, error) {
	query := `SELECT is_identified FROM wallets WHERE wallet_id = $1`
	var isIdentified bool
	err := d.db.QueryRow(query, id).Scan(&isIdentified)
	if err != nil {
		return false, err
	}

	return isIdentified, nil
}
