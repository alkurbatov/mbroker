package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config конфигурация сервиса.
type Config struct {
	// Queues конфигурация очередей сообщений.
	Queues map[string]Queue `mapstructure:"queues"`
}

func New() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read configuration file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal configuration: %w", err)
	}

	return &cfg, nil
}
