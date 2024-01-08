package configs

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Configurator struct {
	envs []string
}

func NewConfigurator() *Configurator {
	var envs = []string{".env", ".env.local"}

	conf := &Configurator{
		envs: envs,
	}

	conf.Load()

	return conf
}

func (c *Configurator) Load() {
	err := godotenv.Load(c.envs...)
	if err != nil {
		log.Fatal("Error loading envs file")
	}
}
