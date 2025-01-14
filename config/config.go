package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath          string
	Port            string
	RpcURL          string
	NFTContractAddr string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		DBPath:          getEnv("DB_PATH", "./data.db"),
		Port:            getEnv("PORT", ":8080"),
		RpcURL:          getEnv("RPC_URL", ""),
		NFTContractAddr: getEnv("NFT_CONTRACT_ADDRESS", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
