package htdfservice

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"

	"github.com/orientwalt/htdf/evm/vm"
	"github.com/orientwalt/htdf/types"

	"github.com/ethereum/go-ethereum/common"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
)

// junying-todo, 2019-11-06
// from x/core/transition.go
// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, homestead bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if len(data) > 0 && homestead {
		gas = params.DefaultMsgGasContractCreation // 53000 -> 60000
	} else {
		gas = params.DefaultMsgGas // 21000 -> 30000
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/ethparams.TxDataNonZeroGas < nz {
			return 0, vm.ErrOutOfGas
		}
		gas += nz * ethparams.TxDataNonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/ethparams.TxDataZeroGas < z {
			return 0, vm.ErrOutOfGas
		}
		gas += z * ethparams.TxDataZeroGas
	}
	return gas, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MsgSend defines a SendFrom message /////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
type MsgSend struct {
	From      sdk.AccAddress
	To        sdk.AccAddress
	Amount    sdk.Coins
	Data      string
	GasPrice  uint64 //unit,  satoshi/gallon
	GasWanted uint64 //unit,  gallon
}

var _ sdk.Msg = MsgSend{}

// NewMsgSend is a constructor function for MsgSend
// Normal Transaction
// Default GasWanted, Default GasPrice
func NewMsgSendDefault(fromaddr sdk.AccAddress, toaddr sdk.AccAddress, amount sdk.Coins) MsgSend {
	return MsgSend{
		From:      fromaddr,
		To:        toaddr,
		Amount:    amount,
		GasPrice:  params.DefaultMinGasPrice,
		GasWanted: params.DefaultMsgGas,
	}
}

// Normal Transaction
// Default GasWanted, Customized GasPrice
func NewMsgSend(fromaddr sdk.AccAddress, toaddr sdk.AccAddress, amount sdk.Coins, gasPrice uint64, gasWanted uint64) MsgSend {
	return MsgSend{
		From:      fromaddr,
		To:        toaddr,
		Amount:    amount,
		GasPrice:  gasPrice,
		GasWanted: gasWanted,
	}
}

// Contract Transaction
func NewMsgSendForData(fromaddr sdk.AccAddress, toaddr sdk.AccAddress, amount sdk.Coins, data string, gasPrice uint64, gasWanted uint64) MsgSend {
	return MsgSend{
		From:      fromaddr,
		To:        toaddr,
		Amount:    amount,
		Data:      data,
		GasPrice:  gasPrice,
		GasWanted: gasWanted,
	}
}

// Route should return the name of the module
func (msg MsgSend) Route() string { return "htdfservice" }

// Type should return the action
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}

	if len(msg.Data) == 0 {
		// classic transfer

		// must have to address
		if msg.To.Empty() {
			return sdk.ErrInvalidAddress(msg.To.String())
		}

		// amount > 0
		if !msg.Amount.IsAllPositive() {
			return sdk.ErrInsufficientCoins("Amount must be positive")
		}

		// junying-todo, 2019-11-12
		if msg.GasWanted < params.DefaultMsgGas {
			// return sdk.ErrOutOfGas(fmt.Sprintf("gaswanted must be greather than %d", params.DefaultMsgGas))
			return sdk.ErrInvalidGas(fmt.Sprintf("gaswanted must be greather than %d", params.DefaultMsgGas))
		}

	} else {
		// junying-todo, 2019-11-12
		inputCode, err := hex.DecodeString(msg.Data)
		if err != nil {
			return sdk.ErrTxDecode("decoding msg.data failed. you should check msg.data")
		}
		//Intrinsic gas calc
		itrsGas, err := IntrinsicGas(inputCode, true)
		if err != nil {
			return sdk.ErrOutOfGas("intrinsic out of gas")
		}
		if itrsGas > msg.GasWanted {
			return sdk.ErrOutOfGas(fmt.Sprintf("gaswanted must be greather than %d to pass validating", itrsGas))
		}

	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSend) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// GetStringAddr defines whose fromaddr is required
// func (msg MsgSend) GetFromAddrStr() string {
// 	return sdk.AccAddress.String(msg.From)
// }

//
func (msg MsgSend) FromAddress() common.Address {
	return types.ToEthAddress(msg.From)
}

// junying-todo, 2019-11-06
func (msg MsgSend) GetData() string {
	return msg.Data
}



// GetMsgSendData implement for sdk.GetMsgDataFunc, yqq 2020-12-03
func GetMsgSendData(msg sdk.Msg)([]byte, error) {
	if msgSend, ok := msg.(MsgSend) ; ok {
		return []byte(msgSend.Data), nil
	}
	return []byte{}, fmt.Errorf("msg is not MsgSend")
}


///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MsgAdd defines a Add message ///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
type MsgAdd struct {
	SystemIssuer sdk.AccAddress
	Amount       sdk.Coins
}

var _ sdk.Msg = MsgAdd{}

// NewMsgAdd is a constructor function for Msgadd
func NewMsgAdd(addr sdk.AccAddress, amount sdk.Coins) MsgAdd {
	return MsgAdd{
		SystemIssuer: addr,
		Amount:       amount,
	}
}

// Route should return the name of the module
func (msg MsgAdd) Route() string { return "htdfservice" }

// Type should return the action
func (msg MsgAdd) Type() string { return "add" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAdd) ValidateBasic() sdk.Error {
	if msg.SystemIssuer.Empty() {
		return sdk.ErrInvalidAddress(msg.SystemIssuer.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Amount must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAdd) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgAdd) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SystemIssuer}
}

// GetStringAddr defines whose fromaddr is required
func (msg MsgAdd) GetSystemIssuerStr() string {
	return sdk.AccAddress.String(msg.SystemIssuer)
}

//
func (msg MsgAdd) GeSystemIssuer() sdk.AccAddress {
	return msg.SystemIssuer
}
