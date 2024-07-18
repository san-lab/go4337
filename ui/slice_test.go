package ui

import (
	"fmt"
	"testing"

	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/state"
)

func TestSlice(t *testing.T) {
	state.InitState()
	abistr := abiutil.TestABI
	abi, _, err := state.ParseABI("AcceptintSlice", abistr)
	if err != nil {

		t.Error(err)
	}
	acceptintslice := abi.Methods["AcceptintSlice"]
	arg := acceptintslice.Inputs[0]
	slice := abiutil.MakeSliceOfType(*arg.Type.Elem, 0, 0)
	slice, err = abiutil.AppendToSlice(slice, uint16(1))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(slice)
}
