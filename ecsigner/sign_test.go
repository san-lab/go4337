package ecsigner

import (
	"fmt"
	"testing"
)

func TestSignature(t *testing.T) {
	msg := []byte{4, 5, 6}
	privkeyhex := "c522c068090d4e888dadbab9967fd81a79a451aff84dce2040df59ad5a6ce1e8"

	s, err := FromHexKey("key", privkeyhex)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := s.Sign(msg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Signature: %x\n", sig)

}
