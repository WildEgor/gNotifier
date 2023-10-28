package configs

import "github.com/joho/godotenv"

type Configurator struct{}

func NewConfigurator() *Configurator {
	return &Configurator{}
}

func (c *Configurator) Load() error {
	err := godotenv.Load(".env", ".env.local")
	return err
}
