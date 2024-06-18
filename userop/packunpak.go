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

func GetUsOpLibPrehash(userOp *PackedUserOp) (hash [32]byte, err error) {
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

func GetUserOpBytesToHash(userOp *PackedUserOp, entryPoint common.Address, chainid uint64) (encoded []byte, err error) {
	h1, err := GetUsOpLibPrehash(userOp)
	args := abi.Arguments{
		{Type: bytes32Ty},
		{Type: addressTy},
		{Type: uint256Ty},
	}
	return args.Pack(h1, entryPoint, big.NewInt(int64(chainid)))
}
