package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type SMSConfig struct {
	BaseURL  string `env:"SMS_BASE_URL"`
	Username string `env:"SMS_USERNAME"`
	Password string `env:"SMS_PASSWORD"`
}

func NewSMSConfig() *SMSConfig {
	cfg := SMSConfig{}

	if err := godotenv.Load(".env", ".env.local"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("[SMSConfig] %+v\n", err)
		}
	}

	return &cfg
}
