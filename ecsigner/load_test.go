package ecsigner

import (
	"fmt"
	"testing"

	"github.com/san-lab/go4337/state"
)

func TestLoad(t *testing.T) {
	// TestLoad is a test for the Load function
	// of the StateStruct
	// This test is not exhaustive and should be
	// expanded to cover more cases
	state.Register(Type, AddECSigner, Unmarshal)
	_, ok := state.Unmarshalers[Type]
	fmt.Println(ok)
	err := state.State.Load()
	if err != nil {
		t.Error(err)
	}
	if len(state.State.Signers) == 0 {
		t.Error("No signers loaded")
	}

}
