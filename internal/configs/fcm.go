package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type FCMConfig struct {
	APIKey     string `env:"FCM_ANDROID_API_KEY"`
	Production bool   `env:"FCM_PRODUCTION"`
	MaxRetry   int    `env:"FCM_MAX_RETRY"`
}

func NewFCMConfig(c *Configurator) *FCMConfig {
	cfg := FCMConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("[FCMConfig] %+v\n", err)
	}

	if cfg.APIKey == "" {
		log.Fatal("[FCMConfig] Failed load Android API key!")
	}

	if cfg.MaxRetry == 0 {
		cfg.MaxRetry = 5
	}

	return &cfg
}
