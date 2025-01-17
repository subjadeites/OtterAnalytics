package config

import (
	"OtterAnalytics/pkg/errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Host          string
	Port          int
	PgsqlHost     string
	PgsqlPort     int
	PgsqlDb       string
	PgsqlUser     string
	PgsqlPassword string
	Token         string
}

var (
	config *Config
	once   sync.Once
)

func initDefaultConfig() *Config {
	return &Config{
		Host:          "localhost",
		Port:          7000,
		PgsqlHost:     "localhost",
		PgsqlPort:     5432,
		PgsqlDb:       "otter_analytics",
		PgsqlUser:     "postgres",
		PgsqlPassword: "password",
		Token:         "",
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func LoadConfig() *Config {
	once.Do(func() {
		config = initDefaultConfig()

		err := godotenv.Load("./.env")
		errors.Must(err, "Error loading .env file")

		config.Host = getEnv("HOST", config.Host)
		config.Port = getEnvAsInt("PORT", config.Port)
		config.PgsqlHost = getEnv("PGSQL_HOST", config.PgsqlHost)
		config.PgsqlPort = getEnvAsInt("PGSQL_PORT", config.PgsqlPort)
		config.PgsqlDb = getEnv("PGSQL_DB", config.PgsqlDb)
		config.PgsqlUser = getEnv("PGSQL_USER", config.PgsqlUser)
		config.PgsqlPassword = getEnv("PGSQL_PASSWORD", config.PgsqlPassword)
		config.Token = getEnv("TOKEN", config.Token)
	})
	return config
}
