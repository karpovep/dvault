package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type (
	ServerConfig struct {
		Port string `yaml:"port" env:"SERVER_PORT" env-default:":9090"`
	}

	LoggerConfig struct {
		Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"debug"`
	}

	Config struct {
		Logger LoggerConfig `yaml:"logger"`
		Server ServerConfig
	}
)

func Init(configFile string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configFile, &cfg)
	if err != nil {
		log.Fatal("Error reading config", err.Error())
	}
	return &cfg
}
