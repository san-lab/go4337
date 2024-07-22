package ui

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

var AddAbiItem = &Item{Label: "Add ABI", Details: "Add a new ABI"}

// If callback is not nil, it will be set and the function will return
func AbisUI(callbackItem *Item) (contract string, err error) {
	for {
		abiitems := []*Item{}
		for _, contract := range state.ListABIs() {
			//We can assume the abi is in the cache
			abiitems = append(abiitems, &Item{Label: contract[0], Details: "Select this ABI", Value: contract[1]})
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

var InputAbiItem = &Item{Label: "Input ABI", Details: "Input ABI manually"}
var ReadAbiFileItem = &Item{Label: "Read ABI from file", Details: "Read ABI from a file"}

func AddAbiUI() {
	items := []*Item{InputAbiItem, ReadAbiFileItem, Back}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Add ABI",
		Items:     items,
		Templates: ItemTemplate,
		Size:      10,
	}

	_, sel, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch sel {
	case Back.Label:
		return
	case InputAbiItem.Label:
		AddAbiManuallyUI()
	case ReadAbiFileItem.Label:
		InputItem := &Item{Label: "Contract Name", Details: "Input the contract name"}
		err := InputNewStringUI(InputItem)
		if err != nil {
			fmt.Println(err)
			return
		}
		contract := InputItem.Value.(string)
		bts, err := SelectFileFromFS("")
		if err != nil {
			fmt.Println(err)
			return
		}
		abistr := string(bts)
		_, _, err = state.ParseABI(contract, abistr)
		if err != nil {
			fmt.Printf("Error %e when parsing abi:\n>>%s<<\n", err, abistr)
		}
	default:
		fmt.Println("Not implemented yet:", sel)
	}
}

func AddAbiManuallyUI() {
	InputItem := &Item{Label: "Contract Name", Details: "Input the contract name"}
	err := InputNewStringUI(InputItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	contract := InputItem.Value.(string)
	abistr, err := MultiLineInput(fmt.Sprintf("Input ABI for %s, add an empty line to finish", contract))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _, err = state.ParseABI(contract, abistr)
	if err != nil {
		fmt.Printf("Error %e when parsing abi:\n>>%s<<\n", err, abistr)
	}
}

var EditAbiItem = &Item{Label: "Edit ABI", Details: "Edit this ABI"}
var DeleteAbiItem = &Item{Label: "Delete ABI", Details: "Delete this ABI"}
var EncodeAbiItem = &Item{Label: "Encode ABI", Details: "Encode this ABI"}

var SelectMethodItem = &Item{Label: "Select a Method"}
var MethodItem = &Item{Label: "Method", Details: "Manage method"}

func AbiUI(contract string, callbackItem *Item) (ret bool) { //should return error?
	var label string
	var proceed bool
	_, _, err := state.GetABI(contract)
	if err != nil {
		label = fmt.Sprintf("Invalid ABI: %v", err)
		proceed = false
	} else {
		label = fmt.Sprintf("ABI for %s", contract)
		proceed = true
	}

	if proceed {
		items := []*Item{}
		if MethodItem.Value != nil && (MethodItem.Value.(*state.MethodCall)).ABIName != contract {
			fmt.Println("Clearing method")
			MethodItem.Value = nil
			MethodItem.DisplayValueString = ""
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
				state.RemoveABI(contract)
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
	method := vmethod.(*state.MethodCall)
	for {
		allSet := true
		items := []*Item{}

		for i, input := range method.ABI.Methods[method.MethodName].Inputs {
			allSet = allSet && method.Params[i] != nil
			items = append(items, &Item{Label: input.Name, DisplayValueString: userop.ParamToString(method.Params[i]), Details: input.Type.String(), Value: input})
		}
		if allSet {
			items = append(items, &Item{Label: "Encode", Details: "Encode this method", Value: method})
		}
		items = append(items, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Set details for " + MethodItem.DisplayValueString + " method",
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
			data, err := userop.EncodeWithParams(method.ABI, method.MethodName, method.Params...)
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
				arg, okk := it.Value.(abi.Argument)
				if !okk {
					fmt.Println("Not an ABI argument")
					continue
				}
				it.Value = method.Params[i]
				err := SetParamUI(it, &arg)
				if err != nil {
					fmt.Println(err)
				} else {
					//t := reflect.TypeOf(v)
					//fmt.Println("Type", t)

					method.Params[i] = it.Value
					//it.DisplayValue = fmt.Sprint(v)
				}
			}
		}
	}
}

func SelectMethodUI(contractName string) {
	items := []*Item{}
	_, methods, err := state.GetABI(contractName)
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
		items = append(items, &Item{Label: k, Details: "select " + v.MethodName, Value: v})
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
				MethodItem.DisplayValueString = it.Label + "()"
				return
			}
		}
	}
}

func Capitalize(s string) string {
	//return s
	return strings.ToUpper(s[:1]) + s[1:]
}

var AppendItem = &Item{Label: "Append", Details: "Append a new value"}

func SetSliceUI(topItem *Item, slice *abi.Argument) error {
	fmt.Println(topItem.Value)
	valueItems := []*Item{}
	v := reflect.ValueOf(topItem.Value)

	// Check if the value is a slice
	// TODO better compare with the type of the slice
	if topItem.Value != nil && v.Kind() == reflect.Slice { //so we have a slice...

		// Iterate over the slice elements
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			valueItems = append(valueItems, &Item{Label: fmt.Sprintf("%s_%v", slice.Name, i), Details: "set " + slice.Type.Elem.String(), Value: element.Interface()})
		}
	}
	loop := true
	for loop {
		items := valueItems
		items = append(items, AppendItem, Set, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     topItem.Label,
			Items:     items,
			Templates: ItemTemplate,
			Size:      10,
		}

		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return err
		}
		switch sel {
		case Back.Label:
			return fmt.Errorf("Back")
		case Set.Label:
			loop = false
		case AppendItem.Label:
			input := &abi.Argument{Type: *slice.Type.Elem}
			newItem := &Item{Label: fmt.Sprintf("%s_%v", slice.Name, len(valueItems)), Details: "set " + slice.Type.Elem.String()}
			err := SetParamUI(newItem, input)
			if err != nil {
				fmt.Println(err)
			}
			valueItems = append(valueItems, newItem)
		default:
			it, ok := GetItem(sel, items)
			if ok {
				err := SetParamUI(it, &abi.Argument{Type: *slice.Type.Elem})
				if err != nil {
					fmt.Println(err)
				}

			}
		}
	}
	//Repackage the values to avoid dependency
	values := abiutil.MakeSliceOfType(*slice.Type.Elem, 0, 0)
	for i := range valueItems {
		var err error
		values, err = abiutil.AppendToSlice(values, valueItems[i].Value)
		if err != nil {
			return fmt.Errorf("Error appending to slice: %v", err)
		}
	}
	topItem.Value = values
	//fmt.Println(topItem.Value)
	return nil

}

func SetTupleUI(item *Item, tuple *abi.Argument) error {
	if tuple.Type.T != abi.TupleTy {
		return fmt.Errorf("Not a tuple")
	}
	valueItems := []*Item{}
	// Expect nil or a tuple
	if item.Value != nil && reflect.TypeOf(item.Value) != tuple.Type.GetType() {
		return fmt.Errorf("Value passed that is not a correct struct")
	}
	// Create a new struct if nil
	if item.Value == nil {
		item.Value = reflect.New(tuple.Type.GetType()).Elem().Interface()
	}

	for i := range tuple.Type.TupleElems {
		name := tuple.Type.GetType().Field(i).Name
		mv := reflect.ValueOf(item.Value).Field(i).Interface()
		nit := &Item{
			Label: name,
			Value: mv,
			//DisplayValue: fmt.Sprint(mv),
			Details: "Set " + tuple.Type.TupleElems[i].String(),
		}
		valueItems = append(valueItems, nit)
	}

	loop := true
	for loop {
		allSet := true
		items := valueItems
		for i := range tuple.Type.TupleElems {

			allSet = allSet && valueItems[i].Value != nil
		}
		if allSet {
			items = append(items, Set)
		}
		items = append(items, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Set tuple values",
			Items:     items,
			Templates: ItemTemplate,
			Size:      10,
		}

		i, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return err
		}
		switch sel {
		case Back.Label:
			return nil
		case "Set":
			loop = false
		default:
			it, ok := GetItem(sel, items)
			if ok {
				err := SetParamUI(it, &abi.Argument{Type: *tuple.Type.TupleElems[i]})
				if err != nil {
					fmt.Println(err)
				}

			}
		}
	}
	//Repackage the values to avoid dependency
	values := []interface{}{}
	for i := range tuple.Type.TupleElems {
		values = append(values, valueItems[i].Value)
	}
	v, err := abiutil.SetTupleValues(tuple, values)
	if err != nil {
		return fmt.Errorf("Error setting tuple values: %v", err)
	}
	item.Value = v
	return nil

}

// Returning the new value as Item.Value. This allows to handle the DisplayValue in a flexible way
func SetParamUI(item *Item, input *abi.Argument) error {

	switch input.Type.T {
	case abi.SliceTy, abi.ArrayTy:
		return SetSliceUI(item, input)
	case abi.TupleTy:
		return SetTupleUI(item, input)

	case abi.AddressTy:
		addr, ok := AddressFromBookUI(item.Label)
		if ok {
			item.Value = *addr
			//item.DisplayValue = addr.String()
			return nil
		}
	case abi.UintTy:
		switch input.Type.Size {
		case 256, 160, 128:
			return InputBigInt(item)
		case 64, 32, 16, 8:
			it := &Item{Label: item.Label}
			err := InputUint(it, input.Type.Size)
			if err != nil {
				return err
			}
			item.Value = it.Value
			//item.DisplayValue = fmt.Sprint(it.Value)
			return nil
		default:
			fmt.Printf("no such uint size: %d\n", input.Type.Size)

		}

	case abi.BytesTy:
		fmt.Println(item.Label)
		bts, err := PotentiallyRecursiveCallDataUI()
		if err != nil {
			return err
		}
		item.Value = bts
		//item.DisplayValue = fmt.Sprintf("0x%x", bts)

	case abi.FixedBytesTy:
		return InputBytes(item, input.Type.Size)
	case abi.BoolTy:
		return InputBool(item)
	case abi.StringTy:
		return InputString(item)

	default:
		fmt.Println("Not implemented yet:", input.Type.String())
	}
	return nil

}