package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	PostgresHost string
	PostgresPort string
}

func LoadCfg() Config {

	var cfg = Config{}

	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST","localhost"))
	cfg.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT","5432"))
	
	return cfg
}

func getOrReturnDefault(env_var,def string) interface{} {
	val, exists := os.LookupEnv(env_var)
	if exists {
		return val
	}
	return def
}