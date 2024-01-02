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

	DbConfig struct {
		Host string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
		Port string `yaml:"port" env:"DB_POST" env-default:"21017"`
		Name string `yaml:"name" env:"DB_NAME" env-default:"dvault"`
		User string `yaml:"user" env:"DB_USER" env-default:"root"`
		Pass string `yaml:"pass" env:"DB_PASS" env-default:"example"`
	}

	Config struct {
		Logger LoggerConfig `yaml:"logger"`
		Server ServerConfig `yaml:"server"`
		Db     DbConfig     `yaml:"db"`
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
