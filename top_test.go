package main

import (
	"fmt"
	"testing"

	"github.com/san-lab/go4337/state"
)

func TestLoadState(t *testing.T) {
	fmt.Println(state.State.Load())
	fmt.Println(state.State.Signers)
}
