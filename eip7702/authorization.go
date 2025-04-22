package eip7702

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
	"github.com/san-lab/go4337/signer"
)

func Sign(sca *types.SetCodeAuthorization, s signer.Signer) error {
	h, err := Hash(sca)
	if err != nil {
		return fmt.Errorf("failed to hash authorization: %w", err)
	}
	sig, err := s.SignHash(h)
	if err != nil {
		return err
	}
	//calculate R and S
	// The sig is in the [R || S || V] format where V is 0 or 1.
	sca.R = *(uint256.NewInt(0).SetBytes(sig[:32]))
	sca.S = *(new(uint256.Int).SetBytes(sig[32:64]))
	sca.V = sig[64]

	return nil

}

func Hash(auth *types.SetCodeAuthorization) ([]byte, error) {
	// Prepare data to hash (ChainID, Address, Nonce)
	rlpData, err := rlp.EncodeToBytes([]interface{}{
		auth.ChainID,
		auth.Address,
		auth.Nonce,
	})
	if err != nil {
		return nil, err
	}

	// Magic constant (used to distinguish from normal transactions)
	magicConstant := []byte{5}

	// Combine magic constant and RLP-encoded data
	finalData := append(magicConstant, rlpData...)

	// Calculate the hash
	hash := crypto.Keccak256(finalData)

	return hash, nil
}

// Serialize encodes the authorization using RLP as specified by EIP-7702.
func Serialize(auth *types.SetCodeAuthorization) ([]byte, error) {
	var buf bytes.Buffer
	err := rlp.Encode(&buf, auth)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
