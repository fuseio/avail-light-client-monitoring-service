package nft

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

func (n *NFTChecker) GetBatchBalance(address string, _ []*big.Int) ([]*big.Int, error) {
	fmt.Printf("Checking NFT balance for address: %s\n", address)

	// We always query tokens 0 to 9.
	numTokens := int64(10)
	tokenIDs := make([]*big.Int, numTokens)
	for i := int64(0); i < numTokens; i++ {
		tokenIDs[i] = big.NewInt(i)
	}

	// Use the ABI function "balanceOfBatch(address[],uint256[])".
	methodID := "0x4e1273f4"

	
	firstOffset := big.NewInt(64)
	firstParamSize := new(big.Int).Mul(big.NewInt(32), big.NewInt(1+numTokens))
	secondOffset := new(big.Int).Add(firstOffset, firstParamSize)

	fmt.Printf("First offset: %s (expected: 40 in hex), second offset: %s (expected: 1a0 in hex)\n",
		firstOffset.Text(16), secondOffset.Text(16))

	data := common.FromHex(methodID)
	data = append(data, common.LeftPadBytes(firstOffset.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(secondOffset.Bytes(), 32)...)

	// First parameter: the "accounts" array.
	data = append(data, common.LeftPadBytes(big.NewInt(numTokens).Bytes(), 32)...)
	addr := common.HexToAddress(address)
	for i := 0; i < int(numTokens); i++ {
		data = append(data, common.LeftPadBytes(addr.Bytes(), 32)...)
	}

	// Second parameter: the "tokenIDs" array.
	data = append(data, common.LeftPadBytes(big.NewInt(numTokens).Bytes(), 32)...)
	for _, id := range tokenIDs {
		data = append(data, common.LeftPadBytes(id.Bytes(), 32)...)
	}

	result, err := n.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &n.contractAddr,
		Data: data,
	}, nil)
	if err != nil {
		fmt.Printf("Contract call error: %v\n", err)
		return nil, fmt.Errorf("contract call failed: %v", err)
	}
	fmt.Printf("Raw result (%d bytes): %x\n", len(result), result)
	if len(result) < 64 {
		return nil, fmt.Errorf("result length is invalid: %d", len(result))
	}

	offset := new(big.Int).SetBytes(result[:32]).Int64()
	fmt.Printf("Decoded offset: %d\n", offset)
	if int64(len(result)) < offset+32 {
		return nil, fmt.Errorf("result too short for length field")
	}

	arrayLen := new(big.Int).SetBytes(result[offset : offset+32]).Int64()
	fmt.Printf("Decoded array length: %d\n", arrayLen)
	if arrayLen != numTokens {
		fmt.Printf("Warning: Expected %d balances, got %d\n", numTokens, arrayLen)
	}

	var totalBalance big.Int
	var balances []*big.Int
	for i := int64(0); i < arrayLen; i++ {
		startPos := offset + 32 + i*32
		if int64(len(result)) < startPos+32 {
			return nil, fmt.Errorf("result too short for index %d", i)
		}
		bal := new(big.Int).SetBytes(result[startPos : startPos+32])
		balances = append(balances, bal)
		totalBalance.Add(&totalBalance, bal)
	}
	fmt.Printf("Total combined balance: %s\n", totalBalance.String())
	return []*big.Int{&totalBalance}, nil
}
