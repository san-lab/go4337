package zksyncera

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func TestDecode(t *testing.T) {
	encTx, _ := hex.DecodeString(EncodedTransaction1[2:])
	ztx := new(ZkSyncTxRLP)
	err := ztx.Decode(encTx)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ztx.Nonce)

	fmt.Println("Signature: ", ztx.CustomSignature)
	fmt.Println("from: ", ztx.From)
	fmt.Println("ChainId: ", ztx.ChainId2)

	tdat, err := ztx.TypedData()
	if err != nil {
		t.Error(err)
	}

	h, _, err := apitypes.TypedDataAndHash(*tdat)
	sig := ztx.CustomSignature
	sig[64] -= 27
	pub, err := crypto.SigToPub(h, sig)
	if err != nil {
		t.Error(err)
	}
	rfrom := crypto.PubkeyToAddress(*pub)
	fmt.Println("Recovered from: ", rfrom.Hex())

}
