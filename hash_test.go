package main

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

func TestPackUserOp(t *testing.T) {
	uop, ok := state.GetUserOp("ThirdOne")
	if !ok {
		t.Error("UserOp not found")
	}
	bt, err := uop.EncodeToHash()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%x\n", bt)

	h1, err := userop.GetUsOpLibPrehashV6(uop)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%x\n", h1)

	bt2, err := userop.GetUserOpBytesToHashV6(uop, common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"), big.NewInt(1))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%x\n", bt2)

	h2, err := userop.GetUserOpHashV6(uop, common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"), big.NewInt(1))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%x\n", h2)

}
