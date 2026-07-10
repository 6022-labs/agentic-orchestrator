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

// AgentEnsOffchainResolverMetaData contains all meta data concerning the AgentEnsOffchainResolver contract.
var AgentEnsOffchainResolverMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"urls_\",\"type\":\"string[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"urls\",\"type\":\"string[]\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"callbackFunction\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"OffchainLookup\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"expires\",\"type\":\"uint64\"}],\"name\":\"SignatureExpired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"UntrustedSigner\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"trusted\",\"type\":\"bool\"}],\"name\":\"SignerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"urls\",\"type\":\"string[]\"}],\"name\":\"UrlsUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"name\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"resolve\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"resolveWithProof\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"trusted\",\"type\":\"bool\"}],\"name\":\"setSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newUrls\",\"type\":\"string[]\"}],\"name\":\"setUrls\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"urls\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AgentEnsOffchainResolverABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentEnsOffchainResolverMetaData.ABI instead.
var AgentEnsOffchainResolverABI = AgentEnsOffchainResolverMetaData.ABI

// AgentEnsOffchainResolver is an auto generated Go binding around an Ethereum contract.
type AgentEnsOffchainResolver struct {
	AgentEnsOffchainResolverCaller     // Read-only binding to the contract
	AgentEnsOffchainResolverTransactor // Write-only binding to the contract
	AgentEnsOffchainResolverFilterer   // Log filterer for contract events
}

// AgentEnsOffchainResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentEnsOffchainResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentEnsOffchainResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentEnsOffchainResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentEnsOffchainResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentEnsOffchainResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentEnsOffchainResolverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentEnsOffchainResolverSession struct {
	Contract     *AgentEnsOffchainResolver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AgentEnsOffchainResolverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentEnsOffchainResolverCallerSession struct {
	Contract *AgentEnsOffchainResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// AgentEnsOffchainResolverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentEnsOffchainResolverTransactorSession struct {
	Contract     *AgentEnsOffchainResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// AgentEnsOffchainResolverRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentEnsOffchainResolverRaw struct {
	Contract *AgentEnsOffchainResolver // Generic contract binding to access the raw methods on
}

// AgentEnsOffchainResolverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentEnsOffchainResolverCallerRaw struct {
	Contract *AgentEnsOffchainResolverCaller // Generic read-only contract binding to access the raw methods on
}

// AgentEnsOffchainResolverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentEnsOffchainResolverTransactorRaw struct {
	Contract *AgentEnsOffchainResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentEnsOffchainResolver creates a new instance of AgentEnsOffchainResolver, bound to a specific deployed contract.
func NewAgentEnsOffchainResolver(address common.Address, backend bind.ContractBackend) (*AgentEnsOffchainResolver, error) {
	contract, err := bindAgentEnsOffchainResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolver{AgentEnsOffchainResolverCaller: AgentEnsOffchainResolverCaller{contract: contract}, AgentEnsOffchainResolverTransactor: AgentEnsOffchainResolverTransactor{contract: contract}, AgentEnsOffchainResolverFilterer: AgentEnsOffchainResolverFilterer{contract: contract}}, nil
}

// NewAgentEnsOffchainResolverCaller creates a new read-only instance of AgentEnsOffchainResolver, bound to a specific deployed contract.
func NewAgentEnsOffchainResolverCaller(address common.Address, caller bind.ContractCaller) (*AgentEnsOffchainResolverCaller, error) {
	contract, err := bindAgentEnsOffchainResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolverCaller{contract: contract}, nil
}

// NewAgentEnsOffchainResolverTransactor creates a new write-only instance of AgentEnsOffchainResolver, bound to a specific deployed contract.
func NewAgentEnsOffchainResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentEnsOffchainResolverTransactor, error) {
	contract, err := bindAgentEnsOffchainResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolverTransactor{contract: contract}, nil
}

// NewAgentEnsOffchainResolverFilterer creates a new log filterer instance of AgentEnsOffchainResolver, bound to a specific deployed contract.
func NewAgentEnsOffchainResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentEnsOffchainResolverFilterer, error) {
	contract, err := bindAgentEnsOffchainResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolverFilterer{contract: contract}, nil
}

// bindAgentEnsOffchainResolver binds a generic wrapper to an already deployed contract.
func bindAgentEnsOffchainResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentEnsOffchainResolverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentEnsOffchainResolver.Contract.AgentEnsOffchainResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.AgentEnsOffchainResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.AgentEnsOffchainResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentEnsOffchainResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentEnsOffchainResolver.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentEnsOffchainResolver.Contract.DEFAULTADMINROLE(&_AgentEnsOffchainResolver.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentEnsOffchainResolver.Contract.DEFAULTADMINROLE(&_AgentEnsOffchainResolver.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _AgentEnsOffchainResolver.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgentEnsOffchainResolver.Contract.GetRoleAdmin(&_AgentEnsOffchainResolver.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgentEnsOffchainResolver.Contract.GetRoleAdmin(&_AgentEnsOffchainResolver.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _AgentEnsOffchainResolver.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgentEnsOffchainResolver.Contract.HasRole(&_AgentEnsOffchainResolver.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgentEnsOffchainResolver.Contract.HasRole(&_AgentEnsOffchainResolver.CallOpts, role, account)
}

// Resolve is a free data retrieval call binding the contract method 0x9061b923.
//
// Solidity: function resolve(bytes name, bytes data) view returns(bytes)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCaller) Resolve(opts *bind.CallOpts, name []byte, data []byte) ([]byte, error) {
	var out []interface{}
	err := _AgentEnsOffchainResolver.contract.Call(opts, &out, "resolve", name, data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Resolve is a free data retrieval call binding the contract method 0x9061b923.
//
// Solidity: function resolve(bytes name, bytes data) view returns(bytes)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) Resolve(name []byte, data []byte) ([]byte, error) {
	return _AgentEnsOffchainResolver.Contract.Resolve(&_AgentEnsOffchainResolver.CallOpts, name, data)
}

// Resolve is a free data retrieval call binding the contract method 0x9061b923.
//
// Solidity: function resolve(bytes name, bytes data) view returns(bytes)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerSession) Resolve(name []byte, data []byte) ([]byte, error) {
	return _AgentEnsOffchainResolver.Contract.Resolve(&_AgentEnsOffchainResolver.CallOpts, name, data)
}

// ResolveWithProof is a free data retrieval call binding the contract method 0xf4d4d2f8.
//
// Solidity: function resolveWithProof(bytes response, bytes extraData) view returns(bytes)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCaller) ResolveWithProof(opts *bind.CallOpts, response []byte, extraData []byte) ([]byte, error) {
	var out []interface{}
	err := _AgentEnsOffchainResolver.contract.Call(opts, &out, "resolveWithProof", response, extraData)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ResolveWithProof is a free data retrieval call binding the contract method 0xf4d4d2f8.
//
// Solidity: function resolveWithProof(bytes response, bytes extraData) view returns(bytes)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) ResolveWithProof(response []byte, extraData []byte) ([]byte, error) {
	return _AgentEnsOffchainResolver.Contract.ResolveWithProof(&_AgentEnsOffchainResolver.CallOpts, response, extraData)
}

// ResolveWithProof is a free data retrieval call binding the contract method 0xf4d4d2f8.
//
// Solidity: function resolveWithProof(bytes response, bytes extraData) view returns(bytes)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerSession) ResolveWithProof(response []byte, extraData []byte) ([]byte, error) {
	return _AgentEnsOffchainResolver.Contract.ResolveWithProof(&_AgentEnsOffchainResolver.CallOpts, response, extraData)
}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCaller) Signers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AgentEnsOffchainResolver.contract.Call(opts, &out, "signers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) Signers(arg0 common.Address) (bool, error) {
	return _AgentEnsOffchainResolver.Contract.Signers(&_AgentEnsOffchainResolver.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerSession) Signers(arg0 common.Address) (bool, error) {
	return _AgentEnsOffchainResolver.Contract.Signers(&_AgentEnsOffchainResolver.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AgentEnsOffchainResolver.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentEnsOffchainResolver.Contract.SupportsInterface(&_AgentEnsOffchainResolver.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentEnsOffchainResolver.Contract.SupportsInterface(&_AgentEnsOffchainResolver.CallOpts, interfaceId)
}

// Urls is a free data retrieval call binding the contract method 0x796676be.
//
// Solidity: function urls(uint256 index) view returns(string)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCaller) Urls(opts *bind.CallOpts, index *big.Int) (string, error) {
	var out []interface{}
	err := _AgentEnsOffchainResolver.contract.Call(opts, &out, "urls", index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Urls is a free data retrieval call binding the contract method 0x796676be.
//
// Solidity: function urls(uint256 index) view returns(string)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) Urls(index *big.Int) (string, error) {
	return _AgentEnsOffchainResolver.Contract.Urls(&_AgentEnsOffchainResolver.CallOpts, index)
}

// Urls is a free data retrieval call binding the contract method 0x796676be.
//
// Solidity: function urls(uint256 index) view returns(string)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverCallerSession) Urls(index *big.Int) (string, error) {
	return _AgentEnsOffchainResolver.Contract.Urls(&_AgentEnsOffchainResolver.CallOpts, index)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.GrantRole(&_AgentEnsOffchainResolver.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.GrantRole(&_AgentEnsOffchainResolver.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.RenounceRole(&_AgentEnsOffchainResolver.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.RenounceRole(&_AgentEnsOffchainResolver.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.RevokeRole(&_AgentEnsOffchainResolver.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.RevokeRole(&_AgentEnsOffchainResolver.TransactOpts, role, account)
}

// SetSigner is a paid mutator transaction binding the contract method 0x31cb6105.
//
// Solidity: function setSigner(address signer, bool trusted) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactor) SetSigner(opts *bind.TransactOpts, signer common.Address, trusted bool) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.contract.Transact(opts, "setSigner", signer, trusted)
}

// SetSigner is a paid mutator transaction binding the contract method 0x31cb6105.
//
// Solidity: function setSigner(address signer, bool trusted) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) SetSigner(signer common.Address, trusted bool) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.SetSigner(&_AgentEnsOffchainResolver.TransactOpts, signer, trusted)
}

// SetSigner is a paid mutator transaction binding the contract method 0x31cb6105.
//
// Solidity: function setSigner(address signer, bool trusted) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactorSession) SetSigner(signer common.Address, trusted bool) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.SetSigner(&_AgentEnsOffchainResolver.TransactOpts, signer, trusted)
}

// SetUrls is a paid mutator transaction binding the contract method 0x6cc895a9.
//
// Solidity: function setUrls(string[] newUrls) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactor) SetUrls(opts *bind.TransactOpts, newUrls []string) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.contract.Transact(opts, "setUrls", newUrls)
}

// SetUrls is a paid mutator transaction binding the contract method 0x6cc895a9.
//
// Solidity: function setUrls(string[] newUrls) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverSession) SetUrls(newUrls []string) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.SetUrls(&_AgentEnsOffchainResolver.TransactOpts, newUrls)
}

// SetUrls is a paid mutator transaction binding the contract method 0x6cc895a9.
//
// Solidity: function setUrls(string[] newUrls) returns()
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverTransactorSession) SetUrls(newUrls []string) (*types.Transaction, error) {
	return _AgentEnsOffchainResolver.Contract.SetUrls(&_AgentEnsOffchainResolver.TransactOpts, newUrls)
}

// AgentEnsOffchainResolverRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverRoleAdminChangedIterator struct {
	Event *AgentEnsOffchainResolverRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *AgentEnsOffchainResolverRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsOffchainResolverRoleAdminChanged)
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
		it.Event = new(AgentEnsOffchainResolverRoleAdminChanged)
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
func (it *AgentEnsOffchainResolverRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsOffchainResolverRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsOffchainResolverRoleAdminChanged represents a RoleAdminChanged event raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*AgentEnsOffchainResolverRoleAdminChangedIterator, error) {

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

	logs, sub, err := _AgentEnsOffchainResolver.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolverRoleAdminChangedIterator{contract: _AgentEnsOffchainResolver.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *AgentEnsOffchainResolverRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _AgentEnsOffchainResolver.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsOffchainResolverRoleAdminChanged)
				if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) ParseRoleAdminChanged(log types.Log) (*AgentEnsOffchainResolverRoleAdminChanged, error) {
	event := new(AgentEnsOffchainResolverRoleAdminChanged)
	if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsOffchainResolverRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverRoleGrantedIterator struct {
	Event *AgentEnsOffchainResolverRoleGranted // Event containing the contract specifics and raw log

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
func (it *AgentEnsOffchainResolverRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsOffchainResolverRoleGranted)
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
		it.Event = new(AgentEnsOffchainResolverRoleGranted)
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
func (it *AgentEnsOffchainResolverRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsOffchainResolverRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsOffchainResolverRoleGranted represents a RoleGranted event raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgentEnsOffchainResolverRoleGrantedIterator, error) {

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

	logs, sub, err := _AgentEnsOffchainResolver.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolverRoleGrantedIterator{contract: _AgentEnsOffchainResolver.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *AgentEnsOffchainResolverRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AgentEnsOffchainResolver.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsOffchainResolverRoleGranted)
				if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) ParseRoleGranted(log types.Log) (*AgentEnsOffchainResolverRoleGranted, error) {
	event := new(AgentEnsOffchainResolverRoleGranted)
	if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsOffchainResolverRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverRoleRevokedIterator struct {
	Event *AgentEnsOffchainResolverRoleRevoked // Event containing the contract specifics and raw log

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
func (it *AgentEnsOffchainResolverRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsOffchainResolverRoleRevoked)
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
		it.Event = new(AgentEnsOffchainResolverRoleRevoked)
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
func (it *AgentEnsOffchainResolverRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsOffchainResolverRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsOffchainResolverRoleRevoked represents a RoleRevoked event raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgentEnsOffchainResolverRoleRevokedIterator, error) {

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

	logs, sub, err := _AgentEnsOffchainResolver.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolverRoleRevokedIterator{contract: _AgentEnsOffchainResolver.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *AgentEnsOffchainResolverRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AgentEnsOffchainResolver.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsOffchainResolverRoleRevoked)
				if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) ParseRoleRevoked(log types.Log) (*AgentEnsOffchainResolverRoleRevoked, error) {
	event := new(AgentEnsOffchainResolverRoleRevoked)
	if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsOffchainResolverSignerUpdatedIterator is returned from FilterSignerUpdated and is used to iterate over the raw logs and unpacked data for SignerUpdated events raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverSignerUpdatedIterator struct {
	Event *AgentEnsOffchainResolverSignerUpdated // Event containing the contract specifics and raw log

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
func (it *AgentEnsOffchainResolverSignerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsOffchainResolverSignerUpdated)
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
		it.Event = new(AgentEnsOffchainResolverSignerUpdated)
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
func (it *AgentEnsOffchainResolverSignerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsOffchainResolverSignerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsOffchainResolverSignerUpdated represents a SignerUpdated event raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverSignerUpdated struct {
	Signer  common.Address
	Trusted bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSignerUpdated is a free log retrieval operation binding the contract event 0xfcaa24b1276bfa7dbf77797c0a984b9df924acbeaabd48cd2f1b0eca379b78fa.
//
// Solidity: event SignerUpdated(address indexed signer, bool trusted)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) FilterSignerUpdated(opts *bind.FilterOpts, signer []common.Address) (*AgentEnsOffchainResolverSignerUpdatedIterator, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _AgentEnsOffchainResolver.contract.FilterLogs(opts, "SignerUpdated", signerRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolverSignerUpdatedIterator{contract: _AgentEnsOffchainResolver.contract, event: "SignerUpdated", logs: logs, sub: sub}, nil
}

// WatchSignerUpdated is a free log subscription operation binding the contract event 0xfcaa24b1276bfa7dbf77797c0a984b9df924acbeaabd48cd2f1b0eca379b78fa.
//
// Solidity: event SignerUpdated(address indexed signer, bool trusted)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) WatchSignerUpdated(opts *bind.WatchOpts, sink chan<- *AgentEnsOffchainResolverSignerUpdated, signer []common.Address) (event.Subscription, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _AgentEnsOffchainResolver.contract.WatchLogs(opts, "SignerUpdated", signerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsOffchainResolverSignerUpdated)
				if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "SignerUpdated", log); err != nil {
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

// ParseSignerUpdated is a log parse operation binding the contract event 0xfcaa24b1276bfa7dbf77797c0a984b9df924acbeaabd48cd2f1b0eca379b78fa.
//
// Solidity: event SignerUpdated(address indexed signer, bool trusted)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) ParseSignerUpdated(log types.Log) (*AgentEnsOffchainResolverSignerUpdated, error) {
	event := new(AgentEnsOffchainResolverSignerUpdated)
	if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "SignerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsOffchainResolverUrlsUpdatedIterator is returned from FilterUrlsUpdated and is used to iterate over the raw logs and unpacked data for UrlsUpdated events raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverUrlsUpdatedIterator struct {
	Event *AgentEnsOffchainResolverUrlsUpdated // Event containing the contract specifics and raw log

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
func (it *AgentEnsOffchainResolverUrlsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsOffchainResolverUrlsUpdated)
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
		it.Event = new(AgentEnsOffchainResolverUrlsUpdated)
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
func (it *AgentEnsOffchainResolverUrlsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsOffchainResolverUrlsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsOffchainResolverUrlsUpdated represents a UrlsUpdated event raised by the AgentEnsOffchainResolver contract.
type AgentEnsOffchainResolverUrlsUpdated struct {
	Urls []string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUrlsUpdated is a free log retrieval operation binding the contract event 0xc320ea45f3697c23c0a846b66e374215bea4a3cfef9a23233e9d9c2a1e7df068.
//
// Solidity: event UrlsUpdated(string[] urls)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) FilterUrlsUpdated(opts *bind.FilterOpts) (*AgentEnsOffchainResolverUrlsUpdatedIterator, error) {

	logs, sub, err := _AgentEnsOffchainResolver.contract.FilterLogs(opts, "UrlsUpdated")
	if err != nil {
		return nil, err
	}
	return &AgentEnsOffchainResolverUrlsUpdatedIterator{contract: _AgentEnsOffchainResolver.contract, event: "UrlsUpdated", logs: logs, sub: sub}, nil
}

// WatchUrlsUpdated is a free log subscription operation binding the contract event 0xc320ea45f3697c23c0a846b66e374215bea4a3cfef9a23233e9d9c2a1e7df068.
//
// Solidity: event UrlsUpdated(string[] urls)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) WatchUrlsUpdated(opts *bind.WatchOpts, sink chan<- *AgentEnsOffchainResolverUrlsUpdated) (event.Subscription, error) {

	logs, sub, err := _AgentEnsOffchainResolver.contract.WatchLogs(opts, "UrlsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsOffchainResolverUrlsUpdated)
				if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "UrlsUpdated", log); err != nil {
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

// ParseUrlsUpdated is a log parse operation binding the contract event 0xc320ea45f3697c23c0a846b66e374215bea4a3cfef9a23233e9d9c2a1e7df068.
//
// Solidity: event UrlsUpdated(string[] urls)
func (_AgentEnsOffchainResolver *AgentEnsOffchainResolverFilterer) ParseUrlsUpdated(log types.Log) (*AgentEnsOffchainResolverUrlsUpdated, error) {
	event := new(AgentEnsOffchainResolverUrlsUpdated)
	if err := _AgentEnsOffchainResolver.contract.UnpackLog(event, "UrlsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// --- RawLog() methods added by go generate ---
func (e AgentEnsOffchainResolverRoleAdminChanged) RawLog() types.Log { return e.Raw }
func (e AgentEnsOffchainResolverRoleGranted) RawLog() types.Log      { return e.Raw }
func (e AgentEnsOffchainResolverRoleRevoked) RawLog() types.Log      { return e.Raw }
func (e AgentEnsOffchainResolverSignerUpdated) RawLog() types.Log    { return e.Raw }
func (e AgentEnsOffchainResolverUrlsUpdated) RawLog() types.Log      { return e.Raw }
