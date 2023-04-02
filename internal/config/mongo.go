package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type MongoConfig struct {
	Host     string `env:"MONGO_HOST"`
	DB       string `env:"MONGO_DB_NAME"`
	Port     int16  `env:"MONGO_PORT"`
	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PASSWORD"`
}

func NewMongoConfig() *MongoConfig {
	cfg := MongoConfig{}

	if err := godotenv.Load(".env", ".env.local"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("%+v\n", err)
		}
	}

	return &cfg
}

func (c *MongoConfig) GetHost() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}
