package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/orientwalt/htdf/types"
)

const (
	//
	FlagEncode = "encode"
	//
	FlagOffline = "offline"
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

// ResultData represents the data returned in an sdk.Result
type ResultData struct {
	ContractAddress ethcmn.Address  `json:"contract_address"`
	Bloom           ethtypes.Bloom  `json:"bloom"`
	Logs            []*ethtypes.Log `json:"logs"`
	Ret             []byte          `json:"ret"`
	TxHash          ethcmn.Hash     `json:"tx_hash"`
}

// String implements fmt.Stringer interface.
func (rd ResultData) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ResultData:
	ContractAddress: %s
	Bloom: %s
	Logs: %v
	Ret: %v
	TxHash: %s
`, rd.ContractAddress.String(), rd.Bloom.Big().String(), rd.Logs, rd.Ret, rd.TxHash.String()))
}
