package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret  string
	ServerAddr string
	EmailSender string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "pass"),
		DBName:     getEnv("DB_NAME", "authdb"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecret"),
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		EmailSender: getEnv("EMAIL_SENDER", "no-reply@example.com"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Error parsing int from env %s: %v", key, err)
		return fallback
	}
	return i
}
