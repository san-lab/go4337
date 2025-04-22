package eip7702

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type SetAuthTx struct {
	ChainID        *big.Int
	Nonce          uint64
	GasTipCap      *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap      *big.Int // a.k.a. maxFeePerGas
	Gas            uint64
	To             *common.Address // do not need this in the context `rlp:"nil"` // nil means contract creation
	Value          *big.Int
	Data           []byte
	AccessList     types.AccessList
	Authorizations []types.SetCodeAuthorization

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

// RLP Encode
func Encode(satx *types.SetCodeTx) ([]byte, error) {
	w := new(bytes.Buffer)
	w.WriteByte(0x04)
	rlp.Encode(w, satx)
	return w.Bytes(), nil
}
