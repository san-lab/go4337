package common

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

// MenuItem is the interface for all menu items — used for heterogeneous slices passed to promptui.
type MenuItem interface {
	fmt.Stringer      // String() returns Label — used by promptui and GetItem
	DisplayValue() string
	LabColon() string
}

type Item[T any] struct {
	Label   string
	Value   T
	Details string
	Display func(T) string // replaces DisplayValueString string
}

func (i *Item[T]) String() string { return i.Label }

func (i *Item[T]) LabColon() string {
	return fmt.Sprintf("%s\t", i.Label)
}

func (i *Item[T]) DisplayValue() string {
	if i == nil {
		return ""
	}
	if i.Display != nil {
		return i.Display(i.Value)
	}
	return ToDisplayString(any(i.Value))
}

func ToDisplayString(v any) string {

	if v == nil {
		return ""
	}

	if v, ok := (v).(bool); ok {
		return fmt.Sprint(v)
	}

	if v, ok := (v).(uint256.Int); ok {
		return fmt.Sprintf("%v \t\t(0x%x)", v, v)
	}

	if v, ok := (v).(uint64); ok {
		return fmt.Sprintf("%v \t\t(0x%x)", v, v)
	}

	if v, ok := (v).(*big.Int); ok {
		return fmt.Sprintf("%v \t\t(0x%x)", v, v)
	}

	if usop, ok := (v).(*userop.UserOperation); ok {
		return fmt.Sprintf("Sender: %s, Nonce: %s", usop.Sender.String(), usop.Nonce)
	}

	if addr, ok := (v).(*common.Address); ok {
		if addr == nil {
			return ""
		}
		return addr.String()
	}

	if str, ok := (v).(fmt.Stringer); ok {
		return ShortString(str.String(), 50)
	}

	if str, ok := (v).(*fmt.Stringer); ok {
		return ShortString((*str).String(), 50)
	}

	rv := reflect.ValueOf(v)
	if rv.IsZero() {
		if reflect.TypeOf(v).ConvertibleTo(reflect.TypeOf(0)) {
			return "0"
		}
		return ""
	}
	var derefv interface{}
	//Dereference, if it is a pointer

	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		derefv = reflect.ValueOf(v).Elem().Interface()
	} else {
		derefv = v
	}

	switch derefv.(type) {
	case string:
		return ShortString(derefv.(string), 40)
	case []byte:
		return ShortHex(derefv.([]byte), 40)
	case [32]byte:
		bt32 := derefv.([32]byte)
		return ShortHex(bt32[:], 40)
	case state.MethodCall:
		method := derefv.(state.MethodCall)
		return method.MethodName
	default:
		return ShortString(fmt.Sprint(derefv), 50)
	}
}

func ShortHex(data []byte, l int) string {

	return ShortString(hex.EncodeToString(data), l)
}

func ShortString(data string, l int) string {
	if len(data) < l+3 {
		return data
	}
	return fmt.Sprintf("%s...%s", data[:l/2], data[len(data)-l/2:])
}

const unicorn = "\U0001F984"

const (
	EXIT = "EXIT"
	BACK = "BACK"
)

// Sentinel items — no meaningful value
var OK = &Item[struct{}]{Label: "OK", Details: "Confirm and proceed"}
var Back = &Item[struct{}]{Label: BACK, Details: "Go back to previous menu"}
var Exit = &Item[struct{}]{Label: EXIT, Details: "Exit the program"}
var Set = &Item[struct{}]{Label: "Set", Details: "Set the value"}

var ItemTemplate = &promptui.SelectTemplates{
	Label:    "{{ . | bold | cyan}}",
	Inactive: `{{if eq .Label "BACK"}}{{.Label | yellow}}{{else if eq .Label "EXIT"}}{{.Label | red}}{{else if eq .Label "Set"}}{{.Label | green}}{{else}}{{ .LabColon }}{{with .DisplayValue}} {{.}}{{end}}{{end}}`,
	Active:   `{{if eq .Label "BACK"}}{{.Label | yellow | bold | underline}}{{else if eq .Label "EXIT"}}{{.Label | red | bold | underline}}{{else if eq .Label "Set"}}{{.Label | green | bold | underline}}{{else}}{{ .LabColon | bold | underline }}{{with .DisplayValue}} {{. | bold}}{{end}}{{end}}`,
	Selected: "{{. | faint}}",
	Details:  "{{ .Details | faint }}",
}

func GetValue(label string, items []MenuItem) (any, bool) {
	i, ok := GetItem(label, items)
	if !ok {
		return nil, false
	}
	// Use reflection to get Value field from concrete type
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName("Value")
	if !f.IsValid() {
		return nil, false
	}
	return f.Interface(), true
}

func GetItem(label string, items []MenuItem) (MenuItem, bool) {
	for _, i := range items {
		if i.String() == label {
			return i, true
		}
	}
	return nil, false
}

// GetTypedItem finds and type-asserts in one step — used for toggle bool items etc.
func GetTypedItem[T any](label string, items []MenuItem) (*Item[T], bool) {
	for _, i := range items {
		if i.String() == label {
			if typed, ok := i.(*Item[T]); ok {
				return typed, true
			}
		}
	}
	return nil, false
}
