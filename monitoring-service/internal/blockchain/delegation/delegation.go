// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegation

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IDelegateRegistryDelegation is an auto generated low-level Go binding around an user-defined struct.
type IDelegateRegistryDelegation struct {
	Type     uint8
	From     common.Address  
	To       common.Address 
	Rights   [32]byte
	Contract common.Address
	TokenId  *big.Int
	Amount   *big.Int
}

// DelegationMetaData contains all meta data concerning the Delegation contract.
var DelegationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"MulticallFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"DelegateAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"DelegateContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DelegateERC1155\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DelegateERC20\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"DelegateERC721\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"}],\"name\":\"checkDelegateForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"}],\"name\":\"checkDelegateForContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"}],\"name\":\"checkDelegateForERC1155\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"}],\"name\":\"checkDelegateForERC20\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"}],\"name\":\"checkDelegateForERC721\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"delegateAll\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"delegateContract\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegateERC1155\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegateERC20\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"delegateERC721\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"hashes\",\"type\":\"bytes32[]\"}],\"name\":\"getDelegationsFromHashes\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIDelegateRegistry.DelegationType\",\"name\":\"type_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structIDelegateRegistry.Delegation[]\",\"name\":\"delegations_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"getIncomingDelegationHashes\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"delegationHashes\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"getIncomingDelegations\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIDelegateRegistry.DelegationType\",\"name\":\"type_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structIDelegateRegistry.Delegation[]\",\"name\":\"delegations_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"getOutgoingDelegationHashes\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"delegationHashes\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"getOutgoingDelegations\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIDelegateRegistry.DelegationType\",\"name\":\"type_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rights\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structIDelegateRegistry.Delegation[]\",\"name\":\"delegations_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"location\",\"type\":\"bytes32\"}],\"name\":\"readSlot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"contents\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"locations\",\"type\":\"bytes32[]\"}],\"name\":\"readSlots\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"contents\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sweep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DelegationABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegationMetaData.ABI instead.
var DelegationABI = DelegationMetaData.ABI

// Delegation is an auto generated Go binding around an Ethereum contract.
type Delegation struct {
	DelegationCaller     // Read-only binding to the contract
	DelegationTransactor // Write-only binding to the contract
	DelegationFilterer   // Log filterer for contract events
}

// DelegationCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegationSession struct {
	Contract     *Delegation       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DelegationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegationCallerSession struct {
	Contract *DelegationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// DelegationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegationTransactorSession struct {
	Contract     *DelegationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// DelegationRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegationRaw struct {
	Contract *Delegation // Generic contract binding to access the raw methods on
}

// DelegationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegationCallerRaw struct {
	Contract *DelegationCaller // Generic read-only contract binding to access the raw methods on
}

// DelegationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegationTransactorRaw struct {
	Contract *DelegationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegation creates a new instance of Delegation, bound to a specific deployed contract.
func NewDelegation(address common.Address, backend bind.ContractBackend) (*Delegation, error) {
	contract, err := bindDelegation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Delegation{DelegationCaller: DelegationCaller{contract: contract}, DelegationTransactor: DelegationTransactor{contract: contract}, DelegationFilterer: DelegationFilterer{contract: contract}}, nil
}

// NewDelegationCaller creates a new read-only instance of Delegation, bound to a specific deployed contract.
func NewDelegationCaller(address common.Address, caller bind.ContractCaller) (*DelegationCaller, error) {
	contract, err := bindDelegation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationCaller{contract: contract}, nil
}

// NewDelegationTransactor creates a new write-only instance of Delegation, bound to a specific deployed contract.
func NewDelegationTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegationTransactor, error) {
	contract, err := bindDelegation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationTransactor{contract: contract}, nil
}

// NewDelegationFilterer creates a new log filterer instance of Delegation, bound to a specific deployed contract.
func NewDelegationFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegationFilterer, error) {
	contract, err := bindDelegation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegationFilterer{contract: contract}, nil
}

// bindDelegation binds a generic wrapper to an already deployed contract.
func bindDelegation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegation *DelegationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegation.Contract.DelegationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegation *DelegationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.Contract.DelegationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegation *DelegationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegation.Contract.DelegationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegation *DelegationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegation *DelegationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegation *DelegationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegation.Contract.contract.Transact(opts, method, params...)
}

// CheckDelegateForAll is a free data retrieval call binding the contract method 0xe839bd53.
//
// Solidity: function checkDelegateForAll(address to, address from, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationCaller) CheckDelegateForAll(opts *bind.CallOpts, to common.Address, from common.Address, rights [32]byte) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "checkDelegateForAll", to, from, rights)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDelegateForAll is a free data retrieval call binding the contract method 0xe839bd53.
//
// Solidity: function checkDelegateForAll(address to, address from, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationSession) CheckDelegateForAll(to common.Address, from common.Address, rights [32]byte) (bool, error) {
	return _Delegation.Contract.CheckDelegateForAll(&_Delegation.CallOpts, to, from, rights)
}

// CheckDelegateForAll is a free data retrieval call binding the contract method 0xe839bd53.
//
// Solidity: function checkDelegateForAll(address to, address from, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationCallerSession) CheckDelegateForAll(to common.Address, from common.Address, rights [32]byte) (bool, error) {
	return _Delegation.Contract.CheckDelegateForAll(&_Delegation.CallOpts, to, from, rights)
}

// CheckDelegateForContract is a free data retrieval call binding the contract method 0x8988eea9.
//
// Solidity: function checkDelegateForContract(address to, address from, address contract_, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationCaller) CheckDelegateForContract(opts *bind.CallOpts, to common.Address, from common.Address, contract_ common.Address, rights [32]byte) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "checkDelegateForContract", to, from, contract_, rights)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDelegateForContract is a free data retrieval call binding the contract method 0x8988eea9.
//
// Solidity: function checkDelegateForContract(address to, address from, address contract_, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationSession) CheckDelegateForContract(to common.Address, from common.Address, contract_ common.Address, rights [32]byte) (bool, error) {
	return _Delegation.Contract.CheckDelegateForContract(&_Delegation.CallOpts, to, from, contract_, rights)
}

// CheckDelegateForContract is a free data retrieval call binding the contract method 0x8988eea9.
//
// Solidity: function checkDelegateForContract(address to, address from, address contract_, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationCallerSession) CheckDelegateForContract(to common.Address, from common.Address, contract_ common.Address, rights [32]byte) (bool, error) {
	return _Delegation.Contract.CheckDelegateForContract(&_Delegation.CallOpts, to, from, contract_, rights)
}

// CheckDelegateForERC1155 is a free data retrieval call binding the contract method 0xb8705875.
//
// Solidity: function checkDelegateForERC1155(address to, address from, address contract_, uint256 tokenId, bytes32 rights) view returns(uint256 amount)
func (_Delegation *DelegationCaller) CheckDelegateForERC1155(opts *bind.CallOpts, to common.Address, from common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "checkDelegateForERC1155", to, from, contract_, tokenId, rights)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CheckDelegateForERC1155 is a free data retrieval call binding the contract method 0xb8705875.
//
// Solidity: function checkDelegateForERC1155(address to, address from, address contract_, uint256 tokenId, bytes32 rights) view returns(uint256 amount)
func (_Delegation *DelegationSession) CheckDelegateForERC1155(to common.Address, from common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte) (*big.Int, error) {
	return _Delegation.Contract.CheckDelegateForERC1155(&_Delegation.CallOpts, to, from, contract_, tokenId, rights)
}

// CheckDelegateForERC1155 is a free data retrieval call binding the contract method 0xb8705875.
//
// Solidity: function checkDelegateForERC1155(address to, address from, address contract_, uint256 tokenId, bytes32 rights) view returns(uint256 amount)
func (_Delegation *DelegationCallerSession) CheckDelegateForERC1155(to common.Address, from common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte) (*big.Int, error) {
	return _Delegation.Contract.CheckDelegateForERC1155(&_Delegation.CallOpts, to, from, contract_, tokenId, rights)
}

// CheckDelegateForERC20 is a free data retrieval call binding the contract method 0xba63c817.
//
// Solidity: function checkDelegateForERC20(address to, address from, address contract_, bytes32 rights) view returns(uint256 amount)
func (_Delegation *DelegationCaller) CheckDelegateForERC20(opts *bind.CallOpts, to common.Address, from common.Address, contract_ common.Address, rights [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "checkDelegateForERC20", to, from, contract_, rights)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CheckDelegateForERC20 is a free data retrieval call binding the contract method 0xba63c817.
//
// Solidity: function checkDelegateForERC20(address to, address from, address contract_, bytes32 rights) view returns(uint256 amount)
func (_Delegation *DelegationSession) CheckDelegateForERC20(to common.Address, from common.Address, contract_ common.Address, rights [32]byte) (*big.Int, error) {
	return _Delegation.Contract.CheckDelegateForERC20(&_Delegation.CallOpts, to, from, contract_, rights)
}

// CheckDelegateForERC20 is a free data retrieval call binding the contract method 0xba63c817.
//
// Solidity: function checkDelegateForERC20(address to, address from, address contract_, bytes32 rights) view returns(uint256 amount)
func (_Delegation *DelegationCallerSession) CheckDelegateForERC20(to common.Address, from common.Address, contract_ common.Address, rights [32]byte) (*big.Int, error) {
	return _Delegation.Contract.CheckDelegateForERC20(&_Delegation.CallOpts, to, from, contract_, rights)
}

// CheckDelegateForERC721 is a free data retrieval call binding the contract method 0xb9f36874.
//
// Solidity: function checkDelegateForERC721(address to, address from, address contract_, uint256 tokenId, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationCaller) CheckDelegateForERC721(opts *bind.CallOpts, to common.Address, from common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "checkDelegateForERC721", to, from, contract_, tokenId, rights)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDelegateForERC721 is a free data retrieval call binding the contract method 0xb9f36874.
//
// Solidity: function checkDelegateForERC721(address to, address from, address contract_, uint256 tokenId, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationSession) CheckDelegateForERC721(to common.Address, from common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte) (bool, error) {
	return _Delegation.Contract.CheckDelegateForERC721(&_Delegation.CallOpts, to, from, contract_, tokenId, rights)
}

// CheckDelegateForERC721 is a free data retrieval call binding the contract method 0xb9f36874.
//
// Solidity: function checkDelegateForERC721(address to, address from, address contract_, uint256 tokenId, bytes32 rights) view returns(bool valid)
func (_Delegation *DelegationCallerSession) CheckDelegateForERC721(to common.Address, from common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte) (bool, error) {
	return _Delegation.Contract.CheckDelegateForERC721(&_Delegation.CallOpts, to, from, contract_, tokenId, rights)
}

// GetDelegationsFromHashes is a free data retrieval call binding the contract method 0x4705ed38.
//
// Solidity: function getDelegationsFromHashes(bytes32[] hashes) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationCaller) GetDelegationsFromHashes(opts *bind.CallOpts, hashes [][32]byte) ([]IDelegateRegistryDelegation, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getDelegationsFromHashes", hashes)

	if err != nil {
		return *new([]IDelegateRegistryDelegation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDelegateRegistryDelegation)).(*[]IDelegateRegistryDelegation)

	return out0, err

}

// GetDelegationsFromHashes is a free data retrieval call binding the contract method 0x4705ed38.
//
// Solidity: function getDelegationsFromHashes(bytes32[] hashes) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationSession) GetDelegationsFromHashes(hashes [][32]byte) ([]IDelegateRegistryDelegation, error) {
	return _Delegation.Contract.GetDelegationsFromHashes(&_Delegation.CallOpts, hashes)
}

// GetDelegationsFromHashes is a free data retrieval call binding the contract method 0x4705ed38.
//
// Solidity: function getDelegationsFromHashes(bytes32[] hashes) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationCallerSession) GetDelegationsFromHashes(hashes [][32]byte) ([]IDelegateRegistryDelegation, error) {
	return _Delegation.Contract.GetDelegationsFromHashes(&_Delegation.CallOpts, hashes)
}

// GetIncomingDelegationHashes is a free data retrieval call binding the contract method 0x063182a5.
//
// Solidity: function getIncomingDelegationHashes(address to) view returns(bytes32[] delegationHashes)
func (_Delegation *DelegationCaller) GetIncomingDelegationHashes(opts *bind.CallOpts, to common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getIncomingDelegationHashes", to)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetIncomingDelegationHashes is a free data retrieval call binding the contract method 0x063182a5.
//
// Solidity: function getIncomingDelegationHashes(address to) view returns(bytes32[] delegationHashes)
func (_Delegation *DelegationSession) GetIncomingDelegationHashes(to common.Address) ([][32]byte, error) {
	return _Delegation.Contract.GetIncomingDelegationHashes(&_Delegation.CallOpts, to)
}

// GetIncomingDelegationHashes is a free data retrieval call binding the contract method 0x063182a5.
//
// Solidity: function getIncomingDelegationHashes(address to) view returns(bytes32[] delegationHashes)
func (_Delegation *DelegationCallerSession) GetIncomingDelegationHashes(to common.Address) ([][32]byte, error) {
	return _Delegation.Contract.GetIncomingDelegationHashes(&_Delegation.CallOpts, to)
}

// GetIncomingDelegations is a free data retrieval call binding the contract method 0x42f87c25.
//
// Solidity: function getIncomingDelegations(address to) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationCaller) GetIncomingDelegations(opts *bind.CallOpts, to common.Address) ([]IDelegateRegistryDelegation, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getIncomingDelegations", to)

	if err != nil {
		return *new([]IDelegateRegistryDelegation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDelegateRegistryDelegation)).(*[]IDelegateRegistryDelegation)

	return out0, err

}

// GetIncomingDelegations is a free data retrieval call binding the contract method 0x42f87c25.
//
// Solidity: function getIncomingDelegations(address to) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationSession) GetIncomingDelegations(to common.Address) ([]IDelegateRegistryDelegation, error) {
	return _Delegation.Contract.GetIncomingDelegations(&_Delegation.CallOpts, to)
}

// GetIncomingDelegations is a free data retrieval call binding the contract method 0x42f87c25.
//
// Solidity: function getIncomingDelegations(address to) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationCallerSession) GetIncomingDelegations(to common.Address) ([]IDelegateRegistryDelegation, error) {
	return _Delegation.Contract.GetIncomingDelegations(&_Delegation.CallOpts, to)
}

// GetOutgoingDelegationHashes is a free data retrieval call binding the contract method 0x01a920a0.
//
// Solidity: function getOutgoingDelegationHashes(address from) view returns(bytes32[] delegationHashes)
func (_Delegation *DelegationCaller) GetOutgoingDelegationHashes(opts *bind.CallOpts, from common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getOutgoingDelegationHashes", from)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetOutgoingDelegationHashes is a free data retrieval call binding the contract method 0x01a920a0.
//
// Solidity: function getOutgoingDelegationHashes(address from) view returns(bytes32[] delegationHashes)
func (_Delegation *DelegationSession) GetOutgoingDelegationHashes(from common.Address) ([][32]byte, error) {
	return _Delegation.Contract.GetOutgoingDelegationHashes(&_Delegation.CallOpts, from)
}

// GetOutgoingDelegationHashes is a free data retrieval call binding the contract method 0x01a920a0.
//
// Solidity: function getOutgoingDelegationHashes(address from) view returns(bytes32[] delegationHashes)
func (_Delegation *DelegationCallerSession) GetOutgoingDelegationHashes(from common.Address) ([][32]byte, error) {
	return _Delegation.Contract.GetOutgoingDelegationHashes(&_Delegation.CallOpts, from)
}

// GetOutgoingDelegations is a free data retrieval call binding the contract method 0x51525e9a.
//
// Solidity: function getOutgoingDelegations(address from) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationCaller) GetOutgoingDelegations(opts *bind.CallOpts, from common.Address) ([]IDelegateRegistryDelegation, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getOutgoingDelegations", from)

	if err != nil {
		return *new([]IDelegateRegistryDelegation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDelegateRegistryDelegation)).(*[]IDelegateRegistryDelegation)

	return out0, err

}

// GetOutgoingDelegations is a free data retrieval call binding the contract method 0x51525e9a.
//
// Solidity: function getOutgoingDelegations(address from) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationSession) GetOutgoingDelegations(from common.Address) ([]IDelegateRegistryDelegation, error) {
	return _Delegation.Contract.GetOutgoingDelegations(&_Delegation.CallOpts, from)
}

// GetOutgoingDelegations is a free data retrieval call binding the contract method 0x51525e9a.
//
// Solidity: function getOutgoingDelegations(address from) view returns((uint8,address,address,bytes32,address,uint256,uint256)[] delegations_)
func (_Delegation *DelegationCallerSession) GetOutgoingDelegations(from common.Address) ([]IDelegateRegistryDelegation, error) {
	return _Delegation.Contract.GetOutgoingDelegations(&_Delegation.CallOpts, from)
}

// ReadSlot is a free data retrieval call binding the contract method 0xe8e834a9.
//
// Solidity: function readSlot(bytes32 location) view returns(bytes32 contents)
func (_Delegation *DelegationCaller) ReadSlot(opts *bind.CallOpts, location [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "readSlot", location)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ReadSlot is a free data retrieval call binding the contract method 0xe8e834a9.
//
// Solidity: function readSlot(bytes32 location) view returns(bytes32 contents)
func (_Delegation *DelegationSession) ReadSlot(location [32]byte) ([32]byte, error) {
	return _Delegation.Contract.ReadSlot(&_Delegation.CallOpts, location)
}

// ReadSlot is a free data retrieval call binding the contract method 0xe8e834a9.
//
// Solidity: function readSlot(bytes32 location) view returns(bytes32 contents)
func (_Delegation *DelegationCallerSession) ReadSlot(location [32]byte) ([32]byte, error) {
	return _Delegation.Contract.ReadSlot(&_Delegation.CallOpts, location)
}

// ReadSlots is a free data retrieval call binding the contract method 0x61451a30.
//
// Solidity: function readSlots(bytes32[] locations) view returns(bytes32[] contents)
func (_Delegation *DelegationCaller) ReadSlots(opts *bind.CallOpts, locations [][32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "readSlots", locations)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// ReadSlots is a free data retrieval call binding the contract method 0x61451a30.
//
// Solidity: function readSlots(bytes32[] locations) view returns(bytes32[] contents)
func (_Delegation *DelegationSession) ReadSlots(locations [][32]byte) ([][32]byte, error) {
	return _Delegation.Contract.ReadSlots(&_Delegation.CallOpts, locations)
}

// ReadSlots is a free data retrieval call binding the contract method 0x61451a30.
//
// Solidity: function readSlots(bytes32[] locations) view returns(bytes32[] contents)
func (_Delegation *DelegationCallerSession) ReadSlots(locations [][32]byte) ([][32]byte, error) {
	return _Delegation.Contract.ReadSlots(&_Delegation.CallOpts, locations)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Delegation *DelegationCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Delegation *DelegationSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Delegation.Contract.SupportsInterface(&_Delegation.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Delegation *DelegationCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Delegation.Contract.SupportsInterface(&_Delegation.CallOpts, interfaceId)
}

// DelegateAll is a paid mutator transaction binding the contract method 0x30ff3140.
//
// Solidity: function delegateAll(address to, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactor) DelegateAll(opts *bind.TransactOpts, to common.Address, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateAll", to, rights, enable)
}

// DelegateAll is a paid mutator transaction binding the contract method 0x30ff3140.
//
// Solidity: function delegateAll(address to, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationSession) DelegateAll(to common.Address, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateAll(&_Delegation.TransactOpts, to, rights, enable)
}

// DelegateAll is a paid mutator transaction binding the contract method 0x30ff3140.
//
// Solidity: function delegateAll(address to, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactorSession) DelegateAll(to common.Address, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateAll(&_Delegation.TransactOpts, to, rights, enable)
}

// DelegateContract is a paid mutator transaction binding the contract method 0xd90e73ab.
//
// Solidity: function delegateContract(address to, address contract_, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactor) DelegateContract(opts *bind.TransactOpts, to common.Address, contract_ common.Address, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateContract", to, contract_, rights, enable)
}

// DelegateContract is a paid mutator transaction binding the contract method 0xd90e73ab.
//
// Solidity: function delegateContract(address to, address contract_, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationSession) DelegateContract(to common.Address, contract_ common.Address, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateContract(&_Delegation.TransactOpts, to, contract_, rights, enable)
}

// DelegateContract is a paid mutator transaction binding the contract method 0xd90e73ab.
//
// Solidity: function delegateContract(address to, address contract_, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactorSession) DelegateContract(to common.Address, contract_ common.Address, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateContract(&_Delegation.TransactOpts, to, contract_, rights, enable)
}

// DelegateERC1155 is a paid mutator transaction binding the contract method 0xab764683.
//
// Solidity: function delegateERC1155(address to, address contract_, uint256 tokenId, bytes32 rights, uint256 amount) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactor) DelegateERC1155(opts *bind.TransactOpts, to common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateERC1155", to, contract_, tokenId, rights, amount)
}

// DelegateERC1155 is a paid mutator transaction binding the contract method 0xab764683.
//
// Solidity: function delegateERC1155(address to, address contract_, uint256 tokenId, bytes32 rights, uint256 amount) payable returns(bytes32 hash)
func (_Delegation *DelegationSession) DelegateERC1155(to common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateERC1155(&_Delegation.TransactOpts, to, contract_, tokenId, rights, amount)
}

// DelegateERC1155 is a paid mutator transaction binding the contract method 0xab764683.
//
// Solidity: function delegateERC1155(address to, address contract_, uint256 tokenId, bytes32 rights, uint256 amount) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactorSession) DelegateERC1155(to common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateERC1155(&_Delegation.TransactOpts, to, contract_, tokenId, rights, amount)
}

// DelegateERC20 is a paid mutator transaction binding the contract method 0x003c2ba6.
//
// Solidity: function delegateERC20(address to, address contract_, bytes32 rights, uint256 amount) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactor) DelegateERC20(opts *bind.TransactOpts, to common.Address, contract_ common.Address, rights [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateERC20", to, contract_, rights, amount)
}

// DelegateERC20 is a paid mutator transaction binding the contract method 0x003c2ba6.
//
// Solidity: function delegateERC20(address to, address contract_, bytes32 rights, uint256 amount) payable returns(bytes32 hash)
func (_Delegation *DelegationSession) DelegateERC20(to common.Address, contract_ common.Address, rights [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateERC20(&_Delegation.TransactOpts, to, contract_, rights, amount)
}

// DelegateERC20 is a paid mutator transaction binding the contract method 0x003c2ba6.
//
// Solidity: function delegateERC20(address to, address contract_, bytes32 rights, uint256 amount) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactorSession) DelegateERC20(to common.Address, contract_ common.Address, rights [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateERC20(&_Delegation.TransactOpts, to, contract_, rights, amount)
}

// DelegateERC721 is a paid mutator transaction binding the contract method 0xb18e2bbb.
//
// Solidity: function delegateERC721(address to, address contract_, uint256 tokenId, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactor) DelegateERC721(opts *bind.TransactOpts, to common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateERC721", to, contract_, tokenId, rights, enable)
}

// DelegateERC721 is a paid mutator transaction binding the contract method 0xb18e2bbb.
//
// Solidity: function delegateERC721(address to, address contract_, uint256 tokenId, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationSession) DelegateERC721(to common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateERC721(&_Delegation.TransactOpts, to, contract_, tokenId, rights, enable)
}

// DelegateERC721 is a paid mutator transaction binding the contract method 0xb18e2bbb.
//
// Solidity: function delegateERC721(address to, address contract_, uint256 tokenId, bytes32 rights, bool enable) payable returns(bytes32 hash)
func (_Delegation *DelegationTransactorSession) DelegateERC721(to common.Address, contract_ common.Address, tokenId *big.Int, rights [32]byte, enable bool) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateERC721(&_Delegation.TransactOpts, to, contract_, tokenId, rights, enable)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_Delegation *DelegationTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_Delegation *DelegationSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Delegation.Contract.Multicall(&_Delegation.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_Delegation *DelegationTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Delegation.Contract.Multicall(&_Delegation.TransactOpts, data)
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_Delegation *DelegationTransactor) Sweep(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "sweep")
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_Delegation *DelegationSession) Sweep() (*types.Transaction, error) {
	return _Delegation.Contract.Sweep(&_Delegation.TransactOpts)
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_Delegation *DelegationTransactorSession) Sweep() (*types.Transaction, error) {
	return _Delegation.Contract.Sweep(&_Delegation.TransactOpts)
}

// DelegationDelegateAllIterator is returned from FilterDelegateAll and is used to iterate over the raw logs and unpacked data for DelegateAll events raised by the Delegation contract.
type DelegationDelegateAllIterator struct {
	Event *DelegationDelegateAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationDelegateAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationDelegateAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationDelegateAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationDelegateAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationDelegateAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationDelegateAll represents a DelegateAll event raised by the Delegation contract.
type DelegationDelegateAll struct {
	From   common.Address
	To     common.Address
	Rights [32]byte
	Enable bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDelegateAll is a free log retrieval operation binding the contract event 0xda3ef6410e30373a9137f83f9781a8129962b6882532b7c229de2e39de423227.
//
// Solidity: event DelegateAll(address indexed from, address indexed to, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) FilterDelegateAll(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DelegationDelegateAllIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "DelegateAll", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DelegationDelegateAllIterator{contract: _Delegation.contract, event: "DelegateAll", logs: logs, sub: sub}, nil
}

// WatchDelegateAll is a free log subscription operation binding the contract event 0xda3ef6410e30373a9137f83f9781a8129962b6882532b7c229de2e39de423227.
//
// Solidity: event DelegateAll(address indexed from, address indexed to, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) WatchDelegateAll(opts *bind.WatchOpts, sink chan<- *DelegationDelegateAll, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "DelegateAll", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationDelegateAll)
				if err := _Delegation.contract.UnpackLog(event, "DelegateAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDelegateAll is a log parse operation binding the contract event 0xda3ef6410e30373a9137f83f9781a8129962b6882532b7c229de2e39de423227.
//
// Solidity: event DelegateAll(address indexed from, address indexed to, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) ParseDelegateAll(log types.Log) (*DelegationDelegateAll, error) {
	event := new(DelegationDelegateAll)
	if err := _Delegation.contract.UnpackLog(event, "DelegateAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationDelegateContractIterator is returned from FilterDelegateContract and is used to iterate over the raw logs and unpacked data for DelegateContract events raised by the Delegation contract.
type DelegationDelegateContractIterator struct {
	Event *DelegationDelegateContract // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationDelegateContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationDelegateContract)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationDelegateContract)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationDelegateContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationDelegateContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationDelegateContract represents a DelegateContract event raised by the Delegation contract.
type DelegationDelegateContract struct {
	From     common.Address
	To       common.Address
	Contract common.Address
	Rights   [32]byte
	Enable   bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateContract is a free log retrieval operation binding the contract event 0x021be15e24de4afc43cfb5d0ba95ca38e0783571e05c12bbe6aece8842ae82df.
//
// Solidity: event DelegateContract(address indexed from, address indexed to, address indexed contract_, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) FilterDelegateContract(opts *bind.FilterOpts, from []common.Address, to []common.Address, contract_ []common.Address) (*DelegationDelegateContractIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var contract_Rule []interface{}
	for _, contract_Item := range contract_ {
		contract_Rule = append(contract_Rule, contract_Item)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "DelegateContract", fromRule, toRule, contract_Rule)
	if err != nil {
		return nil, err
	}
	return &DelegationDelegateContractIterator{contract: _Delegation.contract, event: "DelegateContract", logs: logs, sub: sub}, nil
}

// WatchDelegateContract is a free log subscription operation binding the contract event 0x021be15e24de4afc43cfb5d0ba95ca38e0783571e05c12bbe6aece8842ae82df.
//
// Solidity: event DelegateContract(address indexed from, address indexed to, address indexed contract_, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) WatchDelegateContract(opts *bind.WatchOpts, sink chan<- *DelegationDelegateContract, from []common.Address, to []common.Address, contract_ []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var contract_Rule []interface{}
	for _, contract_Item := range contract_ {
		contract_Rule = append(contract_Rule, contract_Item)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "DelegateContract", fromRule, toRule, contract_Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationDelegateContract)
				if err := _Delegation.contract.UnpackLog(event, "DelegateContract", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDelegateContract is a log parse operation binding the contract event 0x021be15e24de4afc43cfb5d0ba95ca38e0783571e05c12bbe6aece8842ae82df.
//
// Solidity: event DelegateContract(address indexed from, address indexed to, address indexed contract_, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) ParseDelegateContract(log types.Log) (*DelegationDelegateContract, error) {
	event := new(DelegationDelegateContract)
	if err := _Delegation.contract.UnpackLog(event, "DelegateContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationDelegateERC1155Iterator is returned from FilterDelegateERC1155 and is used to iterate over the raw logs and unpacked data for DelegateERC1155 events raised by the Delegation contract.
type DelegationDelegateERC1155Iterator struct {
	Event *DelegationDelegateERC1155 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationDelegateERC1155Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationDelegateERC1155)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationDelegateERC1155)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationDelegateERC1155Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationDelegateERC1155Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationDelegateERC1155 represents a DelegateERC1155 event raised by the Delegation contract.
type DelegationDelegateERC1155 struct {
	From     common.Address
	To       common.Address
	Contract common.Address
	TokenId  *big.Int
	Rights   [32]byte
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateERC1155 is a free log retrieval operation binding the contract event 0x27ab1adc9bca76301ed7a691320766dfa4b4b1aa32c9e05cf789611be7f8c75f.
//
// Solidity: event DelegateERC1155(address indexed from, address indexed to, address indexed contract_, uint256 tokenId, bytes32 rights, uint256 amount)
func (_Delegation *DelegationFilterer) FilterDelegateERC1155(opts *bind.FilterOpts, from []common.Address, to []common.Address, contract_ []common.Address) (*DelegationDelegateERC1155Iterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var contract_Rule []interface{}
	for _, contract_Item := range contract_ {
		contract_Rule = append(contract_Rule, contract_Item)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "DelegateERC1155", fromRule, toRule, contract_Rule)
	if err != nil {
		return nil, err
	}
	return &DelegationDelegateERC1155Iterator{contract: _Delegation.contract, event: "DelegateERC1155", logs: logs, sub: sub}, nil
}

// WatchDelegateERC1155 is a free log subscription operation binding the contract event 0x27ab1adc9bca76301ed7a691320766dfa4b4b1aa32c9e05cf789611be7f8c75f.
//
// Solidity: event DelegateERC1155(address indexed from, address indexed to, address indexed contract_, uint256 tokenId, bytes32 rights, uint256 amount)
func (_Delegation *DelegationFilterer) WatchDelegateERC1155(opts *bind.WatchOpts, sink chan<- *DelegationDelegateERC1155, from []common.Address, to []common.Address, contract_ []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var contract_Rule []interface{}
	for _, contract_Item := range contract_ {
		contract_Rule = append(contract_Rule, contract_Item)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "DelegateERC1155", fromRule, toRule, contract_Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationDelegateERC1155)
				if err := _Delegation.contract.UnpackLog(event, "DelegateERC1155", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDelegateERC1155 is a log parse operation binding the contract event 0x27ab1adc9bca76301ed7a691320766dfa4b4b1aa32c9e05cf789611be7f8c75f.
//
// Solidity: event DelegateERC1155(address indexed from, address indexed to, address indexed contract_, uint256 tokenId, bytes32 rights, uint256 amount)
func (_Delegation *DelegationFilterer) ParseDelegateERC1155(log types.Log) (*DelegationDelegateERC1155, error) {
	event := new(DelegationDelegateERC1155)
	if err := _Delegation.contract.UnpackLog(event, "DelegateERC1155", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationDelegateERC20Iterator is returned from FilterDelegateERC20 and is used to iterate over the raw logs and unpacked data for DelegateERC20 events raised by the Delegation contract.
type DelegationDelegateERC20Iterator struct {
	Event *DelegationDelegateERC20 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationDelegateERC20Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationDelegateERC20)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationDelegateERC20)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationDelegateERC20Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationDelegateERC20Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationDelegateERC20 represents a DelegateERC20 event raised by the Delegation contract.
type DelegationDelegateERC20 struct {
	From     common.Address
	To       common.Address
	Contract common.Address
	Rights   [32]byte
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateERC20 is a free log retrieval operation binding the contract event 0x6ebd000dfc4dc9df04f723f827bae7694230795e8f22ed4af438e074cc982d18.
//
// Solidity: event DelegateERC20(address indexed from, address indexed to, address indexed contract_, bytes32 rights, uint256 amount)
func (_Delegation *DelegationFilterer) FilterDelegateERC20(opts *bind.FilterOpts, from []common.Address, to []common.Address, contract_ []common.Address) (*DelegationDelegateERC20Iterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var contract_Rule []interface{}
	for _, contract_Item := range contract_ {
		contract_Rule = append(contract_Rule, contract_Item)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "DelegateERC20", fromRule, toRule, contract_Rule)
	if err != nil {
		return nil, err
	}
	return &DelegationDelegateERC20Iterator{contract: _Delegation.contract, event: "DelegateERC20", logs: logs, sub: sub}, nil
}

// WatchDelegateERC20 is a free log subscription operation binding the contract event 0x6ebd000dfc4dc9df04f723f827bae7694230795e8f22ed4af438e074cc982d18.
//
// Solidity: event DelegateERC20(address indexed from, address indexed to, address indexed contract_, bytes32 rights, uint256 amount)
func (_Delegation *DelegationFilterer) WatchDelegateERC20(opts *bind.WatchOpts, sink chan<- *DelegationDelegateERC20, from []common.Address, to []common.Address, contract_ []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var contract_Rule []interface{}
	for _, contract_Item := range contract_ {
		contract_Rule = append(contract_Rule, contract_Item)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "DelegateERC20", fromRule, toRule, contract_Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationDelegateERC20)
				if err := _Delegation.contract.UnpackLog(event, "DelegateERC20", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDelegateERC20 is a log parse operation binding the contract event 0x6ebd000dfc4dc9df04f723f827bae7694230795e8f22ed4af438e074cc982d18.
//
// Solidity: event DelegateERC20(address indexed from, address indexed to, address indexed contract_, bytes32 rights, uint256 amount)
func (_Delegation *DelegationFilterer) ParseDelegateERC20(log types.Log) (*DelegationDelegateERC20, error) {
	event := new(DelegationDelegateERC20)
	if err := _Delegation.contract.UnpackLog(event, "DelegateERC20", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationDelegateERC721Iterator is returned from FilterDelegateERC721 and is used to iterate over the raw logs and unpacked data for DelegateERC721 events raised by the Delegation contract.
type DelegationDelegateERC721Iterator struct {
	Event *DelegationDelegateERC721 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationDelegateERC721Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationDelegateERC721)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationDelegateERC721)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationDelegateERC721Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationDelegateERC721Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationDelegateERC721 represents a DelegateERC721 event raised by the Delegation contract.
type DelegationDelegateERC721 struct {
	From     common.Address
	To       common.Address
	Contract common.Address
	TokenId  *big.Int
	Rights   [32]byte
	Enable   bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateERC721 is a free log retrieval operation binding the contract event 0x15e7a1bdcd507dd632d797d38e60cc5a9c0749b9a63097a215c4d006126825c6.
//
// Solidity: event DelegateERC721(address indexed from, address indexed to, address indexed contract_, uint256 tokenId, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) FilterDelegateERC721(opts *bind.FilterOpts, from []common.Address, to []common.Address, contract_ []common.Address) (*DelegationDelegateERC721Iterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var contract_Rule []interface{}
	for _, contract_Item := range contract_ {
		contract_Rule = append(contract_Rule, contract_Item)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "DelegateERC721", fromRule, toRule, contract_Rule)
	if err != nil {
		return nil, err
	}
	return &DelegationDelegateERC721Iterator{contract: _Delegation.contract, event: "DelegateERC721", logs: logs, sub: sub}, nil
}

// WatchDelegateERC721 is a free log subscription operation binding the contract event 0x15e7a1bdcd507dd632d797d38e60cc5a9c0749b9a63097a215c4d006126825c6.
//
// Solidity: event DelegateERC721(address indexed from, address indexed to, address indexed contract_, uint256 tokenId, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) WatchDelegateERC721(opts *bind.WatchOpts, sink chan<- *DelegationDelegateERC721, from []common.Address, to []common.Address, contract_ []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var contract_Rule []interface{}
	for _, contract_Item := range contract_ {
		contract_Rule = append(contract_Rule, contract_Item)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "DelegateERC721", fromRule, toRule, contract_Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationDelegateERC721)
				if err := _Delegation.contract.UnpackLog(event, "DelegateERC721", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDelegateERC721 is a log parse operation binding the contract event 0x15e7a1bdcd507dd632d797d38e60cc5a9c0749b9a63097a215c4d006126825c6.
//
// Solidity: event DelegateERC721(address indexed from, address indexed to, address indexed contract_, uint256 tokenId, bytes32 rights, bool enable)
func (_Delegation *DelegationFilterer) ParseDelegateERC721(log types.Log) (*DelegationDelegateERC721, error) {
	event := new(DelegationDelegateERC721)
	if err := _Delegation.contract.UnpackLog(event, "DelegateERC721", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
