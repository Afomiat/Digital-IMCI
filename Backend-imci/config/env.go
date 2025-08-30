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
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	
	// SMS Configuration (replacing email/SMTP)
	TwilioAccountSID       string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken        string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioFromNumber       string `mapstructure:"TWILIO_FROM_NUMBER"`
	UseMockSMS             bool   `mapstructure:"USE_MOCK_SMS"`
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