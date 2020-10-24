package types

import (
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/codec/types"
	cryptocodec "github.com/orientwalt/htdf/crypto/codec"
	sdk "github.com/orientwalt/htdf/types"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterLegacyAminoCodec registers concrete types on the codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRequestRandom{}, "irismod/random/MsgRequestRandom", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestRandom{},
	)
}
