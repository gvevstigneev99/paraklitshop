package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
	Auth     AuthConfig
}

type AppConfig struct {
	Name string `yaml:"name" env:"APP_NAME" env-default:"shop"`
	Env  string `yaml:"env" env:"APP_ENV" env-default:"local"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"5s"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
	DBName   string `yaml:"dbname" env:"POSTGRES_DB" env-default:"shop_db"`
	SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`
}

type RedisConfig struct {
	Host string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
	DB   int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret" env:"JWT_SECRET"`
	TTL    time.Duration `yaml:"ttl" env:"JWT_TTL" env-default:"24h"`
}

type LogConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"debug"`
}

type AuthConfig struct {
	BuyerPassword  string `yaml:"buyer_password" env:"AUTH_BUYER_PASSWORD"`
	SellerPassword string `yaml:"seller_password" env:"AUTH_SELLER_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = filepath.FromSlash("config/local.yaml")
	}

	// 1) читаем файл (если он есть) 2) поверх него применяем переменные окружения
	if _, err := os.Stat(configPath); err == nil {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			return nil, err
		}
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
