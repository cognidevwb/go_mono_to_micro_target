// Package config reads the customers-service configuration from the
// environment. Every variable is listed in the workspace .env.example.
package config

import (
	"os"
	"strconv"
)

// Config is the whole service configuration.
type Config struct {
	Service     string
	Port        string
	DatabaseURL string
	BrokerURL   string
	MaxConns    int
	// Peers maps a context this service calls to its base URL.
	Peers map[string]string
}

// Load reads the environment, falling back to local-development defaults.
func Load() Config {
	return Config{
		Service:     env("OTEL_SERVICE_NAME", "customers-service"),
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://customers_app:customers_app@localhost:5432/customersdb?sslmode=disable"),
		BrokerURL:   env("BROKER_URL", "nats://localhost:4222"),
		MaxConns:    envInt("DB_MAX_CONNS", 9),
		Peers: map[string]string{
		},
	}
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if n, err := strconv.Atoi(os.Getenv(key)); err == nil && n > 0 {
		return n
	}
	return def
}
