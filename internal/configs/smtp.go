package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type SMTPConfig struct {
	From     string `env:"SMTP_FROM_EMAIL"`
	Host     string `env:"SMTP_HOST"`
	Port     int16  `env:"SMTP_PORT"`
	Username string `env:"SMS_USERNAME"`
	Password string `env:"SMS_PASSWORD"`
}

func NewSMTPConfig(c *Configurator) *SMTPConfig {
	cfg := SMTPConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("[SMTPConfig] %+v\n", err)
	}

	return &cfg
}
