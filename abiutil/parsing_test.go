package abiutil

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestDecodeABI(t *testing.T) {
	abi, abiclrs, err := ParseABIFromString(Subjectstring)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(abiclrs))
	methodName := "AcceptTuple"
	guoh := abi.Methods[methodName]
	fmt.Println(len(guoh.Inputs))
	tupl := guoh.Inputs[0]
	fmt.Println(tupl.Type)
	v := reflect.New(tupl.Type.GetType())
	fmt.Println(v.Interface())
	v.Elem().Field(0).Set(reflect.ValueOf(uint64(4)))
	fmt.Println(v.Interface())
}

func TestSilceOfTuples(t *testing.T) {
	abi, _, err := ParseABIFromString(TestABI)
	if err != nil {
		t.Error(err)
	}
	acceptintslice := abi.Methods["AcceptintSlice"]
	fmt.Println(len(acceptintslice.Inputs))

	input := []uint16{1, 2, 3, 4, 5}

	data, err := abi.Pack("AcceptintSlice", input)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
	var a1 uint16 = 0
	var a2 uint16 = 1
	var a3 uint16 = 2
	var a4 uint16 = 3

	in := []interface{}{a1, a2, a3, a4}
	data, err = abi.Pack("AcceptintSlice", in)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)

}

func TestSlices(t *testing.T) {
	tp, err := abi.NewType("string", "", nil)
	if err != nil {
		t.Error(err)
	}
	slice := MakeSliceOfType(tp, 0, 0)
	slice, err = AppendToSlice(slice, "ok")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(slice)
	slice, err = AppendToSlice(slice, "ok2")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(slice)
	err = SetSliceValue(slice, 0, "ok3")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(slice)
}
