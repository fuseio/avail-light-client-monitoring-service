package config

import (
	"encoding/hex"
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	MongoURI             string
	MongoDB              string
	RpcURL               string
	NFTContractAddr      string
	DelegateContractAddr string
	Rights               []byte
	CheckNFTInterval     int
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // default port
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("MONGO_URI environment variable is required")
	}

	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		return nil, errors.New("MONGO_DB environment variable is required")
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

	rights := os.Getenv("RIGHTS")
	if rights == "" {
		return nil, errors.New("RIGHTS environment variable is required")
	}

	// Remove 0x prefix if present
	if len(rights) >= 2 && rights[:2] == "0x" {
		rights = rights[2:]
	}

	rightsBytes, err := hex.DecodeString(rights)
	if err != nil {
		return nil, errors.New("invalid RIGHTS format")
	}

	checkNFTInterval := os.Getenv("CHECK_NFT_INTERVAL")
	if checkNFTInterval == "" {
		return nil, errors.New("CHECK_NFT_INTERVAL environment variable is required")
	}

	checkNFTIntervalInt, err := strconv.Atoi(checkNFTInterval)
	if err != nil {
		return nil, errors.New("invalid CHECK_NFT_INTERVAL format")
	}

	return &Config{
		Port:                 port,
		MongoURI:             mongoURI,
		MongoDB:              mongoDB,
		RpcURL:               rpcURL,
		NFTContractAddr:      nftContractAddr,
		DelegateContractAddr: delegateContractAddr,
		Rights:               rightsBytes,
		CheckNFTInterval:     checkNFTIntervalInt,
	}, nil
}
