package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	Server     ServerConf
	Storage    StorageConf
	GRPC       GRPCConf
	Migrations MigrationsConf
}

type LoggerConf struct {
	Level    string
	FileName string
}

type ServerConf struct {
	Host string
	Port string
}

type GRPCConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Dsn         string
	StorageType string
}

type MigrationsConf struct {
	AutoMigrate bool
	Dir         string
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}
