package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Env struct {
	LocalServerPort        string `mapstructure:"LOCAL_SERVER_PORT"`
	PostgresDSN            string `mapstructure:"POSTGRES_DSN"`
	JWTSecret              string `mapstructure:"JWT_SECRET"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenExpiryMinute  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_MINUTE"`
	RefreshTokenExpiryDay int    `mapstructure:"REFRESH_TOKEN_EXPIRY_DAY"`

	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`

	MetaWhatsAppAccessToken      string `mapstructure:"META_WHATSAPP_ACCESS_TOKEN"`
	MetaWhatsAppPhoneNumberID    string `mapstructure:"META_WHATSAPP_PHONE_NUMBER_ID"`
	MetaWhatsAppBusinessAccountID string `mapstructure:"META_WHATSAPP_BUSINESS_ACCOUNT_ID"`

	RedisURL string `mapstructure:"REDIS_URL"`

}

func NewEnv() *Env {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var env Env
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can't find the file .env: %v", err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalf("Environment can't be loaded: %v", err)
	}

	return &env
}