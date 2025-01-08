package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func NewConfig(folder, filename string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New("unable to read config file: " + viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.New("unable to unmarshal config file: " + viper.ConfigFileUsed())
	}

	return cfg, nil
}
