package ecsigner

import (
	"fmt"
	"testing"

	"github.com/san-lab/go4337/state"
)

func TestMarshal(t *testing.T) {
	tp := Type
	keybytes := []byte("da3c0f646006521e15b4dfdad3d4c6f1e91caafdfb04c6aa65ea2b6fdc11b4b7")
	esig, err := Unmarshal(keybytes)
	if err != nil {
		t.Error(err)
	}
	bt, err := esig.Marshal()
	fmt.Println(tp)
	fmt.Println(string(bt))

	state.State.Signers = append(state.State.Signers, esig)
	state.State.Save()
}
