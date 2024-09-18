// config/config.go
package config

import (
	"github.com/spf13/viper"
	"log"
)

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type AuthConfig struct {
	Username string
	Password string
}

type Config struct {
	Auth AuthConfig
	JWT  JWTConfig
	Path string
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}
}
