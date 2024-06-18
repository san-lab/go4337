package ui

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
