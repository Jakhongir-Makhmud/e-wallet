package postgres

import (
	"database/sql"
	"e-wallet/storage/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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

	query := `SELECT balance FROM wallets WHERE wallet_id = $1 AND deleted_at IS NULL`
	var balance float64
	err := d.db.QueryRow(query, w.Id).Scan(&balance)
	if err != nil {
		fmt.Println(err, "34")

		return nil, err
	}
	return &models.Wallet{Id: w.Id, Balance: balance}, nil
}

// This method checks wheather wallet exists or not, if not it returns error and Wallet with null fields
func (d Database) CheckWalletExists(w models.Wallet) (*models.Wallet, error) {
	var exists int
	query := `SELECT COUNT(*) FROM wallets WHERE wallet_id = $1 AND deleted_at IS NULL`
	err := d.db.QueryRow(query, w.Id).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if exists != 1 {
		return &models.Wallet{Id: "", Balance: 0}, err
	}

	return d.GetBalance(w)
}

// This method returns total operations and total expense and total income for current month
func (d Database) GetHistory(w models.Wallet) (*models.WalletHistory, error) {
	var (
		walletHistory = models.WalletHistory{}
		now           = time.Now().Month()
	)
	var queryIncome = `SELECT balance, SUM(wi.amount), COUNT(wi.amount)
	FROM wallets w 
	JOIN wallets_incomes wi USING(wallet_id)
	where 
	wallet_id = $1 
	and 
	extract(month from wi.created_at) = $2
	group by wallet_id;
`
	var queryExpense = `SELECT SUM(we.amount), COUNT(we.amount)
	FROM wallets w 
	JOIN wallets_expenses we USING(wallet_id)
	where 
	wallet_id = $1 
	and 
	extract(month from we.created_at) = $2
	group by wallet_id;
`

	month := int(now)

	err := d.db.DB.QueryRow(queryIncome, w.Id, month).Scan(
		&walletHistory.CurrentBalance,
		&walletHistory.TotalIncome,
		&walletHistory.TotalIncomeOperations,
	)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)

		return nil, err
	}
	err = d.db.DB.QueryRow(queryExpense, w.Id, month).Scan(
		&walletHistory.TotalExpense,
		&walletHistory.TotalExpenseOperations,
	)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)

		return nil, err
	}

	walletHistory.Id = w.Id
	walletHistory.TotalOperations = walletHistory.TotalExpenseOperations + walletHistory.TotalIncomeOperations

	return &walletHistory, nil

}

// This method is used to Fill Wallet
func (d Database) FillWallet(w models.WalletFill) (*models.Wallet, error) {
	isIdentified, err := d.isUserIdentified(w.UserId)
	if err != nil {
		fmt.Println(err)
		fmt.Println(isIdentified, "fill wallet 89")
		return nil, err
	}

	wallet, err := d.GetBalance(models.Wallet{Id: w.Id})

	if err != nil {
		fmt.Println(err, "get Balance 96")

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

	queryWallet := `UPDATE wallets SET balance = balance + $2, updated_at = $3 WHERE wallet_id = $1 AND deleted_at IS NULL`
	queryIncome := `INSERT INTO wallets_incomes (income_id,wallet_id,amount,created_at) values ($1,$2,$3,$4)`

	tx, err := d.db.Begin()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	now := time.Now().Format(time.RFC822)

	_, err = tx.Exec(queryWallet, w.Id, w.Amount, now)
	if err != nil {
		fmt.Println(err, "127")
		tx.Rollback()
		return nil, err
	}
	incodeId := uuid.New().String()
	_, err = tx.Exec(queryIncome, incodeId, w.Id, w.Amount, now)
	if err != nil {
		fmt.Println(err, "134")
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		fmt.Println(err)
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
	query := `SELECT is_identified FROM users WHERE user_id = $1`
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
		fmt.Println(err)

		return false, err
	}
	if count == 0 {
		return true, nil
	}
	return false, nil

}

// This method to create new user
func (d Database) NewUser(u models.User) (*models.User, error) {
	// insert query for not identified users
	query := `INSERT INTO users (user_id,first_name,last_name,email,created_at) VALUES (
		$1,$2,$3,$4,$5
	)`

	// insert query for identified users, who provides phone number
	queryIden := `INSERT INTO users (user_id,first_name,last_name,email,phone,is_identified,created_at) VALUES (
		$1,$2,$3,$4,$5,$6,$7
	)`

	var (
		now = time.Now().Format(time.RFC822)
		err error
	)
	// First statement is for not identidied users, second for indentified
	if u.PhoneNumber == "" {
		_, err = d.db.Exec(query, u.Id, u.FirstName, u.LastName, u.Email, now)
	} else {
		_, err = d.db.Exec(queryIden, u.Id, u.FirstName, u.LastName, u.Email, u.PhoneNumber, true, now)
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// This method creates new wallet
func (d Database) NewWallet(nw models.NewWallet) (*models.Wallet, error) {
	isIden, err := d.isUserIdentified(nw.UserId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	now := time.Now().Format(time.RFC822)
	query := `INSERT INTO wallets (wallet_id,user_id,is_identified,created_at) VALUES (
		$1,$2,$3,$4
	)`
	_, err = d.db.Exec(query, nw.WalletId, nw.UserId, isIden, now)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &models.Wallet{Id: nw.WalletId, Balance: 0}, nil
}
