package staking

import (
	v1Sdk "github.com/orientwalt/htdf/app/v1/staking/types"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/staking/keeper"
)

// Inflate every block, update inflation parameters once per hour
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	params.UnbondingTime = v1Sdk.DefaultUnbondingTime
	k.SetParams(ctx, params)
}
