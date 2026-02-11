package config

type Config struct {
	TelegramToken   string `env:"TELEGRAM_APITOKEN,required"`
	AdminTelegramID int64  `env:"TELEGRAM_ADMINID,required"`
	Port            string `env:"PORT,required"`

	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
}
