package configs

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Configs struct {
	Debug bool `env:"DEBUG" env-default:"false"`
	NATS  NATS
	HTTP  HTTP
}

type NATS struct {
	URL string `env:"NATS_URL"`
}

type HTTP struct {
	Host string `env:"HTTP_HOST"`
	Port int    `env:"HTTP_PORT" env-default:"8001"`
}

func MustConfigs() *Configs {
	cfg := Configs{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to read configs: %s", err.Error())
	}

	return &cfg
}
