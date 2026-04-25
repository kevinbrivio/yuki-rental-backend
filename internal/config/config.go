package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App     AppConfig
	DB      DBConfig
	Session SessionConfig
	CORS    CORSConfig
}

type AppConfig struct {
	Port int
	Host string
	Env  string
}

type DBConfig struct {
	Port     int
	Host     string
	Password string
	Name     string
	SSLMode  string
	User     string
}

type SessionConfig struct {
	Secret []byte
	MaxAge time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

// --- METHODS -------------
func (db DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

func (a AppConfig) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a AppConfig) IsProd() bool {
	return a.Env == "production"
}

func envOrDefault(key, fallback string) string {
	envStr := os.Getenv(key)
	if envStr != "" {
		return envStr
	}
	
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	envStr := os.Getenv(key)
	if envStr != "" {
		// convert str to int
		envInt, err := strconv.Atoi(envStr)
		if err != nil {
			return fallback
		}
		
		return envInt
	}
	
	return fallback
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Port: envIntOrDefault("PORT", 8080),
			Host: envOrDefault("HOST", "0.0.0.0"),
			Env: envOrDefault("ENV", "development"),
		},
		DB: DBConfig{
			Host: envOrDefault("DB_HOST", "localhost"),
			Port: envIntOrDefault("DB_PORT", 5432),
			User: envOrDefault("DB_USER", "yuki"),
			Password: envOrDefault("DB_PASSWORD", "yuki_secret"),
			Name: envOrDefault("DB_NAME", "yuki_rental"),
			SSLMode: envOrDefault("DB_SSLMODE", "disable"),
		},
		Session: SessionConfig{
			MaxAge: time.Duration(envIntOrDefault("SESSION_MAX_AGE", 86400)) * time.Second,
		},
	}
	
	// Read SESSION_SECRET
	session := os.Getenv("SESSION_SECRET")
	if session == "" {
		if cfg.App.IsProd() {
			return nil, errors.New("SESSION_SECRET is required in production")
		}
		session = "evwev-we9vwevbe2e9v2bv929vnsnvs992993n--12#4vn"
	} 
	
	cfg.Session.Secret = []byte(session)
	
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins != "" {
		 cfg.CORS.AllowedOrigins = strings.Split(origins, ",")
	} else {
		cfg.CORS.AllowedOrigins= []string{"http://localhost:3000", "http://localhost:5173"}
	}
	
	return cfg, nil
}
