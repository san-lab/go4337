// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package baseaccv6

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

// PackedUserOperation is an auto generated low-level Go binding around an user-defined struct.
type PackedUserOperation struct {
	Sender             common.Address
	Nonce              *big.Int
	InitCode           []byte
	CallData           []byte
	AccountGasLimits   [32]byte
	PreVerificationGas *big.Int
	GasFees            [32]byte
	PaymasterAndData   []byte
	Signature          []byte
}

// Baseaccv6MetaData contains all meta data concerning the Baseaccv6 contract.
var Baseaccv6MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"entryPoint\",\"outputs\":[{\"internalType\":\"contractIEntryPoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"accountGasLimits\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"gasFees\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structPackedUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"missingAccountFunds\",\"type\":\"uint256\"}],\"name\":\"validateUserOp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// Baseaccv6ABI is the input ABI used to generate the binding from.
// Deprecated: Use Baseaccv6MetaData.ABI instead.
var Baseaccv6ABI = Baseaccv6MetaData.ABI

// Baseaccv6 is an auto generated Go binding around an Ethereum contract.
type Baseaccv6 struct {
	Baseaccv6Caller     // Read-only binding to the contract
	Baseaccv6Transactor // Write-only binding to the contract
	Baseaccv6Filterer   // Log filterer for contract events
}

// Baseaccv6Caller is an auto generated read-only Go binding around an Ethereum contract.
type Baseaccv6Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Baseaccv6Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Baseaccv6Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Baseaccv6Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Baseaccv6Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Baseaccv6Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Baseaccv6Session struct {
	Contract     *Baseaccv6        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Baseaccv6CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Baseaccv6CallerSession struct {
	Contract *Baseaccv6Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// Baseaccv6TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Baseaccv6TransactorSession struct {
	Contract     *Baseaccv6Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Baseaccv6Raw is an auto generated low-level Go binding around an Ethereum contract.
type Baseaccv6Raw struct {
	Contract *Baseaccv6 // Generic contract binding to access the raw methods on
}

// Baseaccv6CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Baseaccv6CallerRaw struct {
	Contract *Baseaccv6Caller // Generic read-only contract binding to access the raw methods on
}

// Baseaccv6TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Baseaccv6TransactorRaw struct {
	Contract *Baseaccv6Transactor // Generic write-only contract binding to access the raw methods on
}

// NewBaseaccv6 creates a new instance of Baseaccv6, bound to a specific deployed contract.
func NewBaseaccv6(address common.Address, backend bind.ContractBackend) (*Baseaccv6, error) {
	contract, err := bindBaseaccv6(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Baseaccv6{Baseaccv6Caller: Baseaccv6Caller{contract: contract}, Baseaccv6Transactor: Baseaccv6Transactor{contract: contract}, Baseaccv6Filterer: Baseaccv6Filterer{contract: contract}}, nil
}

// NewBaseaccv6Caller creates a new read-only instance of Baseaccv6, bound to a specific deployed contract.
func NewBaseaccv6Caller(address common.Address, caller bind.ContractCaller) (*Baseaccv6Caller, error) {
	contract, err := bindBaseaccv6(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Baseaccv6Caller{contract: contract}, nil
}

// NewBaseaccv6Transactor creates a new write-only instance of Baseaccv6, bound to a specific deployed contract.
func NewBaseaccv6Transactor(address common.Address, transactor bind.ContractTransactor) (*Baseaccv6Transactor, error) {
	contract, err := bindBaseaccv6(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Baseaccv6Transactor{contract: contract}, nil
}

// NewBaseaccv6Filterer creates a new log filterer instance of Baseaccv6, bound to a specific deployed contract.
func NewBaseaccv6Filterer(address common.Address, filterer bind.ContractFilterer) (*Baseaccv6Filterer, error) {
	contract, err := bindBaseaccv6(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Baseaccv6Filterer{contract: contract}, nil
}

// bindBaseaccv6 binds a generic wrapper to an already deployed contract.
func bindBaseaccv6(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Baseaccv6MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Baseaccv6 *Baseaccv6Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Baseaccv6.Contract.Baseaccv6Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Baseaccv6 *Baseaccv6Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Baseaccv6.Contract.Baseaccv6Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Baseaccv6 *Baseaccv6Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Baseaccv6.Contract.Baseaccv6Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Baseaccv6 *Baseaccv6CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Baseaccv6.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Baseaccv6 *Baseaccv6TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Baseaccv6.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Baseaccv6 *Baseaccv6TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Baseaccv6.Contract.contract.Transact(opts, method, params...)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_Baseaccv6 *Baseaccv6Caller) EntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Baseaccv6.contract.Call(opts, &out, "entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_Baseaccv6 *Baseaccv6Session) EntryPoint() (common.Address, error) {
	return _Baseaccv6.Contract.EntryPoint(&_Baseaccv6.CallOpts)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_Baseaccv6 *Baseaccv6CallerSession) EntryPoint() (common.Address, error) {
	return _Baseaccv6.Contract.EntryPoint(&_Baseaccv6.CallOpts)
}

// GetNonce is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_Baseaccv6 *Baseaccv6Caller) GetNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Baseaccv6.contract.Call(opts, &out, "getNonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonce is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_Baseaccv6 *Baseaccv6Session) GetNonce() (*big.Int, error) {
	return _Baseaccv6.Contract.GetNonce(&_Baseaccv6.CallOpts)
}

// GetNonce is a free data retrieval call binding the contract method 0xd087d288.
//
// Solidity: function getNonce() view returns(uint256)
func (_Baseaccv6 *Baseaccv6CallerSession) GetNonce() (*big.Int, error) {
	return _Baseaccv6.Contract.GetNonce(&_Baseaccv6.CallOpts)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x19822f7c.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (_Baseaccv6 *Baseaccv6Transactor) ValidateUserOp(opts *bind.TransactOpts, userOp PackedUserOperation, userOpHash [32]byte, missingAccountFunds *big.Int) (*types.Transaction, error) {
	return _Baseaccv6.contract.Transact(opts, "validateUserOp", userOp, userOpHash, missingAccountFunds)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x19822f7c.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (_Baseaccv6 *Baseaccv6Session) ValidateUserOp(userOp PackedUserOperation, userOpHash [32]byte, missingAccountFunds *big.Int) (*types.Transaction, error) {
	return _Baseaccv6.Contract.ValidateUserOp(&_Baseaccv6.TransactOpts, userOp, userOpHash, missingAccountFunds)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x19822f7c.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (_Baseaccv6 *Baseaccv6TransactorSession) ValidateUserOp(userOp PackedUserOperation, userOpHash [32]byte, missingAccountFunds *big.Int) (*types.Transaction, error) {
	return _Baseaccv6.Contract.ValidateUserOp(&_Baseaccv6.TransactOpts, userOp, userOpHash, missingAccountFunds)
}
