package slashing

import (
	"github.com/orientwalt/htdf/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "htdf/MsgUnjail", nil)
}

var cdcEmpty = codec.New()
