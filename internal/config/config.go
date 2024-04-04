package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	yamlPath = "configs"
	yamlFile = "config"
	envFile  = ".env"
)

type Config struct {
	Redis      RedisConfig
	HTTPServer HTTPServerConfig
	Cache      CacheConfig
}

type CacheConfig struct {
	Expiration time.Duration
}

type RedisConfig struct {
	Password string
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   int    `mapstructure:"db_name"`
}

type HTTPServerConfig struct {
	Port           string        `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
}

func InitConfig() (*Config, error) {
	if err := read(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return fmt.Errorf("unmarshal redis config: %v", err)
	}

	if err := viper.UnmarshalKey("cache", &cfg.Cache); err != nil {
		return fmt.Errorf("unmarshal redis config: %v", err)
	}

	if err := viper.UnmarshalKey("http_server", &cfg.HTTPServer); err != nil {
		return fmt.Errorf("unmarshal http server config: %v", err)
	}

	return nil
}

func read() error {
	viper.AddConfigPath(yamlPath)
	viper.SetConfigName(yamlFile)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read yaml file: %v: %v", envFile, err)
	}

	return nil
}
