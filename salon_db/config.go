package main

import (
	"os"
)

type Config struct {
	DBURL  string
	PORT   string
	APIKey string
}

func LoadConfig() Config {
	return Config{
		DBURL:  os.Getenv("DB_POSTGRES_URL"),
		PORT:   os.Getenv("SERVER_PORT"),
		APIKey: os.Getenv("API_KEY"),
	}
}
