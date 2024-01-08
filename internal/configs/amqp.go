package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type AMQPConfig struct {
	URI string `env:"AMQP_URI"`
}

func NewAMQPConfig(
	c *Configurator,
) *AMQPConfig {
	cfg := AMQPConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}
