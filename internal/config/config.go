package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	StoragePath string `yaml:"StoragePath"`
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

//вопрос log который тут использут мой из main или другой = нет это другой, видно это через import.
// 2 варианта или в main обработать ошибку (этот) или сэда передать log
