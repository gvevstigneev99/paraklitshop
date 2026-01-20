package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"APP_ENV" env-default:"local"`
	ServerPort string `yaml:"server_port" env:"SERVER_PORT" env-default:"8080"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
