package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	authtypes "github.com/orientwalt/htdf/x/auth/types"

	"github.com/orientwalt/htdf/modules/htlc/types"
)

// Keeper defines the HTLC keeper
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.Marshaler
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

// NewKeeper creates a new HTLC Keeper instance
func NewKeeper(
	cdc codec.Marshaler,
	key sdk.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	// ensure the HTLC module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// GetHTLCAccount returns the HTLC module account
func (k Keeper) GetHTLCAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
