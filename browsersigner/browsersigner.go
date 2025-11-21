package browsersigner

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/state"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/userop"
)

func Init() {
	state.Register(SType, AddBrowserSigner, Unmarshal)
}

type BrowserSigner struct {
	SignerAddress common.Address
	name          string
}

const SType = "Browser_Plugin"

/*
type Signer interface {
	SignMessage([]byte) ([]byte, error)
	SignEIP712([]byte) ([]byte, error)
	SignHash([]byte) ([]byte, error) //without any decorations
	Name() string
	String() string //Details
	Type() string
	Marshal() ([]byte, error)
	GetKey() any
}*/

func (bs *BrowserSigner) SignMessage([]byte) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (bs *BrowserSigner) SignHash([]byte) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (bs *BrowserSigner) SignEIP712(uop *userop.UserOperation, chainId *big.Int, entrypoint common.Address) ([]byte, error) {
	return SignEIP712Way(uop, chainId, entrypoint)
}

func (bs *BrowserSigner) SignUserop(uop *userop.UserOperation, chainId *big.Int, entrypoint common.Address) ([]byte, error) {
	return SignEIP712Way(uop, chainId, entrypoint)
}

func (bs *BrowserSigner) Name() string {
	return bs.name
}

func (bs *BrowserSigner) Type() string {
	return SType
}

func (bs *BrowserSigner) GetKey() any {
	return nil
}

func (bs *BrowserSigner) String() string {
	return "mm" + bs.SignerAddress.Hex()
}

func (bs *BrowserSigner) Marshal() ([]byte, error) {
	return []byte(bs.Type() + ";" + bs.name + ":" + bs.SignerAddress.Hex()), nil
}

var nameCounter = 0

func Unmarshal(bts []byte) (signer.Signer, error) {
	nameAndAddress := string(bts)
	var name, hexaddr string
	terms := strings.Split(nameAndAddress, ":")
	if len(terms) == 1 {
		name = "unnamed" + fmt.Sprint(nameCounter)
		nameCounter++
		hexaddr = terms[0]
	} else {
		name = terms[0]
		hexaddr = terms[1]
	}

	bs := new(BrowserSigner)

	bs.name = name
	bs.SignerAddress = common.HexToAddress(hexaddr)
	return bs, nil
}

// UI
func AddBrowserSigner() (err error) {
	nit := &Item{Label: "Signer name"}
	err = InputNewStringUI(nit)
	if err != nil {
		return
	}
	name, ok := nit.Value.(string)
	if !ok {
		return fmt.Errorf("invalid input for signer name")
	}

	it := &Item{Label: "Input Browser Signer Address"}
	err = InputBytes(it, -1)
	if err != nil {
		return
	}
	bt, ok := it.Value.([]byte)
	if !ok || len(bt) == 0 {
		return fmt.Errorf("invalid input")
	}

	brs := new(BrowserSigner)
	brs.name = name
	brs.SignerAddress = common.BytesToAddress(bt)
	state.AddSigner(brs)
	fmt.Println("Added browser signer", brs.Name(), brs.String())
	return

}
