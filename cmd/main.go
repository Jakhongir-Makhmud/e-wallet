package main

import (
	// "github.com/go-redis/redis"
	"e-wallet/api"
	"e-wallet/api/auth"
	"e-wallet/config"
	"e-wallet/storage"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

func main() {

	cfg := config.LoadCfg()
	psqlCred := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	cfg.PostgresHost,
	cfg.PostgresPort,
	cfg.PostgresUser,
	cfg.PostgresPass,
	cfg.PostgresDB,
)
	psql,err := sqlx.Connect("postgres",psqlCred)
	if err != nil {
		panic(err)
	}
	redis := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",cfg.RedisHost,cfg.RedisPort),
		DB: 0,
		Password: "",
	})

	storagePool := storage.NewStorage(psql)

	
	server := api.New(api.Options{
		Cfg: cfg,
		Repo: storagePool.Storage(),
		Redis: redis,
		Auth: auth.Auth{Cfg: cfg},
	})


	err = server.Run(cfg.Port)
	if err != nil {
		panic(err)
	}

}
