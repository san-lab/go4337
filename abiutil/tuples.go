package abiutil

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func decapitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] >= 97 {
		return s
	}
	return string(s[0]+32) + s[1:]
}

func SetTupleValues(tuple *abi.Argument, values []interface{}) (interface{}, error) {
	if len(values) != len(tuple.Type.TupleElems) {
		return nil, fmt.Errorf("expected %d values, got %d", len(tuple.Type.TupleElems), len(values))
	}
	v := reflect.New(tuple.Type.GetType())
	//v := reflect.ValueOf(input)
	for i := range tuple.Type.TupleElems {
		value := values[i]
		//If value is an *Address, change to to Address
		//TODO this is a hack, should be done in SetParamUI
		if value != nil {
			if reflect.TypeOf(value).String() == "*common.Address" {
				value = reflect.ValueOf(value).Elem().Interface()
			}
			v.Elem().Field(i).Set(reflect.ValueOf(value))
		}

	}

	return v.Elem().Interface(), nil
}

// CloneStruct clones a struct using reflection
func CloneStruct(input interface{}) (interface{}, error) {
	// Get the reflect.Value of the input
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a pointer to a struct")
	}

	// Dereference the pointer to get the struct value
	v = v.Elem()

	// Create a new instance of the struct type
	newStruct := reflect.New(v.Type()).Elem()

	// Copy each field from the original struct to the new instance
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		newStruct.Field(i).Set(fieldValue)
	}

	// Return a pointer to the new struct
	return newStruct.Addr().Interface(), nil
}
