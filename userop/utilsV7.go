package userop

import "math/big"

func (mop *UserOperation) GetRequiredPrefund() *big.Int {
	requiredGas := big.NewInt(0)
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.VerificationGasLimit)))
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.CallGasLimit)))
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.PaymasterVerificationGasLimit)))
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.PaymasterPostOpGasLimit)))
	requiredGas.Add(requiredGas, big.NewInt(int64(mop.PreVerificationGas)))
	requiredGas.Mul(requiredGas, big.NewInt(int64(mop.MaxFeePerGas)))
	//requiredGas.Mul(requiredGas, big.NewInt(int64(mop.MaxPriorityFeePerGas)))
	/*requiredGas.Add(requiredGas, mop.CallGasLimit)
	requiredGas.Add(requiredGas, mop.PaymasterVerificationGasLimit)
	requiredGas.Add(requiredGas, mop.PaymasterPostOpGasLimit)
	requiredGas.Add(requiredGas, mop.PreVerificationGas)
	requiredGas.Mul(requiredGas, mop.MaxFeePerGas)*/
	return requiredGas
}
