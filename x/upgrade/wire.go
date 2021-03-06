package upgrade

import (
	"github.com/orientwalt/htdf/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&VersionInfo{}, "htdf/upgrade/VersionInfo", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
