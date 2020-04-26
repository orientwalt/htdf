package mint

import (
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// Not Register mint codec in app, deprecated now
	//cdc.RegisterConcrete(Minter{}, "htdf/mint/Minter", nil)
	cdc.RegisterConcrete(&Params{}, "htdf/mint/Params", nil)
	cdc.RegisterConcrete(&sdk.Dec{}, "htdf/mint/rewards", nil)
	// cdc.RegisterConcrete(&Params{}, "mint/Params", nil)
	// cdc.RegisterConcrete(&BlockReward{}, "htdf/mint/BlockReward", nil)
	// cdc.RegisterConcrete(&sdk.Dec{}, "types/Dec", nil)
	// cdc.RegisterConcrete(&sdk.Int{}, "types/Int", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
