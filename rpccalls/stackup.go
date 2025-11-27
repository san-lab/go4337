package rpccalls

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

const StackupAPIDefURL = "https://api.stackup.sh/v1/paymaster/"

const StackUpPMPayTemplate = "{\"jsonrpc\": \"2.0\",\"id\": 1,\"method\": \"pm_sponsorUserOperation\",\"params\": [%s,\"0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789\",{\"type\": \"payg\"}]}"
const StackUpSendOpTemplate = "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_sendUserOperation\",\"params\":[%s,\"0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789\" ]}"

type StackUpPMPayGResult struct {
	PaymasterAndData     string `json:"paymasterAndData"`
	PreVerificationGas   string `json:"preVerificationGas"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	CallGasLimit         string `json:"callGasLimit"`
}

func StackUpPMPayCall(url, key string, usop *userop.UserOpForApiV6) (*StackUpPMPayGResult, error) {
	ar := &APIRequest{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "pm_sponsorUserOperation",
		Params:  []interface{}{usop, "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789", map[string]string{"type": "payg"}},
	}
	result := &StackUpPMPayGResult{}
	_, err := ApiCall(url, key, ar, result)
	if err != nil {
		return nil, fmt.Errorf("API Call error: %w", err)
	}
	return result, nil

}

func IncorporateStackUpPMResToUserOp(usop *userop.UserOperation, res *StackUpPMPayGResult) error {
	pma := common.HexToAddress(res.PaymasterAndData[:42])
	usop.Paymaster = &pma
	usop.PaymasterData, _ = hex.DecodeString(res.PaymasterAndData[42:])

	usop.PreVerificationGas = ConvHexOrZero(res.PreVerificationGas)
	usop.VerificationGasLimit = ConvHexOrZero(res.VerificationGasLimit)
	usop.CallGasLimit = ConvHexOrZero(res.CallGasLimit)
	return nil
}

func ConvHexOrZero(s string) uint64 {
	i, err := strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
	if err != nil && state.DEBUG {
		log.Println(err)
	}
	return i
}
