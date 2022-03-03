package storage

import (
	"github.com/jmoiron/sqlx"
	"e-wallet/storage/repo"
	"e-wallet/storage/postgres"
)

type storagePool struct {
	db *sqlx.DB
	storage *postgres.Database
}

type StorageI interface {
	Storage() repo.Repo
}

func NewStorage(db *sqlx.DB) *storagePool {
	return &storagePool{
		db: db,
		storage: postgres.NewDatabase(db),
	}
}

func (s *storagePool) Storage() *postgres.Database {
	return s.storage
}

