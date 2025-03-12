package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type AppConfig struct {
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`
}

// NewAppConfig returns a new instance of AppConfig
func NewAppConfig() *AppConfig {
	var c AppConfig
	err := env.Parse(&c)
	if err != nil {
		log.Fatalf("error loading app config: %v", err)
	}

	return &c
}
