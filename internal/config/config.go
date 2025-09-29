package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string
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
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}

	return &cfg
}
