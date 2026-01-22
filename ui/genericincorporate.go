package ui

import (
	"encoding/json"
	"fmt"
	"reflect"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/userop"
)

// FieldInfo holds all data needed to build one UI item and apply the value.
// (This struct remains unchanged)
type FieldInfo struct {
	FieldName    string
	CurrentValue interface{}
	NewValue     interface{}
	IsGasLimit   bool
}

// IncorporateRecommendedValuesUI presents a promptui select screen to incorporate
// non-nil values from newValues into the UserOperation.
func IncorporateRecommendedValuesUI(usop *userop.UserOperation, rawNewValues any) {
	newValues := MapResultToGasAndPaymasterValues(rawNewValues)
	b, _ := json.MarshalIndent(newValues, " ", " ")
	fmt.Println("new values:", string(b))
	if state.DEBUG {
		fmt.Println(newValues.Paymaster)
	}
	// 1. Build the list of actionable items using reflection
	fields := getApplicableFields(usop, newValues)
	if len(fields) == 0 {
		fmt.Println("No new recommended values to incorporate.")
		return
	}

	var items []*common.Item
	// Map keyed by the *full unique label* string to look up the FieldInfo
	fieldMap := make(map[string]*FieldInfo)

	IncorporateAllItem := &common.Item{Label: "Incorporate All"}
	items = append(items, IncorporateAllItem)

	for _, fi := range fields {
		// Create the full unique label for the promptui item
		label := fmt.Sprintf("%s (Current: %s | New: %s)", fi.FieldName, common.ToDisplayString(fi.CurrentValue), common.ToDisplayString(fi.NewValue))

		item := &common.Item{
			Label: label,
			Value: false, // UI toggle state
			// We MUST NOT set UserData as it doesn't exist on common.Item
		}
		items = append(items, item)
		//fmt.Printf("Added item %s for the field %s\n", item.Label, fi.FieldName)
		// Map the unique label back to the FieldInfo
		fieldMap[label] = &fi
	}

	items = append(items, common.Set, common.Back) // Add action buttons

	// 2. Main UI Loop
	for {
		spr := promptui.Select{Label: "Incorporate Recommended Values", Items: items, Templates: common.ItemTemplate, Size: 10}
		_, sel, err := spr.Run() // sel is the *label* string

		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}

		switch sel {
		case common.Back.Label:
			return

		case common.Set.Label, IncorporateAllItem.Label:
			all := (sel == IncorporateAllItem.Label)

			// Apply selected or all parameters
			for _, item := range items[1 : len(items)-2] { // Iterate over dynamic field items (skip All, Set, Back)
				//fmt.Println("xyz", item.Label, fieldMap[item.Label])
				// Only process if selected OR if 'Incorporate All' was chosen
				if item.Value.(bool) || all {
					// Use the item's label (which is unique) to look up the field info
					fieldInfo := fieldMap[item.Label]

					if fieldInfo != nil {
						applyValueToUserOp(usop, fieldInfo)
					}
				}
			}

			// state.Save() // Assuming this is still necessary after applying values
			return

		default:
			// Toggle selection for a field
			it, _ := common.GetItem(sel, items)
			if it.Label != IncorporateAllItem.Label && it.Label != common.Set.Label && it.Label != common.Back.Label {
				it.Value = !it.Value.(bool)
			}
		}
	}
}

// NOTE: getApplicableFields and applyValueToUserOp remain unchanged
// as their logic is correct and decoupled from the Item structure.
// They are included below for completeness of the file.

// --- Helper Functions (Unchanged) ---

// getApplicableFields uses reflection to find non-nil fields in newValues and
// pair them with the corresponding current value from the usop.
func getApplicableFields(usop *userop.UserOperation, newValues *GasAndPaymasterValues) []FieldInfo {
	var fields []FieldInfo

	vNew := reflect.ValueOf(newValues).Elem() // Dereference the newValues pointer
	vUsop := reflect.ValueOf(usop).Elem()

	for i := 0; i < vNew.NumField(); i++ {
		newField := vNew.Field(i)
		newFieldType := vNew.Type().Field(i)
		fieldName := newFieldType.Name

		// Skip if the field is not a pointer, and it's zero/empty (this is important for []byte)
		if newField.Kind() != reflect.Ptr && newField.IsZero() {
			continue
		}

		// Check if the new value is provided (i.e., not nil) - only for pointers
		if newField.Kind() == reflect.Ptr && newField.IsNil() {
			continue
		}

		var actualNewValue interface{}

		// FIX: Correctly handle pointers vs. non-pointers (slices)
		if newField.Kind() == reflect.Ptr {
			// This branch handles all *uint64 and *common.Address fields
			actualNewValue = newField.Elem().Interface()
		} else {
			// This branch handles PaymasterData ([]byte) and any other non-pointer fields
			actualNewValue = newField.Interface()
		}

		// Get the current value from the UserOperation
		usopField := vUsop.FieldByName(fieldName)
		if !usopField.IsValid() {
			continue
		}

		// Special handling for PaymasterData and CallData for display
		var usopCurrentValue interface{}
		if fieldName == "PaymasterData" || fieldName == "CallData" {
			// Only show prefix for byte slices
			usopCurrentValue = fmt.Sprintf("0x%x...", usopField.Bytes())
		} else {
			usopCurrentValue = usopField.Interface()
		}

		fi := FieldInfo{
			FieldName:    fieldName,
			CurrentValue: usopCurrentValue,
			NewValue:     actualNewValue,
			IsGasLimit: // ... (Gas limit tag check remains unchanged)
			newFieldType.Tag.Get("json") == "callGasLimit" ||
				newFieldType.Tag.Get("json") == "verificationGasLimit" ||
				newFieldType.Tag.Get("json") == "preVerificationGas" ||
				newFieldType.Tag.Get("json") == "paymasterVerificationGasLimit" ||
				newFieldType.Tag.Get("json") == "paymasterPostOpGasLimit",
		}

		fields = append(fields, fi)
	}

	return fields
}

func applyValueToUserOp(usop *userop.UserOperation, fi *FieldInfo) {
	//fmt.Println("applying:", fi)
	vUsop := reflect.ValueOf(usop).Elem()
	field := vUsop.FieldByName(fi.FieldName)
	fmt.Println(fi.FieldName, field.Type())

	if field.IsValid() && field.CanSet() {
		newValue := reflect.ValueOf(fi.NewValue)

		switch field.Kind() {
		case reflect.Uint64:
			if newValue.Kind() == reflect.Uint64 {
				field.SetUint(newValue.Uint())
			}
		case reflect.Slice:
			if newValue.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Uint8 {
				field.SetBytes(newValue.Bytes())
			}
		case reflect.Struct:
			if field.Type().String() == "common.Address" {
				if s, ok := fi.NewValue.(string); ok {
					addr := ecommon.HexToAddress(s)
					field.Set(reflect.ValueOf(addr))
				}
			} else {
				fmt.Println("!!! Unknown Struct:", field.Type(), "!!!")
			}

		case reflect.Ptr:
			if field.Type().String() == "*common.Address" {
				if s, ok := fi.NewValue.(string); ok {
					addr := ecommon.HexToAddress(s)
					field.Set(reflect.ValueOf(&addr))
					return
				} else if as, ok := fi.NewValue.(ecommon.Address); ok {
					field.Set(reflect.ValueOf(&as))
				} else {
					fmt.Println("not a string:", fi.NewValue)
					fmt.Printf("NewValue dynamic type: %T\n", fi.NewValue)
				}
			}

			if field.Type().ConvertibleTo(newValue.Type()) {
				field.Set(newValue)
			}
		default:
			if newValue.Type().ConvertibleTo(field.Type()) {
				field.Set(newValue.Convert(field.Type()))
			}
		}
	} else {
		fmt.Printf("Warning: Cannot set field %s on UserOperation.\n", fi.FieldName)
	}
}

type GasAndPaymasterValues struct {
	// Gas Limits
	CallGasLimit                  *uint64 `json:"callGasLimit,omitempty"`
	VerificationGasLimit          *uint64 `json:"verificationGasLimit,omitempty"`
	PreVerificationGas            *uint64 `json:"preVerificationGas,omitempty"`
	PaymasterVerificationGasLimit *uint64 `json:"paymasterVerificationGasLimit,omitempty"`
	PaymasterPostOpGasLimit       *uint64 `json:"paymasterPostOpGasLimit,omitempty"`

	// Fees
	MaxFeePerGas         *uint64 `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *uint64 `json:"maxPriorityFeePerGas,omitempty"`

	// Paymaster Data
	Paymaster     *ecommon.Address `json:"paymaster,omitempty"`
	PaymasterData []byte           `json:"paymasterData,omitempty"`
}

// MapResultToGasAndPaymasterValues generically converts any RPC result struct
// into the standardized GasAndPaymasterValues struct using reflection.
func MapResultToGasAndPaymasterValues(result any) *GasAndPaymasterValues {

	gasValues := &GasAndPaymasterValues{}

	vResult := reflect.ValueOf(result)

	if vResult.Kind() == reflect.Ptr {
		vResult = vResult.Elem()
	}
	if vResult.Kind() != reflect.Struct {
		fmt.Printf("Input for mapping is not a struct: %v\n", vResult.Kind())
		return gasValues
	}

	vGasValues := reflect.ValueOf(gasValues).Elem()

	// Iterate over the fields of the *input* result struct
	for i := 0; i < vResult.NumField(); i++ {
		fieldResult := vResult.Field(i)
		fieldTypeResult := vResult.Type().Field(i)
		fieldName := fieldTypeResult.Name

		// Find the corresponding field in GasAndPaymasterValues
		fieldGasValues := vGasValues.FieldByName(fieldName)

		if !fieldGasValues.IsValid() || !fieldGasValues.CanSet() {
			continue // Skip non-existent or unexported fields
		}

		// Check if the input field is nil (e.g., if the input struct also uses pointers)
		if fieldResult.Kind() == reflect.Ptr && fieldResult.IsNil() {
			continue
		}

		// --- Type-Specific Handling ---

		// 1. Handle Paymaster (*common.Address)
		if fieldName == "Paymaster" {
			if fieldResult.Kind() == reflect.String {
				addrStr := fieldResult.String()
				if addrStr == "" || addrStr == (ecommon.Address{}).Hex() {
					continue
				}
				addr := ecommon.HexToAddress(addrStr)
				// Create and set pointer to common.Address
				fieldGasValues.Set(reflect.ValueOf(&addr))
			} else if fieldResult.Type() == reflect.TypeOf((*ecommon.Address)(nil)) {
				// Input is already a *common.Address (e.g., from another Paymaster struct)
				fieldGasValues.Set(fieldResult)
			}
			continue
		}

		// 2. Handle PaymasterData ([]byte)
		if fieldName == "PaymasterData" {
			if fieldResult.Kind() == reflect.Slice && fieldResult.Type().Elem().Kind() == reflect.Uint8 {
				// Input is already []byte
				if len(fieldResult.Bytes()) > 0 {
					fieldGasValues.Set(fieldResult)
				}
			} else if fieldResult.Kind() == reflect.String {
				// Input is a hex string (e.g., "0x...")
				dataStr := fieldResult.String()
				if dataStr == "" {
					continue
				}
				dataBytes := ecommon.FromHex(dataStr)
				if len(dataBytes) > 0 {
					fieldGasValues.Set(reflect.ValueOf(dataBytes))
				}
			}
			continue
		}

		// 3. Handle uint64 fields (Gas Limits & Fees)
		// This is the original logic for numerical values
		var valueToSet uint64

		switch fieldResult.Kind() {
		case reflect.Uint64:
			valueToSet = fieldResult.Uint()

		case reflect.String:
			strVal := fieldResult.String()
			if strVal == "" || strVal == "0x0" {
				valueToSet = 0
			}
			valueToSet = common.ConvHexOrZero(strVal)

		default:
			// Ignore all other types for numerical fields
			continue
		}

		// Set the *uint64 field

		ptrValue := new(uint64)
		*ptrValue = valueToSet
		fieldGasValues.Set(reflect.ValueOf(ptrValue))

	}

	return gasValues
}
