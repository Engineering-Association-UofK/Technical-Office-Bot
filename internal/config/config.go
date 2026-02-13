package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var App Config

type Config struct {
	TelegramToken string `env:"TELEGRAM_API_TOKEN,required"`
	Port          string `env:"PORT,required"`

	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

	JwtSecret string `env:"JWT_SECRET"`
}

func Load() error {
	godotenv.Load()
	return env.Parse(&App)
}
