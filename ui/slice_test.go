package ui

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/state"
	ucommon "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/userop"
)

func TestSlice(t *testing.T) {
	//state.InitState()
	abistr := abiutil.TestABI
	abia, err := state.ParseABI("AcceptintSlice", abistr)
	if err != nil {

		t.Error(err)
	}
	abi := abia.ABI
	acceptintslice := abi.Methods["AcceptintSlice"]
	arg := acceptintslice.Inputs[0]
	slice := abiutil.MakeSliceOfType(*arg.Type.Elem, 0, 0)
	slice, err = abiutil.AppendToSlice(slice, uint16(1))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(slice)
}

func TestDisplayValue(t *testing.T) {
	//state.InitState()
	it := new(ucommon.Item)
	it.Value = common.Address{}
	fmt.Println(it.DisplayValue())
	it.Value = &it.Value
	fmt.Println(it.DisplayValue())
	usop := new(userop.UserOperation)
	it.Value = usop.Sender
	fmt.Println(it.Value == nil)
	fmt.Println(reflect.ValueOf(it.Value).IsNil())
	//fmt.Println(elem)
	fmt.Println(">>" + it.DisplayValue() + "<<")
}
