package types

import (
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/codec/types"
)

// RegisterLegacyAminoCodec registers the sdk message type.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.base.v1beta1.Msg", (*Msg)(nil))
}
