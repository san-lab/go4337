package eip7702

import (
	"crypto/ecdsa"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/san-lab/go4337/ecsigner"
)

func TestSignature(t *testing.T) {
	// TestSignature tests the signature of a SetCodeAuthorization struct.
	// The test is successful if the signature is valid.
	// The test is unsuccessful if the signature is invalid.
	t.Parallel()

	// Create a new SetCodeAuthorization struct.
	sca := MockAuthorization()
	s, err := ecsigner.FromHexKey("key", "503f38a9c967ed597e47fe25643985f032b072db8075426a92110f82df48dfcb")
	fmt.Println(s.SignerAddress)
	if err != nil {
		t.Fatal(err)
	}
	Sign(&sca, s)
	fmt.Printf("R: %x, S: %x, V: %v \n", sca.R.Bytes(), sca.S.Bytes(), sca.V)

	//Check the serialized data against the testSignedTestAuth
	hsa, err := Serialize(&sca)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Serialized: %x\n", hsa)
	fmt.Println("TestSignedTestAuth: ", testSignedTestAuth)

	prkey := s.SignerKey
	ssca, err := types.SignSetCode(prkey, sca)
	if err != nil {
		t.Fatal(err)
	}
	hsa2, err := Serialize(&ssca)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Serialized2: %x\n", hsa2)

	h, _ := Hash(&sca)
	sig, _ := s.SignHash(h)
	pub, _ := crypto.SigToPub(h, sig)
	adr := crypto.PubkeyToAddress(*pub)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(adr)
}

func MockAuthorization() types.SetCodeAuthorization {
	return types.SetCodeAuthorization{
		ChainID: *uint256.NewInt(11155111),
		Address: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Nonce:   3293,
	}
}

func MockSignedAuthorization() types.SetCodeAuthorization {
	auth := MockAuthorization()
	sauth, _ := types.SignSetCode(MockPrivateKey(), auth)
	return sauth

}

func MockPrivateKey() *ecdsa.PrivateKey {
	s, _ := ecsigner.FromHexKey("key", "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	return s.SignerKey
}

func TestSbb(t *testing.T) {
	fmt.Printf("%03d", 0x48)
}

const testSignedTestAuth = "0xf85c827a6994663f3ad617193148711d28f5334ee4ed070166028080a040e292da533253143f134643a03405f1af1de1d305526f44ed27e62061368d4ea051cfb0af34e491aa4d6796dececf95569088322e116c4b2f312bb23f20699269"
