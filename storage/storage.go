package storage

import (
	"e-wallet/storage/postgres"
	"e-wallet/storage/repo"

	"github.com/jmoiron/sqlx"
)

type storagePool struct {
	db      *sqlx.DB
	storage *postgres.Database
}

type StorageI interface {
	Storage() repo.Repo
}

func NewStorage(db *sqlx.DB) *storagePool {
	return &storagePool{
		db:      db,
		storage: postgres.NewDatabase(db),
	}
}

func (s *storagePool) Storage() *postgres.Database {
	return s.storage
}
