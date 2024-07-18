package state

import (
	"testing"

	"github.com/san-lab/go4337/abiutil"
)

func TestParseABI(t *testing.T) {
	InitState()
	abistr := abiutil.Subjectstring
	_, _, err := ParseABI("AcceptTuple", abistr)
	if err != nil {

		t.Error(err)
	}
}
