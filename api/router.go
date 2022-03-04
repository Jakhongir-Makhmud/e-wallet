package api

import (
	"e-wallet/api/auth"
	"e-wallet/api/handlers"
	"e-wallet/config"
	"e-wallet/storage/repo"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Options struct {
	Cfg   config.Config
	Repo  repo.Repo
	Redis *redis.Client
	Auth  auth.Auth
}
// @title E-Wallet 
// @securitydefinitions.oauth2.accessCode Digest
// @in header
// @name Authorization

func New(options Options) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(options.Auth.Auth)

	handler := handlers.NewHandler(options.Cfg, options.Repo, options.Redis, options.Auth)

	router.POST("/check/wallet/exist",handler.CheckWalletExists)
	router.POST("/wallet/balance",handler.GetBalance)
	router.POST("/wallet/history",handler.GetHistory)
	router.POST("/wallet/fill",handler.FillWallet)

	return router
}
