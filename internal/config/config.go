package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbPath string
	AppEnv string
}

func Load() (*Config, error) {
	godotenv.Load()

	dbPath := os.Getenv("DB_PATH")
	appEnv := os.Getenv("APP_ENV")

	if dbPath == "" {
		return nil, errors.New("DB_URL must be set")
	}

	if appEnv == "" {
		return nil, errors.New("APP_ENV must be set")
	}

	cfg := Config{
		DbPath: dbPath,
		AppEnv: appEnv,
	}

	return &cfg, nil
}
