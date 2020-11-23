package slashing

import (
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/x/slashing/types"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(types.MsgUnjail{}, "htdf/MsgUnjail", nil)
}

var cdcEmpty = codec.New()
