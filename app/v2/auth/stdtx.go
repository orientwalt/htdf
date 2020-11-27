package auth

import (
	"fmt"

	"github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
	xauth "github.com/orientwalt/htdf/x/auth"
)

// junying-todo, 2019-11-22
// copy from  x/auth/stdtx.go and changed to normal function
func ValidateFeeV2(tx xauth.StdTx) sdk.Error {
	// junying-todo, 2019-11-13
	// MinGasPrice Checking
	if tx.Fee.GasPrice < params.DefaultMinGasPrice {
		return sdk.ErrGasPriceTooLow(fmt.Sprintf("gasprice must be greater than %d%s", params.DefaultMinGasPrice, sdk.DefaultDenom))
	}

	// junying-todo, 2019-11-13
	// Validate Msgs &
	// Check MinGas for staking txs
	var msgs = tx.Msgs
	if msgs == nil || len(msgs) == 0 {
		return sdk.ErrUnknownRequest("Tx.GetMsgs() must return at least one message in list")
	}
	count := 0
	for _, msg := range msgs {
		// Validate the Msg.
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}
		if msg.Route() == "htdfservice" {
			count = count + 1
		}
	}
	// only MsgSends or only OtherTypes in a Tx
	if count > 0 && count != len(msgs) {
		return sdk.ErrUnknownRequest("can't mix htdfservice msg with other-type msgs in a Tx")
	}
	// one MsgSend in one Tx
	if count > 1 {
		return sdk.ErrUnknownRequest("can't include more than one htdfservice msgs in a Tx")
	}

	// Checking minimum gaswanted condition for transactions
	minTxGasWanted := uint64(len(msgs)) * params.DefaultMsgGas
	if tx.Fee.GasWanted < minTxGasWanted {
		return sdk.ErrInvalidGas(fmt.Sprintf("Tx[count(msgs)=%d] gaswanted must be greater than %d", len(msgs), minTxGasWanted))
	}

	// fix issue #9 and issue #10
	// Checking maximum gaswanted condition for transactions
	if tx.Fee.GasWanted > params.TxGasLimit {
		return sdk.ErrInvalidGas(fmt.Sprintf("GasWanted[%d]  of Tx could not excess TxGasLimit[%d]", tx.Fee.GasWanted, params.TxGasLimit))
	}
	return nil
}
