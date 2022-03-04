package repo

import "e-wallet/storage/models"

type Repo interface {
	GetBalance(w models.Wallet) (*models.Wallet, error)
	CheckWalletExists(w models.Wallet) (*models.Wallet, error)
	GetHistory(w models.Wallet) (*models.WalletHistory, error)
	FillWallet(w models.WalletFill) (*models.Wallet, error)
	CheckUserById(id string) (bool,error)
}
