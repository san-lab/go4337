package rpccalls

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/userop"
)

type SponsorUserOperationReq struct {
	Userop     userop.UserOpForApiV78
	Entrypoint common.Address
	Context    Payment
}

type Payment struct {
	Token             common.Address
	SponsorshipPolicy string
}
