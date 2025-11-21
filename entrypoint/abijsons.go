package entrypoint

import "github.com/ethereum/go-ethereum/common"

const EntryPointAddressV6 = "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"
const EntryPointAddressV7 = "0x0000000071727De22E5E9d8BAf0edAc6f37da032"
const EntryPointAddressV8 = "0x4337084D9E255Ff0702461CF8895CE9E3b5Ff108"

var E6Address common.Address
var E7Address common.Address
var E8Address common.Address

func init() {
	E6Address = common.HexToAddress(EntryPointAddressV6)
	E7Address = common.HexToAddress(EntryPointAddressV7)
	E8Address = common.HexToAddress(EntryPointAddressV8)
}

// These are generated from the eth-infinitism implementation
const EntryPointV6AbiJson = `{
	"compiler": {
		"version": "0.8.26+commit.8a97fa7a"
	},
	"language": "Solidity",
	"output": {
		"abi": [
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "preOpGas",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "paid",
						"type": "uint256"
					},
					{
						"internalType": "uint48",
						"name": "validAfter",
						"type": "uint48"
					},
					{
						"internalType": "uint48",
						"name": "validUntil",
						"type": "uint48"
					},
					{
						"internalType": "bool",
						"name": "targetSuccess",
						"type": "bool"
					},
					{
						"internalType": "bytes",
						"name": "targetResult",
						"type": "bytes"
					}
				],
				"name": "ExecutionResult",
				"type": "error"
			},
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "opIndex",
						"type": "uint256"
					},
					{
						"internalType": "string",
						"name": "reason",
						"type": "string"
					}
				],
				"name": "FailedOp",
				"type": "error"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					}
				],
				"name": "SenderAddressResult",
				"type": "error"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "aggregator",
						"type": "address"
					}
				],
				"name": "SignatureValidationFailed",
				"type": "error"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "preOpGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "prefund",
								"type": "uint256"
							},
							{
								"internalType": "bool",
								"name": "sigFailed",
								"type": "bool"
							},
							{
								"internalType": "uint48",
								"name": "validAfter",
								"type": "uint48"
							},
							{
								"internalType": "uint48",
								"name": "validUntil",
								"type": "uint48"
							},
							{
								"internalType": "bytes",
								"name": "paymasterContext",
								"type": "bytes"
							}
						],
						"internalType": "struct IEntryPoint.ReturnInfo",
						"name": "returnInfo",
						"type": "tuple"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "stake",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "unstakeDelaySec",
								"type": "uint256"
							}
						],
						"internalType": "struct IStakeManager.StakeInfo",
						"name": "senderInfo",
						"type": "tuple"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "stake",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "unstakeDelaySec",
								"type": "uint256"
							}
						],
						"internalType": "struct IStakeManager.StakeInfo",
						"name": "factoryInfo",
						"type": "tuple"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "stake",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "unstakeDelaySec",
								"type": "uint256"
							}
						],
						"internalType": "struct IStakeManager.StakeInfo",
						"name": "paymasterInfo",
						"type": "tuple"
					}
				],
				"name": "ValidationResult",
				"type": "error"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "preOpGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "prefund",
								"type": "uint256"
							},
							{
								"internalType": "bool",
								"name": "sigFailed",
								"type": "bool"
							},
							{
								"internalType": "uint48",
								"name": "validAfter",
								"type": "uint48"
							},
							{
								"internalType": "uint48",
								"name": "validUntil",
								"type": "uint48"
							},
							{
								"internalType": "bytes",
								"name": "paymasterContext",
								"type": "bytes"
							}
						],
						"internalType": "struct IEntryPoint.ReturnInfo",
						"name": "returnInfo",
						"type": "tuple"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "stake",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "unstakeDelaySec",
								"type": "uint256"
							}
						],
						"internalType": "struct IStakeManager.StakeInfo",
						"name": "senderInfo",
						"type": "tuple"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "stake",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "unstakeDelaySec",
								"type": "uint256"
							}
						],
						"internalType": "struct IStakeManager.StakeInfo",
						"name": "factoryInfo",
						"type": "tuple"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "stake",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "unstakeDelaySec",
								"type": "uint256"
							}
						],
						"internalType": "struct IStakeManager.StakeInfo",
						"name": "paymasterInfo",
						"type": "tuple"
					},
					{
						"components": [
							{
								"internalType": "address",
								"name": "aggregator",
								"type": "address"
							},
							{
								"components": [
									{
										"internalType": "uint256",
										"name": "stake",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "unstakeDelaySec",
										"type": "uint256"
									}
								],
								"internalType": "struct IStakeManager.StakeInfo",
								"name": "stakeInfo",
								"type": "tuple"
							}
						],
						"internalType": "struct IEntryPoint.AggregatorStakeInfo",
						"name": "aggregatorInfo",
						"type": "tuple"
					}
				],
				"name": "ValidationResultWithAggregation",
				"type": "error"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "factory",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "paymaster",
						"type": "address"
					}
				],
				"name": "AccountDeployed",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [],
				"name": "BeforeExecution",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "totalDeposit",
						"type": "uint256"
					}
				],
				"name": "Deposited",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "aggregator",
						"type": "address"
					}
				],
				"name": "SignatureAggregatorChanged",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "totalStaked",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "unstakeDelaySec",
						"type": "uint256"
					}
				],
				"name": "StakeLocked",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "withdrawTime",
						"type": "uint256"
					}
				],
				"name": "StakeUnlocked",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "withdrawAddress",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "StakeWithdrawn",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "paymaster",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "bool",
						"name": "success",
						"type": "bool"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "actualGasCost",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "actualGasUsed",
						"type": "uint256"
					}
				],
				"name": "UserOperationEvent",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "bytes",
						"name": "revertReason",
						"type": "bytes"
					}
				],
				"name": "UserOperationRevertReason",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "withdrawAddress",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "Withdrawn",
				"type": "event"
			},
			{
				"inputs": [],
				"name": "SIG_VALIDATION_FAILED",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "bytes",
						"name": "initCode",
						"type": "bytes"
					},
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"internalType": "bytes",
						"name": "paymasterAndData",
						"type": "bytes"
					}
				],
				"name": "_validateSenderAndPaymaster",
				"outputs": [],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint32",
						"name": "unstakeDelaySec",
						"type": "uint32"
					}
				],
				"name": "addStake",
				"outputs": [],
				"stateMutability": "payable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					}
				],
				"name": "balanceOf",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					}
				],
				"name": "depositTo",
				"outputs": [],
				"stateMutability": "payable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					}
				],
				"name": "deposits",
				"outputs": [
					{
						"internalType": "uint112",
						"name": "deposit",
						"type": "uint112"
					},
					{
						"internalType": "bool",
						"name": "staked",
						"type": "bool"
					},
					{
						"internalType": "uint112",
						"name": "stake",
						"type": "uint112"
					},
					{
						"internalType": "uint32",
						"name": "unstakeDelaySec",
						"type": "uint32"
					},
					{
						"internalType": "uint48",
						"name": "withdrawTime",
						"type": "uint48"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					}
				],
				"name": "getDepositInfo",
				"outputs": [
					{
						"components": [
							{
								"internalType": "uint112",
								"name": "deposit",
								"type": "uint112"
							},
							{
								"internalType": "bool",
								"name": "staked",
								"type": "bool"
							},
							{
								"internalType": "uint112",
								"name": "stake",
								"type": "uint112"
							},
							{
								"internalType": "uint32",
								"name": "unstakeDelaySec",
								"type": "uint32"
							},
							{
								"internalType": "uint48",
								"name": "withdrawTime",
								"type": "uint48"
							}
						],
						"internalType": "struct IStakeManager.DepositInfo",
						"name": "info",
						"type": "tuple"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"internalType": "uint192",
						"name": "key",
						"type": "uint192"
					}
				],
				"name": "getNonce",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "bytes",
						"name": "initCode",
						"type": "bytes"
					}
				],
				"name": "getSenderAddress",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "address",
								"name": "sender",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "nonce",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "initCode",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "callData",
								"type": "bytes"
							},
							{
								"internalType": "uint256",
								"name": "callGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "verificationGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "preVerificationGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxPriorityFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "paymasterAndData",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct UserOperation",
						"name": "userOp",
						"type": "tuple"
					}
				],
				"name": "getUserOpHash",
				"outputs": [
					{
						"internalType": "bytes32",
						"name": "",
						"type": "bytes32"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"components": [
									{
										"internalType": "address",
										"name": "sender",
										"type": "address"
									},
									{
										"internalType": "uint256",
										"name": "nonce",
										"type": "uint256"
									},
									{
										"internalType": "bytes",
										"name": "initCode",
										"type": "bytes"
									},
									{
										"internalType": "bytes",
										"name": "callData",
										"type": "bytes"
									},
									{
										"internalType": "uint256",
										"name": "callGasLimit",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "verificationGasLimit",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "preVerificationGas",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "maxFeePerGas",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "maxPriorityFeePerGas",
										"type": "uint256"
									},
									{
										"internalType": "bytes",
										"name": "paymasterAndData",
										"type": "bytes"
									},
									{
										"internalType": "bytes",
										"name": "signature",
										"type": "bytes"
									}
								],
								"internalType": "struct UserOperation[]",
								"name": "userOps",
								"type": "tuple[]"
							},
							{
								"internalType": "contract IAggregator",
								"name": "aggregator",
								"type": "address"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct IEntryPoint.UserOpsPerAggregator[]",
						"name": "opsPerAggregator",
						"type": "tuple[]"
					},
					{
						"internalType": "address payable",
						"name": "beneficiary",
						"type": "address"
					}
				],
				"name": "handleAggregatedOps",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "address",
								"name": "sender",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "nonce",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "initCode",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "callData",
								"type": "bytes"
							},
							{
								"internalType": "uint256",
								"name": "callGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "verificationGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "preVerificationGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxPriorityFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "paymasterAndData",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct UserOperation[]",
						"name": "ops",
						"type": "tuple[]"
					},
					{
						"internalType": "address payable",
						"name": "beneficiary",
						"type": "address"
					}
				],
				"name": "handleOps",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint192",
						"name": "key",
						"type": "uint192"
					}
				],
				"name": "incrementNonce",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "bytes",
						"name": "callData",
						"type": "bytes"
					},
					{
						"components": [
							{
								"components": [
									{
										"internalType": "address",
										"name": "sender",
										"type": "address"
									},
									{
										"internalType": "uint256",
										"name": "nonce",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "callGasLimit",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "verificationGasLimit",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "preVerificationGas",
										"type": "uint256"
									},
									{
										"internalType": "address",
										"name": "paymaster",
										"type": "address"
									},
									{
										"internalType": "uint256",
										"name": "maxFeePerGas",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "maxPriorityFeePerGas",
										"type": "uint256"
									}
								],
								"internalType": "struct EntryPoint.MemoryUserOp",
								"name": "mUserOp",
								"type": "tuple"
							},
							{
								"internalType": "bytes32",
								"name": "userOpHash",
								"type": "bytes32"
							},
							{
								"internalType": "uint256",
								"name": "prefund",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "contextOffset",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "preOpGas",
								"type": "uint256"
							}
						],
						"internalType": "struct EntryPoint.UserOpInfo",
						"name": "opInfo",
						"type": "tuple"
					},
					{
						"internalType": "bytes",
						"name": "context",
						"type": "bytes"
					}
				],
				"name": "innerHandleOp",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "actualGasCost",
						"type": "uint256"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					},
					{
						"internalType": "uint192",
						"name": "",
						"type": "uint192"
					}
				],
				"name": "nonceSequenceNumber",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "address",
								"name": "sender",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "nonce",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "initCode",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "callData",
								"type": "bytes"
							},
							{
								"internalType": "uint256",
								"name": "callGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "verificationGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "preVerificationGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxPriorityFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "paymasterAndData",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct UserOperation",
						"name": "op",
						"type": "tuple"
					},
					{
						"internalType": "address",
						"name": "target",
						"type": "address"
					},
					{
						"internalType": "bytes",
						"name": "targetCallData",
						"type": "bytes"
					}
				],
				"name": "simulateHandleOp",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "address",
								"name": "sender",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "nonce",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "initCode",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "callData",
								"type": "bytes"
							},
							{
								"internalType": "uint256",
								"name": "callGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "verificationGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "preVerificationGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxPriorityFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "paymasterAndData",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct UserOperation",
						"name": "userOp",
						"type": "tuple"
					}
				],
				"name": "simulateValidation",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "unlockStake",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address payable",
						"name": "withdrawAddress",
						"type": "address"
					}
				],
				"name": "withdrawStake",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address payable",
						"name": "withdrawAddress",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "withdrawAmount",
						"type": "uint256"
					}
				],
				"name": "withdrawTo",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"stateMutability": "payable",
				"type": "receive"
			}
		],
		"devdoc": {
			"errors": {
				"FailedOp(uint256,string)": [
					{
						"params": {
							"opIndex": "- index into the array of ops to the failed one (in simulateValidation, this is always zero)",
							"reason": "- revert reason      The string starts with a unique code \"AAmn\", where \"m\" is \"1\" for factory, \"2\" for account and \"3\" for paymaster issues,      so a failure can be attributed to the correct entity.   Should be caught in off-chain handleOps simulation and not happen on-chain.   Useful for mitigating DoS attempts against batchers or for troubleshooting of factory/account/paymaster reverts."
						}
					}
				],
				"ValidationResult((uint256,uint256,bool,uint48,uint48,bytes),(uint256,uint256),(uint256,uint256),(uint256,uint256))": [
					{
						"params": {
							"factoryInfo": "stake information about the factory (if any)",
							"paymasterInfo": "stake information about the paymaster (if any)",
							"returnInfo": "gas and time-range returned values",
							"senderInfo": "stake information about the sender"
						}
					}
				],
				"ValidationResultWithAggregation((uint256,uint256,bool,uint48,uint48,bytes),(uint256,uint256),(uint256,uint256),(uint256,uint256),(address,(uint256,uint256)))": [
					{
						"params": {
							"aggregatorInfo": "signature aggregation info (if the account requires signature aggregator)      bundler MUST use it to verify the signature, or reject the UserOperation",
							"factoryInfo": "stake information about the factory (if any)",
							"paymasterInfo": "stake information about the paymaster (if any)",
							"returnInfo": "gas and time-range returned values",
							"senderInfo": "stake information about the sender"
						}
					}
				]
			},
			"events": {
				"AccountDeployed(bytes32,address,address,address)": {
					"params": {
						"factory": "the factory used to deploy this account (in the initCode)",
						"paymaster": "the paymaster used by this UserOp",
						"sender": "the account that is deployed",
						"userOpHash": "the userOp that deployed this account. UserOperationEvent will follow."
					}
				},
				"UserOperationRevertReason(bytes32,address,uint256,bytes)": {
					"params": {
						"nonce": "the nonce used in the request",
						"revertReason": "- the return bytes from the (reverted) call to \"callData\".",
						"sender": "the sender of this request",
						"userOpHash": "the request unique identifier."
					}
				}
			},
			"kind": "dev",
			"methods": {
				"addStake(uint32)": {
					"params": {
						"unstakeDelaySec": "the new lock duration before the deposit can be withdrawn."
					}
				},
				"getDepositInfo(address)": {
					"returns": {
						"info": "- full deposit information of given account"
					}
				},
				"getNonce(address,uint192)": {
					"params": {
						"key": "the high 192 bit of the nonce",
						"sender": "the account address"
					},
					"returns": {
						"nonce": "a full nonce to pass for next UserOp with this sender."
					}
				},
				"getSenderAddress(bytes)": {
					"params": {
						"initCode": "the constructor code to be passed into the UserOperation."
					}
				},
				"handleAggregatedOps(((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[],address,bytes)[],address)": {
					"params": {
						"beneficiary": "the address to receive the fees",
						"opsPerAggregator": "the operations to execute, grouped by aggregator (or address(0) for no-aggregator accounts)"
					}
				},
				"handleOps((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[],address)": {
					"params": {
						"beneficiary": "the address to receive the fees",
						"ops": "the operations to execute"
					}
				},
				"simulateHandleOp((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes),address,bytes)": {
					"params": {
						"op": "the UserOperation to simulate",
						"target": "if nonzero, a target address to call after userop simulation. If called, the targetSuccess and targetResult        are set to the return from that call.",
						"targetCallData": "callData to pass to target address"
					}
				},
				"simulateValidation((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes))": {
					"details": "this method always revert. Successful result is ValidationResult error. other errors are failures.The node must also verify it doesn't use banned opcodes, and that it doesn't reference storage outside the account's data.",
					"params": {
						"userOp": "the user operation to validate."
					}
				},
				"withdrawStake(address)": {
					"params": {
						"withdrawAddress": "the address to send withdrawn value."
					}
				},
				"withdrawTo(address,uint256)": {
					"params": {
						"withdrawAddress": "the address to send withdrawn value.",
						"withdrawAmount": "the amount to withdraw."
					}
				}
			},
			"version": 1
		},
		"userdoc": {
			"errors": {
				"ExecutionResult(uint256,uint256,uint48,uint48,bool,bytes)": [
					{
						"notice": "return value of simulateHandleOp"
					}
				],
				"FailedOp(uint256,string)": [
					{
						"notice": "a custom revert error of handleOps, to identify the offending op.  NOTE: if simulateValidation passes successfully, there should be no reason for handleOps to fail on it."
					}
				],
				"SenderAddressResult(address)": [
					{
						"notice": "return value of getSenderAddress"
					}
				],
				"SignatureValidationFailed(address)": [
					{
						"notice": "error case when a signature aggregator fails to verify the aggregated signature it had created."
					}
				],
				"ValidationResult((uint256,uint256,bool,uint48,uint48,bytes),(uint256,uint256),(uint256,uint256),(uint256,uint256))": [
					{
						"notice": "Successful result from simulateValidation."
					}
				],
				"ValidationResultWithAggregation((uint256,uint256,bool,uint48,uint48,bytes),(uint256,uint256),(uint256,uint256),(uint256,uint256),(address,(uint256,uint256)))": [
					{
						"notice": "Successful result from simulateValidation, if the account returns a signature aggregator"
					}
				]
			},
			"events": {
				"AccountDeployed(bytes32,address,address,address)": {
					"notice": "account \"sender\" was deployed."
				},
				"BeforeExecution()": {
					"notice": "an event emitted by handleOps(), before starting the execution loop. any event emitted before this event, is part of the validation."
				},
				"SignatureAggregatorChanged(address)": {
					"notice": "signature aggregator used by the following UserOperationEvents within this bundle."
				},
				"StakeLocked(address,uint256,uint256)": {
					"notice": "Emitted when stake or unstake delay are modified"
				},
				"StakeUnlocked(address,uint256)": {
					"notice": "Emitted once a stake is scheduled for withdrawal"
				},
				"UserOperationRevertReason(bytes32,address,uint256,bytes)": {
					"notice": "An event emitted if the UserOperation \"callData\" reverted with non-zero length"
				}
			},
			"kind": "user",
			"methods": {
				"SIG_VALIDATION_FAILED()": {
					"notice": "for simulation purposes, validateUserOp (and validatePaymasterUserOp) must return this value in case of signature failure, instead of revert."
				},
				"_validateSenderAndPaymaster(bytes,address,bytes)": {
					"notice": "Called only during simulation. This function always reverts to prevent warm/cold storage differentiation in simulation vs execution."
				},
				"addStake(uint32)": {
					"notice": "add to the account's stake - amount and delay any pending unstake is first cancelled."
				},
				"balanceOf(address)": {
					"notice": "return the deposit (for gas payment) of the account"
				},
				"depositTo(address)": {
					"notice": "add to the deposit of the given account"
				},
				"deposits(address)": {
					"notice": "maps paymaster to their deposits and stakes"
				},
				"getNonce(address,uint192)": {
					"notice": "Return the next nonce for this sender. Within a given key, the nonce values are sequenced (starting with zero, and incremented by one on each userop) But UserOp with different keys can come with arbitrary order."
				},
				"getSenderAddress(bytes)": {
					"notice": "Get counterfactual sender address.  Calculate the sender contract address that will be generated by the initCode and salt in the UserOperation. this method always revert, and returns the address in SenderAddressResult error"
				},
				"getUserOpHash((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes))": {
					"notice": "generate a request Id - unique identifier for this request. the request ID is a hash over the content of the userOp (except the signature), the entrypoint and the chainid."
				},
				"handleAggregatedOps(((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[],address,bytes)[],address)": {
					"notice": "Execute a batch of UserOperation with Aggregators"
				},
				"handleOps((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[],address)": {
					"notice": "Execute a batch of UserOperations. no signature aggregator is used. if any account requires an aggregator (that is, it returned an aggregator when performing simulateValidation), then handleAggregatedOps() must be used instead."
				},
				"incrementNonce(uint192)": {
					"notice": "Manually increment the nonce of the sender. This method is exposed just for completeness.. Account does NOT need to call it, neither during validation, nor elsewhere, as the EntryPoint will update the nonce regardless. Possible use-case is call it with various keys to \"initialize\" their nonces to one, so that future UserOperations will not pay extra for the first transaction with a given key."
				},
				"innerHandleOp(bytes,((address,uint256,uint256,uint256,uint256,address,uint256,uint256),bytes32,uint256,uint256,uint256),bytes)": {
					"notice": "inner function to handle a UserOperation. Must be declared \"external\" to open a call context, but it can only be called by handleOps."
				},
				"nonceSequenceNumber(address,uint192)": {
					"notice": "The next valid sequence number for a given nonce key."
				},
				"simulateHandleOp((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes),address,bytes)": {
					"notice": "simulate full execution of a UserOperation (including both validation and target execution) this method will always revert with \"ExecutionResult\". it performs full validation of the UserOperation, but ignores signature error. an optional target address is called after the userop succeeds, and its value is returned (before the entire call is reverted) Note that in order to collect the the success/failure of the target call, it must be executed with trace enabled to track the emitted events."
				},
				"simulateValidation((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes))": {
					"notice": "Simulate a call to account.validateUserOp and paymaster.validatePaymasterUserOp."
				},
				"unlockStake()": {
					"notice": "attempt to unlock the stake. the value can be withdrawn (using withdrawStake) after the unstake delay."
				},
				"withdrawStake(address)": {
					"notice": "withdraw from the (unlocked) stake. must first call unlockStake and wait for the unstakeDelay to pass"
				},
				"withdrawTo(address,uint256)": {
					"notice": "withdraw from the deposit."
				}
			},
			"version": 1
		}
	},
	"settings": {
		"compilationTarget": {
			"contracts/erc4337v6/contracts/core/EntryPoint.sol": "EntryPoint"
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
		"@openzeppelin/contracts/security/ReentrancyGuard.sol": {
			"keccak256": "0xa535a5df777d44e945dd24aa43a11e44b024140fc340ad0dfe42acf4002aade1",
			"license": "MIT",
			"urls": [
				"bzz-raw://41319e7f621f2dc3733511332c4fd032f8e32ad2aa7fd6f665c19741d9941a34",
				"dweb:/ipfs/QmcYR3bd862GD1Bc7jwrU9bGxrhUu5na1oP964bDCu2id1"
			]
		},
		"contracts/erc4337v6/contracts/core/EntryPoint.sol": {
			"keccak256": "0x04f86318b47f052d7308795ffae6ecec0d023d2458b4e17751b89a0e4acfcdc6",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://c9f6e359c8dbe875ad974d3a0fb7b3d62319a6b115c44bac1e4587ae2ad2edaf",
				"dweb:/ipfs/QmTSWTov2rUeYk8cwzrtsd3uVXokCYok4gMiZ1sPs9tycH"
			]
		},
		"contracts/erc4337v6/contracts/core/Helpers.sol": {
			"keccak256": "0x591c87519f7155d1909210276b77925ab2722a99b7b5d5649aecc36ebbdb045a",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://69643e83f68e6a13d5075c7565bfce326673b0bd98c432033c4603ea84835746",
				"dweb:/ipfs/QmSwSzjYyV7qudi5vvsmzHMG2Z4YJZxX51RRXXVCLaNcEU"
			]
		},
		"contracts/erc4337v6/contracts/core/NonceManager.sol": {
			"keccak256": "0xa17a4a6fde70088ab18ffe6df830f3efa31f1cd0e1a7160336c96e3c94984d25",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://b38615df9f80c56282b72888e9ba1eb1a9413fa67a0dbf094deda7af9feb38e7",
				"dweb:/ipfs/QmSzcXetEJRH4UHuUmZiSgX6bFgfqHWfmyuxVnh4NosMk1"
			]
		},
		"contracts/erc4337v6/contracts/core/SenderCreator.sol": {
			"keccak256": "0x44b9449fec82d6cdfb01d52fdd5a72f90099c651316123810cf9633f00b018c2",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://a9c0487390e72638681d175c45bc92425c802fffdca4bd0ae8457782ee284612",
				"dweb:/ipfs/QmVbzuehCUWJWqEHyMWuc6cRVbxfcMdFsmGL9o4Wz7WY2x"
			]
		},
		"contracts/erc4337v6/contracts/core/StakeManager.sol": {
			"keccak256": "0x21aa0956382bd000b1b8c3b1d19ca6ebcd6c9029eebb19c612fb38ee5dd2430a",
			"license": "GPL-3.0-only",
			"urls": [
				"bzz-raw://0a625c8795354d9f429367f9c1d14eb8af7db9c7f2c2a2033e2066ced76a573a",
				"dweb:/ipfs/Qmd1j6UarUg54q1G2HCNCLQz8XGVZR1qxX7eQ6cytHpQPN"
			]
		},
		"contracts/erc4337v6/contracts/interfaces/IAccount.sol": {
			"keccak256": "0x556a0e5980de18e90b115553ed502408155ba35f58642823010d9288047bc418",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://a0f420134b79596db8737173c7b933ae0a33059e107b6327c43aa40d4744a9e4",
				"dweb:/ipfs/QmRo8s1AhXmEMV7uPYnbpYwU19e9Bk6jmYBJTiPx3Fo85W"
			]
		},
		"contracts/erc4337v6/contracts/interfaces/IAggregator.sol": {
			"keccak256": "0x060e9ddb0152250c269ba0640dc5753834ac44cf182a2837d508c0c529cae26a",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://20ed837bc5909c89ff1910246bf245a5dad6840aa939382e1694964eb7dbd37b",
				"dweb:/ipfs/QmTMybRq5yyghPDDs1ZCNAVB9sSJ4WHe6Q9mejuKPTAdNP"
			]
		},
		"contracts/erc4337v6/contracts/interfaces/IEntryPoint.sol": {
			"keccak256": "0x3a90bf308819ed125fa4202f880999caff8a8686633b8ddb79a30ca240d5b8f8",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://d2d21cc92c2fdab2b58d21bc25d4cd0e8c284b922528a186b087b818d54bc6cf",
				"dweb:/ipfs/QmT1qrfuBjsv2rmRCDn8mgPXHp94hARJwzbcDuBLDTbFWd"
			]
		},
		"contracts/erc4337v6/contracts/interfaces/INonceManager.sol": {
			"keccak256": "0x509871e6c63663cdcc3eba19920fe84e991f38b289b1377ac3c3a6d9f22d7e12",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://00fe21b4349b24c50df60e1a705179293982bd9e7a32b78d4bac9620f89e7fe2",
				"dweb:/ipfs/QmSFFYGfUwQbVa6hASjU7YxTvgi2HkfrPr4X5oPHscHg8b"
			]
		},
		"contracts/erc4337v6/contracts/interfaces/IPaymaster.sol": {
			"keccak256": "0x36858ba8685024974f533530420688da3454d29996ebc42e410673a1ed2ec456",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://616cbcf51778b1961b7f20a547bec7efae6d1d565df0f651926241ed8bde9ad8",
				"dweb:/ipfs/QmaVsgffUUmeUJYgStvRr8cNZ1LBbrc3FYNLW4JT1dVLia"
			]
		},
		"contracts/erc4337v6/contracts/interfaces/IStakeManager.sol": {
			"keccak256": "0xd227b02888cd4ac68daebcdfd992ec00f9fff66fa3b3bb16f656cd582fa3480f",
			"license": "GPL-3.0-only",
			"urls": [
				"bzz-raw://b389da4714a138be63704a576a482505eab2855e263b38a93706395d8d42e7c3",
				"dweb:/ipfs/QmeeAZpdHwUXxqP8pxA7GNtoCGBmmH4FaqLLwScVKGxtxZ"
			]
		},
		"contracts/erc4337v6/contracts/interfaces/UserOperation.sol": {
			"keccak256": "0x61374003361059087fdcf17967a7bba052badeaf5c7f0ae689166f8aafd3a45c",
			"license": "GPL-3.0",
			"urls": [
				"bzz-raw://6ff83c59432e733bf6304dda27cd4b0f34401917dd535e2669cc842d2d26568c",
				"dweb:/ipfs/QmPJbHU5TAjHqUTZzAcicEeG2nknmwCN43L4EW9LHbknTN"
			]
		},
		"contracts/erc4337v6/contracts/utils/Exec.sol": {
			"keccak256": "0x5b232117afbc2939f3ffc92745614867e9e1d475a3e1e5443adae13c200174f1",
			"license": "LGPL-3.0-only",
			"urls": [
				"bzz-raw://62e7365379a06ead7b47637945bcaee095d51aab1d3ac00ddec69443e6cbe9fe",
				"dweb:/ipfs/QmctG3aw4U3KMSMeJKoLJ1NJewjMWfppnd1m3kxNTe39Uy"
			]
		}
	},
	"version": 1
}`

const EntryPointV7AbiJson = `[
			{
				"inputs": [
					{
						"internalType": "bool",
						"name": "success",
						"type": "bool"
					},
					{
						"internalType": "bytes",
						"name": "ret",
						"type": "bytes"
					}
				],
				"name": "DelegateAndRevert",
				"type": "error"
			},
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "opIndex",
						"type": "uint256"
					},
					{
						"internalType": "string",
						"name": "reason",
						"type": "string"
					}
				],
				"name": "FailedOp",
				"type": "error"
			},
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "opIndex",
						"type": "uint256"
					},
					{
						"internalType": "string",
						"name": "reason",
						"type": "string"
					},
					{
						"internalType": "bytes",
						"name": "inner",
						"type": "bytes"
					}
				],
				"name": "FailedOpWithRevert",
				"type": "error"
			},
			{
				"inputs": [
					{
						"internalType": "bytes",
						"name": "returnData",
						"type": "bytes"
					}
				],
				"name": "PostOpReverted",
				"type": "error"
			},
			{
				"inputs": [],
				"name": "ReentrancyGuardReentrantCall",
				"type": "error"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					}
				],
				"name": "SenderAddressResult",
				"type": "error"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "aggregator",
						"type": "address"
					}
				],
				"name": "SignatureValidationFailed",
				"type": "error"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "factory",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "paymaster",
						"type": "address"
					}
				],
				"name": "AccountDeployed",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [],
				"name": "BeforeExecution",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "totalDeposit",
						"type": "uint256"
					}
				],
				"name": "Deposited",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "bytes",
						"name": "revertReason",
						"type": "bytes"
					}
				],
				"name": "PostOpRevertReason",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "aggregator",
						"type": "address"
					}
				],
				"name": "SignatureAggregatorChanged",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "totalStaked",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "unstakeDelaySec",
						"type": "uint256"
					}
				],
				"name": "StakeLocked",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "withdrawTime",
						"type": "uint256"
					}
				],
				"name": "StakeUnlocked",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "withdrawAddress",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "StakeWithdrawn",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "paymaster",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "bool",
						"name": "success",
						"type": "bool"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "actualGasCost",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "actualGasUsed",
						"type": "uint256"
					}
				],
				"name": "UserOperationEvent",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					}
				],
				"name": "UserOperationPrefundTooLow",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "bytes",
						"name": "revertReason",
						"type": "bytes"
					}
				],
				"name": "UserOperationRevertReason",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "address",
						"name": "withdrawAddress",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "Withdrawn",
				"type": "event"
			},
			{
				"inputs": [
					{
						"internalType": "uint32",
						"name": "unstakeDelaySec",
						"type": "uint32"
					}
				],
				"name": "addStake",
				"outputs": [],
				"stateMutability": "payable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					}
				],
				"name": "balanceOf",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "target",
						"type": "address"
					},
					{
						"internalType": "bytes",
						"name": "data",
						"type": "bytes"
					}
				],
				"name": "delegateAndRevert",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					}
				],
				"name": "depositTo",
				"outputs": [],
				"stateMutability": "payable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					}
				],
				"name": "deposits",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "deposit",
						"type": "uint256"
					},
					{
						"internalType": "bool",
						"name": "staked",
						"type": "bool"
					},
					{
						"internalType": "uint112",
						"name": "stake",
						"type": "uint112"
					},
					{
						"internalType": "uint32",
						"name": "unstakeDelaySec",
						"type": "uint32"
					},
					{
						"internalType": "uint48",
						"name": "withdrawTime",
						"type": "uint48"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					}
				],
				"name": "getDepositInfo",
				"outputs": [
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "deposit",
								"type": "uint256"
							},
							{
								"internalType": "bool",
								"name": "staked",
								"type": "bool"
							},
							{
								"internalType": "uint112",
								"name": "stake",
								"type": "uint112"
							},
							{
								"internalType": "uint32",
								"name": "unstakeDelaySec",
								"type": "uint32"
							},
							{
								"internalType": "uint48",
								"name": "withdrawTime",
								"type": "uint48"
							}
						],
						"internalType": "struct IStakeManager.DepositInfo",
						"name": "info",
						"type": "tuple"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"internalType": "uint192",
						"name": "key",
						"type": "uint192"
					}
				],
				"name": "getNonce",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "bytes",
						"name": "initCode",
						"type": "bytes"
					}
				],
				"name": "getSenderAddress",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "address",
								"name": "sender",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "nonce",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "initCode",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "callData",
								"type": "bytes"
							},
							{
								"internalType": "bytes32",
								"name": "accountGasLimits",
								"type": "bytes32"
							},
							{
								"internalType": "uint256",
								"name": "preVerificationGas",
								"type": "uint256"
							},
							{
								"internalType": "bytes32",
								"name": "gasFees",
								"type": "bytes32"
							},
							{
								"internalType": "bytes",
								"name": "paymasterAndData",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct PackedUserOperation",
						"name": "userOp",
						"type": "tuple"
					}
				],
				"name": "getUserOpHash",
				"outputs": [
					{
						"internalType": "bytes32",
						"name": "",
						"type": "bytes32"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"components": [
									{
										"internalType": "address",
										"name": "sender",
										"type": "address"
									},
									{
										"internalType": "uint256",
										"name": "nonce",
										"type": "uint256"
									},
									{
										"internalType": "bytes",
										"name": "initCode",
										"type": "bytes"
									},
									{
										"internalType": "bytes",
										"name": "callData",
										"type": "bytes"
									},
									{
										"internalType": "bytes32",
										"name": "accountGasLimits",
										"type": "bytes32"
									},
									{
										"internalType": "uint256",
										"name": "preVerificationGas",
										"type": "uint256"
									},
									{
										"internalType": "bytes32",
										"name": "gasFees",
										"type": "bytes32"
									},
									{
										"internalType": "bytes",
										"name": "paymasterAndData",
										"type": "bytes"
									},
									{
										"internalType": "bytes",
										"name": "signature",
										"type": "bytes"
									}
								],
								"internalType": "struct PackedUserOperation[]",
								"name": "userOps",
								"type": "tuple[]"
							},
							{
								"internalType": "contract IAggregator",
								"name": "aggregator",
								"type": "address"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct IEntryPoint.UserOpsPerAggregator[]",
						"name": "opsPerAggregator",
						"type": "tuple[]"
					},
					{
						"internalType": "address payable",
						"name": "beneficiary",
						"type": "address"
					}
				],
				"name": "handleAggregatedOps",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"components": [
							{
								"internalType": "address",
								"name": "sender",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "nonce",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "initCode",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "callData",
								"type": "bytes"
							},
							{
								"internalType": "bytes32",
								"name": "accountGasLimits",
								"type": "bytes32"
							},
							{
								"internalType": "uint256",
								"name": "preVerificationGas",
								"type": "uint256"
							},
							{
								"internalType": "bytes32",
								"name": "gasFees",
								"type": "bytes32"
							},
							{
								"internalType": "bytes",
								"name": "paymasterAndData",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct PackedUserOperation[]",
						"name": "ops",
						"type": "tuple[]"
					},
					{
						"internalType": "address payable",
						"name": "beneficiary",
						"type": "address"
					}
				],
				"name": "handleOps",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint192",
						"name": "key",
						"type": "uint192"
					}
				],
				"name": "incrementNonce",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "bytes",
						"name": "callData",
						"type": "bytes"
					},
					{
						"components": [
							{
								"components": [
									{
										"internalType": "address",
										"name": "sender",
										"type": "address"
									},
									{
										"internalType": "uint256",
										"name": "nonce",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "verificationGasLimit",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "callGasLimit",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "paymasterVerificationGasLimit",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "paymasterPostOpGasLimit",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "preVerificationGas",
										"type": "uint256"
									},
									{
										"internalType": "address",
										"name": "paymaster",
										"type": "address"
									},
									{
										"internalType": "uint256",
										"name": "maxFeePerGas",
										"type": "uint256"
									},
									{
										"internalType": "uint256",
										"name": "maxPriorityFeePerGas",
										"type": "uint256"
									}
								],
								"internalType": "struct EntryPoint.MemoryUserOp",
								"name": "mUserOp",
								"type": "tuple"
							},
							{
								"internalType": "bytes32",
								"name": "userOpHash",
								"type": "bytes32"
							},
							{
								"internalType": "uint256",
								"name": "prefund",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "contextOffset",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "preOpGas",
								"type": "uint256"
							}
						],
						"internalType": "struct EntryPoint.UserOpInfo",
						"name": "opInfo",
						"type": "tuple"
					},
					{
						"internalType": "bytes",
						"name": "context",
						"type": "bytes"
					}
				],
				"name": "innerHandleOp",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "actualGasCost",
						"type": "uint256"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					},
					{
						"internalType": "uint192",
						"name": "",
						"type": "uint192"
					}
				],
				"name": "nonceSequenceNumber",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "bytes4",
						"name": "interfaceId",
						"type": "bytes4"
					}
				],
				"name": "supportsInterface",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "unlockStake",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address payable",
						"name": "withdrawAddress",
						"type": "address"
					}
				],
				"name": "withdrawStake",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address payable",
						"name": "withdrawAddress",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "withdrawAmount",
						"type": "uint256"
					}
				],
				"name": "withdrawTo",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"stateMutability": "payable",
				"type": "receive"
			}
		]`

const EntryPointV8AbiJson = `[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [
			{
				"internalType": "bool",
				"name": "success",
				"type": "bool"
			},
			{
				"internalType": "bytes",
				"name": "ret",
				"type": "bytes"
			}
		],
		"name": "DelegateAndRevert",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "withdrawAddress",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "revertReason",
				"type": "bytes"
			}
		],
		"name": "DepositWithdrawalFailed",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "Eip7702SenderNotDelegate",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "Eip7702SenderWithoutCode",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "opIndex",
				"type": "uint256"
			},
			{
				"internalType": "string",
				"name": "reason",
				"type": "string"
			}
		],
		"name": "FailedOp",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "opIndex",
				"type": "uint256"
			},
			{
				"internalType": "string",
				"name": "reason",
				"type": "string"
			},
			{
				"internalType": "bytes",
				"name": "inner",
				"type": "bytes"
			}
		],
		"name": "FailedOpWithRevert",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "beneficiary",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "revertData",
				"type": "bytes"
			}
		],
		"name": "FailedSendToBeneficiary",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "currentDeposit",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "withdrawAmount",
				"type": "uint256"
			}
		],
		"name": "InsufficientDeposit",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "InternalFunction",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "beneficiary",
				"type": "address"
			}
		],
		"name": "InvalidBeneficiary",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "paymaster",
				"type": "address"
			}
		],
		"name": "InvalidPaymaster",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "paymasterAndDataLength",
				"type": "uint256"
			}
		],
		"name": "InvalidPaymasterData",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "dataLength",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "pmSignatureLength",
				"type": "uint256"
			}
		],
		"name": "InvalidPaymasterSignatureLength",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "InvalidShortString",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "msgValue",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "currentStake",
				"type": "uint256"
			}
		],
		"name": "InvalidStake",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "newUnstakeDelaySec",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "currentUnstakeDelaySec",
				"type": "uint256"
			}
		],
		"name": "InvalidUnstakeDelay",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "currentStake",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "unstakeDelaySec",
				"type": "uint256"
			},
			{
				"internalType": "bool",
				"name": "staked",
				"type": "bool"
			}
		],
		"name": "NotStaked",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "bytes",
				"name": "returnData",
				"type": "bytes"
			}
		],
		"name": "PostOpReverted",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "Reentrancy",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "SenderAddressResult",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "aggregator",
				"type": "address"
			}
		],
		"name": "SignatureValidationFailed",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "withdrawTime",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "blockTimestamp",
				"type": "uint256"
			}
		],
		"name": "StakeNotUnlocked",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "withdrawAddress",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "revertReason",
				"type": "bytes"
			}
		],
		"name": "StakeWithdrawalFailed",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "str",
				"type": "string"
			}
		],
		"name": "StringTooLong",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "withdrawTime",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "blockTimestamp",
				"type": "uint256"
			}
		],
		"name": "WithdrawalNotDue",
		"type": "error"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "userOpHash",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "factory",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "paymaster",
				"type": "address"
			}
		],
		"name": "AccountDeployed",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [],
		"name": "BeforeExecution",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "totalDeposit",
				"type": "uint256"
			}
		],
		"name": "Deposited",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [],
		"name": "EIP712DomainChanged",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "userOpHash",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "delegate",
				"type": "address"
			}
		],
		"name": "EIP7702AccountInitialized",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "userOpHash",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "unusedFactory",
				"type": "address"
			}
		],
		"name": "IgnoredInitCode",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "userOpHash",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "nonce",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "revertReason",
				"type": "bytes"
			}
		],
		"name": "PostOpRevertReason",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "aggregator",
				"type": "address"
			}
		],
		"name": "SignatureAggregatorChanged",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "totalStaked",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "unstakeDelaySec",
				"type": "uint256"
			}
		],
		"name": "StakeLocked",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "withdrawTime",
				"type": "uint256"
			}
		],
		"name": "StakeUnlocked",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "withdrawAddress",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "StakeWithdrawn",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "userOpHash",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "paymaster",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "nonce",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "bool",
				"name": "success",
				"type": "bool"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "actualGasCost",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "actualGasUsed",
				"type": "uint256"
			}
		],
		"name": "UserOperationEvent",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "userOpHash",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "nonce",
				"type": "uint256"
			}
		],
		"name": "UserOperationPrefundTooLow",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "userOpHash",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "nonce",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "bytes",
				"name": "revertReason",
				"type": "bytes"
			}
		],
		"name": "UserOperationRevertReason",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "withdrawAddress",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "Withdrawn",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "uint32",
				"name": "unstakeDelaySec",
				"type": "uint32"
			}
		],
		"name": "addStake",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "balanceOf",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "target",
				"type": "address"
			},
			{
				"internalType": "bytes",
				"name": "data",
				"type": "bytes"
			}
		],
		"name": "delegateAndRevert",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "depositTo",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "eip712Domain",
		"outputs": [
			{
				"internalType": "bytes1",
				"name": "fields",
				"type": "bytes1"
			},
			{
				"internalType": "string",
				"name": "name",
				"type": "string"
			},
			{
				"internalType": "string",
				"name": "version",
				"type": "string"
			},
			{
				"internalType": "uint256",
				"name": "chainId",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "verifyingContract",
				"type": "address"
			},
			{
				"internalType": "bytes32",
				"name": "salt",
				"type": "bytes32"
			},
			{
				"internalType": "uint256[]",
				"name": "extensions",
				"type": "uint256[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getCurrentUserOpHash",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "getDepositInfo",
		"outputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "deposit",
						"type": "uint256"
					},
					{
						"internalType": "bool",
						"name": "staked",
						"type": "bool"
					},
					{
						"internalType": "uint112",
						"name": "stake",
						"type": "uint112"
					},
					{
						"internalType": "uint32",
						"name": "unstakeDelaySec",
						"type": "uint32"
					},
					{
						"internalType": "uint48",
						"name": "withdrawTime",
						"type": "uint48"
					}
				],
				"internalType": "struct IStakeManager.DepositInfo",
				"name": "info",
				"type": "tuple"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getDomainSeparatorV4",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "sender",
				"type": "address"
			},
			{
				"internalType": "uint192",
				"name": "key",
				"type": "uint192"
			}
		],
		"name": "getNonce",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "nonce",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getPackedUserOpTypeHash",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes",
				"name": "initCode",
				"type": "bytes"
			}
		],
		"name": "getSenderAddress",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"components": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					},
					{
						"internalType": "bytes",
						"name": "initCode",
						"type": "bytes"
					},
					{
						"internalType": "bytes",
						"name": "callData",
						"type": "bytes"
					},
					{
						"internalType": "bytes32",
						"name": "accountGasLimits",
						"type": "bytes32"
					},
					{
						"internalType": "uint256",
						"name": "preVerificationGas",
						"type": "uint256"
					},
					{
						"internalType": "bytes32",
						"name": "gasFees",
						"type": "bytes32"
					},
					{
						"internalType": "bytes",
						"name": "paymasterAndData",
						"type": "bytes"
					},
					{
						"internalType": "bytes",
						"name": "signature",
						"type": "bytes"
					}
				],
				"internalType": "struct PackedUserOperation",
				"name": "userOp",
				"type": "tuple"
			}
		],
		"name": "getUserOpHash",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"components": [
					{
						"components": [
							{
								"internalType": "address",
								"name": "sender",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "nonce",
								"type": "uint256"
							},
							{
								"internalType": "bytes",
								"name": "initCode",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "callData",
								"type": "bytes"
							},
							{
								"internalType": "bytes32",
								"name": "accountGasLimits",
								"type": "bytes32"
							},
							{
								"internalType": "uint256",
								"name": "preVerificationGas",
								"type": "uint256"
							},
							{
								"internalType": "bytes32",
								"name": "gasFees",
								"type": "bytes32"
							},
							{
								"internalType": "bytes",
								"name": "paymasterAndData",
								"type": "bytes"
							},
							{
								"internalType": "bytes",
								"name": "signature",
								"type": "bytes"
							}
						],
						"internalType": "struct PackedUserOperation[]",
						"name": "userOps",
						"type": "tuple[]"
					},
					{
						"internalType": "contract IAggregator",
						"name": "aggregator",
						"type": "address"
					},
					{
						"internalType": "bytes",
						"name": "signature",
						"type": "bytes"
					}
				],
				"internalType": "struct IEntryPoint.UserOpsPerAggregator[]",
				"name": "opsPerAggregator",
				"type": "tuple[]"
			},
			{
				"internalType": "address payable",
				"name": "beneficiary",
				"type": "address"
			}
		],
		"name": "handleAggregatedOps",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"components": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "nonce",
						"type": "uint256"
					},
					{
						"internalType": "bytes",
						"name": "initCode",
						"type": "bytes"
					},
					{
						"internalType": "bytes",
						"name": "callData",
						"type": "bytes"
					},
					{
						"internalType": "bytes32",
						"name": "accountGasLimits",
						"type": "bytes32"
					},
					{
						"internalType": "uint256",
						"name": "preVerificationGas",
						"type": "uint256"
					},
					{
						"internalType": "bytes32",
						"name": "gasFees",
						"type": "bytes32"
					},
					{
						"internalType": "bytes",
						"name": "paymasterAndData",
						"type": "bytes"
					},
					{
						"internalType": "bytes",
						"name": "signature",
						"type": "bytes"
					}
				],
				"internalType": "struct PackedUserOperation[]",
				"name": "ops",
				"type": "tuple[]"
			},
			{
				"internalType": "address payable",
				"name": "beneficiary",
				"type": "address"
			}
		],
		"name": "handleOps",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint192",
				"name": "key",
				"type": "uint192"
			}
		],
		"name": "incrementNonce",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes",
				"name": "callData",
				"type": "bytes"
			},
			{
				"components": [
					{
						"components": [
							{
								"internalType": "address",
								"name": "sender",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "nonce",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "verificationGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "callGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "paymasterVerificationGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "paymasterPostOpGasLimit",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "preVerificationGas",
								"type": "uint256"
							},
							{
								"internalType": "address",
								"name": "paymaster",
								"type": "address"
							},
							{
								"internalType": "uint256",
								"name": "maxFeePerGas",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "maxPriorityFeePerGas",
								"type": "uint256"
							}
						],
						"internalType": "struct EntryPoint.MemoryUserOp",
						"name": "mUserOp",
						"type": "tuple"
					},
					{
						"internalType": "bytes32",
						"name": "userOpHash",
						"type": "bytes32"
					},
					{
						"internalType": "uint256",
						"name": "prefund",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "contextOffset",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "preOpGas",
						"type": "uint256"
					}
				],
				"internalType": "struct EntryPoint.UserOpInfo",
				"name": "opInfo",
				"type": "tuple"
			},
			{
				"internalType": "bytes",
				"name": "context",
				"type": "bytes"
			}
		],
		"name": "innerHandleOp",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "actualGasCost",
				"type": "uint256"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"internalType": "uint192",
				"name": "",
				"type": "uint192"
			}
		],
		"name": "nonceSequenceNumber",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "senderCreator",
		"outputs": [
			{
				"internalType": "contract ISenderCreator",
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
				"internalType": "bytes4",
				"name": "interfaceId",
				"type": "bytes4"
			}
		],
		"name": "supportsInterface",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "unlockStake",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address payable",
				"name": "withdrawAddress",
				"type": "address"
			}
		],
		"name": "withdrawStake",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address payable",
				"name": "withdrawAddress",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "withdrawAmount",
				"type": "uint256"
			}
		],
		"name": "withdrawTo",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"stateMutability": "payable",
		"type": "receive"
	}
]`
