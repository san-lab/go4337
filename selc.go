// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

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

// SelcMetaData contains all meta data concerning the Selc contract.
var SelcMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes4[]\",\"name\":\"_sels\",\"type\":\"bytes4[]\"}],\"name\":\"setSelectors\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SelcABI is the input ABI used to generate the binding from.
// Deprecated: Use SelcMetaData.ABI instead.
var SelcABI = SelcMetaData.ABI

// Selc is an auto generated Go binding around an Ethereum contract.
type Selc struct {
	SelcCaller     // Read-only binding to the contract
	SelcTransactor // Write-only binding to the contract
	SelcFilterer   // Log filterer for contract events
}

// SelcCaller is an auto generated read-only Go binding around an Ethereum contract.
type SelcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SelcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SelcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SelcSession struct {
	Contract     *Selc             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SelcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SelcCallerSession struct {
	Contract *SelcCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SelcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SelcTransactorSession struct {
	Contract     *SelcTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SelcRaw is an auto generated low-level Go binding around an Ethereum contract.
type SelcRaw struct {
	Contract *Selc // Generic contract binding to access the raw methods on
}

// SelcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SelcCallerRaw struct {
	Contract *SelcCaller // Generic read-only contract binding to access the raw methods on
}

// SelcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SelcTransactorRaw struct {
	Contract *SelcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSelc creates a new instance of Selc, bound to a specific deployed contract.
func NewSelc(address common.Address, backend bind.ContractBackend) (*Selc, error) {
	contract, err := bindSelc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Selc{SelcCaller: SelcCaller{contract: contract}, SelcTransactor: SelcTransactor{contract: contract}, SelcFilterer: SelcFilterer{contract: contract}}, nil
}

// NewSelcCaller creates a new read-only instance of Selc, bound to a specific deployed contract.
func NewSelcCaller(address common.Address, caller bind.ContractCaller) (*SelcCaller, error) {
	contract, err := bindSelc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SelcCaller{contract: contract}, nil
}

// NewSelcTransactor creates a new write-only instance of Selc, bound to a specific deployed contract.
func NewSelcTransactor(address common.Address, transactor bind.ContractTransactor) (*SelcTransactor, error) {
	contract, err := bindSelc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SelcTransactor{contract: contract}, nil
}

// NewSelcFilterer creates a new log filterer instance of Selc, bound to a specific deployed contract.
func NewSelcFilterer(address common.Address, filterer bind.ContractFilterer) (*SelcFilterer, error) {
	contract, err := bindSelc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SelcFilterer{contract: contract}, nil
}

// bindSelc binds a generic wrapper to an already deployed contract.
func bindSelc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SelcMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Selc *SelcRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Selc.Contract.SelcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Selc *SelcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Selc.Contract.SelcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Selc *SelcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Selc.Contract.SelcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Selc *SelcCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Selc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Selc *SelcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Selc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Selc *SelcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Selc.Contract.contract.Transact(opts, method, params...)
}

// SetSelectors is a paid mutator transaction binding the contract method 0x9030799d.
//
// Solidity: function setSelectors(bytes4[] _sels) returns()
func (_Selc *SelcTransactor) SetSelectors(opts *bind.TransactOpts, _sels [][4]byte) (*types.Transaction, error) {
	return _Selc.contract.Transact(opts, "setSelectors", _sels)
}

// SetSelectors is a paid mutator transaction binding the contract method 0x9030799d.
//
// Solidity: function setSelectors(bytes4[] _sels) returns()
func (_Selc *SelcSession) SetSelectors(_sels [][4]byte) (*types.Transaction, error) {
	return _Selc.Contract.SetSelectors(&_Selc.TransactOpts, _sels)
}

// SetSelectors is a paid mutator transaction binding the contract method 0x9030799d.
//
// Solidity: function setSelectors(bytes4[] _sels) returns()
func (_Selc *SelcTransactorSession) SetSelectors(_sels [][4]byte) (*types.Transaction, error) {
	return _Selc.Contract.SetSelectors(&_Selc.TransactOpts, _sels)
}
