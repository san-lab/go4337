package abiutil

import "encoding/json"

type ABIStruct struct {
	Compiler struct {
		Version string `json:"version"`
	} `json:"compiler"`
	Language string `json:"language"`
	Output   struct {
		Abi json.RawMessage `json:"abi"`
		/*
			[]struct {
				Inputs []struct {
					InternalType string `json:"internalType"`
					Name         string `json:"name"`
					Type         string `json:"type"`
				} `json:"inputs,omitempty"`
				StateMutability string `json:"stateMutability,omitempty"`
				Type            string `json:"type"`
				Name            string `json:"name,omitempty"`
				Anonymous       bool   `json:"anonymous,omitempty"`
				Outputs         []struct {
					InternalType string `json:"internalType"`
					Name         string `json:"name"`
					Type         string `json:"type"`
				} `json:"outputs,omitempty"`
			}*/
		/*
			Devdoc struct {
				Errors struct {
					ECDSAInvalidSignature []struct {
						Details string `json:"details"`
					} `json:"ECDSAInvalidSignature()"`
					ECDSAInvalidSignatureLengthUint256 []struct {
						Details string `json:"details"`
					} `json:"ECDSAInvalidSignatureLength(uint256)"`
					ECDSAInvalidSignatureSBytes32 []struct {
						Details string `json:"details"`
					} `json:"ECDSAInvalidSignatureS(bytes32)"`
					ERC1967InvalidImplementationAddress []struct {
						Details string `json:"details"`
					} `json:"ERC1967InvalidImplementation(address)"`
					ERC1967NonPayable []struct {
						Details string `json:"details"`
					} `json:"ERC1967NonPayable()"`
					InvalidInitialization []struct {
						Details string `json:"details"`
					} `json:"InvalidInitialization()"`
					NotInitializing []struct {
						Details string `json:"details"`
					} `json:"NotInitializing()"`
					UUPSUnauthorizedCallContext []struct {
						Details string `json:"details"`
					} `json:"UUPSUnauthorizedCallContext()"`
					UUPSUnsupportedProxiableUUIDBytes32 []struct {
						Details string `json:"details"`
					} `json:"UUPSUnsupportedProxiableUUID(bytes32)"`
				} `json:"errors"`
				Events struct {
					InitializedUint64 struct {
						Details string `json:"details"`
					} `json:"Initialized(uint64)"`
					UpgradedAddress struct {
						Details string `json:"details"`
					} `json:"Upgraded(address)"`
				} `json:"events"`
				Kind    string `json:"kind"`
				Methods struct {
					ExecuteAddressUint256Bytes struct {
						Params struct {
							Dest  string `json:"dest"`
							Func  string `json:"func"`
							Value string `json:"value"`
						} `json:"params"`
					} `json:"execute(address,uint256,bytes)"`
					ExecuteBatchAddressUint256Bytes struct {
						Details string `json:"details"`
						Params  struct {
							Dest  string `json:"dest"`
							Func  string `json:"func"`
							Value string `json:"value"`
						} `json:"params"`
					} `json:"executeBatch(address[],uint256[],bytes[])"`
					InitializeAddress struct {
						Details string `json:"details"`
						Params  struct {
							AnOwner string `json:"anOwner"`
						} `json:"params"`
					} `json:"initialize(address)"`
					ProxiableUUID struct {
						Details string `json:"details"`
					} `json:"proxiableUUID()"`
					SupportsInterfaceBytes4 struct {
						Details string `json:"details"`
					} `json:"supportsInterface(bytes4)"`
					UpgradeToAndCallAddressBytes struct {
						CustomOzUpgradesUnsafeAllowReachable string `json:"custom:oz-upgrades-unsafe-allow-reachable"`
						Details                              string `json:"details"`
					} `json:"upgradeToAndCall(address,bytes)"`
					ValidateUserOpAddressUint256BytesBytesBytes32Uint256Bytes32BytesBytesBytes32Uint256 struct {
						Details string `json:"details"`
						Params  struct {
							MissingAccountFunds string `json:"missingAccountFunds"`
							UserOp              string `json:"userOp"`
							UserOpHash          string `json:"userOpHash"`
						} `json:"params"`
						Returns struct {
							ValidationData string `json:"validationData"`
						} `json:"returns"`
					} `json:"validateUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes),bytes32,uint256)"`
					WithdrawDepositToAddressUint256 struct {
						Params struct {
							Amount          string `json:"amount"`
							WithdrawAddress string `json:"withdrawAddress"`
						} `json:"params"`
					} `json:"withdrawDepositTo(address,uint256)"`
				} `json:"methods"`
				Version int `json:"version"`
			} `json:"devdoc"`
		*/
		/*
			Userdoc struct {
				Kind    string `json:"kind"`
				Methods struct {
					AddDeposit struct {
						Notice string `json:"notice"`
					} `json:"addDeposit()"`
					EntryPoint struct {
						Notice string `json:"notice"`
					} `json:"entryPoint()"`
					ExecuteAddressUint256Bytes struct {
						Notice string `json:"notice"`
					} `json:"execute(address,uint256,bytes)"`
					ExecuteBatchAddressUint256Bytes struct {
						Notice string `json:"notice"`
					} `json:"executeBatch(address[],uint256[],bytes[])"`
					GetDeposit struct {
						Notice string `json:"notice"`
					} `json:"getDeposit()"`
					GetNonce struct {
						Notice string `json:"notice"`
					} `json:"getNonce()"`
					ValidateUserOpAddressUint256BytesBytesBytes32Uint256Bytes32BytesBytesBytes32Uint256 struct {
						Notice string `json:"notice"`
					} `json:"validateUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes),bytes32,uint256)"`
					WithdrawDepositToAddressUint256 struct {
						Notice string `json:"notice"`
					} `json:"withdrawDepositTo(address,uint256)"`
				} `json:"methods"`
				Notice  string `json:"notice"`
				Version int    `json:"version"`
			} `json:"userdoc"`
		*/
	} `json:"output"`
	/*
		Settings struct {
			CompilationTarget struct {
				ContractsContractsSamplesSimpleAccountSol string `json:"contracts/contracts/samples/SimpleAccount.sol"`
			} `json:"compilationTarget"`
			EvmVersion string `json:"evmVersion"`
			Libraries  struct {
			} `json:"libraries"`
			Metadata struct {
				BytecodeHash string `json:"bytecodeHash"`
			} `json:"metadata"`
			Optimizer struct {
				Enabled bool `json:"enabled"`
				Runs    int  `json:"runs"`
			} `json:"optimizer"`
			Remappings []interface{} `json:"remappings"`
		} `json:"settings"`
	*/
	/*
		Sources struct {
			OpenzeppelinContractsInterfacesDraftIERC1822Sol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/interfaces/draft-IERC1822.sol"`
			OpenzeppelinContractsProxyERC1967ERC1967UtilsSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/proxy/ERC1967/ERC1967Utils.sol"`
			OpenzeppelinContractsProxyBeaconIBeaconSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/proxy/beacon/IBeacon.sol"`
			OpenzeppelinContractsProxyUtilsInitializableSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/proxy/utils/Initializable.sol"`
			OpenzeppelinContractsProxyUtilsUUPSUpgradeableSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol"`
			OpenzeppelinContractsTokenERC1155IERC1155ReceiverSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/token/ERC1155/IERC1155Receiver.sol"`
			OpenzeppelinContractsTokenERC721IERC721ReceiverSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol"`
			OpenzeppelinContractsUtilsAddressSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/utils/Address.sol"`
			OpenzeppelinContractsUtilsStorageSlotSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/utils/StorageSlot.sol"`
			OpenzeppelinContractsUtilsStringsSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/utils/Strings.sol"`
			OpenzeppelinContractsUtilsCryptographyECDSASol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/utils/cryptography/ECDSA.sol"`
			OpenzeppelinContractsUtilsCryptographyMessageHashUtilsSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol"`
			OpenzeppelinContractsUtilsIntrospectionIERC165Sol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/utils/introspection/IERC165.sol"`
			OpenzeppelinContractsUtilsMathMathSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"@openzeppelin/contracts/utils/math/Math.sol"`
			ContractsContractsCoreBaseAccountSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/core/BaseAccount.sol"`
			ContractsContractsCoreHelpersSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/core/Helpers.sol"`
			ContractsContractsCoreUserOperationLibSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/core/UserOperationLib.sol"`
			ContractsContractsInterfacesIAccountSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/interfaces/IAccount.sol"`
			ContractsContractsInterfacesIAggregatorSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/interfaces/IAggregator.sol"`
			ContractsContractsInterfacesIEntryPointSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/interfaces/IEntryPoint.sol"`
			ContractsContractsInterfacesINonceManagerSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/interfaces/INonceManager.sol"`
			ContractsContractsInterfacesIStakeManagerSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/interfaces/IStakeManager.sol"`
			ContractsContractsInterfacesPackedUserOperationSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/interfaces/PackedUserOperation.sol"`
			ContractsContractsSamplesSimpleAccountSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/samples/SimpleAccount.sol"`
			ContractsContractsSamplesCallbackTokenCallbackHandlerSol struct {
				Keccak256 string   `json:"keccak256"`
				License   string   `json:"license"`
				Urls      []string `json:"urls"`
			} `json:"contracts/contracts/samples/callback/TokenCallbackHandler.sol"`
		} `json:"sources"`
	*/
	Version int `json:"version"`
}

type FullABIJsonRemix struct {
	Deploy struct {
		VM struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"VM:-"`
		Main1 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"main:1"`
		Ropsten3 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"ropsten:3"`
		Rinkeby4 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"rinkeby:4"`
		Kovan42 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"kovan:42"`
		Goerli5 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"goerli:5"`
		Custom struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"Custom"`
	} `json:"deploy"`
	Data struct {
		Bytecode struct {
			FunctionDebugData struct {
				Num23 struct {
					EntryPoint     interface{} `json:"entryPoint"`
					ID             int         `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"@_23"`
				AbiDecodeTupleTStringMemoryPtrFromMemory struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"abi_decode_tuple_t_string_memory_ptr_fromMemory"`
				ArrayDataslotStringStorage struct {
					EntryPoint     interface{} `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"array_dataslot_string_storage"`
				CleanUpBytearrayEndSlotsStringStorage struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"clean_up_bytearray_end_slots_string_storage"`
				CopyByteArrayToStorageFromTStringMemoryPtrToTStringStorage struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"copy_byte_array_to_storage_from_t_string_memory_ptr_to_t_string_storage"`
				ExtractByteArrayLength struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"extract_byte_array_length"`
				ExtractUsedPartAndSetLengthOfShortByteArray struct {
					EntryPoint     interface{} `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"extract_used_part_and_set_length_of_short_byte_array"`
				PanicError0X41 struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"panic_error_0x41"`
			} `json:"functionDebugData"`
			GeneratedSources []struct {
				Ast struct {
					NativeSrc  string `json:"nativeSrc"`
					NodeType   string `json:"nodeType"`
					Src        string `json:"src"`
					Statements []struct {
						NativeSrc  string        `json:"nativeSrc"`
						NodeType   string        `json:"nodeType"`
						Src        string        `json:"src"`
						Statements []interface{} `json:"statements,omitempty"`
						Body       struct {
							NativeSrc  string `json:"nativeSrc"`
							NodeType   string `json:"nodeType"`
							Src        string `json:"src"`
							Statements []struct {
								Expression struct {
									Arguments []struct {
										Kind      string `json:"kind,omitempty"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
										Type      string `json:"type,omitempty"`
										Value     string `json:"value,omitempty"`
										Arguments []struct {
											Kind      string `json:"kind"`
											NativeSrc string `json:"nativeSrc"`
											NodeType  string `json:"nodeType"`
											Src       string `json:"src"`
											Type      string `json:"type"`
											Value     string `json:"value"`
										} `json:"arguments,omitempty"`
										FunctionName struct {
											Name      string `json:"name"`
											NativeSrc string `json:"nativeSrc"`
											NodeType  string `json:"nodeType"`
											Src       string `json:"src"`
										} `json:"functionName,omitempty"`
									} `json:"arguments"`
									FunctionName struct {
										Name      string `json:"name"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
									} `json:"functionName"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
								} `json:"expression"`
								NativeSrc string `json:"nativeSrc"`
								NodeType  string `json:"nodeType"`
								Src       string `json:"src"`
							} `json:"statements"`
						} `json:"body,omitempty"`
						Name       string `json:"name,omitempty"`
						Parameters []struct {
							Name      string `json:"name"`
							NativeSrc string `json:"nativeSrc"`
							NodeType  string `json:"nodeType"`
							Src       string `json:"src"`
							Type      string `json:"type"`
						} `json:"parameters,omitempty"`
						ReturnVariables []struct {
							Name      string `json:"name"`
							NativeSrc string `json:"nativeSrc"`
							NodeType  string `json:"nodeType"`
							Src       string `json:"src"`
							Type      string `json:"type"`
						} `json:"returnVariables,omitempty"`
					} `json:"statements"`
				} `json:"ast"`
				Contents string `json:"contents"`
				ID       int    `json:"id"`
				Language string `json:"language"`
				Name     string `json:"name"`
			} `json:"generatedSources"`
			LinkReferences struct {
			} `json:"linkReferences"`
			Object    string `json:"object"`
			Opcodes   string `json:"opcodes"`
			SourceMap string `json:"sourceMap"`
		} `json:"bytecode"`
		DeployedBytecode struct {
			FunctionDebugData struct {
				AddressToBool7 struct {
					EntryPoint     interface{} `json:"entryPoint"`
					ID             int         `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"@addressToBool_7"`
				Greet31 struct {
					EntryPoint     int `json:"entryPoint"`
					ID             int `json:"id"`
					ParameterSlots int `json:"parameterSlots"`
					ReturnSlots    int `json:"returnSlots"`
				} `json:"@greet_31"`
				SetGreeting41 struct {
					EntryPoint     int `json:"entryPoint"`
					ID             int `json:"id"`
					ParameterSlots int `json:"parameterSlots"`
					ReturnSlots    int `json:"returnSlots"`
				} `json:"@setGreeting_41"`
				TestOverflows69 struct {
					EntryPoint     int `json:"entryPoint"`
					ID             int `json:"id"`
					ParameterSlots int `json:"parameterSlots"`
					ReturnSlots    int `json:"returnSlots"`
				} `json:"@testOverflows_69"`
				AbiDecodeTupleTAddress struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"abi_decode_tuple_t_address"`
				AbiDecodeTupleTStringMemoryPtr struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"abi_decode_tuple_t_string_memory_ptr"`
				AbiDecodeTupleTUint256TUint256 struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"abi_decode_tuple_t_uint256t_uint256"`
				AbiEncodeTupleTBoolToTBoolFromStackReversed struct {
					EntryPoint     interface{} `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"abi_encode_tuple_t_bool__to_t_bool__fromStack_reversed"`
				AbiEncodeTupleTStringMemoryPtrToTStringMemoryPtrFromStackReversed struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"abi_encode_tuple_t_string_memory_ptr__to_t_string_memory_ptr__fromStack_reversed"`
				AbiEncodeTupleTUint256TUint256ToTUint256TUint256FromStackReversed struct {
					EntryPoint     interface{} `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"abi_encode_tuple_t_uint256_t_uint256__to_t_uint256_t_uint256__fromStack_reversed"`
				ArrayDataslotStringStorage struct {
					EntryPoint     interface{} `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"array_dataslot_string_storage"`
				CheckedMulTUint256 struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"checked_mul_t_uint256"`
				CheckedSubTUint256 struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"checked_sub_t_uint256"`
				CleanUpBytearrayEndSlotsStringStorage struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"clean_up_bytearray_end_slots_string_storage"`
				CopyByteArrayToStorageFromTStringMemoryPtrToTStringStorage struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"copy_byte_array_to_storage_from_t_string_memory_ptr_to_t_string_storage"`
				ExtractByteArrayLength struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"extract_byte_array_length"`
				ExtractUsedPartAndSetLengthOfShortByteArray struct {
					EntryPoint     interface{} `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"extract_used_part_and_set_length_of_short_byte_array"`
				PanicError0X11 struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"panic_error_0x11"`
				PanicError0X41 struct {
					EntryPoint     int         `json:"entryPoint"`
					ID             interface{} `json:"id"`
					ParameterSlots int         `json:"parameterSlots"`
					ReturnSlots    int         `json:"returnSlots"`
				} `json:"panic_error_0x41"`
			} `json:"functionDebugData"`
			GeneratedSources []struct {
				Ast struct {
					NativeSrc  string `json:"nativeSrc"`
					NodeType   string `json:"nodeType"`
					Src        string `json:"src"`
					Statements []struct {
						NativeSrc  string        `json:"nativeSrc"`
						NodeType   string        `json:"nodeType"`
						Src        string        `json:"src"`
						Statements []interface{} `json:"statements,omitempty"`
						Body       struct {
							NativeSrc  string `json:"nativeSrc"`
							NodeType   string `json:"nodeType"`
							Src        string `json:"src"`
							Statements []struct {
								Body struct {
									NativeSrc  string `json:"nativeSrc"`
									NodeType   string `json:"nodeType"`
									Src        string `json:"src"`
									Statements []struct {
										Expression struct {
											Arguments []struct {
												Kind      string `json:"kind"`
												NativeSrc string `json:"nativeSrc"`
												NodeType  string `json:"nodeType"`
												Src       string `json:"src"`
												Type      string `json:"type"`
												Value     string `json:"value"`
											} `json:"arguments"`
											FunctionName struct {
												Name      string `json:"name"`
												NativeSrc string `json:"nativeSrc"`
												NodeType  string `json:"nodeType"`
												Src       string `json:"src"`
											} `json:"functionName"`
											NativeSrc string `json:"nativeSrc"`
											NodeType  string `json:"nodeType"`
											Src       string `json:"src"`
										} `json:"expression"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
									} `json:"statements"`
								} `json:"body,omitempty"`
								Condition struct {
									Arguments []struct {
										Arguments []struct {
											Name      string `json:"name"`
											NativeSrc string `json:"nativeSrc"`
											NodeType  string `json:"nodeType"`
											Src       string `json:"src"`
										} `json:"arguments,omitempty"`
										FunctionName struct {
											Name      string `json:"name"`
											NativeSrc string `json:"nativeSrc"`
											NodeType  string `json:"nodeType"`
											Src       string `json:"src"`
										} `json:"functionName,omitempty"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
										Kind      string `json:"kind,omitempty"`
										Type      string `json:"type,omitempty"`
										Value     string `json:"value,omitempty"`
									} `json:"arguments"`
									FunctionName struct {
										Name      string `json:"name"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
									} `json:"functionName"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
								} `json:"condition,omitempty"`
								NativeSrc string `json:"nativeSrc"`
								NodeType  string `json:"nodeType"`
								Src       string `json:"src"`
								Value     struct {
									Kind      string `json:"kind"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
									Type      string `json:"type"`
									Value     string `json:"value"`
								} `json:"value,omitempty"`
								Variables []struct {
									Name      string `json:"name"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
									Type      string `json:"type"`
								} `json:"variables,omitempty"`
								Value0 struct {
									Arguments []struct {
										Name      string `json:"name"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
									} `json:"arguments"`
									FunctionName struct {
										Name      string `json:"name"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
									} `json:"functionName"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
								} `json:"value,omitempty"`
								VariableNames []struct {
									Name      string `json:"name"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
								} `json:"variableNames,omitempty"`
								Value1 struct {
									Name      string `json:"name"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
								} `json:"value,omitempty"`
								Value2 struct {
									Arguments []struct {
										Arguments []struct {
											Name      string `json:"name,omitempty"`
											NativeSrc string `json:"nativeSrc"`
											NodeType  string `json:"nodeType"`
											Src       string `json:"src"`
											Kind      string `json:"kind,omitempty"`
											Type      string `json:"type,omitempty"`
											Value     string `json:"value,omitempty"`
										} `json:"arguments"`
										FunctionName struct {
											Name      string `json:"name"`
											NativeSrc string `json:"nativeSrc"`
											NodeType  string `json:"nodeType"`
											Src       string `json:"src"`
										} `json:"functionName"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
									} `json:"arguments"`
									FunctionName struct {
										Name      string `json:"name"`
										NativeSrc string `json:"nativeSrc"`
										NodeType  string `json:"nodeType"`
										Src       string `json:"src"`
									} `json:"functionName"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
								} `json:"value,omitempty"`
								Value3 struct {
									Name      string `json:"name"`
									NativeSrc string `json:"nativeSrc"`
									NodeType  string `json:"nodeType"`
									Src       string `json:"src"`
								} `json:"value,omitempty"`
							} `json:"statements"`
						} `json:"body,omitempty"`
						Name       string `json:"name,omitempty"`
						Parameters []struct {
							Name      string `json:"name"`
							NativeSrc string `json:"nativeSrc"`
							NodeType  string `json:"nodeType"`
							Src       string `json:"src"`
							Type      string `json:"type"`
						} `json:"parameters,omitempty"`
						ReturnVariables []struct {
							Name      string `json:"name"`
							NativeSrc string `json:"nativeSrc"`
							NodeType  string `json:"nodeType"`
							Src       string `json:"src"`
							Type      string `json:"type"`
						} `json:"returnVariables,omitempty"`
					} `json:"statements"`
				} `json:"ast"`
				Contents string `json:"contents"`
				ID       int    `json:"id"`
				Language string `json:"language"`
				Name     string `json:"name"`
			} `json:"generatedSources"`
			ImmutableReferences struct {
			} `json:"immutableReferences"`
			LinkReferences struct {
			} `json:"linkReferences"`
			Object    string `json:"object"`
			Opcodes   string `json:"opcodes"`
			SourceMap string `json:"sourceMap"`
		} `json:"deployedBytecode"`
		GasEstimates struct {
			Creation struct {
				CodeDepositCost string `json:"codeDepositCost"`
				ExecutionCost   string `json:"executionCost"`
				TotalCost       string `json:"totalCost"`
			} `json:"creation"`
			External struct {
				AddressToBoolAddress        string `json:"addressToBool(address)"`
				Greet                       string `json:"greet()"`
				SetGreetingString           string `json:"setGreeting(string)"`
				TestOverflowsUint256Uint256 string `json:"testOverflows(uint256,uint256)"`
			} `json:"external"`
		} `json:"gasEstimates"`
		MethodIdentifiers struct {
			AddressToBoolAddress        string `json:"addressToBool(address)"`
			Greet                       string `json:"greet()"`
			SetGreetingString           string `json:"setGreeting(string)"`
			TestOverflowsUint256Uint256 string `json:"testOverflows(uint256,uint256)"`
		} `json:"methodIdentifiers"`
	} `json:"data"`
	Abi string `json:"abi"`
}
