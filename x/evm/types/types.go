package types

import (
	"encoding/json"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/orientwalt/htdf/types"
)

const (
	//
	FlagEncode = "encode"
	//
	FlagOffline = "offline"

	TxReceiptStatusSuccess = uint(1)
	TxReceiptStatusFail    = uint(0)
)

//
type SendTxResp struct {
	ErrCode         sdk.CodeType `json:"code"`
	ErrMsg          string       `json:"message"`
	ContractAddress string       `json:"contract_address"`
	EvmOutput       string       `json:"evm_output"`
}

//
func (result SendTxResp) String() string {
	result.ErrMsg = sdk.GetErrMsg(result.ErrCode)
	data, _ := json.Marshal(&result)
	return string(data)
}

// GasInfo returns the gas limit, gas consumed and gas refunded from the EVM transition
// execution
type GasInfo struct {
	GasLimit    uint64
	GasConsumed uint64
	GasRefunded uint64
}

// ExecutionResult represents what's returned from a transition
type ExecutionResult struct {
	Logs    []*ethtypes.Log
	Bloom   *big.Int
	Result  *sdk.Result
	GasInfo GasInfo
}
