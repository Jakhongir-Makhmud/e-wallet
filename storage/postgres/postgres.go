package postgres

import (
	"e-wallet/storage/models"
	"fmt"
	"strings"
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
func (d Database) GetHistory(w models.Wallet) (*models.WalletHistory, error) {
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

// This method is used to Fill Wallet
func (d Database) FillWallet(w models.WalletFill) (*models.Wallet, error) {
	isIdentified, err := d.isUserIdentified(w.Id)
	if err != nil {
		return nil, err
	}

	wallet, err := d.GetBalance(models.Wallet{Id: w.Id})

	if err != nil {
		return nil, err
	}

	currentBalance := wallet.Balance + w.Amount

	if !isIdentified {
		if currentBalance > LimitNotIden {
			return nil, fmt.Errorf("you are not identified user, so your limit is %v", LimitNotIden)
		}
	}

	if currentBalance > Limit {
		return nil, fmt.Errorf("your limit is %v", Limit)
	}

	if !isIdentified {
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
// This mehtod is used to check user via email address wheather we have such user or not, if yes the result is true
func (d Database) CheckUser(email string) (bool, error) {
	query := `SELECT COUNT(*) WHERE email = $1`
	email = strings.TrimSpace(email)
	var isExists int
	err := d.db.QueryRow(query, email).Scan(&isExists)

	if err != nil {
		return false, err
	}

	if isExists == 1 {
		return true, nil
	}

	return false, nil
}

// This method checks where user is identified or not
func (d Database) isUserIdentified(id string) (bool, error) {
	query := `SELECT is_identified FROM wallets WHERE wallet_id = $1`
	var isIdentified bool
	err := d.db.QueryRow(query, id).Scan(&isIdentified)
	if err != nil {
		return false, err
	}

	return isIdentified, nil
}

// This method checks wheather do we have such user or not, if not result is false
func (d Database) CheckUserById(id string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE user_id = $1 AND deleted_at IS NULL`
	err := d.db.QueryRow(query, id).Scan(&count)

	if err != nil {
		return false, err
	}
	if count == 0 {
		return true, nil
	}
	return false, nil

}

// This method to create new user
func (d Database) NewUser(u models.User) (*models.User, error) {

	query := `INSERT INTO users (user_id,first_name,last_name,email,created_at) VALUES (
		$1,$2,$3,$4,$5
	)`
	var (
		now = time.Now().Format(time.RFC3339)
		id  string
	)
	err := d.db.QueryRow(query, u.Id, u, u.FirstName, u.LastName, u.Email, now).Scan(&id)
	if err != nil {
		return nil, err
	}
	if id != "" {
		return &u, nil
	}
	return nil, nil
}

// This method creates new wallet
func (d Database) NewWallet(nw models.NewWallet) (*models.Wallet,error) {
	isIden,err := d.isUserIdentified(nw.UserId)
	if err != nil {
			return nil,err
	}
	now := time.Now().Format(time.RFC3339)
	query := `INSERT INTO wallets (wallet_id,user_id,is_identified,created_at) VALUES (
		$1,$2,$3,$4
	)`
	_,err = d.db.Exec(query,nw.WalletId,nw.UserId,isIden,now)
	
	if err != nil {
		return nil,err
	}

	return &models.Wallet{Id: nw.WalletId,Balance: 0},nil
}