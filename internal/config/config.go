package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Dsn       string
	JwtSecret string
}

func LoadEnv() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
