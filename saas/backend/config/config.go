package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Port        int
	Environment string
	LogLevel    string
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	Expiry     time.Duration
	RefreshExp time.Duration
}

type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnvAsInt("PORT", 8080),
			Environment: getEnv("ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://ts2go:ts2go_dev_password@localhost:5432/ts2go_saas?sslmode=disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://:ts2go_redis_password@localhost:6379/0"),
			Password: getEnv("REDIS_PASSWORD", "ts2go_redis_password"),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			Expiry:     getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),
			RefreshExp: getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},
		Storage: StorageConfig{
			Endpoint:  getEnv("STORAGE_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("STORAGE_ACCESS_KEY", "ts2go"),
			SecretKey: getEnv("STORAGE_SECRET_KEY", "ts2go_minio_password"),
			Bucket:    getEnv("STORAGE_BUCKET", "ts2go-files"),
			UseSSL:    getEnvAsBool("STORAGE_USE_SSL", false),
		},
	}

	// Validate required fields
	if cfg.JWT.Secret == "your-super-secret-jwt-key-change-in-production" && cfg.Server.Environment == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return cfg, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
