package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type AMQPConfig struct {
	URI string `env:"AMQP_URI"`
}

func NewAMQPConfig() *AMQPConfig {
	cfg := AMQPConfig{}

	if err := godotenv.Load(".env", ".env.local"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("%+v\n", err)
		}
	}

	return &cfg
}
