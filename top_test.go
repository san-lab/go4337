package main

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/san-lab/go4337/state"
)

func TestLoadState(t *testing.T) {
	//ecsigner.Init()
	//state.InitState()
	mabi, _, err := state.GetABI("Subject")

	if err != nil {
		t.Error(err)
	}
	methodName := "AcceptTuple"
	guoh := mabi.Methods[methodName]
	fmt.Println(len(guoh.Inputs))
	tup := guoh.Inputs[0]
	fmt.Println(tup.Type)
	fmt.Println(tup.Type.T == abi.TupleTy)
	fmt.Println(tup.Type.TupleElems)
	fmt.Println(tup.Type.TupleElems[0])
	inps := guoh.Inputs
	fmt.Println(inps)
	type TupleArgs struct{ A, B abi.Argument }
	var uint64Ty, _ = abi.NewType("uint64", "", nil)
	enc, err := inps.Pack(TupleArgs{
		A: abi.Argument{Type: uint64Ty, Name: "first"},
		B: abi.Argument{Type: uint64Ty, Name: "second"},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%x\n", enc)

}
