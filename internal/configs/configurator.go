package configs

import "github.com/joho/godotenv"

type Configurator struct {
	envs []string
}

func NewConfigurator() *Configurator {
	var envs = []string{".env", ".env.local"}

	return &Configurator{
		envs: envs,
	}
}

func (c *Configurator) Load() error {
	err := godotenv.Load(c.envs...)
	return err
}
