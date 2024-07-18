package abiutil

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func ParseABIFromString(abistr string) (ABI *abi.ABI, clearABIString string, err error) {
	//First try to fit into a full ABI-Metadata json
	abistuct := new(ABIStruct)
	abistr = Sanitize(abistr, -1)
	err = json.Unmarshal([]byte(abistr), abistuct)
	if err == nil {
		abistr = string(abistuct.Output.Abi)
	}

	//else try the whole block
	a, err := abi.JSON(strings.NewReader(abistr))
	return &a, abistr, err

}

func Sanitize(abi string, l int) string {
	// Remove newlines and leading/trailing whitespace
	san := strings.TrimSpace(strings.ReplaceAll(abi, "\n", ""))
	san = strings.ReplaceAll(san, "\t", "")
	if l > 0 && len(san) > l {
		return san[:l] + "..."
	}
	return san
}

var Subjectstring = `[
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
				"name": "TupeSet",
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
						"internalType": "struct Subject.mtuple",
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
						"internalType": "bytes",
						"name": "calldat",
						"type": "bytes"
					},
					{
						"internalType": "address",
						"name": "_object",
						"type": "address"
					}
				],
				"name": "CallObject",
				"outputs": [
					{
						"internalType": "bytes",
						"name": "",
						"type": "bytes"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "GetCallMeBytecode",
				"outputs": [
					{
						"internalType": "bytes",
						"name": "",
						"type": "bytes"
					}
				],
				"stateMutability": "pure",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "_owner",
						"type": "address"
					}
				],
				"name": "GetCallOwner",
				"outputs": [
					{
						"internalType": "bytes",
						"name": "",
						"type": "bytes"
					}
				],
				"stateMutability": "pure",
				"type": "function"
			}
		]`
