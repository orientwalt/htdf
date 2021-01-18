package types

import (
	"github.com/orientwalt/htdf/codec"
)

var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "htdfservice/send", nil)
	// cdc.RegisterConcrete(MsgAdd{}, "htdfservice/add", nil)
}

// Evm module events
const (
	EventTypeMsgSend = TypeMsgSend

	AttributeKeyContractAddress = "contract"
	AttributeKeyRecipient       = "recipient"
	AttributeValueCategory      = ModuleName
)
