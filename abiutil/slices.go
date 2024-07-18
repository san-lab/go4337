package abiutil

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func MakeSliceOfType(solType abi.Type, length, cap int) interface{} {
	var elementType reflect.Type
	elementType = solType.GetType()
	sliceType := reflect.SliceOf(elementType)
	sliceValue := reflect.MakeSlice(sliceType, length, cap)
	return sliceValue.Interface()
}

func AppendToSlice(slice interface{}, value interface{}) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	valueValue := reflect.ValueOf(value)

	// Check if the provided slice is indeed a slice
	if sliceValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("provided interface is not a slice")
	}

	// Check if the value is assignable to the slice element type
	elemType := sliceValue.Type().Elem()
	if !valueValue.Type().AssignableTo(elemType) {
		return nil, fmt.Errorf("cannot assign value of type %s to slice of type %s", valueValue.Type(), elemType)
	}

	// Append the value to the slice
	sliceValue = reflect.Append(sliceValue, valueValue)
	return sliceValue.Interface(), nil
}

func SetSliceValue(slice interface{}, index int, value interface{}) error {
	sliceValue := reflect.ValueOf(slice)
	valueValue := reflect.ValueOf(value)

	// Check if the provided slice is indeed a slice
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("provided interface is not a slice")
	}

	// Check if the index is within bounds
	if index < 0 || index >= sliceValue.Len() {
		return fmt.Errorf("index out of bounds")
	}

	// Check if the value is assignable to the slice element type
	elemType := sliceValue.Type().Elem()
	if !valueValue.Type().AssignableTo(elemType) {
		return fmt.Errorf("cannot assign value of type %s to slice of type %s", valueValue.Type(), elemType)
	}

	// Set the value at the specified index
	sliceValue.Index(index).Set(valueValue)
	return nil
}

const TestABI = `{
	"compiler": {
		"version": "0.8.26+commit.8a97fa7a"
	},
	"language": "Solidity",
	"output": {
		"abi": [
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"name": "IntReceived",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"name": "IntsAccepted",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "string",
						"name": "",
						"type": "string"
					}
				],
				"name": "StringReceived",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "uint64",
						"name": "",
						"type": "uint64"
					},
					{
						"indexed": false,
						"internalType": "uint64",
						"name": "",
						"type": "uint64"
					}
				],
				"name": "TupleSet",
				"type": "event"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "uint64",
								"name": "first",
								"type": "uint64"
							},
							{
								"internalType": "uint64",
								"name": "second",
								"type": "uint64"
							}
						],
						"internalType": "struct TestEncodings.mtuple",
						"name": "tuple",
						"type": "tuple"
					}
				],
				"name": "AcceptTuple",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint16[]",
						"name": "ints",
						"type": "uint16[]"
					}
				],
				"name": "AcceptintSlice",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "bytes",
								"name": "bts",
								"type": "bytes"
							},
							{
								"internalType": "uint32",
								"name": "someint",
								"type": "uint32"
							}
						],
						"internalType": "struct TestEncodings.btuple",
						"name": "mybtuple",
						"type": "tuple"
					}
				],
				"name": "BytesAndIntTuple",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "uint64",
								"name": "first",
								"type": "uint64"
							},
							{
								"internalType": "uint64",
								"name": "second",
								"type": "uint64"
							}
						],
						"internalType": "struct TestEncodings.mtuple[]",
						"name": "tuples",
						"type": "tuple[]"
					}
				],
				"name": "SliceOfTuples",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "string",
						"name": "mesg",
						"type": "string"
					},
					{
						"internalType": "uint256",
						"name": "myint",
						"type": "uint256"
					}
				],
				"name": "StringAndInt",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			}
		],
		"devdoc": {
			"kind": "dev",
			"methods": {},
			"version": 1
		},
		"userdoc": {
			"kind": "user",
			"methods": {},
			"version": 1
		}
	},
	"settings": {
		"compilationTarget": {
			"contracts/contracts/samples/TestEncodings.sol": "TestEncodings"
		},
		"evmVersion": "cancun",
		"libraries": {},
		"metadata": {
			"bytecodeHash": "ipfs"
		},
		"optimizer": {
			"enabled": true,
			"runs": 200
		},
		"remappings": []
	},
	"sources": {
		"contracts/contracts/samples/TestEncodings.sol": {
			"keccak256": "0x6593cc7c3ec8adda7844579fd58f2eb936261a8e9ca8f0ed0d1b622760081796",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://2d9088f61a34223b313ebaebece402923cbe90f45a49d08e534842fda78e047c",
				"dweb:/ipfs/QmTt9iKxiipNmBdEMwSCBMWxAdDgWsxMAvwGfcSYPotbxF"
			]
		}
	},
	"version": 1
}`
