package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
	}

	JWT struct {
		Secret          string
		AccessDuration  time.Duration
		RefreshDuration time.Duration
	}

	BackendPort     string
	SaltLength      int
	MinPasswordSize int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file didn't found")
	}

	cfg := &Config{}

	cfg.BackendPort = os.Getenv("BACKEND_PORT")

	saltLength, err := strconv.Atoi(os.Getenv("SALT_LENGTH"))
	if err != nil {
		log.Fatal("invalid SALT_LENGTH")
	}
	cfg.SaltLength = saltLength

	minPasswordSize, err := strconv.Atoi(os.Getenv("MIN_PASSWORD_SIZE"))
	if err != nil {
		log.Fatal("invalid MIN_PASSWORD_SIZE")
	}
	cfg.MinPasswordSize = minPasswordSize

	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Database = os.Getenv("POSTGRES_DB")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	access, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		log.Fatal("invalid ACCESS_TOKEN_DURATION")
	}
	cfg.JWT.AccessDuration = access
	refresh, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		log.Fatal("invalid REFRESH_TOKEN_DURATION")
	}
	cfg.JWT.RefreshDuration = refresh

	return cfg
}
