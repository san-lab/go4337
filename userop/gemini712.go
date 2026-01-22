package userop

import (
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Constants defined by ERC-4337 v0.7
var (
	// TYPE_HASH is the Keccak256 of the type definition string
	// PackedUserOperation(address sender,uint256 nonce,bytes initCode,bytes callData,bytes32 accountGasLimits,uint256 preVerificationGas,bytes32 gasFees,bytes paymasterAndData)
	TYPE_HASH = crypto.Keccak256Hash([]byte("PackedUserOperation(address sender,uint256 nonce,bytes initCode,bytes callData,bytes32 accountGasLimits,uint256 preVerificationGas,bytes32 gasFees,bytes paymasterAndData)"))
)

// SignUserOp computes the EIP-712 hash and signs it with the private key.
func SignUserOp(userOp *UserOperation, entryPoint common.Address, chainID *big.Int, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// 1. Get the EIP-712 digest (the hash to be signed)
	digest, err := GetUserOpHashV712(userOp, entryPoint, chainID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("userOpHash: 0x%x\n", digest)

	// 2. Sign the digest
	signature, err := crypto.Sign(digest.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	// 3. Add the 'V' value adjustment (standard Geth crypto.Sign returns V as 0 or 1, usually needs +27)
	// Note: ERC-4337 validation usually handles standard ECDSA recovery (v=27/28).
	if signature[64] < 27 {
		signature[64] += 27
	}

	return signature, nil
}

// GetUserOpHashV712 calculates the EIP-712 digest: Keccak256("\x19\x01" ‖ domainSeparator ‖ hashStruct(userOp))
func GetUserOpHashV712(op *UserOperation, entryPoint common.Address, chainID *big.Int) (common.Hash, error) {
	domainSeparator, err := buildDomainSeparator(entryPoint, chainID)
	if err != nil {
		return common.Hash{}, err
	}
	fmt.Printf("dom sep: 0x%x\n", domainSeparator)

	structHash := buildStructHash(op)

	// EIP-712 encoding: keccak256("\x19\x01" + domainSeparator + structHash)
	data := []byte{0x19, 0x01}
	data = append(data, domainSeparator.Bytes()...)
	data = append(data, structHash.Bytes()...)

	return crypto.Keccak256Hash(data), nil
}

// buildStructHash calculates hashStruct(PackedUserOperation)
func buildStructHash(op *UserOperation) common.Hash {
	// Logic for v0.7 / v8:
	// H(typeHash ‖ sender ‖ nonce ‖ keccak(initCode) ‖ keccak(callData) ‖ accountGasLimits ‖ preVerifGas ‖ gasFees ‖ keccak(paymasterAndData))

	// Prepare the complex fields
	initCode := getInitCode(op)
	paymasterAndData := getPaymasterAndData(op)
	accountGasLimits := packUint128s(op.VerificationGasLimit, op.CallGasLimit)
	gasFees := packUint128s(op.MaxPriorityFeePerGas, op.MaxFeePerGas)

	// Buffer to hold encoded data
	// 32 bytes per field * 9 fields
	encoded := make([]byte, 0, 32*9)

	// 1. TypeHash
	encoded = append(encoded, TYPE_HASH.Bytes()...)
	fmt.Printf("type hash: 0x%x\n", TYPE_HASH)
	// 2. Sender (padded to 32 bytes)
	encoded = append(encoded, common.LeftPadBytes(op.Sender.Bytes(), 32)...)

	// 3. Nonce (padded to 32 bytes)
	//nonceBytes := make([]byte, 32)
	//binary.BigEndian.PutUint64(nonceBytes[24:], op.Nonce) // Assuming nonce fits in uint64, otherwise use big.Int
	nonceBytes := op.Nonce.To32Bytes()
	encoded = append(encoded, nonceBytes...)

	// 4. Keccak(initCode)
	initCodeHash := crypto.Keccak256Hash(initCode)
	encoded = append(encoded, initCodeHash.Bytes()...)

	// 5. Keccak(callData)
	callDataHash := crypto.Keccak256Hash(op.CallData)
	encoded = append(encoded, callDataHash.Bytes()...)

	// 6. AccountGasLimits (bytes32)
	encoded = append(encoded, accountGasLimits[:]...)

	// 7. PreVerificationGas (uint256)
	preVerifBytes := make([]byte, 32)
	binary.BigEndian.PutUint64(preVerifBytes[24:], op.PreVerificationGas)
	encoded = append(encoded, preVerifBytes...)

	// 8. GasFees (bytes32)
	encoded = append(encoded, gasFees[:]...)

	// 9. Keccak(paymasterAndData)
	pmHash := crypto.Keccak256Hash(paymasterAndData)
	encoded = append(encoded, pmHash.Bytes()...)

	return crypto.Keccak256Hash(encoded)
}

// buildDomainSeparator constructs the EIP-712 domain separator
func buildDomainSeparator(entryPoint common.Address, chainID *big.Int) (common.Hash, error) {
	// Domain TypeHash: Keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)")
	domainTypeHash := crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))

	nameHash := crypto.Keccak256Hash([]byte("ERC4337"))
	versionHash := crypto.Keccak256Hash([]byte("1"))

	encoded := make([]byte, 0, 32*4)
	encoded = append(encoded, domainTypeHash.Bytes()...)
	encoded = append(encoded, nameHash.Bytes()...)
	encoded = append(encoded, versionHash.Bytes()...)
	encoded = append(encoded, common.BigToHash(chainID).Bytes()...)
	encoded = append(encoded, common.LeftPadBytes(entryPoint.Bytes(), 32)...)

	return crypto.Keccak256Hash(encoded), nil
}

// ---------------- Helpers for Packing v0.7 Fields ----------------

// packUint128s packs two uint64s into a single bytes32 (High << 128 | Low)
// Used for AccountGasLimits (VerificationGasLimit, CallGasLimit)
// and GasFees (MaxPriorityFeePerGas, MaxFeePerGas)
func packUint128s(high, low uint64) [32]byte {
	var packed [32]byte
	binary.BigEndian.PutUint64(packed[8:16], high)
	binary.BigEndian.PutUint64(packed[24:32], low)
	return packed
}

// getInitCode concatenates Factory address + FactoryData
func getInitCode(op *UserOperation) []byte {
	if op.Factory == nil || *op.Factory == (common.Address{}) {
		return []byte{}
	}
	// 20 bytes address + bytes data
	return append(op.Factory.Bytes(), op.FactoryData...)
}

// getPaymasterAndData concatenates Paymaster + Gas Limits + Data
// Format: [Paymaster (20 bytes)] [PaymasterVerificationGasLimit (16 bytes)] [PaymasterPostOpGasLimit (16 bytes)] [PaymasterData]
func getPaymasterAndData(op *UserOperation) []byte {
	if op.Paymaster == nil || *op.Paymaster == (common.Address{}) {
		return []byte{}
	}

	// Capacity: 20 + 16 + 16 + len(data)
	result := make([]byte, 0, 20+32+len(op.PaymasterData))

	// 1. Paymaster Address (20 bytes)
	result = append(result, op.Paymaster.Bytes()...)

	// 2. PaymasterVerificationGasLimit (16 bytes / uint128)
	verifGasBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(verifGasBytes[8:], op.PaymasterVerificationGasLimit)
	result = append(result, verifGasBytes...)

	// 3. PaymasterPostOpGasLimit (16 bytes / uint128)
	postOpGasBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(postOpGasBytes[8:], op.PaymasterPostOpGasLimit)
	result = append(result, postOpGasBytes...)

	// 4. Paymaster Data
	result = append(result, op.PaymasterData...)

	return result
}
