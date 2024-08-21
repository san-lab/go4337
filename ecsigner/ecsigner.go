package ecsigner

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui"
)

func Init() {
	state.Register(Type, AddECSigner, Unmarshal)
}

type ECSigner struct {
	SignerAddress common.Address
	SignerKey     *ecdsa.PrivateKey
}

const Type = "ECDSA plain"

func (ecsigner *ECSigner) Type() string {
	return Type
}

func (ecsigner *ECSigner) Sign(mssg []byte) ([]byte, error) {
	mhash := signer.ToEthSignedMessageHash(mssg)
	sig, err := crypto.Sign(mhash[:], ecsigner.SignerKey)
	if err != nil {
		return nil, fmt.Errorf("sign error: %v", err)
	}
	sig[64] += 27
	return sig, nil
}

func (ecsigner *ECSigner) String() string {
	return ecsigner.SignerAddress.String()
}

func AddECSigner() (err error) {
	it := &ui.Item{Label: "Input new ECDSA private key in HEX"}
	err = ui.InputBytes(it, -1)
	if err != nil {
		return
	}
	bt, ok := it.Value.([]byte)
	if !ok || len(bt) == 0 {
		return fmt.Errorf("invalid input")
	}
	//fmt.Printf("Got key %x\n", bt)
	privkey, err := crypto.ToECDSA(bt)
	if err != nil {
		return fmt.Errorf("invalid key: %v", err)
	}
	ecsigner := &ECSigner{SignerKey: privkey}
	ecsigner.SignerAddress = common.BytesToAddress(crypto.PubkeyToAddress(privkey.PublicKey).Bytes())
	state.State.Signers = append(state.State.Signers, ecsigner)
	fmt.Println("Added signer", ecsigner.String())
	state.State.Save()
	return

}

func (ecsigner *ECSigner) Marshal() ([]byte, error) {
	return []byte(Type + ";" + hex.EncodeToString(ecsigner.SignerKey.D.Bytes())), nil
}

func Unmarshal(bt []byte) (signer.Signer, error) {
	hexkey := string(bt)
	privkey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %v", err)
	}
	ecsigner := new(ECSigner)
	ecsigner.SignerKey = privkey
	ecsigner.SignerAddress = common.BytesToAddress(crypto.PubkeyToAddress(privkey.PublicKey).Bytes())
	return ecsigner, nil
}

// KeyContainer function
func (ecsigner *ECSigner) GetKey() *ecdsa.PrivateKey {
	return ecsigner.SignerKey
}
