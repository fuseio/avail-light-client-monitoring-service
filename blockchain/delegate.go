package blockchain

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed abi/delegation_registry.json
var delegationRegistryABI []byte

type DelegationRegistry struct {
	client    *ethclient.Client
	address   common.Address
	parsedABI abi.ABI
}

func NewDelegationRegistry(rpcURL, contractAddress string) (*DelegationRegistry, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(delegationRegistryABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse delegation registry ABI: %v", err)
	}

	return &DelegationRegistry{
		client:    client,
		address:   common.HexToAddress(contractAddress),
		parsedABI: parsedABI,
	}, nil
}

func (d *DelegationRegistry) CheckDelegateForToken(delegate, vault, contract common.Address, tokenID *big.Int) (bool, error) {
	data, err := d.parsedABI.Pack("checkDelegateForToken", delegate, vault, contract, tokenID)
	if err != nil {
		return false, fmt.Errorf("failed to pack data for checkDelegateForToken: %v", err)
	}

	result, err := d.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &d.address,
		Data: data,
	}, nil)
	if err != nil {
		return false, fmt.Errorf("contract call failed for checkDelegateForToken: %v", err)
	}

	var out bool
	if err := d.parsedABI.UnpackIntoInterface(&out, "checkDelegateForToken", result); err != nil {
		return false, fmt.Errorf("failed to unpack checkDelegateForToken result: %v", err)
	}

	return out, nil
}

func (d *DelegationRegistry) CheckDelegateForContract(delegate, vault, contract common.Address) (bool, error) {
	data, err := d.parsedABI.Pack("checkDelegateForContract", delegate, vault, contract)
	if err != nil {
		return false, fmt.Errorf("failed to pack data for checkDelegateForContract: %v", err)
	}

	result, err := d.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &d.address,
		Data: data,
	}, nil)
	if err != nil {
		return false, fmt.Errorf("contract call failed for checkDelegateForContract: %v", err)
	}

	var out bool
	if err := d.parsedABI.UnpackIntoInterface(&out, "checkDelegateForContract", result); err != nil {
		return false, fmt.Errorf("failed to unpack checkDelegateForContract result: %v", err)
	}

	return out, nil
}

func (d *DelegationRegistry) CheckDelegateForAll(delegate, vault common.Address) (bool, error) {
	data, err := d.parsedABI.Pack("checkDelegateForAll", delegate, vault)
	if err != nil {
		return false, fmt.Errorf("failed to pack data for checkDelegateForAll: %v", err)
	}

	result, err := d.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &d.address,
		Data: data,
	}, nil)
	if err != nil {
		return false, fmt.Errorf("contract call failed for checkDelegateForAll: %v", err)
	}

	var out bool
	if err := d.parsedABI.UnpackIntoInterface(&out, "checkDelegateForAll", result); err != nil {
		return false, fmt.Errorf("failed to unpack checkDelegateForAll result: %v", err)
	}

	return out, nil
}

func (d *DelegationRegistry) CheckDelegateForERC1155(delegate, vault, contract common.Address, tokenID *big.Int, rights [32]byte) (*big.Int, error) {
	data, err := d.parsedABI.Pack("checkDelegateForERC1155", delegate, vault, contract, tokenID, rights)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data for checkDelegateForERC1155: %v", err)
	}

	result, err := d.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &d.address,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed for checkDelegateForERC1155: %v", err)
	}

	var out *big.Int
	if err := d.parsedABI.UnpackIntoInterface(&out, "checkDelegateForERC1155", result); err != nil {
		return nil, fmt.Errorf("failed to unpack checkDelegateForERC1155 result: %v", err)
	}

	return out, nil
}
