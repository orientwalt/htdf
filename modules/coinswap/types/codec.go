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
}

// RegisterLegacyAminoCodec registers concrete types on the codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSwapOrder{}, "irismod/coinswap/MsgSwapOrder", nil)
	cdc.RegisterConcrete(&MsgAddLiquidity{}, "irismod/coinswap/MsgAddLiquidity", nil)
	cdc.RegisterConcrete(&MsgRemoveLiquidity{}, "irismod/coinswap/MsgRemoveLiquidity", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSwapOrder{},
		&MsgAddLiquidity{},
		&MsgRemoveLiquidity{},
	)
}
