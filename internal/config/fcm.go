package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type FCMConfig struct {
	APIKey string `env:"FCM_ANDROID_API_KEY"`
}

func NewFCMConfig() *FCMConfig {
	cfg := FCMConfig{}

	if err := godotenv.Load(".env", ".env.local"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("[FCMConfig] %+v\n", err)
		}

		if cfg.APIKey == "" {
			log.Fatal("[FCMConfig] Failed load Android API key!")
		}
	}

	return &cfg
}
