package state

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/userop"
	"github.com/san-lab/go4337/zksyncera"
)

// Some Types that are automatically abi-generated from entrypoin contract abi's
// they are here, because their instantiation depends on parsing of the abi definitions
var UserOpV6Type reflect.Type
var PackedUserOpV7Type reflect.Type

const GasLimitOffsetDefault = 1000000

type StateStruct struct {
	AddressBooks map[string]*AddressBook
	//Default Gas Costs?
	Signers        map[string]signer.Signer `json:"-"`
	SignersRaw     []string
	ABIArts        map[string]*AbiArtifacts //ABI strings memorized
	UserOps        map[string]*userop.UserOperation
	ChainID        uint64
	Dictionaries   map[string]map[string]string
	RPCEndpoints   map[string]*RPCEndpoint
	GasLimitOffset uint64
	ERA712s        map[string]*zksyncera.ZkSyncTxRLP
	SetAuthTxs     map[string]*types.SetCodeTx
	SetAuthAuths   map[string]*types.SetCodeAuthorization
}

//type ABIs map[string]string

var state *StateStruct

var stateMux = &sync.Mutex{}

var DEBUG = false

func Log(a ...interface{}) {
	if DEBUG {
		fmt.Println(a...)
	}
}

const EntrypointV6 = "EntrypointV6"
const EntrypointV7 = "EntrypointV7"

type AbiArtifacts struct {
	ABIString       string
	ABI             *abi.ABI `json:"-"`
	DeployBytecode  []byte
	ExecuteBytecode []byte
	MethodCalls     map[string]*MethodCall `json:"-"`
}

type RPCEndpoint struct {
	Name    string
	URL     string
	ChainId *big.Int
	APIKey  string
}

func init() {
	stateMux.Lock()
	defer stateMux.Unlock()
	state = new(StateStruct)
	err := state.Load()
	if err != nil {
		fmt.Println(err)
	}
	if state.AddressBooks == nil {
		state.AddressBooks = make(map[string]*AddressBook)
		state.AddressBooks[Sender] = &AddressBook{}
		state.AddressBooks[Paymaster] = &AddressBook{}
		state.AddressBooks[CustomEntrypoint] = &AddressBook{}
	}
	if state.ABIArts == nil {
		state.ABIArts = make(map[string]*AbiArtifacts)
	}
	if state.UserOps == nil {
		state.UserOps = make(map[string]*userop.UserOperation)
	}
	if state.RPCEndpoints == nil {
		state.RPCEndpoints = make(map[string]*RPCEndpoint)
	}
	if state.Dictionaries == nil {
		state.Dictionaries = make(map[string]map[string]string)
	}
	if state.ERA712s == nil {
		state.ERA712s = make(map[string]*zksyncera.ZkSyncTxRLP)
	}
	if state.SetAuthTxs == nil {
		state.SetAuthTxs = make(map[string]*types.SetCodeTx)
	}
	if state.SetAuthAuths == nil {
		state.SetAuthAuths = make(map[string]*types.SetCodeAuthorization)
	}

	//Add the Entrypoint abis
	v6arts, err := ParseABI(EntrypointV6, entrypoint.EntryPointV6AbiJson)
	if err != nil {
		fmt.Println(err)
	} else {
		UserOpV6Type = v6arts.ABI.Methods["getUserOpHash"].Inputs[0].Type.GetType()

	}
	v7arts, err := ParseABI(EntrypointV7, entrypoint.EntryPointV7AbiJson)
	if err != nil {
		fmt.Println(err)
	} else {
		PackedUserOpV7Type = v7arts.ABI.Methods["getUserOpHash"].Inputs[0].Type.GetType()
	}

	if state.GasLimitOffset == 0 {
		state.GasLimitOffset = GasLimitOffsetDefault
	}

}

func GetGasLimitOffset() uint64 {
	return state.GasLimitOffset
}

func SetGasLimitOffset(offset uint64) {
	state.GasLimitOffset = offset
	state.Save()
}

type AddressBook map[string]*common.Address

const Sender = "Sender"
const Paymaster = "Paymaster"
const CustomEntrypoint = "Custom Entrypoint"

func GetAddressBooks() []string {
	chapters := []string{}
	for k := range state.AddressBooks {
		chapters = append(chapters, k)
	}
	sort.Strings(chapters)
	return chapters
}

func GetAddressBook(label string) (*AddressBook, bool) {
	ab, ok := state.AddressBooks[label]
	if !ok || ab == nil {
		ab = &AddressBook{}
		state.AddressBooks[label] = ab
	}
	return ab, ok
}

func (ab *AddressBook) Get(name string) (*common.Address, bool) {
	addr, ok := (*ab)[name]
	return addr, ok
}

func (ab *AddressBook) Add(name string, addr *common.Address) {
	(*ab)[name] = addr
	state.Save()
}

func (ab *AddressBook) Remove(addr *common.Address) {
	//Find the address index
	for n, a := range *ab {
		if a.Hex() == addr.Hex() {
			delete(*ab, n)
			break
		}
	}

	state.Save()
}

func (ab *AddressBook) RemoveByName(name string) bool {
	_, ok := (*ab)[name]
	delete(*ab, name)
	state.Save()
	return ok
}

func (ab *AddressBook) Rename(oldname, newname string) bool {
	addr, ok := (*ab)[oldname]
	if !ok {
		return false
	}
	delete(*ab, oldname)
	(*ab)[newname] = addr
	state.Save()
	return true
}

func (ab *AddressBook) Keys() []string {
	keys := []string{}
	for k := range *ab {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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
	bt, err := json.MarshalIndent(state, "", "  ")
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
	err = json.Unmarshal(bt, state)
	if err != nil {
		return err
	}
	for _, raw := range st.SignersRaw {
		terms := strings.SplitN(raw, `;`, 2)
		if len(terms) != 2 {
			fmt.Println("Invalid signer data:", raw)
			continue
		}
		//fmt.Println("Unmarshaling signer:", terms[0])
		su, ok := Unmarshalers[terms[0]]
		if !ok {
			//fmt.Println("Unknown signer type:", terms[0])
			UnmarshalledSignersBuffer[terms[0]] = append(UnmarshalledSignersBuffer[terms[0]], []byte(terms[1])) //Store for later
			continue
		}
		s, err := su([]byte(terms[1]))
		if err != nil {
			fmt.Println("Unmarshal error:", err)
			continue
		}
		state.Signers[s.Name()] = s

	}
	for k, v := range st.ABIArts {
		ParseABI(k, v.ABIString)
	}
	//fmt.Println("Loaded state", State)
	return nil

}

var UnmarshalledSignersBuffer = map[string][][]byte{}

var SignerTypes = map[string]signer.AddSigner{}
var Unmarshalers = map[string]signer.Unmarshal{}

func Register(signerType string, add signer.AddSigner, unmarshal signer.Unmarshal) {
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
		AddSigner(s)
	}
	if clearBuffer {
		delete(UnmarshalledSignersBuffer, signerType)
	}

}

// var abicache = map[string]*abi.ABI{}                  //Transient cache
// var methodcache = map[string]map[string]*MethodCall{} //ContractName.MethodName

type MethodCall struct {
	MethodName string
	ABIName    string
	ABI        *abi.ABI
	Params     []interface{}
}

func ParseABI(contractname, abiString string) (*AbiArtifacts, error) {

	abi, clrabistr, err := abiutil.ParseABIFromString(abiString)
	if err != nil {
		fmt.Println(clrabistr)
		return nil, fmt.Errorf("Error parsing ABI: %w", err)
	}

	aart, ok := state.ABIArts[contractname]
	if !ok {
		aart = &AbiArtifacts{ABIString: clrabistr}

		state.ABIArts[contractname] = aart
	}

	aart.ABI = abi
	if len(aart.MethodCalls) != len(abi.Methods) {
		methodCalls := map[string]*MethodCall{}
		for l, m := range abi.Methods {
			params := make([]interface{}, len(m.Inputs))
			methodCalls[m.Name] = &MethodCall{MethodName: l, ABIName: m.Name, ABI: abi, Params: params}
		}
		aart.MethodCalls = methodCalls
	}

	//abicache[contractname] = abi

	return aart, nil
}

func RemoveABI(name string) error {
	delete(state.ABIArts, name)
	return state.Save()
}

func GetABI(name string) (*AbiArtifacts, error) {
	aart, ok := state.ABIArts[name]
	if !ok {
		return nil, fmt.Errorf("ABI not found")
	}
	return aart, nil
}

func ListABIs() [][]string {
	abinames := []string{}
	for k := range state.ABIArts {
		abinames = append(abinames, k)
	}
	//sort the names
	sort.Strings(abinames)
	abis := [][]string{}
	for _, k := range abinames {
		abis = append(abis, []string{k, state.ABIArts[k].ABIString})
	}
	return abis
}

func Save() error {
	return state.Save()
}

func Load() error {
	return state.Load()
}

func AddSigner(sig signer.Signer) {
	if state.Signers == nil {
		state.Signers = make(map[string]signer.Signer)
	}
	state.Signers[sig.Name()] = sig
	state.Save()
}

func GetSigners() []string {
	keys := []string{}
	for k := range state.Signers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetSigner(name string) signer.Signer {

	return state.Signers[name]
}

func RemoveSigner(name string) {
	delete(state.Signers, name)
	state.Save()
}

func GetDictionary(name string) map[string]string {
	dict, ok := state.Dictionaries[name]
	if !ok {
		dict = make(map[string]string)
		state.Dictionaries[name] = dict
	}
	return dict
}

func GetDictionaries() []string {
	dicts := []string{}
	for k := range state.Dictionaries {
		dicts = append(dicts, k)
	}
	sort.Strings(dicts)
	return dicts
}

const ApiKeysLabel = "ApiKeys"

func ListApiKeys() []string {
	ApiKeys := GetDictionary(ApiKeysLabel)
	keys := []string{}
	for k := range ApiKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetApiKey(name string) string {
	return GetDictionary(ApiKeysLabel)[name]
}

func AddApiKey(name, key string) {
	GetDictionary(ApiKeysLabel)[name] = key
	state.Save()
}

const ApiEndpointsLabel = "ApiEndpoints"

func ListApiEndpoints() []string {
	ApiEndpoints := GetDictionary(ApiEndpointsLabel)
	keys := []string{}
	for k := range ApiEndpoints {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetApiEndpoint(name string) string {
	return GetDictionary(ApiEndpointsLabel)[name]
}

func AddApiEndpoint(name, url string) {
	GetDictionary(ApiEndpointsLabel)[name] = url
	state.Save()
}

func SetChainId(id any) {
	uid, ok := id.(uint64)
	if ok {
		state.ChainID = uid
		state.Save()

	}
}

func GetChainId() uint64 {
	return state.ChainID
}

func GetRPCEndpoints() map[string]*RPCEndpoint {
	return state.RPCEndpoints
}

func GetRPCEndpoint(name string) *RPCEndpoint {
	ep, ok := state.RPCEndpoints[name]
	if !ok {
		return nil
	}
	return ep
}

func RemoveRPCEndpoint(name string) {
	delete(state.RPCEndpoints, name)
	state.Save()
}

func AddRPCEndpoint(name, url string, chainId *big.Int) {
	state.RPCEndpoints[name] = &RPCEndpoint{Name: name, URL: url, ChainId: chainId}
	state.Save()
}

const MethodTemplatesLabel = "MethodTemplates"

func ListMethodTemplates() []string {
	templates := GetDictionary(MethodTemplatesLabel)
	keys := []string{}
	for k := range templates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetMethodTemplate(name string) string {
	return GetDictionary(MethodTemplatesLabel)[name]
}

func AddMethodTemplate(name, template string) {
	GetDictionary(MethodTemplatesLabel)[name] = template
	state.Save()
}

func AddUserOp(name string, uop *userop.UserOperation) {
	state.UserOps[name] = uop
	state.Save()
}

func GetUserOp(name string) (*userop.UserOperation, bool) {
	u, ok := state.UserOps[name]
	return u, ok
}

func GetUserOps() map[string]*userop.UserOperation {
	return state.UserOps
}

func RemoveUserOp(name string) {
	delete(state.UserOps, name)
	state.Save()
}

// Dies silently if the dictionary does not exist
func AddToDictionatyWithIndice(dict, label, value string) {
	d, ok := state.Dictionaries[dict]
	if !ok {
		return
	}
	idx := 1
	for {
		lab := fmt.Sprintf("%s_%0d", label, idx)
		if _, ok := d[lab]; !ok {
			d[lab] = value
			state.Save()
			return
		}
		idx++
	}

}

func AddSetAuthTx(name string, tx *types.SetCodeTx) {
	state.SetAuthTxs[name] = tx
	state.Save()
}

func GetSetAuthTx(name string) (*types.SetCodeTx, bool) {
	tx, ok := state.SetAuthTxs[name]
	return tx, ok
}

func RemoveSetAuthTx(name string) {
	delete(state.SetAuthTxs, name)
	state.Save()
}

func ListSetAuthTxs() []string {
	keys := []string{}
	for k := range state.SetAuthTxs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func AddSetAuthAuth(name string, auth *types.SetCodeAuthorization) {
	state.SetAuthAuths[name] = auth
	state.Save()
}

func GetSetAuthAuth(name string) (*types.SetCodeAuthorization, bool) {
	auth, ok := state.SetAuthAuths[name]
	return auth, ok
}

func RemoveSetAuthAuth(name string) {
	delete(state.SetAuthAuths, name)
	state.Save()
}

func ListSetAuthAuths() []string {
	keys := []string{}
	for k := range state.SetAuthAuths {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func GetZkERA712(name string) (*zksyncera.ZkSyncTxRLP, bool) {
	era, ok := state.ERA712s[name]
	return era, ok
}

func AddZkERA712(name string, era *zksyncera.ZkSyncTxRLP) {
	state.ERA712s[name] = era
	state.Save()
}

func RemoveZkERA712(name string) {
	delete(state.ERA712s, name)
	state.Save()
}

func ListZkERA712s() []string {
	keys := []string{}
	for k := range state.ERA712s {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
