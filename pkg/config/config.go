package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost               string
	DBPort               string
	DBName               string
	DBUser               string
	DBPassword           string
	ServicePort          string
	HQEndPoint           string
	HQBasicAuthorization string
}

func Load() (*Config, error) {
	cfg := &Config{
		DBHost:               os.Getenv("DB_HOST"),
		DBPort:               os.Getenv("DB_PORT"),
		DBName:               os.Getenv("DB_NAME"),
		DBUser:               os.Getenv("DB_USER"),
		DBPassword:           os.Getenv("DB_PASSWORD"),
		ServicePort:          os.Getenv("SERVICE_PORT"),
		HQEndPoint:           os.Getenv("HQ_END_POINT"),
		HQBasicAuthorization: os.Getenv("HQ_BASIC_AUTHORIZATION"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBPort == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.ServicePort == "" {
		return fmt.Errorf("SERVICE_PORT is required")
	}
	if c.HQEndPoint == "" {
		return fmt.Errorf("HQ_END_POINT is required")
	}
	if c.HQBasicAuthorization == "" {
		return fmt.Errorf("HQ_BASIC_AUTHORIZATION is required")
	}
	return nil
}
