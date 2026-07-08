package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBPath      string
	DatabaseURL string // PostgreSQL DSN — si se define, tiene prioridad sobre DBPath
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		Port:        getEnv("PORT", ":8080"),
		DBPath:      getEnv("DB_PATH", "garra.db"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
