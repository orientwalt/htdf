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

// RegisterLegacyAminoCodec registers the necessary x/bank interfaces and concrete types
// on the provided Amino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateFeed{}, "irismod/oracle/MsgCreateFeed", nil)
	cdc.RegisterConcrete(&MsgStartFeed{}, "irismod/oracle/MsgStartFeed", nil)
	cdc.RegisterConcrete(&MsgPauseFeed{}, "irismod/oracle/MsgPauseFeed", nil)
	cdc.RegisterConcrete(&MsgEditFeed{}, "irismod/oracle/MsgEditFeed", nil)

	cdc.RegisterConcrete(&Feed{}, "irismod/oracle/Feed", nil)
	cdc.RegisterConcrete(&FeedContext{}, "irismod/oracle/FeedContext", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateFeed{},
		&MsgStartFeed{},
		&MsgPauseFeed{},
		&MsgEditFeed{},
	)
}
