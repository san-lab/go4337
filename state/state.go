package state

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/signer"
)

type StateStruct struct {
	AddressBooks map[string]*AddressBook
	//Default Gas Costs?
	Signers    []signer.Signer `json:"-"`
	SignersRaw []string
	ABIs       map[string]string
}

var State *StateStruct

var stateMux = &sync.Mutex{}

func InitState() {
	stateMux.Lock()
	defer stateMux.Unlock()
	State = new(StateStruct)
	err := State.Load()
	if err != nil {
		fmt.Println(err)
		State.AddressBooks = make(map[string]*AddressBook)

		State.AddressBooks[Sender] = &AddressBook{}
		State.AddressBooks[Paymaster] = &AddressBook{}
		State.Signers = []signer.Signer{}
		State.ABIs = make(map[string]string)
	}
}

type AddressBook []*common.Address

const Sender = "Sender"
const Paymaster = "Paymaster"

func GetAddressBook(label string) (*AddressBook, bool) {
	ab, ok := State.AddressBooks[label]
	if !ok {
		ab = &AddressBook{}
		State.AddressBooks[label] = ab
	}
	return ab, true
}

func (ab *AddressBook) Add(addr *common.Address) {
	*ab = append(*ab, addr)
	State.Save()
}

func (ab *AddressBook) Remove(addr *common.Address) {
	//Find the address index
	i := 0
	for i = 0; i < len(*ab); i++ {
		if (*ab)[i].Hex() == addr.Hex() {
			break
		}
	}
	if i == len(*ab) {
		return
	}
	//Remove the address
	*ab = append((*ab)[:i], (*ab)[i+1:]...)
	State.Save()
}

var StateFile = "state.json"

func (st *StateStruct) Save() {
	st.SignersRaw = []string{}
	for _, s := range st.Signers {
		bt, err := s.Marshal()
		if err != nil {
			fmt.Println(err)
			continue
		}
		st.SignersRaw = append(st.SignersRaw, string(bt))
	}
	bt, err := json.MarshalIndent(State, "", "  ")
	if err != nil {
		fmt.Println("Error saving:", err)
		return
	}
	err = os.WriteFile(StateFile, bt, 0644)
	if err != nil {
		fmt.Println(err)
	}

}

func (st *StateStruct) Load() error {
	bt, err := os.ReadFile(StateFile)
	if err != nil {
		return err

	}
	err = json.Unmarshal(bt, State)
	if err != nil {
		return err
	}
	for _, raw := range st.SignersRaw {
		terms := strings.SplitN(raw, `;`, 2)
		if len(terms) != 2 {
			fmt.Println("Invalid signer data:", raw)
			continue
		}
		fmt.Println("Unmarshaling signer:", terms[0])
		su, ok := Unmarshalers[terms[0]]
		if !ok {
			fmt.Println("Unknown signer type:", terms[0])
			continue
		}
		s, err := su([]byte(terms[1]))
		if err != nil {
			fmt.Println("Unmarshal error:", err)
			continue
		}
		st.Signers = append(st.Signers, s)

	}
	return nil

}

var SignerTypes = map[string]signer.AddSigner{}
var Unmarshalers = map[string]func([]byte) (signer.Signer, error){}

func Register(signerType string, add signer.AddSigner, unmarshal func([]byte) (signer.Signer, error)) {
	stateMux.Lock()
	defer stateMux.Unlock()
	SignerTypes[signerType] = add
	Unmarshalers[signerType] = unmarshal

}
