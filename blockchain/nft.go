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

func NewNFTChecker(rpcURL, contractAddress string) (*NFTChecker, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	return &NFTChecker{
		client:       client,
		contractAddr: common.HexToAddress(contractAddress),
	}, nil
}

func (n *NFTChecker) HasNFT(address string, tokenID *big.Int) (bool, error) {
	// ERC1155 balanceOf function signature
	balanceOfSig := "0x00fdd58e" // balanceOf(address,uint256)

	addr := common.HexToAddress(address)

	// Pack the address and token ID
	paddedAddr := common.LeftPadBytes(addr.Bytes(), 32)
	paddedTokenID := common.LeftPadBytes(tokenID.Bytes(), 32)

	data := append(common.FromHex(balanceOfSig), append(paddedAddr, paddedTokenID...)...)

	result, err := n.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &n.contractAddr,
		Data: data,
	}, nil)

	if err != nil {
		return false, fmt.Errorf("contract call failed: %v", err)
	}

	balance := new(big.Int).SetBytes(result)
	return balance.Cmp(big.NewInt(0)) > 0, nil
}

func (n *NFTChecker) GetContractAddress() common.Address {
	return n.contractAddr
}
