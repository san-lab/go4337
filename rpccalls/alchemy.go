package rpccalls

// USerOpV6 - standard

type AlchemyV7UserOp struct {
	Sender                        string `json:"sender"`
	Nonce                         string `json:"nonce"`
	Factory                       string `json:"factory"`
	FactoryData                   string `json:"factoryData"`
	CallData                      string `json:"callData"`
	CallGasLimit                  string `json:"callGasLimit"`
	VerificationGasLimit          string `json:"verificationGasLimit"`
	PreVerificationGas            string `json:"preVerificationGas"`
	MaxFeePerGas                  string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas          string `json:"maxPriorityFeePerGas"`
	PaymasterVerificationGasLimit string `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       string `json:"paymasterPostOpGasLimit"`
	Signature                     string `json:"signature"`
	Paymaster                     string `json:"paymaster"`
	PaymasterData                 string `json:"paymasterData"`
}

type AlchemyEstimateGasCostResponse struct {
	PreVerificationGas   string `json:"preVerificationGas"`
	CallGasLimit         string `json:"callGasLimit"`
	VerificationGasLimit string `json:"verificationGasLimit"`
}

/*

eth_sendUserOperation
Submits a user operation to a Bundler. If the request is successful, the endpoint will return a user operation hash that the caller can use to look up the status of the user operation. If it fails, or another error occurs, an error code and description will be returned.
eth_estimateUserOperationGas
Estimates the gas values for a user operation. It returns the preVerificationGas, verificationGasLimit, and callGasLimit values associated with the provided user operation.
eth_getUserOperationByHash
Returns a user operation based on the given user operation hash. It returns the user operation along with extra information including what block/transaction it was included in. If the operation has not yet been included, it will return null.
eth_getUserOperationReceipt
Returns a user operation receipt ( metadata associated with the given user operation ) based on the given user operation hash. It returns null if the user operation has not yet been included.
rundler_maxPriorityFeePerGas
Returns a fee per gas that is an estimate of how much users should set as a priority fee in UOs for Rundler endpoints.
eth_supportedEntryPoints
Returns a list of Entrypoint contract addresses supported by the bunder endpoints.

*/
