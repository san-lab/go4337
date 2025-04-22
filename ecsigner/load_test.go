package ecsigner

import (
	"crypto/ecdsa"
	"fmt"
	"testing"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/common"
)

func TestLoad(t *testing.T) {
	// TestLoad is a test for the Load function
	// of the StateStruct
	// This test is not exhaustive and should be
	// expanded to cover more cases
	state.Register(Type, AddECSigner, Unmarshal)
	_, ok := state.Unmarshalers[Type]
	fmt.Println(ok)
	err := state.Load()
	if err != nil {
		t.Error(err)
	}
	if len(state.GetSigners()) == 0 {
		t.Error("No signers loaded")
	}

}

func TestDisplayValue(t *testing.T) {
	it := &common.Item{Label: "Input new ECDSA private key in HEX"}

	it.Value = ECSigner{SignerAddress: ecommon.Address{0x01}, SignerKey: &ecdsa.PrivateKey{}}
	fmt.Println(it.DisplayValue())

}
