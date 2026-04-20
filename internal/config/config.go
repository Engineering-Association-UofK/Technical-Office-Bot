package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var App Config

type Config struct {
	TelegramToken   string `env:"TELEGRAM_API_TOKEN,required"`
	TelegramChannel int64  `env:"TELEGRAM_CHANNEL_ID,required"`
	Port            string `env:"PORT,required"`

	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

	JwtSecret    string `env:"JWT_SECRET,required"`
	BackupDir    string `env:"BACKUP_DIR,required"`
	BackupSecret string `env:"BACKUP_SECRET,required"`

	UserName string `env:"USERNAME,required"`
	Password string `env:"PASSWORD,required"`
	Host     string `env:"HOST,required"`
}

func Load() error {
	godotenv.Load()
	return env.Parse(&App)
}
