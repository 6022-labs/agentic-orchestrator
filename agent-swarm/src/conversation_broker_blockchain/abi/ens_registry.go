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

// EnsRegistryMetaData contains all meta data concerning the EnsRegistry contract.
var EnsRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"resolver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EnsRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use EnsRegistryMetaData.ABI instead.
var EnsRegistryABI = EnsRegistryMetaData.ABI

// EnsRegistry is an auto generated Go binding around an Ethereum contract.
type EnsRegistry struct {
	EnsRegistryCaller     // Read-only binding to the contract
	EnsRegistryTransactor // Write-only binding to the contract
	EnsRegistryFilterer   // Log filterer for contract events
}

// EnsRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type EnsRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EnsRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EnsRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EnsRegistrySession struct {
	Contract     *EnsRegistry      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EnsRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EnsRegistryCallerSession struct {
	Contract *EnsRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// EnsRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EnsRegistryTransactorSession struct {
	Contract     *EnsRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// EnsRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type EnsRegistryRaw struct {
	Contract *EnsRegistry // Generic contract binding to access the raw methods on
}

// EnsRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EnsRegistryCallerRaw struct {
	Contract *EnsRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// EnsRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EnsRegistryTransactorRaw struct {
	Contract *EnsRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEnsRegistry creates a new instance of EnsRegistry, bound to a specific deployed contract.
func NewEnsRegistry(address common.Address, backend bind.ContractBackend) (*EnsRegistry, error) {
	contract, err := bindEnsRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EnsRegistry{EnsRegistryCaller: EnsRegistryCaller{contract: contract}, EnsRegistryTransactor: EnsRegistryTransactor{contract: contract}, EnsRegistryFilterer: EnsRegistryFilterer{contract: contract}}, nil
}

// NewEnsRegistryCaller creates a new read-only instance of EnsRegistry, bound to a specific deployed contract.
func NewEnsRegistryCaller(address common.Address, caller bind.ContractCaller) (*EnsRegistryCaller, error) {
	contract, err := bindEnsRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnsRegistryCaller{contract: contract}, nil
}

// NewEnsRegistryTransactor creates a new write-only instance of EnsRegistry, bound to a specific deployed contract.
func NewEnsRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*EnsRegistryTransactor, error) {
	contract, err := bindEnsRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnsRegistryTransactor{contract: contract}, nil
}

// NewEnsRegistryFilterer creates a new log filterer instance of EnsRegistry, bound to a specific deployed contract.
func NewEnsRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*EnsRegistryFilterer, error) {
	contract, err := bindEnsRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnsRegistryFilterer{contract: contract}, nil
}

// bindEnsRegistry binds a generic wrapper to an already deployed contract.
func bindEnsRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EnsRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnsRegistry *EnsRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnsRegistry.Contract.EnsRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnsRegistry *EnsRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnsRegistry.Contract.EnsRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnsRegistry *EnsRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnsRegistry.Contract.EnsRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnsRegistry *EnsRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnsRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnsRegistry *EnsRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnsRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnsRegistry *EnsRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnsRegistry.Contract.contract.Transact(opts, method, params...)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 node) view returns(address)
func (_EnsRegistry *EnsRegistryCaller) Resolver(opts *bind.CallOpts, node [32]byte) (common.Address, error) {
	var out []interface{}
	err := _EnsRegistry.contract.Call(opts, &out, "resolver", node)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 node) view returns(address)
func (_EnsRegistry *EnsRegistrySession) Resolver(node [32]byte) (common.Address, error) {
	return _EnsRegistry.Contract.Resolver(&_EnsRegistry.CallOpts, node)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(bytes32 node) view returns(address)
func (_EnsRegistry *EnsRegistryCallerSession) Resolver(node [32]byte) (common.Address, error) {
	return _EnsRegistry.Contract.Resolver(&_EnsRegistry.CallOpts, node)
}
