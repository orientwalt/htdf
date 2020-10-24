package std

import (
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/codec/types"
	cryptocodec "github.com/orientwalt/htdf/crypto/codec"
	sdk "github.com/orientwalt/htdf/types"
	txtypes "github.com/orientwalt/htdf/types/tx"
	vesting "github.com/orientwalt/htdf/x/auth/vesting/types"
)

// RegisterLegacyAminoCodec registers types with the Amino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)
	vesting.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers Interfaces from sdk/types, vesting, crypto, tx.
func RegisterInterfaces(interfaceRegistry types.InterfaceRegistry) {
	sdk.RegisterInterfaces(interfaceRegistry)
	txtypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	vesting.RegisterInterfaces(interfaceRegistry)
}
