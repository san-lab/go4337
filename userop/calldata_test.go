package userop

import (
	"math/big"
)

const tupleABI = `[{
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
			}]`

const testAbi = `[
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "_uint",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "_address",
						"type": "address"
					}
				],
				"name": "IntAndaddr",
				"type": "event"
			},
			{
				"inputs": [],
				"name": "CallMe",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "_address",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "_uint",
						"type": "uint256"
					}
				],
				"name": "addressAndInt",
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
				"name": "counter",
				"outputs": [
					{
						"internalType": "uint64",
						"name": "",
						"type": "uint64"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "counter1",
				"outputs": [
					{
						"internalType": "uint16",
						"name": "",
						"type": "uint16"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "counter2",
				"outputs": [
					{
						"internalType": "uint16",
						"name": "",
						"type": "uint16"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "_256",
						"type": "uint256"
					},
					{
						"internalType": "uint64",
						"name": "_64",
						"type": "uint64"
					}
				],
				"name": "int256Andint64",
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
						"internalType": "uint256",
						"name": "_uint",
						"type": "uint256"
					},
					{
						"internalType": "address",
						"name": "_address",
						"type": "address"
					}
				],
				"name": "intAndAddress",
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
				"name": "owner",
				"outputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					}
				],
				"stateMutability": "view",
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
				"name": "setOwner",
				"outputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			}
		]`

const calldataIntInt17and42 = `ce4023f20000000000000000000000000000000000000000000000000000000000000011000000000000000000000000000000000000000000000000000000000000002a`

const testAddress = `0x9F2b8EAA0cb96bc709482eBdcB8f18dFB12D3133`

var testBigInt = big.NewInt(143)

const callDataAddrAndInt = `06cdc3c10000000000000000000000009f2b8eaa0cb96bc709482ebdcb8f18dfb12d3133000000000000000000000000000000000000000000000000000000000000008f`
