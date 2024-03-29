package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	Port    string `env:"APP_PORT"`
	Mode    string `env:"APP_MODE"`
	GoEnv   string `env:"GO_ENV"`
	Version string `env:"VERSION"`
}

func NewAppConfig(
	c *Configurator,
) *AppConfig {
	cfg := AppConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	if cfg.GoEnv == "" {
		cfg.GoEnv = "local"
	}

	if cfg.Version == "" {
		cfg.Version = "local"
	}

	return &cfg
}

func (ac AppConfig) IsProduction() bool {
	if ac.Mode == "develop" {
		return false
	}

	return true
}
