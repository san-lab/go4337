package eip712

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func HashTypedData(data *apitypes.TypedData) ([]byte, error) {
	domain, err := data.HashStruct("EIP712Domain", data.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to get hash of typed data domain: %w", err)
	}
	dataHash, err := data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash of typed message: %w", err)
	}
	prefixedData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domain), string(dataHash)))
	prefixedDataHash := crypto.Keccak256(prefixedData)
	return prefixedDataHash, nil
}
