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
		big.NewInt(int64(puop.Nonce)),
		UnsafeSliceToBytes32(crypto.Keccak256(puop.InitCode)),
		UnsafeSliceToBytes32(crypto.Keccak256(puop.CallData)),
		puop.AccountGasLimits,
		big.NewInt(int64(puop.PreVerificationGas)),
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
	return arguments.Pack(
		uop.Sender,
		big.NewInt(int64(uop.Nonce)),
		UnsafeSliceToBytes32(crypto.Keccak256(uop.InitData())),
		UnsafeSliceToBytes32(crypto.Keccak256(uop.CallData)),
		big.NewInt(int64(uop.CallGasLimit)),
		big.NewInt(int64(uop.VerificationGasLimit)),
		big.NewInt(int64(uop.PreVerificationGas)),
		big.NewInt(int64(uop.MaxFeePerGas)),
		big.NewInt(int64(uop.MaxPriorityFeePerGas)),
		UnsafeSliceToBytes32(crypto.Keccak256(uop.PaymasterAndData())),
	)
}

func GetUsOpLibPrehash(userOp *PackedUserOp) (hash [32]byte, err error) {
	enc1, err := userOp.EncodeToHash()
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
func GetUserOpHash(userOp *PackedUserOp, entryPoint common.Address, chainid uint64) (hash [32]byte, err error) {

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
