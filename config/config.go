package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	PostgresHost string
	PostgresPort int
	PostgresUser string
	PostgresPass string
	PostgresDB string
	RedisHost string
	RedisPort int
	Port         string
	SecretKey []byte
}

func LoadCfg() Config {

	var cfg = Config{}
	
	cfg.Port = cast.ToString(getOrReturnDefault("HOST",":8080"))
	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	cfg.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_USER", "jakhongir"))
	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1"))
	cfg.PostgresDB = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "ewallet"))
	cfg.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	cfg.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))
	cfg.SecretKey = []byte(getOrReturnDefault("SECRET_KEY","secret").(string))

	return cfg
}

func getOrReturnDefault(env_var string, def interface{}) interface{} {
	val, exists := os.LookupEnv(env_var)
	if exists {
		return val
	}
	return def
}
