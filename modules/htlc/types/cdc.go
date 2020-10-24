package types

import (
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/codec/types"
	cryptocodec "github.com/orientwalt/htdf/crypto/codec"
	sdk "github.com/orientwalt/htdf/types"
)

// ModuleCdc defines the module codec
var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateHTLC{}, "irismod/htlc/MsgCreateHTLC", nil)
	cdc.RegisterConcrete(&MsgClaimHTLC{}, "irismod/htlc/MsgClaimHTLC", nil)
	cdc.RegisterConcrete(&MsgRefundHTLC{}, "irismod/htlc/MsgRefundHTLC", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateHTLC{},
		&MsgClaimHTLC{},
		&MsgRefundHTLC{},
	)
}
