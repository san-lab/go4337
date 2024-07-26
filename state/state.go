package state

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/userop"
)

// Some Types that are automatically abi-generated from entrypoin contract abi's
// they are here, because their instantiation depends on parsing of the abi definitions
var UserOpV6Type reflect.Type
var PackedUserOpV7Type reflect.Type

type StateStruct struct {
	AddressBooks map[string]*AddressBook
	//Default Gas Costs?
	Signers    []signer.Signer `json:"-"`
	SignersRaw []string
	ABIs       map[string]string //ABI strings memorized
	UserOps    map[string]*userop.UserOperation
	ChainID    uint64
}

//type ABIs map[string]string

var State *StateStruct

var stateMux = &sync.Mutex{}

const EntrypointV6 = "EntrypointV6"
const EntrypointV7 = "EntrypointV7"

func init() {
	stateMux.Lock()
	defer stateMux.Unlock()
	State = new(StateStruct)
	err := State.Load()
	if err != nil {
		fmt.Println(err)
	}
	if State.AddressBooks == nil {
		State.AddressBooks = make(map[string]*AddressBook)
		State.AddressBooks[Sender] = &AddressBook{}
		State.AddressBooks[Paymaster] = &AddressBook{}
		State.AddressBooks[CustomEntrypoint] = &AddressBook{}
	}
	if State.ABIs == nil {
		State.ABIs = make(map[string]string)
	}
	if State.UserOps == nil {
		State.UserOps = make(map[string]*userop.UserOperation)
	}

	//Add the Entrypoin abis
	v6, _, err := ParseABI(EntrypointV6, entrypoint.EntryPointV6AbiJson)
	if err != nil {
		fmt.Println(err)
	} else {
		UserOpV6Type = v6.Methods["getUserOpHash"].Inputs[0].Type.GetType()

	}
	v7, _, err := ParseABI(EntrypointV7, entrypoint.EntryPointV7AbiJson)
	if err != nil {
		fmt.Println(err)
	} else {
		PackedUserOpV7Type = v7.Methods["getUserOpHash"].Inputs[0].Type.GetType()
	}

}

type AddressBook []*common.Address

const Sender = "Sender"
const Paymaster = "Paymaster"
const CustomEntrypoint = "Custom Entrypoint"

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

func (st *StateStruct) Save() error {
	st.SignersRaw = []string{}
	for _, s := range st.Signers {
		bt, err := s.Marshal()
		if err != nil {
			fmt.Println(err)
			continue
		}
		st.SignersRaw = append(st.SignersRaw, string(bt))
	}
	//Also save the undecoded ones
	for k, v := range UnmarshalledSignersBuffer {
		for _, raw := range v {
			st.SignersRaw = append(st.SignersRaw, k+";"+string(raw))
		}
	}
	bt, err := json.MarshalIndent(State, "", "  ")
	if err != nil {
		return fmt.Errorf("Error saving: %w", err)

	}
	err = os.WriteFile(StateFile, bt, 0644)
	if err != nil {
		fmt.Errorf("Error saving: %w", err)
	}
	return nil

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
			UnmarshalledSignersBuffer[terms[0]] = append(UnmarshalledSignersBuffer[terms[0]], []byte(terms[1])) //Store for later
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

var UnmarshalledSignersBuffer = map[string][][]byte{}

var SignerTypes = map[string]signer.AddSigner{}
var Unmarshalers = map[string]func([]byte) (signer.Signer, error){}

func Register(signerType string, add signer.AddSigner, unmarshal func([]byte) (signer.Signer, error)) {
	stateMux.Lock()
	defer stateMux.Unlock()
	SignerTypes[signerType] = add
	Unmarshalers[signerType] = unmarshal
	clearBuffer := true
	for _, raw := range UnmarshalledSignersBuffer[signerType] {
		s, err := unmarshal(raw)
		if err != nil {
			fmt.Println("Unmarshal error:", err)
			clearBuffer = false
			continue
		}
		State.Signers = append(State.Signers, s)
	}
	if clearBuffer {
		delete(UnmarshalledSignersBuffer, signerType)
	}

}

var abicache = map[string]*abi.ABI{}                  //Transient cache
var methodcache = map[string]map[string]*MethodCall{} //ContractName.MethodName

type MethodCall struct {
	MethodName string
	ABIName    string
	ABI        *abi.ABI
	Params     []interface{}
}

func ParseABI(contractname, abiString string) (*abi.ABI, map[string]*MethodCall, error) {

	abi, clrabistr, err := abiutil.ParseABIFromString(abiString)
	if err != nil {
		fmt.Println(clrabistr)
		return abi, nil, fmt.Errorf("Error parsing ABI: %w", err)
	}
	methodCalls := map[string]*MethodCall{}
	for l, m := range abi.Methods {
		params := make([]interface{}, len(m.Inputs))
		methodCalls[m.Name] = &MethodCall{MethodName: l, ABIName: m.Name, ABI: abi, Params: params}
	}
	abicache[contractname] = abi
	methodcache[contractname] = methodCalls
	State.ABIs[contractname] = clrabistr
	return abi, methodCalls, nil
}

func RemoveABI(name string) error {
	delete(State.ABIs, name)
	return State.Save()
}

func GetABI(name string) (*abi.ABI, map[string]*MethodCall, error) {
	abistring, ok := State.ABIs[name]
	if !ok {
		return nil, nil, fmt.Errorf("ABI not found")
	}
	cabi, ok := abicache[name]
	if ok {
		return cabi, methodcache[name], nil

	}

	return ParseABI(name, abistring)
}

func ListABIs() [][]string {
	abinames := []string{}
	for k := range State.ABIs {
		abinames = append(abinames, k)
	}
	//sort the names
	sort.Strings(abinames)
	abis := [][]string{}
	for _, k := range abinames {
		abis = append(abis, []string{k, State.ABIs[k]})
	}
	return abis
}
