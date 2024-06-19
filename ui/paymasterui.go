package ui

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/state"
)

var AddPaymasterItem = &Item{Label: "Add Paymaster", Details: "Add a new paymaster address"}

func PaymasterUI() (*common.Address, bool) {

	addr, ok := AddressFromBookUI(state.Paymaster)
	if ok {
		PaymasterItem.Value = addr
		PaymasterItem.DisplayValue = addr.String()
	}
	return addr, ok
}
