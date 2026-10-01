package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpireHours int
}

func LoadConfig() (*Config, error) {
	// Load .env nếu file tồn tại.
	// Khi deploy Docker/cloud, environment variables
	// có thể được cung cấp trực tiếp.
	_ = godotenv.Load()

	jwtExpireHours := 24

	if value := os.Getenv("JWT_EXPIRE_HOURS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRE_HOURS: %w", err)
		}

		jwtExpireHours = parsed
	}

	config := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME", "library_db"),

		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpireHours: jwtExpireHours,
	}

	return config, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
