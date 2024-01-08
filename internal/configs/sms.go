package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type SMSConfig struct {
	BaseURL  string `env:"SMS_BASE_URL"`
	Username string `env:"SMS_USERNAME"`
	Password string `env:"SMS_PASSWORD"`
}

func NewSMSConfig(c *Configurator) *SMSConfig {
	cfg := SMSConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("[SMSConfig] %+v\n", err)
	}

	return &cfg
}
