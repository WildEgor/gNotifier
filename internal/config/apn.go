package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type APNConfig struct {
	KeyPath    string `env:"APN_KEY_PATH"`
	KeyBase64  string `env:"APN_KEY_BASE64"`
	KeyType    string `env:"APN_KEY_TYPE"`
	KeyID      string `env:"APN_KEY_ID"`
	TeamID     string `env:"APN_TEAM_ID"`
	Password   string `env:"APN_PASSWORD"`
	Production bool   `env:"APN_PRODUCTION"`
}

func NewAPNConfig() *FCMConfig {
	cfg := FCMConfig{}

	if err := godotenv.Load(".env", ".env.local"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("[APNConfig] %+v\n", err)
		}
	}

	return &cfg
}
