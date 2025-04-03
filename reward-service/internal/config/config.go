package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	MongoDB struct {
		URI      string
		Database string
		Timeout  time.Duration
	}
	HTTP struct {
		Address string
		Timeout time.Duration
	}
	Reward struct {
		Hour   int
		Minute int
	}
}

func Load() (*Config, error) {
	config := &Config{}
	
	config.MongoDB.URI = getEnv("MONGO_URI", "mongodb://localhost:27017")
	config.MongoDB.Database = getEnv("MONGO_DB", "rewards")
	config.MongoDB.Timeout = time.Duration(getEnvInt("MONGO_TIMEOUT_SECONDS", 10)) * time.Second
	
	config.HTTP.Address = getEnv("HTTP_ADDR", ":8081")
	config.HTTP.Timeout = time.Duration(getEnvInt("HTTP_TIMEOUT_SECONDS", 30)) * time.Second
	
	config.Reward.Hour = getEnvInt("REWARD_HOUR", 15)
	config.Reward.Minute = getEnvInt("REWARD_MINUTE", 41)
	
	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}