package api

import (
	"e-wallet/api/handlers"
	"e-wallet/config"
	"e-wallet/storage/repo"

	"github.com/gin-gonic/gin"
)


type Options struct {
	Cfg config.Config
	Repo repo.Repo
	// auth
}

func New(options Options) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())

	handler := handlers.NewHandler(options.Cfg,options.Repo)

	router.GET("/",)

}