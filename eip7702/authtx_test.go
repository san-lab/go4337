package eip7702

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/san-lab/go4337/ecsigner"
)

func TestEncoding(t *testing.T) {
	to := common.HexToAddress("0x663F3ad617193148711d28f5334eE4Ed07016602")
	sadt := types.SetCodeTx{
		ChainID:    uint256.NewInt(1115511),
		Nonce:      1,
		GasTipCap:  uint256.NewInt(10_000_000_000),
		GasFeeCap:  uint256.NewInt(30_000_000_000),
		Gas:        100_000,
		To:         to,
		Value:      uint256.NewInt(0),
		Data:       []byte{},
		AccessList: types.AccessList{},
		AuthList: []types.SetCodeAuthorization{
			MockSignedAuthorization(),
		},
		V: uint256.NewInt(0x1b),
		R: uint256.NewInt(8),
		S: uint256.NewInt(9),
	}
	enc, err := Encode(&sadt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Encoded: %x\n", enc)

	signer, err := ecsigner.FromHexKey("key", "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	if err != nil {
		t.Fatal(err)
	}
	stx, err := types.SignNewTx(signer.SignerKey, types.NewPragueSigner(sadt.ChainID.ToBig()), &sadt)
	if err != nil {
		t.Fatal(err)
	}

	w := new(bytes.Buffer)
	stx.EncodeRLP(w)
	enc1 := w.Bytes()
	t.Logf("Encoded1: %x\n", enc1)

	j, err := json.MarshalIndent(sadt, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("JSON: %s\n", j)

	sadt2 := types.SetCodeTx{}
	err = json.Unmarshal(j, &sadt2)
	if err != nil {
		t.Fatal(err)
	}

	stx2, err := types.SignNewTx(signer.SignerKey, types.NewPragueSigner(sadt2.ChainID.ToBig()), &sadt2)
	if err != nil {
		t.Fatal(err)
	}

	w.Reset()
	stx2.EncodeRLP(w)
	enc2 := w.Bytes()
	fmt.Println(bytes.Compare(enc1, enc2))

}
