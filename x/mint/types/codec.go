package types

import (
	"github.com/orientwalt/htdf/codec"
	cryptocodec "github.com/orientwalt/htdf/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
