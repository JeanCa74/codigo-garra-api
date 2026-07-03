package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DBPath string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		Port:   getEnv("PORT", ":8080"),
		DBPath: getEnv("DB_PATH", "garra.db"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
