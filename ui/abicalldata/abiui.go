package abicalldata

import (
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/userop"
)

var AddAbiItem = &common.Item[struct{}]{Label: "Add ABI", Details: "Add a new ABI"}

// If callback is not nil, it will be set and the function will return
func SelectAbiUI(label string) (contract string, ok bool) {

	abiitems := []common.MenuItem{}
	for _, contract := range state.ListABIs() {
		//We can assume the abi is in the cache
		abiitems = append(abiitems, &common.Item[string]{Label: contract[0], Details: "Select this ABI", Value: contract[1]})
	}
	abiitems = append(abiitems, common.Back)
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     label,
		Items:     abiitems,
		Templates: common.ItemTemplate,
		Size:      10,
	}

	_, contract, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	switch contract {
	case common.Back.Label:
		return "", false
	default:
		return contract, true
	}
}

func AbisUI(callbackItem *common.Item[[]byte]) (contract string, err error) {
	for {
		abiitems := []common.MenuItem{}
		for _, contract := range state.ListABIs() {
			//We can assume the abi is in the cache
			abiitems = append(abiitems, &common.Item[string]{Label: contract[0], Details: "Select this ABI", Value: contract[1]})
		}
		abiitems = append(abiitems, AddAbiItem, common.Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Select an ABI",
			Items:     abiitems,
			Templates: common.ItemTemplate,
			Size:      10,
		}

		_, contract, err = prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch contract {
		case common.Back.Label:
			return
		case AddAbiItem.Label:
			AddAbiUI()
		default:
			it, ok := common.GetItem(contract, abiitems)
			if ok {
				AbiUI(it.String(), callbackItem)
				return contract, nil
			}
		}
	}
}

var InputAbiItem = &common.Item[struct{}]{Label: "Input ABI", Details: "Input ABI manually"}
var ReadAbiFileItem = &common.Item[struct{}]{Label: "Read ABI from file", Details: "Read ABI from a file"}

func AddAbiUI() {
	items := []common.MenuItem{InputAbiItem, ReadAbiFileItem, common.Back}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Add ABI",
		Items:     items,
		Templates: common.ItemTemplate,
		Size:      10,
	}

	_, sel, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch sel {
	case common.Back.Label:
		return
	case InputAbiItem.Label:
		AddAbiManuallyUI()
	case ReadAbiFileItem.Label:
		InputItem := &common.Item[string]{Label: "Contract Name", Details: "Input the contract name"}
		err := common.InputNewStringUI(InputItem)
		if err != nil {
			fmt.Println(err)
			return
		}
		contract := InputItem.Value
		bts, err := common.SelectFileFromFS("")
		if err != nil {
			fmt.Println(err)
			return
		}
		abistr := string(bts)
		_, err = state.ParseABI(contract, abistr)
		if err != nil {
			fmt.Printf("Error %e when parsing abi:\n>>%s<<\n", err, common.ShortString(abistr, 100))
		}
	default:
		fmt.Println("Not implemented yet:", sel)
	}
}

func AddAbiManuallyUI() {
	InputItem := &common.Item[string]{Label: "Contract Name", Details: "Input the contract name"}
	err := common.InputNewStringUI(InputItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	contract := InputItem.Value
	abistr, err := common.MultiLineInput(fmt.Sprintf("Input ABI for %s, add an empty line to finish", contract))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = state.ParseABI(contract, abistr)
	if err != nil {
		fmt.Printf("Error %e when parsing abi:\n>>%s<<\n", err, abistr)
	}
}

var EditAbiItem = &common.Item[struct{}]{Label: "Edit ABI", Details: "Edit this ABI"}
var DeleteAbiItem = &common.Item[struct{}]{Label: "Delete ABI", Details: "Delete this ABI"}
var EncodeAbiItem = &common.Item[struct{}]{Label: "Encode ABI", Details: "Encode this ABI"}

var SelectMethodItem = &common.Item[struct{}]{Label: "Select a Method"}
var ConstructorItem = &common.Item[struct{}]{Label: "Constructor", Details: "Encode constructor parameters"}
var MethodItem = &common.Item[*state.MethodCall]{Label: "Method", Details: "Manage method"}
var DeployBytecodeItem = &common.Item[[]byte]{Label: "Deploy Bytecode", Details: "Deploy Bytecode"}
var ExecuteBytecodeItem = &common.Item[[]byte]{Label: "Execute Bytecode", Details: "Execute Bytecode"}
var SelectrorsItem = &common.Item[struct{}]{Label: "Selectors", Details: "Manage selectors"}

func AbiUI(contract string, callbackItem *common.Item[[]byte]) (ret bool) { //should return error?
	var label string
	var proceed bool
	abiart, err := state.GetABI(contract)
	if err != nil {
		label = fmt.Sprintf("Invalid ABI: %v", err)
		proceed = false
	} else {
		label = fmt.Sprintf("ABI for %s", contract)
		proceed = true
		DeployBytecodeItem.Value = abiart.DeployBytecode
		ExecuteBytecodeItem.Value = abiart.ExecuteBytecode

	}

	if proceed {
		items := []common.MenuItem{}
		if MethodItem.Value != nil && MethodItem.Value.ABIName != contract {
			fmt.Println("Clearing method")
			MethodItem.Value = nil
			MethodItem.Display = nil
		}
		items = append(items, EditAbiItem, DeleteAbiItem, SelectMethodItem, MethodItem, ConstructorItem, DeployBytecodeItem,
			ExecuteBytecodeItem, SelectrorsItem, common.Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     label,
			Items:     items,
			Templates: common.ItemTemplate,
			Size:      10,
		}
		for {
			_, sel, err := prompt.Run()
			if err != nil {
				fmt.Println(err)
				return
			}
			switch sel {
			case common.Back.Label:
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
			case ConstructorItem.Label:
				if EncodeConstructorParamsUI(callbackItem, abiart.ABI) {
					return true
				}

			case DeployBytecodeItem.Label:
				err := common.InputBytes(DeployBytecodeItem, -1)
				if err != nil {
					fmt.Println(err)
				} else {
					abiart.DeployBytecode = DeployBytecodeItem.Value
				}
				state.Save()

			case ExecuteBytecodeItem.Label:
				err := common.InputBytes(ExecuteBytecodeItem, -1)
				if err != nil {
					fmt.Println(err)
				} else {
					abiart.ExecuteBytecode = ExecuteBytecodeItem.Value
				}
				state.Save()
			case SelectrorsItem.Label:
				sels := make([][4]byte, len(abiart.ABI.Methods))
				i := 0
				first := true
				fmt.Print("[")
				for _, m := range abiart.ABI.Methods {
					if !first {
						fmt.Print(", ")
					}
					fmt.Printf(`"0x%x"`, m.ID)

					copy(sels[i][:], m.ID[:4])
					i++
					first = false
				}
				fmt.Println("]")
				abiselct, err := abi.NewType("bytes4[]", "", nil)

				if err != nil {
					fmt.Println(err)
					return
				}
				abiselc := abi.Arguments{{Type: abiselct}}
				bt, err := abiselc.Pack(sels)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("Encoded Selectors: 0x%x\n", bt)

			default:
				fmt.Println("Not implemented yet:", sel)
			}
		}
	}
	return
}

func MethodUI(callbackItem *common.Item[[]byte]) (ret bool) {
	vmethod := MethodItem.Value
	if vmethod == nil {
		fmt.Println("No method selected")
		return
	}
	method := vmethod
	// Compute the display label for prompt
	methodDisplayLabel := ""
	if MethodItem.Display != nil {
		methodDisplayLabel = MethodItem.Display(method)
	} else {
		methodDisplayLabel = method.MethodName
	}
	for {
		allSet := true
		items := []common.MenuItem{}

		for i, input := range method.ABI.Methods[method.MethodName].Inputs {
			allSet = allSet && method.Params[i] != nil
			paramDisplay := userop.ParamToString(method.Params[i])
			items = append(items, &common.Item[any]{Label: fmt.Sprintf("%v. %s", i, input.Name), Display: func(v any) string { return paramDisplay }, Details: input.Type.String(), Value: input})
		}
		if allSet {
			items = append(items, &common.Item[*state.MethodCall]{Label: "Encode", Details: "Encode this method", Value: method})
		}
		items = append(items, common.Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Set details for " + methodDisplayLabel + " method",
			Items:     items,
			Templates: common.ItemTemplate,
			Size:      10,
		}

		i, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case common.Back.Label:
			return
		case "Encode":
			data, err := userop.EncodeWithParams(method.ABI, method.MethodName, method.Params...)
			if err != nil {
				fmt.Println("Errrrorrr!!!!!", err)
				continue
			}
			CurrentReturnType = method.ABI.Methods[method.MethodName].Outputs
			if callbackItem != nil {
				callbackItem.Value = data
				return true //a bit controversial hack
			}
			fmt.Printf("Encoded call: 0x%x\n", data)
		default:
			it, ok := common.GetItem(sel, items)
			if ok {
				// We need to get the underlying *Item[any] to get the ABI argument
				typedIt, ok2 := it.(*common.Item[any])
				if !ok2 {
					fmt.Println("Not an ABI argument item")
					continue
				}
				arg, okk := typedIt.Value.(abi.Argument)
				if !okk {
					fmt.Println("Not an ABI argument")
					continue
				}
				// Use a temporary any-typed item for SetParamUI
				paramItem := &common.Item[any]{Label: typedIt.Label, Value: method.Params[i]}
				err := SetParamUI(paramItem, &arg)
				if err != nil {
					fmt.Println(err)
				} else {
					method.Params[i] = paramItem.Value
				}
			}
		}
	}
}

var CurrentReturnType abi.Arguments

func SelectMethodUI(contractName string) {
	items := []common.MenuItem{}
	abiart, err := state.GetABI(contractName)
	if err != nil {
		fmt.Println(err)
		return
	}
	methods := abiart.MethodCalls
	keys := make([]string, 0, len(methods))
	for k := range methods {
		keys = append(keys, k)
	}
	//Sort the keys
	sort.Strings(keys)
	for _, k := range keys {
		v := methods[k]
		dInputs := DisplayInputsTypes(v.ABI.Methods[v.MethodName].Inputs)
		items = append(items, &common.Item[*state.MethodCall]{Label: k, Details: "select " + v.MethodName, Value: v,
			Display: func(mc *state.MethodCall) string { return dInputs }})
	}
	items = append(items, common.Back)
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select a method",
		Items:     items,
		Templates: common.ItemTemplate,
		Size:      10,
	}
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case common.Back.Label:
			return
		default:
			it, ok := common.GetItem(sel, items)
			if ok {
				typedIt, ok2 := it.(*common.Item[*state.MethodCall])
				if !ok2 {
					continue
				}
				MethodItem.Value = typedIt.Value
				mc := typedIt.Value
				m := abiart.ABI.Methods[mc.MethodName]
				MethodItem.Display = func(v *state.MethodCall) string { return v.MethodName + "()" }
				MethodItem.Details = fmt.Sprintf("0x%x, %s", m.ID, m.Sig)
				return
			}
		}
	}
}

func Capitalize(s string) string {
	//return s
	return strings.ToUpper(s[:1]) + s[1:]
}

var AppendItem = &common.Item[struct{}]{Label: "Append", Details: "Append a new value"}
var SelectorsFromABIItem = &common.Item[struct{}]{Label: "Selectors from ABI", Details: "Get all method selectors from the ABI"}

func SetSliceUI(topItem *common.Item[any], slice *abi.Argument) error {
	fmt.Println(topItem.Value)
	valueItems := []*common.Item[any]{}
	v := reflect.ValueOf(topItem.Value)

	// Check if the value is a slice
	// TODO better compare with the type of the slice
	if topItem.Value != nil && v.Kind() == reflect.Slice { //so we have a slice...

		// Iterate over the slice elements
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			valueItems = append(valueItems, &common.Item[any]{Label: fmt.Sprintf("%s_%v", slice.Name, i), Details: "set " + slice.Type.Elem.String(), Value: element.Interface()})
		}
	}
	//detect it this is a slice of bytes4
	bytes4 := false
	if slice.Type.Elem.T == abi.FixedBytesTy && slice.Type.Elem.Size == 4 {
		bytes4 = true
	}

	loop := true
	for loop {
		var items []common.MenuItem
		if bytes4 {
			items = append(items, SelectorsFromABIItem)
		}
		for _, vi := range valueItems {
			items = append(items, vi)
		}
		items = append(items, AppendItem, common.Set, common.Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     topItem.Label,
			Items:     items,
			Templates: common.ItemTemplate,
			Size:      10,
		}

		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return err
		}
		switch sel {
		case common.Back.Label:
			return fmt.Errorf("Back")
		case common.Set.Label:
			loop = false
		case SelectorsFromABIItem.Label:
			if SelectorsFromAbiUI(topItem) {
				return nil
			}

		case AppendItem.Label:
			input := &abi.Argument{Type: *slice.Type.Elem}
			newItem := &common.Item[any]{Label: fmt.Sprintf("%s_%v", slice.Name, len(valueItems)), Details: "set " + slice.Type.Elem.String()}
			err := SetParamUI(newItem, input)
			if err != nil {
				fmt.Println(err)
			}
			valueItems = append(valueItems, newItem)
		default:
			it, ok := common.GetItem(sel, items)
			if ok {
				typedIt, ok2 := it.(*common.Item[any])
				if ok2 {
					err := SetParamUI(typedIt, &abi.Argument{Type: *slice.Type.Elem})
					if err != nil {
						fmt.Println(err)
					}
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

func SetTupleUI(item *common.Item[any], tuple *abi.Argument) error {
	if tuple.Type.T != abi.TupleTy {
		return fmt.Errorf("Not a tuple")
	}
	valueItems := []*common.Item[any]{}

	//---------------Not a known tuple, so a generic approach----------------

	// Expect nil or a tuple
	if item.Value != nil && reflect.TypeOf(item.Value) != tuple.Type.GetType() {
		fmt.Println("Value passed that is not a correct struct")
		fmt.Println(item.Value)
		item.Value = nil
	}
	// Create a new struct if nil
	if item.Value == nil {
		item.Value = reflect.New(tuple.Type.GetType()).Elem().Interface()
	}

	for i := range tuple.Type.TupleElems {
		name := tuple.Type.GetType().Field(i).Name
		mv := reflect.ValueOf(item.Value).Field(i).Interface()
		nit := &common.Item[any]{
			Label:   name,
			Value:   mv,
			Details: "Set " + tuple.Type.TupleElems[i].String(),
		}
		valueItems = append(valueItems, nit)
	}

	loop := true
	for loop {
		allSet := true
		items := []common.MenuItem{}
		for i := range tuple.Type.TupleElems {
			items = append(items, valueItems[i])
			allSet = allSet && valueItems[i].Value != nil
		}
		if allSet {
			items = append(items, common.Set)
		}
		items = append(items, common.Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Set tuple values",
			Items:     items,
			Templates: common.ItemTemplate,
			Size:      10,
		}

		i, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return err
		}
		switch sel {
		case common.Back.Label:
			return nil
		case "Set":
			loop = false
		default:
			it, ok := common.GetItem(sel, items)
			if ok {
				typedIt, ok2 := it.(*common.Item[any])
				if ok2 {
					err := SetParamUI(typedIt, &abi.Argument{Type: *tuple.Type.TupleElems[i]})
					if err != nil {
						fmt.Println(err)
					}
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

// Returning the new value as common.Item.Value. This allows to handle the DisplayValue in a flexible way
func SetParamUI(item *common.Item[any], input *abi.Argument) error {

	switch input.Type.T {
	case abi.SliceTy, abi.ArrayTy:
		return SetSliceUI(item, input)
	case abi.TupleTy:
		return SetTupleUI(item, input)

	case abi.AddressTy:
		_, addr, ok := common.AddressFromBookUI(item.Label)
		if ok {
			item.Value = *addr
			return nil
		}
	case abi.UintTy:
		switch input.Type.Size {
		case 256, 160, 128:
			bigItem := &common.Item[*big.Int]{Label: item.Label}
			_, err := common.InputBigInt(bigItem)
			if err != nil {
				return err
			}
			item.Value = bigItem.Value
			return nil
		case 64, 32, 16, 8:
			u64Item := &common.Item[uint64]{Label: item.Label}
			_, err := common.InputUint(u64Item, input.Type.Size)
			if err != nil {
				return err
			}
			item.Value = u64Item.Value
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

	case abi.FixedBytesTy:
		// For fixed-size bytes, we use a []byte item and store as []byte
		bytesItem := &common.Item[[]byte]{Label: item.Label}
		err := common.InputBytes(bytesItem, input.Type.Size)
		if err != nil {
			return err
		}
		item.Value = bytesItem.Value
	case abi.BoolTy:
		boolItem := &common.Item[bool]{Label: item.Label}
		err := common.InputBool(boolItem)
		if err != nil {
			return err
		}
		item.Value = boolItem.Value
	case abi.StringTy:
		strItem := &common.Item[string]{Label: item.Label}
		err := common.InputNewStringUI(strItem)
		if err != nil {
			return err
		}
		item.Value = strItem.Value

	default:
		fmt.Println("Not implemented yet:", input.Type.String())
	}
	return nil

}

func DisplayInputsTypes(inputs abi.Arguments) string {
	types := []string{}
	for _, input := range inputs {
		types = append(types, input.Type.String())
	}
	return "(" + strings.Join(types, ", ") + ")"
}

func EncodeConstructorParamsUI(callBackItem *common.Item[[]byte], mabi *abi.ABI) (ret bool) {
	if mabi == nil {
		fmt.Println("No ABI")
		return
	}
	constr := mabi.Constructor
	params := constr.Inputs
	values := make([]interface{}, len(params))
	encodeItem := &common.Item[struct{}]{Label: "Encode", Details: "Encode these parameters"}
	encodeWithBytecodeItem := &common.Item[struct{}]{Label: "Encode with Bytecode", Details: "Encode with Bytecode"}
	for {
		allSet := true
		items := []common.MenuItem{}

		for i, input := range constr.Inputs {
			allSet = allSet && values[i] != nil
			vStr := fmt.Sprint(values[i])
			items = append(items, &common.Item[any]{Label: input.Name, Details: input.Type.String(), Value: input, Display: func(v any) string { return vStr }})
		}
		if allSet {
			items = append(items, encodeItem)
			if DeployBytecodeItem.Value != nil {

				items = append(items, encodeWithBytecodeItem)
			}

		}
		items = append(items, common.Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Set parameters of the constructor",
			Items:     items,
			Templates: common.ItemTemplate,
			Size:      10,
		}

		i, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case common.Back.Label:
			return
		case encodeItem.Label, encodeWithBytecodeItem.Label:
			data, err := userop.EncodeWithParams(mabi, "", values...)
			if err != nil {
				fmt.Println("Error encoding parameters", err)
				continue
			}
			if sel == encodeWithBytecodeItem.Label {
				btc := DeployBytecodeItem.Value
				data = append(btc, data...)
			}
			if callBackItem != nil {
				callBackItem.Value = data
				ret = true
				return //a bit controversial hack
			}
			fmt.Printf("Encoded call: 0x%x\n", data)

		default:
			it, ok := common.GetItem(sel, items)
			if ok {
				typedIt, ok2 := it.(*common.Item[any])
				if !ok2 {
					fmt.Println("Not an ABI argument item")
					continue
				}
				arg, okk := typedIt.Value.(abi.Argument)
				if !okk {
					fmt.Println("Not an ABI argument")
					continue
				}
				paramItem := &common.Item[any]{Label: typedIt.Label, Value: params[i]}
				err := SetParamUI(paramItem, &arg)
				if err != nil {
					fmt.Println(err)
				} else {
					values[i] = paramItem.Value
				}
			}
		}
	}

}

func SelectorsFromAbiUI(callbackItem *common.Item[any]) bool {

	contract, ok := SelectAbiUI("From which abi do you want to get the Selectors?")
	if !ok {
		return false
	}
	abiart, err := state.GetABI(contract)
	if err != nil {
		fmt.Println(err)
		return false
	}
	callbackItem.Value = abiutil.GetSelectorsAsSlice(abiart.ABI)
	return true
}
