package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
)

const (
	// TypeMsgEthereumTx defines the type string of an Ethereum tranasction
	TypeMsgSend = "send"
)

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
func NewMsgSend(fromaddr, toaddr sdk.AccAddress, amount sdk.Coins, gasPrice, gasWanted uint64) MsgSend {
	return MsgSend{
		From:      fromaddr,
		To:        toaddr,
		Amount:    amount,
		GasPrice:  gasPrice,
		GasWanted: gasWanted,
	}
}

// Contract Transaction
func NewMsgSendForData(fromaddr, toaddr sdk.AccAddress, amount sdk.Coins, data string, gasPrice, gasWanted uint64) MsgSend {
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
			return sdk.ErrOutOfGas(fmt.Sprintf("gaswanted must be greather than %d", params.DefaultMsgGas))
		}

	} else {
		// junying-todo, 2019-11-12
		inputCode, err := hex.DecodeString(msg.Data)
		if err != nil {
			return sdk.ErrTxDecode("decoding msg.data failed. you should check msg.data")
		}
		//Intrinsic gas calc
		itrsGas, err := IntrinsicGas(inputCode, msg.To == nil, true)
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
	return sdk.ToEthAddress(msg.From)
}

func (msg MsgSend) ToAddress() common.Address {
	return sdk.ToEthAddress(msg.To)
}

// junying-todo, 2019-11-06
func (msg MsgSend) GetData() string {
	return msg.Data
}

// GetMsgs returns a single MsgSend as an sdk.Msg.
func (msg MsgSend) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// Fee returns gasprice * gaslimit.
func (msg MsgSend) Fee() *big.Int {
	return new(big.Int).SetUint64(msg.GasWanted * msg.GasPrice)
}

func (msg *MsgSend) ChainID() *big.Int {
	return deriveChainID(msg.Data.V)
}

// deriveChainID derives the chain id from the given v parameter
func deriveChainID(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
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
