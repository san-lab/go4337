package common

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

type Item struct {
	Label              string
	Value              interface{}
	Details            string
	DisplayValueString string //could be a function, but I am lazy
}

func (i *Item) DisplayValue() string {
	if i.DisplayValueString != "" {
		return ShortString(i.DisplayValueString, 72)
	}
	if i.Value == nil {
		return ""
	}

	if v, ok := (i.Value).(bool); ok {
		return fmt.Sprint(v)
	}

	if v, ok := (i.Value).(uint256.Int); ok {
		return v.String()
	}

	if usop, ok := (i.Value).(*userop.UserOperation); ok {
		return fmt.Sprintf("Sender: %s, Nonce: %d", usop.Sender.String(), usop.Nonce)
	}

	if addr, ok := (i.Value).(*common.Address); ok {
		if addr == nil {
			return ""
		}
		return addr.String()
	}

	if str, ok := (i.Value).(fmt.Stringer); ok {
		return ShortString(str.String(), 50)
	}

	if str, ok := (i.Value).(*fmt.Stringer); ok {
		return ShortString((*str).String(), 50)
	}

	if reflect.ValueOf(i.Value).IsZero() {
		if reflect.TypeOf(i.Value).ConvertibleTo(reflect.TypeOf(0)) {
			return "0"
		}
		return ""
	}
	var derefv interface{}
	//Dereference, if it is a pointer

	if reflect.TypeOf(i.Value).Kind() == reflect.Ptr {
		derefv = reflect.ValueOf(i.Value).Elem().Interface()
	} else {
		derefv = i.Value
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

func (i *Item) String() string {
	return i.Label
}

const unicorn = "\U0001F984"

const (
	EXIT = "EXIT"
	BACK = "BACK"
)

// Rewrite the above as variable instantiation witt OK=&Item{Label: "OK"} pattern
var OK = &Item{Label: "OK", Details: "Confirm and proceed"}
var Back = &Item{Label: BACK, Details: "Go back to previous menu"}
var Exit = &Item{Label: EXIT, Details: "Exit the program"}
var Set = &Item{Label: "Set", Details: "Set the value"}

var ItemTemplate = &promptui.SelectTemplates{
	Label:    "{{ . | bold | cyan}}",
	Inactive: `{{if eq .Label "BACK"}}{{.Label | yellow}}{{else if eq .Label "EXIT"}}{{.Label | red}}{{else if eq .Label "Set"}}{{.Label | green}}{{else}}{{ .Label }}{{with .DisplayValue}}: {{.}}{{end}}{{end}}`,
	Active:   `{{if eq .Label "BACK"}}{{.Label | yellow | bold | underline}}{{else if eq .Label "EXIT"}}{{.Label | red | bold | underline}}{{else if eq .Label "Set"}}{{.Label | green | bold | underline}}{{else}}{{ .Label | bold | underline }}{{with .DisplayValue}}: {{. | bold}}{{end}}{{end}}`,
	Selected: "{{. | faint}}",
	Details:  "{{ .Details | faint }}",
}

func GetValue(label string, items []*Item) (interface{}, bool) {
	it, ok := GetItem(label, items)
	if !ok {
		return nil, false
	}
	return it.Value, true
}

func GetItem(label string, items []*Item) (*Item, bool) {
	for _, i := range items {
		if i.Label == label {
			return i, true
		}
	}
	return nil, false
}
