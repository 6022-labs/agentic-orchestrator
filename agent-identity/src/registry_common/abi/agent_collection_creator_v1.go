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

// AgentCollectionCreatorV1MetaData contains all meta data concerning the AgentCollectionCreatorV1 contract.
var AgentCollectionCreatorV1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialCollectionsManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialCollectionsDescriptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialCollectionImplementation\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedDeployment\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldDescriptor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newDescriptor\",\"type\":\"address\"}],\"name\":\"DefaultCollectionsDescriptorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collectionImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collectionsDescriptor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collectionsManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"moderator\",\"type\":\"address\"}],\"name\":\"createCollection\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"collectionAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newDescriptor\",\"type\":\"address\"}],\"name\":\"updateCollectionsDescriptor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AgentCollectionCreatorV1ABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentCollectionCreatorV1MetaData.ABI instead.
var AgentCollectionCreatorV1ABI = AgentCollectionCreatorV1MetaData.ABI

// AgentCollectionCreatorV1 is an auto generated Go binding around an Ethereum contract.
type AgentCollectionCreatorV1 struct {
	AgentCollectionCreatorV1Caller     // Read-only binding to the contract
	AgentCollectionCreatorV1Transactor // Write-only binding to the contract
	AgentCollectionCreatorV1Filterer   // Log filterer for contract events
}

// AgentCollectionCreatorV1Caller is an auto generated read-only Go binding around an Ethereum contract.
type AgentCollectionCreatorV1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentCollectionCreatorV1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentCollectionCreatorV1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentCollectionCreatorV1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentCollectionCreatorV1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentCollectionCreatorV1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentCollectionCreatorV1Session struct {
	Contract     *AgentCollectionCreatorV1 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AgentCollectionCreatorV1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentCollectionCreatorV1CallerSession struct {
	Contract *AgentCollectionCreatorV1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// AgentCollectionCreatorV1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentCollectionCreatorV1TransactorSession struct {
	Contract     *AgentCollectionCreatorV1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// AgentCollectionCreatorV1Raw is an auto generated low-level Go binding around an Ethereum contract.
type AgentCollectionCreatorV1Raw struct {
	Contract *AgentCollectionCreatorV1 // Generic contract binding to access the raw methods on
}

// AgentCollectionCreatorV1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentCollectionCreatorV1CallerRaw struct {
	Contract *AgentCollectionCreatorV1Caller // Generic read-only contract binding to access the raw methods on
}

// AgentCollectionCreatorV1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentCollectionCreatorV1TransactorRaw struct {
	Contract *AgentCollectionCreatorV1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentCollectionCreatorV1 creates a new instance of AgentCollectionCreatorV1, bound to a specific deployed contract.
func NewAgentCollectionCreatorV1(address common.Address, backend bind.ContractBackend) (*AgentCollectionCreatorV1, error) {
	contract, err := bindAgentCollectionCreatorV1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentCollectionCreatorV1{AgentCollectionCreatorV1Caller: AgentCollectionCreatorV1Caller{contract: contract}, AgentCollectionCreatorV1Transactor: AgentCollectionCreatorV1Transactor{contract: contract}, AgentCollectionCreatorV1Filterer: AgentCollectionCreatorV1Filterer{contract: contract}}, nil
}

// NewAgentCollectionCreatorV1Caller creates a new read-only instance of AgentCollectionCreatorV1, bound to a specific deployed contract.
func NewAgentCollectionCreatorV1Caller(address common.Address, caller bind.ContractCaller) (*AgentCollectionCreatorV1Caller, error) {
	contract, err := bindAgentCollectionCreatorV1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentCollectionCreatorV1Caller{contract: contract}, nil
}

// NewAgentCollectionCreatorV1Transactor creates a new write-only instance of AgentCollectionCreatorV1, bound to a specific deployed contract.
func NewAgentCollectionCreatorV1Transactor(address common.Address, transactor bind.ContractTransactor) (*AgentCollectionCreatorV1Transactor, error) {
	contract, err := bindAgentCollectionCreatorV1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentCollectionCreatorV1Transactor{contract: contract}, nil
}

// NewAgentCollectionCreatorV1Filterer creates a new log filterer instance of AgentCollectionCreatorV1, bound to a specific deployed contract.
func NewAgentCollectionCreatorV1Filterer(address common.Address, filterer bind.ContractFilterer) (*AgentCollectionCreatorV1Filterer, error) {
	contract, err := bindAgentCollectionCreatorV1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentCollectionCreatorV1Filterer{contract: contract}, nil
}

// bindAgentCollectionCreatorV1 binds a generic wrapper to an already deployed contract.
func bindAgentCollectionCreatorV1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentCollectionCreatorV1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentCollectionCreatorV1.Contract.AgentCollectionCreatorV1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.AgentCollectionCreatorV1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.AgentCollectionCreatorV1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentCollectionCreatorV1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentCollectionCreatorV1.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentCollectionCreatorV1.Contract.DEFAULTADMINROLE(&_AgentCollectionCreatorV1.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentCollectionCreatorV1.Contract.DEFAULTADMINROLE(&_AgentCollectionCreatorV1.CallOpts)
}

// CollectionImplementation is a free data retrieval call binding the contract method 0x0f3b901f.
//
// Solidity: function collectionImplementation() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Caller) CollectionImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentCollectionCreatorV1.contract.Call(opts, &out, "collectionImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollectionImplementation is a free data retrieval call binding the contract method 0x0f3b901f.
//
// Solidity: function collectionImplementation() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) CollectionImplementation() (common.Address, error) {
	return _AgentCollectionCreatorV1.Contract.CollectionImplementation(&_AgentCollectionCreatorV1.CallOpts)
}

// CollectionImplementation is a free data retrieval call binding the contract method 0x0f3b901f.
//
// Solidity: function collectionImplementation() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1CallerSession) CollectionImplementation() (common.Address, error) {
	return _AgentCollectionCreatorV1.Contract.CollectionImplementation(&_AgentCollectionCreatorV1.CallOpts)
}

// CollectionsDescriptor is a free data retrieval call binding the contract method 0xab1fdaef.
//
// Solidity: function collectionsDescriptor() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Caller) CollectionsDescriptor(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentCollectionCreatorV1.contract.Call(opts, &out, "collectionsDescriptor")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollectionsDescriptor is a free data retrieval call binding the contract method 0xab1fdaef.
//
// Solidity: function collectionsDescriptor() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) CollectionsDescriptor() (common.Address, error) {
	return _AgentCollectionCreatorV1.Contract.CollectionsDescriptor(&_AgentCollectionCreatorV1.CallOpts)
}

// CollectionsDescriptor is a free data retrieval call binding the contract method 0xab1fdaef.
//
// Solidity: function collectionsDescriptor() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1CallerSession) CollectionsDescriptor() (common.Address, error) {
	return _AgentCollectionCreatorV1.Contract.CollectionsDescriptor(&_AgentCollectionCreatorV1.CallOpts)
}

// CollectionsManager is a free data retrieval call binding the contract method 0x37c11b2a.
//
// Solidity: function collectionsManager() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Caller) CollectionsManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentCollectionCreatorV1.contract.Call(opts, &out, "collectionsManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollectionsManager is a free data retrieval call binding the contract method 0x37c11b2a.
//
// Solidity: function collectionsManager() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) CollectionsManager() (common.Address, error) {
	return _AgentCollectionCreatorV1.Contract.CollectionsManager(&_AgentCollectionCreatorV1.CallOpts)
}

// CollectionsManager is a free data retrieval call binding the contract method 0x37c11b2a.
//
// Solidity: function collectionsManager() view returns(address)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1CallerSession) CollectionsManager() (common.Address, error) {
	return _AgentCollectionCreatorV1.Contract.CollectionsManager(&_AgentCollectionCreatorV1.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _AgentCollectionCreatorV1.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgentCollectionCreatorV1.Contract.GetRoleAdmin(&_AgentCollectionCreatorV1.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgentCollectionCreatorV1.Contract.GetRoleAdmin(&_AgentCollectionCreatorV1.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _AgentCollectionCreatorV1.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgentCollectionCreatorV1.Contract.HasRole(&_AgentCollectionCreatorV1.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgentCollectionCreatorV1.Contract.HasRole(&_AgentCollectionCreatorV1.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AgentCollectionCreatorV1.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentCollectionCreatorV1.Contract.SupportsInterface(&_AgentCollectionCreatorV1.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentCollectionCreatorV1.Contract.SupportsInterface(&_AgentCollectionCreatorV1.CallOpts, interfaceId)
}

// CreateCollection is a paid mutator transaction binding the contract method 0xf652a195.
//
// Solidity: function createCollection(string name, string symbol, string description, address admin, address moderator) returns(address collectionAddress)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Transactor) CreateCollection(opts *bind.TransactOpts, name string, symbol string, description string, admin common.Address, moderator common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.contract.Transact(opts, "createCollection", name, symbol, description, admin, moderator)
}

// CreateCollection is a paid mutator transaction binding the contract method 0xf652a195.
//
// Solidity: function createCollection(string name, string symbol, string description, address admin, address moderator) returns(address collectionAddress)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) CreateCollection(name string, symbol string, description string, admin common.Address, moderator common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.CreateCollection(&_AgentCollectionCreatorV1.TransactOpts, name, symbol, description, admin, moderator)
}

// CreateCollection is a paid mutator transaction binding the contract method 0xf652a195.
//
// Solidity: function createCollection(string name, string symbol, string description, address admin, address moderator) returns(address collectionAddress)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1TransactorSession) CreateCollection(name string, symbol string, description string, admin common.Address, moderator common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.CreateCollection(&_AgentCollectionCreatorV1.TransactOpts, name, symbol, description, admin, moderator)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.GrantRole(&_AgentCollectionCreatorV1.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.GrantRole(&_AgentCollectionCreatorV1.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.RenounceRole(&_AgentCollectionCreatorV1.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1TransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.RenounceRole(&_AgentCollectionCreatorV1.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.RevokeRole(&_AgentCollectionCreatorV1.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.RevokeRole(&_AgentCollectionCreatorV1.TransactOpts, role, account)
}

// UpdateCollectionsDescriptor is a paid mutator transaction binding the contract method 0x798931a3.
//
// Solidity: function updateCollectionsDescriptor(address newDescriptor) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Transactor) UpdateCollectionsDescriptor(opts *bind.TransactOpts, newDescriptor common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.contract.Transact(opts, "updateCollectionsDescriptor", newDescriptor)
}

// UpdateCollectionsDescriptor is a paid mutator transaction binding the contract method 0x798931a3.
//
// Solidity: function updateCollectionsDescriptor(address newDescriptor) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Session) UpdateCollectionsDescriptor(newDescriptor common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.UpdateCollectionsDescriptor(&_AgentCollectionCreatorV1.TransactOpts, newDescriptor)
}

// UpdateCollectionsDescriptor is a paid mutator transaction binding the contract method 0x798931a3.
//
// Solidity: function updateCollectionsDescriptor(address newDescriptor) returns()
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1TransactorSession) UpdateCollectionsDescriptor(newDescriptor common.Address) (*types.Transaction, error) {
	return _AgentCollectionCreatorV1.Contract.UpdateCollectionsDescriptor(&_AgentCollectionCreatorV1.TransactOpts, newDescriptor)
}

// AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdatedIterator is returned from FilterDefaultCollectionsDescriptorUpdated and is used to iterate over the raw logs and unpacked data for DefaultCollectionsDescriptorUpdated events raised by the AgentCollectionCreatorV1 contract.
type AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdatedIterator struct {
	Event *AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated // Event containing the contract specifics and raw log

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
func (it *AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated)
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
		it.Event = new(AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated)
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
func (it *AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated represents a DefaultCollectionsDescriptorUpdated event raised by the AgentCollectionCreatorV1 contract.
type AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated struct {
	OldDescriptor common.Address
	NewDescriptor common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDefaultCollectionsDescriptorUpdated is a free log retrieval operation binding the contract event 0x1463407b87e87141e8c721b5e40c2784ff753056ff95f28e08a54ffa2d03befc.
//
// Solidity: event DefaultCollectionsDescriptorUpdated(address oldDescriptor, address newDescriptor)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) FilterDefaultCollectionsDescriptorUpdated(opts *bind.FilterOpts) (*AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdatedIterator, error) {

	logs, sub, err := _AgentCollectionCreatorV1.contract.FilterLogs(opts, "DefaultCollectionsDescriptorUpdated")
	if err != nil {
		return nil, err
	}
	return &AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdatedIterator{contract: _AgentCollectionCreatorV1.contract, event: "DefaultCollectionsDescriptorUpdated", logs: logs, sub: sub}, nil
}

// WatchDefaultCollectionsDescriptorUpdated is a free log subscription operation binding the contract event 0x1463407b87e87141e8c721b5e40c2784ff753056ff95f28e08a54ffa2d03befc.
//
// Solidity: event DefaultCollectionsDescriptorUpdated(address oldDescriptor, address newDescriptor)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) WatchDefaultCollectionsDescriptorUpdated(opts *bind.WatchOpts, sink chan<- *AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated) (event.Subscription, error) {

	logs, sub, err := _AgentCollectionCreatorV1.contract.WatchLogs(opts, "DefaultCollectionsDescriptorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated)
				if err := _AgentCollectionCreatorV1.contract.UnpackLog(event, "DefaultCollectionsDescriptorUpdated", log); err != nil {
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

// ParseDefaultCollectionsDescriptorUpdated is a log parse operation binding the contract event 0x1463407b87e87141e8c721b5e40c2784ff753056ff95f28e08a54ffa2d03befc.
//
// Solidity: event DefaultCollectionsDescriptorUpdated(address oldDescriptor, address newDescriptor)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) ParseDefaultCollectionsDescriptorUpdated(log types.Log) (*AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated, error) {
	event := new(AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated)
	if err := _AgentCollectionCreatorV1.contract.UnpackLog(event, "DefaultCollectionsDescriptorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentCollectionCreatorV1RoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the AgentCollectionCreatorV1 contract.
type AgentCollectionCreatorV1RoleAdminChangedIterator struct {
	Event *AgentCollectionCreatorV1RoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *AgentCollectionCreatorV1RoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentCollectionCreatorV1RoleAdminChanged)
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
		it.Event = new(AgentCollectionCreatorV1RoleAdminChanged)
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
func (it *AgentCollectionCreatorV1RoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentCollectionCreatorV1RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentCollectionCreatorV1RoleAdminChanged represents a RoleAdminChanged event raised by the AgentCollectionCreatorV1 contract.
type AgentCollectionCreatorV1RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*AgentCollectionCreatorV1RoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _AgentCollectionCreatorV1.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &AgentCollectionCreatorV1RoleAdminChangedIterator{contract: _AgentCollectionCreatorV1.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *AgentCollectionCreatorV1RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _AgentCollectionCreatorV1.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentCollectionCreatorV1RoleAdminChanged)
				if err := _AgentCollectionCreatorV1.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) ParseRoleAdminChanged(log types.Log) (*AgentCollectionCreatorV1RoleAdminChanged, error) {
	event := new(AgentCollectionCreatorV1RoleAdminChanged)
	if err := _AgentCollectionCreatorV1.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentCollectionCreatorV1RoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the AgentCollectionCreatorV1 contract.
type AgentCollectionCreatorV1RoleGrantedIterator struct {
	Event *AgentCollectionCreatorV1RoleGranted // Event containing the contract specifics and raw log

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
func (it *AgentCollectionCreatorV1RoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentCollectionCreatorV1RoleGranted)
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
		it.Event = new(AgentCollectionCreatorV1RoleGranted)
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
func (it *AgentCollectionCreatorV1RoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentCollectionCreatorV1RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentCollectionCreatorV1RoleGranted represents a RoleGranted event raised by the AgentCollectionCreatorV1 contract.
type AgentCollectionCreatorV1RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgentCollectionCreatorV1RoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgentCollectionCreatorV1.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgentCollectionCreatorV1RoleGrantedIterator{contract: _AgentCollectionCreatorV1.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *AgentCollectionCreatorV1RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgentCollectionCreatorV1.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentCollectionCreatorV1RoleGranted)
				if err := _AgentCollectionCreatorV1.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) ParseRoleGranted(log types.Log) (*AgentCollectionCreatorV1RoleGranted, error) {
	event := new(AgentCollectionCreatorV1RoleGranted)
	if err := _AgentCollectionCreatorV1.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentCollectionCreatorV1RoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the AgentCollectionCreatorV1 contract.
type AgentCollectionCreatorV1RoleRevokedIterator struct {
	Event *AgentCollectionCreatorV1RoleRevoked // Event containing the contract specifics and raw log

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
func (it *AgentCollectionCreatorV1RoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentCollectionCreatorV1RoleRevoked)
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
		it.Event = new(AgentCollectionCreatorV1RoleRevoked)
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
func (it *AgentCollectionCreatorV1RoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentCollectionCreatorV1RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentCollectionCreatorV1RoleRevoked represents a RoleRevoked event raised by the AgentCollectionCreatorV1 contract.
type AgentCollectionCreatorV1RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgentCollectionCreatorV1RoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgentCollectionCreatorV1.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgentCollectionCreatorV1RoleRevokedIterator{contract: _AgentCollectionCreatorV1.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *AgentCollectionCreatorV1RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgentCollectionCreatorV1.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentCollectionCreatorV1RoleRevoked)
				if err := _AgentCollectionCreatorV1.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentCollectionCreatorV1 *AgentCollectionCreatorV1Filterer) ParseRoleRevoked(log types.Log) (*AgentCollectionCreatorV1RoleRevoked, error) {
	event := new(AgentCollectionCreatorV1RoleRevoked)
	if err := _AgentCollectionCreatorV1.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// --- RawLog() methods added by go generate ---
func (e AgentCollectionCreatorV1DefaultCollectionsDescriptorUpdated) RawLog() types.Log { return e.Raw }
func (e AgentCollectionCreatorV1RoleAdminChanged) RawLog() types.Log { return e.Raw }
func (e AgentCollectionCreatorV1RoleGranted) RawLog() types.Log { return e.Raw }
func (e AgentCollectionCreatorV1RoleRevoked) RawLog() types.Log { return e.Raw }
