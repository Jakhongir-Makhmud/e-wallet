package handlers

import (
	"e-wallet/config"
	"e-wallet/storage/repo"
)

type handlers struct {
	// auth
	cfg  config.Config
	repo repo.Repo
}

func NewHandler(cfg config.Config, repo repo.Repo) *handlers {
	return &handlers{
		cfg:  cfg,
		repo: repo,
	}
}
