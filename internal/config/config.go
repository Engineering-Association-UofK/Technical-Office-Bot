package config

type Config struct {
	TelegramToken   string `env:"TELEGRAM_APITOKEN,required"`
	AdminTelegramID int64  `env:"TELEGRAM_ADMINID,required"`
}
