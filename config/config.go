package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DbConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Driver   string
}

type Config struct {
	DbConf
	ApiConfig
}

type ApiConfig struct {
	ApiPort string
}

func (c *Config) InitialConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load env %v", err)
	}

	// Database config
	c.DbConf = DbConf{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = ApiConfig{
		ApiPort: os.Getenv("API_PORT"),
	}

	if c.Host == "" || c.Port == "" || c.User == "" || c.Password == "" || c.Database == "" || c.Driver == "" {
		return fmt.Errorf("missing required env variables")
	}
	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.InitialConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}
