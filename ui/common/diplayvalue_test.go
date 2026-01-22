package common

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/san-lab/go4337/userop"
)

func TestDiplay(t *testing.T) {
	it := &Item{}
	u, _ := uint256.FromDecimal("123")
	it.Value = u
	fmt.Println(it.DisplayValue())

	it.Value = *u
	fmt.Println(it.DisplayValue())

	u2 := userop.FullNonce(big.NewInt(1), uint64(7))
	it2 := new(Item)
	it2.Value = u2
	fmt.Println("i", it2.DisplayValue())

}
