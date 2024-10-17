package ecsigner

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"

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
	name          string
	SignerAddress common.Address
	SignerKey     *ecdsa.PrivateKey
}

const Type = "ECDSA plain"

func (ecsigner *ECSigner) Type() string {
	return Type
}

func (ecsigner *ECSigner) Name() string {
	return ecsigner.name
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
	nit := &ui.Item{Label: "Signer name"}
	err = ui.InputNewStringUI(nit)
	if err != nil {
		return
	}
	name, ok := nit.Value.(string)
	if !ok {
		return fmt.Errorf("invalid input for signer name")
	}

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
	ecsigner.name = name
	ecsigner.SignerAddress = common.BytesToAddress(crypto.PubkeyToAddress(privkey.PublicKey).Bytes())
	//state.State.Signers = append(state.State.Signers, ecsigner)
	state.AddSigner(ecsigner)
	fmt.Println("Added signer", ecsigner.Name(), ecsigner.String())
	return

}

func (ecsigner *ECSigner) Marshal() ([]byte, error) {
	return []byte(Type + ";" + ecsigner.name + ":" + hex.EncodeToString(ecsigner.SignerKey.D.Bytes())), nil
}

var nameCounter = 0

func Unmarshal(bt []byte) (signer.Signer, error) {
	nameAndHexkey := string(bt)
	var name, hexkey string
	terms := strings.Split(nameAndHexkey, ":")
	if len(terms) == 1 {
		name = "unnamed" + fmt.Sprint(nameCounter)
		nameCounter++
		hexkey = terms[0]
	} else {
		name = terms[0]
		hexkey = terms[1]
	}
	privkey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %v", err)
	}
	ecsigner := new(ECSigner)
	ecsigner.SignerKey = privkey
	ecsigner.name = name
	ecsigner.SignerAddress = common.BytesToAddress(crypto.PubkeyToAddress(privkey.PublicKey).Bytes())
	return ecsigner, nil
}

// KeyContainer function
func (ecsigner *ECSigner) GetKey() *ecdsa.PrivateKey {
	return ecsigner.SignerKey
}

func (ecsigner ECSigner) SignHash(inhash []byte) ([]byte, error) {
	return crypto.Sign(inhash, ecsigner.SignerKey)
}
