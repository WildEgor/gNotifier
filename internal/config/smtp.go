package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type SMTPConfig struct {
	From     string `env:"SMTP_FROM_EMAIL"`
	Host     string `env:"SMTP_HOST"`
	Port     int16  `env:"SMTP_PORT"`
	Username string `env:"SMS_USERNAME"`
	Password string `env:"SMS_PASSWORD"`
}

func NewSMTPConfig() *SMTPConfig {
	cfg := SMTPConfig{}

	if err := godotenv.Load(".env", ".env.local"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("[SMTPConfig] %+v\n", err)
		}
	}

	return &cfg
}
