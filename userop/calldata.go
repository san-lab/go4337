package userop

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/san-lab/go4337/state"
)

func EncodeWithParams(abis abi.ABI, method string, params ...interface{}) (data []byte, err error) {
	data, err = abis.Pack(method, params...)
	return
}

type Method struct {
	Name    string
	ABIName string
	ABI     abi.ABI
	Params  []interface{}
}

var abicache = map[string]abi.ABI{}
var methodcache = map[string]map[string]*Method{}

func GetABI(name string) (abi.ABI, map[string]*Method, error) {
	cabi, ok := abicache[name]
	if ok {
		return cabi, methodcache[name], nil

	}

	abistring, ok := state.State.ABIs[name]
	if !ok {
		return abi.ABI{}, nil, fmt.Errorf("ABI not found")
	}
	return ParseABI(abistring, name)
}

func ParseABI(as, methodname string) (abi.ABI, map[string]*Method, error) {

	abi, err := abi.JSON(strings.NewReader(as))
	if err != nil {
		return abi, nil, err
	}
	methods := map[string]*Method{}
	for l, m := range abi.Methods {
		params := make([]interface{}, len(m.Inputs))
		methods[m.Name] = &Method{Name: l, ABIName: m.Name, ABI: abi, Params: params}
	}
	abicache[methodname] = abi
	methodcache[methodname] = methods
	return abi, methods, nil
}
