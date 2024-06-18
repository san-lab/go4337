package ui

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

var AddAbiItem = &Item{Label: "Add ABI", Details: "Add a new ABI"}

// If callback is not nil, it will be set and the function will return
func AbisUI(callbackItem *Item) (contract string, err error) {
	for {
		abiitems := []*Item{}
		for contract, abi := range state.State.ABIs {
			abiitems = append(abiitems, &Item{Label: contract, Details: "Select this ABI", Value: abi, DisplayValue: sanitize(abi, 50)})
		}
		abiitems = append(abiitems, AddAbiItem, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Select an ABI",
			Items:     abiitems,
			Templates: ItemTemplate,
			Size:      10,
		}

		_, contract, err = prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch contract {
		case Back.Label:
			return
		case AddAbiItem.Label:
			AddAbiUI()
		default:
			it, ok := GetItem(contract, abiitems)
			if ok {
				AbiUI(it.Label, callbackItem)
				return contract, nil
			}
		}
	}
}

func AddAbiUI() {
	contract, err := InputNewStringUI("Contract Name")
	if err != nil {
		fmt.Println(err)
		return
	}
	abistr, err := MultiLineInput("Input ABI, add an empty line to finish")
	if err != nil {
		fmt.Println(err)
		return
	}
	abistr = sanitize(abistr, -1)
	//First try to fit into a full ABI-Metadata json
	abistuct := new(ABIStruct)
	err = json.Unmarshal([]byte(abistr), abistuct)
	if err == nil {
		abistr = string(abistuct.Output.Abi)
	}

	//else try the whole block
	abi.JSON(strings.NewReader(abistr))
	cABI, err := abi.JSON(strings.NewReader(abistr))
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range cABI.Methods {
		fmt.Println(k, v)
	}
	state.State.ABIs[contract] = sanitize(abistr, -1)
	state.State.Save()
}

var EditAbiItem = &Item{Label: "Edit ABI", Details: "Edit this ABI"}
var DeleteAbiItem = &Item{Label: "Delete ABI", Details: "Delete this ABI"}
var EncodeAbiItem = &Item{Label: "Encode ABI", Details: "Encode this ABI"}

var SelectMethodItem = &Item{Label: "Select a Method"}
var MethodItem = &Item{Label: "Method", Details: "Manage method"}

func AbiUI(contract string, callbackItem *Item) (ret bool) { //should return error?
	var label string
	var proceed bool
	_, _, err := userop.GetABI(contract)
	if err != nil {
		label = fmt.Sprintf("Invalid ABI: %v", err)
		proceed = false
	} else {
		label = fmt.Sprintf("ABI for %s", contract)
		proceed = true
	}

	if proceed {
		items := []*Item{}
		if MethodItem.Value != nil && (MethodItem.Value.(*userop.Method)).ABIName != contract {
			fmt.Println("Clearing method")
			MethodItem.Value = nil
			MethodItem.DisplayValue = ""
		}
		items = append(items, EditAbiItem, DeleteAbiItem, SelectMethodItem, MethodItem, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     label,
			Items:     items,
			Templates: ItemTemplate,
			Size:      10,
		}
		for {
			_, sel, err := prompt.Run()
			if err != nil {
				fmt.Println(err)
				return
			}
			switch sel {
			case Back.Label:
				return
			case EditAbiItem.Label:
				fmt.Println("Edit ABI")
			case DeleteAbiItem.Label:
				delete(state.State.ABIs, contract)
				return
			case SelectMethodItem.Label:
				SelectMethodUI(contract)
			case MethodItem.Label:
				if MethodUI(callbackItem) {
					return true
				}

			default:
				fmt.Println("Not implemented yet:", sel)
			}
		}
	}
	return
}

func MethodUI(callbackItem *Item) (ret bool) {
	vmethod := MethodItem.Value
	if vmethod == nil {
		fmt.Println("No method selected")
		return
	}
	method := vmethod.(*userop.Method)

	for {
		allSet := true
		items := []*Item{}
		for i, input := range method.ABI.Methods[method.Name].Inputs {
			allSet = allSet && method.Params[i] != nil
			items = append(items, &Item{Label: input.Name, DisplayValue: userop.ParamToString(method.Params[i]), Details: input.Type.String(), Value: input})
		}
		if allSet {
			items = append(items, &Item{Label: "Encode", Details: "Encode this method", Value: method})
		}
		items = append(items, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Set details for " + MethodItem.DisplayValue + " method",
			Items:     items,
			Templates: ItemTemplate,
			Size:      10,
		}

		i, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case "Encode":
			data, err := userop.EncodeWithParams(method.ABI, method.Name, method.Params...)
			if err != nil {
				fmt.Println("Errrrorrr!!!!!", err)
				continue
			}
			if callbackItem != nil {
				callbackItem.Value = data
				return true //a bit controversial hack
			}
			fmt.Printf("Dencoded call: 0x%x\n", data)
		default:
			it, ok := GetItem(sel, items)
			if ok {
				v, err := SetParamUI(it.Label, it.Value.(abi.Argument))
				if err != nil {
					fmt.Println(err)
				} else {
					//t := reflect.TypeOf(v)
					//fmt.Println("Type", t)
					it.Value = v
					method.Params[i] = v
					//it.DisplayValue = fmt.Sprint(v)
				}
			}
		}
	}
}

func SelectMethodUI(contract string) {
	items := []*Item{}
	_, methods, err := userop.GetABI(contract)
	if err != nil {
		fmt.Println(err)
		return
	}
	keys := make([]string, 0, len(methods))
	for k := range methods {
		keys = append(keys, k)
	}
	//Sort the keys
	sort.Strings(keys)
	for _, k := range keys {
		v := methods[k]
		items = append(items, &Item{Label: k, Details: "select " + v.Name, Value: v})
	}
	items = append(items, Back)
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select a method",
		Items:     items,
		Templates: ItemTemplate,
		Size:      10,
	}
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case Back.Label:
			return
		default:
			it, ok := GetItem(sel, items)
			if ok {
				MethodItem.Value = it.Value
				MethodItem.DisplayValue = it.Label + "()"
				return
			}
		}
	}
}

func SetParamUI(label string, input abi.Argument) (interface{}, error) {
	switch input.Type.T {
	case abi.AddressTy:
		return InputNewAddressUI(label)
	case abi.UintTy:
		switch input.Type.Size {
		case 256, 160, 128:
			return InputBigInt(label)
		case 64, 32, 16, 8:
			it := &Item{Label: label}
			err := InputUint(it, input.Type.Size)
			if err != nil {
				return nil, err
			}

			return it.Value, nil
		default:
			fmt.Printf("no such uint size: %d\n", input.Type.Size)

		}

	case abi.BytesTy:

		return PotentiallyRecursiveCallDataUI()

	case abi.FixedBytesTy:
		item := &Item{Label: fmt.Sprintf("Bytes%d", input.Type.Size)}
		err := InputBytes(item, input.Type.Size)
		return item.Value, err

	default:
		fmt.Println("Not implemented yet:", input.Type.String())
	}
	return nil, nil

}
