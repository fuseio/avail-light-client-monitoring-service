package config

import (
	"errors"
	"os"
)

type Config struct {
	Port                 string
	DBPath               string
	RpcURL               string
	NFTContractAddr      string
	DelegateContractAddr string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // default port
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return nil, errors.New("DB_PATH environment variable is required")
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		return nil, errors.New("RPC_URL environment variable is required")
	}

	nftContractAddr := os.Getenv("NFT_CONTRACT_ADDRESS")
	if nftContractAddr == "" {
		return nil, errors.New("NFT_CONTRACT_ADDRESS environment variable is required")
	}

	delegateContractAddr := os.Getenv("DELEGATE_CONTRACT_ADDRESS")
	if delegateContractAddr == "" {
		return nil, errors.New("DELEGATE_CONTRACT_ADDRESS environment variable is required")
	}

	return &Config{
		Port:                 port,
		DBPath:               dbPath,
		RpcURL:               rpcURL,
		NFTContractAddr:      nftContractAddr,
		DelegateContractAddr: delegateContractAddr,
	}, nil
}
