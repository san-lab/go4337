package userop

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var addressTy, _ = abi.NewType("address", "", nil)
var uint256Ty, _ = abi.NewType("uint256", "", nil)
var uint64Ty, _ = abi.NewType("uint64", "", nil)
var uint48Ty, _ = abi.NewType("uint48", "", nil)
var bytesTy, _ = abi.NewType("bytes", "", nil)
var bytes32Ty, _ = abi.NewType("bytes32", "", nil)

func (puop *PackedUserOp) EncodeToHash() ([]byte, error) {
	arguments := abi.Arguments{
		{Type: addressTy}, //sender
		{Type: uint256Ty}, //nonce
		{Type: bytes32Ty}, //initCode
		{Type: bytes32Ty}, //callData
		{Type: bytes32Ty}, //accountGasLimits
		{Type: uint256Ty}, //preVerificationGas
		{Type: bytes32Ty}, //gasFees
		{Type: bytes32Ty}, //paymasterAndData
		//{Type: bytesTy},   //signature
	}
	return arguments.Pack(
		puop.Sender,
		puop.Nonce,
		UnsafeSliceToBytes32(crypto.Keccak256(puop.InitCode)),
		UnsafeSliceToBytes32(crypto.Keccak256(puop.CallData)),
		puop.AccountGasLimits,
		puop.PreVerificationGas,
		puop.GasFees,
		UnsafeSliceToBytes32(crypto.Keccak256(puop.PaymasterAndData)),
	)
}

func (uop *UserOperation) EncodeToHash() ([]byte, error) {
	arguments := abi.Arguments{
		{Type: addressTy}, //sender
		{Type: uint256Ty}, //nonce
		{Type: bytes32Ty}, //initCode
		{Type: bytes32Ty}, //callData
		{Type: uint256Ty}, //callGasLimit
		{Type: uint256Ty}, //verificationGasLimit
		{Type: uint256Ty}, //preVerificationGas
		{Type: uint256Ty}, //maxFeePerGas
		{Type: uint256Ty}, //maxPriorityFeePerGas
		{Type: bytes32Ty}, //paymasterAndData
		//{Type: bytesTy},   //signature
	}
	u6 := uop.MarshalV6UserOp()
	return arguments.Pack(
		u6.Sender,
		u6.Nonce,
		UnsafeSliceToBytes32(crypto.Keccak256(u6.InitCode)),
		UnsafeSliceToBytes32(crypto.Keccak256(u6.CallData)),
		u6.CallGasLimit,
		u6.VerificationGasLimit,
		u6.PreVerificationGas,
		u6.MaxFeePerGas,
		u6.MaxPriorityFeePerGas,
		UnsafeSliceToBytes32(crypto.Keccak256(u6.PaymasterAndData)),
	)
}

func GetUsOpLibPrehash(pUserOp *PackedUserOp) (hash [32]byte, err error) {
	enc1, err := pUserOp.EncodeToHash()
	if err != nil {
		err = fmt.Errorf("encode error: %v", err)
		return
	}
	return UnsafeSliceToBytes32(crypto.Keccak256(enc1)), nil
}

func GetUsOpLibPrehashV6(userOp *UserOperation) (hash [32]byte, err error) {
	enc1, err := userOp.EncodeToHash()
	if err != nil {
		err = fmt.Errorf("encode error: %v", err)
		return
	}
	return UnsafeSliceToBytes32(crypto.Keccak256(enc1)), nil
}

/*
keccak256(abi.encode(UserOperationLib.hash(userOp), address(this), block.chainid));
*/
func GetUserOpHashV7(userOp *PackedUserOp, entryPoint common.Address, chainid uint64) (hash [32]byte, err error) {

	enc2, err := GetUserOpBytesToHash(userOp, entryPoint, chainid)
	if err != nil {
		err = fmt.Errorf("pack error: %v", err)
		return
	}
	return UnsafeSliceToBytes32(crypto.Keccak256(enc2)), nil

}

/*
keccak256(abi.encode(UserOperationLib.hash(userOp), address(this), block.chainid));
*/
func GetUserOpHashV6(userOp *UserOperation, entryPoint common.Address, chainid uint64) (hash [32]byte, err error) {
	enc2, err := GetUserOpBytesToHashV6(userOp, entryPoint, chainid)
	if err != nil {
		err = fmt.Errorf("pack error: %v", err)
		return
	}
	return UnsafeSliceToBytes32(crypto.Keccak256(enc2)), nil

}

func GetUserOpBytesToHash(userOp *PackedUserOp, entryPoint common.Address, chainid uint64) (encoded []byte, err error) {
	h1, err := GetUsOpLibPrehash(userOp)
	args := abi.Arguments{
		{Type: bytes32Ty},
		{Type: addressTy},
		{Type: uint256Ty},
	}
	return args.Pack(h1, entryPoint, big.NewInt(int64(chainid)))
}

func GetUserOpBytesToHashV6(userOp *UserOperation, entryPoint common.Address, chainid uint64) (encoded []byte, err error) {
	h1, err := GetUsOpLibPrehashV6(userOp)
	args := abi.Arguments{
		{Type: bytes32Ty},
		{Type: addressTy},
		{Type: uint256Ty},
	}
	return args.Pack(h1, entryPoint, big.NewInt(int64(chainid)))
}

type UserOpForApiV6 struct {
	Sender               string `json:"sender"`
	Nonce                string `json:"nonce"`
	InitCode             string `json:"initCode"`
	CallData             string `json:"callData"`
	CallGasLimit         string `json:"callGasLimit"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	PreVerificationGas   string `json:"preVerificationGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	PaymasterAndData     string `json:"paymasterAndData"`
	Signature            string `json:"signature"`
}

func (uop *UserOperation) ToUserOpForApiV6() *UserOpForApiV6 {
	sender := uop.Sender
	if sender == nil {
		sender = &common.Address{}
	}
	return &UserOpForApiV6{
		Sender:               sender.Hex(),
		Nonce:                fmt.Sprintf("0x%x", uop.Nonce),
		InitCode:             fmt.Sprintf("0x%x", uop.InitData()),
		CallData:             fmt.Sprintf("0x%x", uop.CallData),
		CallGasLimit:         fmt.Sprintf("0x%x", uop.CallGasLimit),
		VerificationGasLimit: fmt.Sprintf("0x%x", uop.VerificationGasLimit),
		PreVerificationGas:   fmt.Sprintf("0x%x", uop.PreVerificationGas),
		MaxFeePerGas:         fmt.Sprintf("0x%x", uop.MaxFeePerGas),
		MaxPriorityFeePerGas: fmt.Sprintf("0x%x", uop.MaxPriorityFeePerGas),
		PaymasterAndData:     fmt.Sprintf("0x%x", uop.PaymasterAndData()),
		Signature:            fmt.Sprintf("0x%x", uop.Signature),
	}
}
