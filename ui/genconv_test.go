package ui

import (
	"fmt"
	"testing"

	"github.com/san-lab/go4337/rpccalls"
)

// Example of your version-specific result struct

// Other fields are nil
func TestGasResults(t *testing.T) {
	egr6 := rpccalls.EthEstimateUserOperationGasResultV6{CallGasLimit: 7, PreVerificationGas: 8, MaxFeePerGas: "0x22"}
	gapv := MapResultToGasAndPaymasterValues(egr6)
	fmt.Println(*gapv.CallGasLimit, *gapv.PreVerificationGas, *gapv.MaxFeePerGas)

	egr7 := rpccalls.EthEstimateUserOperationGasResultV7{CallGasLimit: "0x07", PreVerificationGas: "0x08"}
	gapv = MapResultToGasAndPaymasterValues(&egr7)
	fmt.Println(*gapv.CallGasLimit, *gapv.PreVerificationGas)
}
