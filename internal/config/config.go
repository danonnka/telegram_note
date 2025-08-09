package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	StoragePath string `yaml:"StoragePath"`
	Token       string `env:"TOKEN"`
}

func MustLoad() *Config {
	configPath := "config/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("файл конфига не найден: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("ошибка чтения конфига: %s", err)
	}
	return &cfg
}
