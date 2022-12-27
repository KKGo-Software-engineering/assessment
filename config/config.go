package config

import "os"

type Config struct {
	Database string
	Port     string
}

func NewConfig() *Config {
	return &Config{
		Database: os.Getenv("DATABASE_URL"),
		Port:     os.Getenv("PORT"),
	}
}
