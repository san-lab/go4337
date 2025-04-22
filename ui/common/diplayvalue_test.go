package common

import (
	"fmt"
	"testing"

	"github.com/holiman/uint256"
)

func TestDiplay(t *testing.T) {
	it := &Item{}
	u, _ := uint256.FromDecimal("123")
	it.Value = u
	fmt.Println(it.DisplayValue())

	it.Value = *u
	fmt.Println(it.DisplayValue())
}
