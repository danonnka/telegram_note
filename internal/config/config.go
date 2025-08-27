package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	StoragePath string `yaml:"StoragePath"`
	Token       string `env:"TOKEN"`
}

func MustLoad() (*Config, error) {
	configPath := "config/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл конфига не найден: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка чтения конфига: %s", err)
	}
	return &cfg, nil
}
