package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Dsn       string
	JwtSecret string
}

func LoadEnv() *Config {

	// Load .env file if it exists (optional for deployment)
	godotenv.Load()

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=gotickets port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		Dsn:       dsn,
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}
