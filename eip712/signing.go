package eip712

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func NewSigner() {
	sig := core.SignerAPI{}
	fmt.Println(sig)
	sig.SignTypedData(context.Background(), common.NewMixedcaseAddress(common.Address{}), apitypes.TypedData{})
}
