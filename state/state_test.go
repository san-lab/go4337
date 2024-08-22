package state

import (
	"fmt"
	"testing"

	"github.com/san-lab/go4337/abiutil"
)

func TestPersistBytecode(t *testing.T) {
	testABIStr := abiutil.Subjectstring
	a1, err := ParseABI("test", testABIStr)
	if err != nil {
		t.Error(err)
	}
	a2, err := GetABI("test")
	if err != nil {
		t.Error(err)
	}
	if a1 != a2 {
		fmt.Println("Not equal")
	}
	a1.ExecuteBytecode = []byte{0x01, 0x02, 0x03}
	err = state.Save()
	if err != nil {
		t.Error(err)
	}
	err = state.Load()
	if err != nil {
		t.Error(err)
	}
	a3, err := GetABI("test")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a3.ExecuteBytecode)
}
