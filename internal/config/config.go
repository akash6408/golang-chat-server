package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type HTTPServer struct {
	Address string `yaml:"address" env:"HTTP_ADDRESS"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE"`
}

type JWTConfig struct {
	Secret          string `yaml:"secret" env:"JWT_SECRET"`
	Issuer          string `yaml:"issuer" env:"JWT_ISSUER"`
	AccessTTLMinute int    `yaml:"access_ttl_minutes" env:"JWT_ACCESS_TTL_MINUTES"`
}

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	DBConfig   `yaml:"db" env:"ENV" env-required:"true"`
	JWT        JWTConfig `yaml:"jwt"`
}

func MustLoad() *Config {
	// Load variables from .env if present (non-fatal if missing)
	_ = godotenv.Load()

	// Read configuration strictly from environment variables
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("can not read env: %s", err.Error())
	}

	if cfg.Env == "" {
		cfg.Env = "dev"
	}

	return &cfg
}
