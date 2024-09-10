package ui

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/state"
)

func TestIdentifyUsOpV6(t *testing.T) {
	abistr := entrypoint.EntryPointV6AbiJson
	abia, err := state.ParseABI("EntrypointV6", abistr)
	if err != nil {
		t.Error(err)
	}
	abi := abia.ABI
	simHops, ok := abi.Methods["simulateHandleOp"]
	if !ok {
		t.Error("Method not found")
	}

	fmt.Println(simHops.Inputs[0].Type.GetType())
	t1 := simHops.Inputs[0].Type.GetType()
	opHash := abi.Methods["getUserOpHash"]
	t2 := opHash.Inputs[0].Type.GetType()
	fmt.Println(t1 == t2)

	var v3 interface{}
	v3 = &abiutil.UserOperationV6{}
	t3 := reflect.ValueOf(v3).Type()
	fmt.Println(t1 == t3)
	fmt.Println(t1.String())
	fmt.Println("----")
	fmt.Println(t3.String())

}

func TestParseBigInt(t *testing.T) {
	st := []string{"1", ".2K", "300", "4K", ".05M", "0xaA34"}

	for _, s := range st {
		i, err := ParseBigInt(s)
		fmt.Printf("s: %s\n, i: %v\n, h: %x,  err: %v\n", s, i, i, err)
	}
}
