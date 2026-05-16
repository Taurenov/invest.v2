package config

import "os"

type Config struct {
	APIAddr       string
	APIToken      string
	JWTSecret     string
	DatabaseURL   string
	RedisURL      string
	EngineHTTP    string
	AllowCORS     bool
}

func Load() Config {
	return Config{
		APIAddr:     env("API_ADDR", ":8080"),
		APIToken:    env("API_TOKEN", "dev-token"),
		JWTSecret:   env("JWT_SECRET", "dev-jwt-secret-change-in-production"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		EngineHTTP:  env("ENGINE_HTTP_URL", "http://127.0.0.1:50052"),
		AllowCORS:   os.Getenv("CORS") != "0",
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
