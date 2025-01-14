package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NFTChecker struct {
	client       *ethclient.Client
	contractAddr common.Address
}

func NewNFTChecker(rpcURL string, contractAddress string) (*NFTChecker, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	return &NFTChecker{
		client:       client,
		contractAddr: common.HexToAddress(contractAddress),
	}, nil
}

func (n *NFTChecker) HasNFT(address string) (bool, error) {
	// Basic ERC721 balanceOf function signature
	balanceOfSig := "balanceOf(address)"

	addr := common.HexToAddress(address)

	data, err := n.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &n.contractAddr,
		Data: common.FromHex(balanceOfSig + fmt.Sprintf("%064x", addr.Big())),
	}, nil)

	if err != nil {
		return false, fmt.Errorf("contract call failed: %v", err)
	}

	balance := new(big.Int).SetBytes(data)
	return balance.Cmp(big.NewInt(0)) > 0, nil
}
