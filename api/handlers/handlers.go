package handlers

import (
	"e-wallet/api/auth"
	"e-wallet/config"
	"e-wallet/storage/repo"

	"github.com/go-redis/redis/v8"
)

type handlers struct {
	auth auth.Auth
	cfg  config.Config
	repo repo.Repo
	redis *redis.Client
}

func NewHandler(cfg config.Config, repo repo.Repo,redis *redis.Client,auth auth.Auth) *handlers {
	return &handlers{
		cfg:  cfg,
		repo: repo,
		redis: redis,
		auth: auth,
	}
}
