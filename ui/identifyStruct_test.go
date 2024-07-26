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
	abi, _, err := state.ParseABI("EntrypointV6", abistr)
	if err != nil {
		t.Error(err)
	}
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
