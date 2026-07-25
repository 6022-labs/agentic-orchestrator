// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

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

// EnsTextResolverMetaData contains all meta data concerning the EnsTextResolver contract.
var EnsTextResolverMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"text\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EnsTextResolverABI is the input ABI used to generate the binding from.
// Deprecated: Use EnsTextResolverMetaData.ABI instead.
var EnsTextResolverABI = EnsTextResolverMetaData.ABI

// EnsTextResolver is an auto generated Go binding around an Ethereum contract.
type EnsTextResolver struct {
	EnsTextResolverCaller     // Read-only binding to the contract
	EnsTextResolverTransactor // Write-only binding to the contract
	EnsTextResolverFilterer   // Log filterer for contract events
}

// EnsTextResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type EnsTextResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsTextResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EnsTextResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsTextResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EnsTextResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsTextResolverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EnsTextResolverSession struct {
	Contract     *EnsTextResolver  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EnsTextResolverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EnsTextResolverCallerSession struct {
	Contract *EnsTextResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// EnsTextResolverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EnsTextResolverTransactorSession struct {
	Contract     *EnsTextResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// EnsTextResolverRaw is an auto generated low-level Go binding around an Ethereum contract.
type EnsTextResolverRaw struct {
	Contract *EnsTextResolver // Generic contract binding to access the raw methods on
}

// EnsTextResolverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EnsTextResolverCallerRaw struct {
	Contract *EnsTextResolverCaller // Generic read-only contract binding to access the raw methods on
}

// EnsTextResolverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EnsTextResolverTransactorRaw struct {
	Contract *EnsTextResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEnsTextResolver creates a new instance of EnsTextResolver, bound to a specific deployed contract.
func NewEnsTextResolver(address common.Address, backend bind.ContractBackend) (*EnsTextResolver, error) {
	contract, err := bindEnsTextResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EnsTextResolver{EnsTextResolverCaller: EnsTextResolverCaller{contract: contract}, EnsTextResolverTransactor: EnsTextResolverTransactor{contract: contract}, EnsTextResolverFilterer: EnsTextResolverFilterer{contract: contract}}, nil
}

// NewEnsTextResolverCaller creates a new read-only instance of EnsTextResolver, bound to a specific deployed contract.
func NewEnsTextResolverCaller(address common.Address, caller bind.ContractCaller) (*EnsTextResolverCaller, error) {
	contract, err := bindEnsTextResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnsTextResolverCaller{contract: contract}, nil
}

// NewEnsTextResolverTransactor creates a new write-only instance of EnsTextResolver, bound to a specific deployed contract.
func NewEnsTextResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*EnsTextResolverTransactor, error) {
	contract, err := bindEnsTextResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnsTextResolverTransactor{contract: contract}, nil
}

// NewEnsTextResolverFilterer creates a new log filterer instance of EnsTextResolver, bound to a specific deployed contract.
func NewEnsTextResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*EnsTextResolverFilterer, error) {
	contract, err := bindEnsTextResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnsTextResolverFilterer{contract: contract}, nil
}

// bindEnsTextResolver binds a generic wrapper to an already deployed contract.
func bindEnsTextResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EnsTextResolverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnsTextResolver *EnsTextResolverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnsTextResolver.Contract.EnsTextResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnsTextResolver *EnsTextResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnsTextResolver.Contract.EnsTextResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnsTextResolver *EnsTextResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnsTextResolver.Contract.EnsTextResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnsTextResolver *EnsTextResolverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnsTextResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnsTextResolver *EnsTextResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnsTextResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnsTextResolver *EnsTextResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnsTextResolver.Contract.contract.Transact(opts, method, params...)
}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(bytes32 node, string key) view returns(string)
func (_EnsTextResolver *EnsTextResolverCaller) Text(opts *bind.CallOpts, node [32]byte, key string) (string, error) {
	var out []interface{}
	err := _EnsTextResolver.contract.Call(opts, &out, "text", node, key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(bytes32 node, string key) view returns(string)
func (_EnsTextResolver *EnsTextResolverSession) Text(node [32]byte, key string) (string, error) {
	return _EnsTextResolver.Contract.Text(&_EnsTextResolver.CallOpts, node, key)
}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(bytes32 node, string key) view returns(string)
func (_EnsTextResolver *EnsTextResolverCallerSession) Text(node [32]byte, key string) (string, error) {
	return _EnsTextResolver.Contract.Text(&_EnsTextResolver.CallOpts, node, key)
}
