package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	HTTPPort      string
	AllowedOrigin string
	SeedData      bool
	SeedAdmin     bool
	JWT           JWTConfig
	AdminSeed     AdminSeedConfig
	Database      DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type AdminSeedConfig struct {
	Name     string
	Email    string
	Password string
}

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		log.Printf("warning: %v", err)
	}

	cfg := Config{
		AppEnv:        getRequiredEnv("APP_ENV", "development"),
		HTTPPort:      getRequiredEnv("HTTP_PORT", "8080"),
		AllowedOrigin: getRequiredEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
		SeedData:      getEnvBool("SEED_DATA", false),
		SeedAdmin:     getEnvBool("SEED_ADMIN", false),
		JWT: JWTConfig{
			Secret:          getRequiredEnv("JWT_SECRET", "change-me-development-secret"),
			ExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		},
		AdminSeed: AdminSeedConfig{
			Name:     getRequiredEnv("ADMIN_NAME", "Administrator"),
			Email:    getRequiredEnv("ADMIN_EMAIL", "admin@example.com"),
			Password: os.Getenv("ADMIN_PASSWORD"),
		},
		Database: DatabaseConfig{
			Host:     getRequiredEnv("DB_HOST", "localhost"),
			Port:     getRequiredEnv("DB_PORT", "5434"),
			User:     getRequiredEnv("DB_USER", "macends"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     getRequiredEnv("DB_NAME", "live_coding_db"),
			SSLMode:  getRequiredEnv("DB_SSLMODE", "disable"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	requiredValues := map[string]string{
		"HTTP_PORT":  c.HTTPPort,
		"DB_HOST":    c.Database.Host,
		"DB_PORT":    c.Database.Port,
		"DB_USER":    c.Database.User,
		"DB_NAME":    c.Database.Name,
		"DB_SSLMODE": c.Database.SSLMode,
		"JWT_SECRET": c.JWT.Secret,
	}

	for key, value := range requiredValues {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("missing required config: %s", key)
		}
	}

	if c.JWT.ExpirationHours <= 0 {
		return fmt.Errorf("invalid required config: JWT_EXPIRATION_HOURS must be greater than 0")
	}

	if c.SeedAdmin && strings.TrimSpace(c.AdminSeed.Password) == "" {
		return fmt.Errorf("missing required config: ADMIN_PASSWORD")
	}

	return nil
}

func loadDotEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cwd, cwdErr := os.Getwd()
			if cwdErr != nil {
				return fmt.Errorf(".env file not found and current directory could not be read: %w", cwdErr)
			}

			return fmt.Errorf(".env file not found in %s; using OS environment variables and defaults", cwd)
		}

		return fmt.Errorf("failed to load .env file: %w", err)
	}

	return nil
}

func getRequiredEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	return strings.TrimSpace(value)
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsedValue
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsedValue
}
