package userop

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (mop *UserOperation) GetRequiredPrefund() *big.Int {
	requiredGas := big.NewInt(0)
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.VerificationGasLimit)))
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.CallGasLimit)))
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.PaymasterVerificationGasLimit)))
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.PaymasterPostOpGasLimit)))
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.PreVerificationGas)))
	requiredGas.Mul(requiredGas, big.NewInt(int64(mop.MaxFeePerGas)))
	//requiredGas.Mul(requiredGas, big.NewInt(int64(mop.MaxPriorityFeePerGas)))
	/*requiredGas.Add(requiredGas, mop.CallGasLimit)
	requiredGas.Add(requiredGas, mop.PaymasterVerificationGasLimit)
	requiredGas.Add(requiredGas, mop.PaymasterPostOpGasLimit)
	requiredGas.Add(requiredGas, mop.PreVerificationGas)
	requiredGas.Mul(requiredGas, mop.MaxFeePerGas)*/
	return requiredGas
}

const PAYMASTER_VALIDATION_GAS_OFFSET = 20
const PAYMASTER_POSTOP_GAS_OFFSET = 36
const PAYMASTER_DATA_OFFSET = 52

func GetPaymasterV7Hash(puop *PackedUserOp, paymaster *common.Address, chainId, validUntil, validAfter uint64) ([]byte, []byte, error) {
	//can't use userOp.hash(), since it contains also the paymasterAndData itself.
	/*
	   address sender = userOp.getSender();
	   return
	       keccak256(
	       abi.encode(
	           sender,
	           userOp.nonce,
	           keccak256(userOp.initCode),
	           keccak256(userOp.callData),
	           userOp.accountGasLimits,
	           uint256(bytes32(userOp.paymasterAndData[PAYMASTER_VALIDATION_GAS_OFFSET : PAYMASTER_DATA_OFFSET])),
	           userOp.preVerificationGas,
	           userOp.gasFees,
	           block.chainid,
	           address(this),
	           validUntil,
	           validAfter
	       )
	   );
	*/
	if len(puop.PaymasterAndData) < PAYMASTER_DATA_OFFSET {
		return nil, nil, fmt.Errorf("PaymasterAndData too short")
	}
	if paymaster == nil {
		return nil, nil, fmt.Errorf("Paymaster address is nil")
	}
	args := abi.Arguments{
		{Type: addressTy}, //sender
		{Type: uint256Ty}, //nonce
		{Type: bytes32Ty}, //initCode
		{Type: bytes32Ty}, //callData
		{Type: bytes32Ty}, //accountGasLimits
		{Type: uint256Ty}, //PaymasterValidationGas
		{Type: uint256Ty}, //preVerificationGas
		{Type: bytes32Ty}, //gasFees
		{Type: uint256Ty}, //chainid
		{Type: addressTy}, //Paymaster's address
		{Type: uint48Ty},  //validUntil
		{Type: uint48Ty},  //validAfter
	}

	bts, err := args.Pack(
		puop.Sender,
		big.NewInt(int64(puop.Nonce)),
		UnsafeSliceToBytes32(crypto.Keccak256(puop.InitCode)),
		UnsafeSliceToBytes32(crypto.Keccak256(puop.CallData)),
		puop.AccountGasLimits,
		new(big.Int).SetBytes(puop.PaymasterAndData[PAYMASTER_VALIDATION_GAS_OFFSET:PAYMASTER_DATA_OFFSET]),
		big.NewInt(int64(puop.PreVerificationGas)),
		puop.GasFees,
		big.NewInt(int64(chainId)),
		*paymaster,
		big.NewInt(int64(validUntil)),
		big.NewInt(int64(validAfter)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Pack error in GetPaymasterV7Hash: %v", err)
	}

	return bts, crypto.Keccak256(bts), nil

}
