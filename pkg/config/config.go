package config

import "os"

type Config struct {
	database string
	port     string
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable: " + name)
	}
	return v
}

func NewConfig() *Config {
	return &Config{
		database: getenv("DATABASE_URL"),
		port:     getenv("PORT"),
	}
}

func (c Config) Database() string {
	return c.database
}

func (c Config) Port() string {
	return c.port
}
