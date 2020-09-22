package mint

import (
	"os"
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/stretchr/testify/require"

	dbm "github.com/tendermint/tm-db"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	"github.com/orientwalt/htdf/x/bank"
	"github.com/orientwalt/htdf/x/params"
	"github.com/orientwalt/htdf/x/staking"
	stakekeeper "github.com/orientwalt/htdf/x/staking/keeper"
)

type testInput struct {
	ctx        sdk.Context
	cdc        *codec.Codec
	mintKeeper Keeper
}

func createTestCodec() *codec.Codec {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	return cdc
}

func newTestInput(t *testing.T) testInput {
	cdc := createTestCodec()
	db := dbm.NewMemDB()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyFeeCollection := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyMint := sdk.NewKVStoreKey(StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyFeeCollection, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	feeCollectionKeeper := auth.NewFeeCollectionKeeper(cdc, keyFeeCollection)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)

	stakingKeeper := staking.NewKeeper(
		cdc, keyStaking, tkeyStaking, bankKeeper, paramsKeeper.Subspace(staking.DefaultParamspace), staking.DefaultCodespace, stakekeeper.NopMetrics(),
	)

	mintKeeper := NewKeeper(
		cdc, keyMint, paramsKeeper.Subspace(DefaultParamspace), &stakingKeeper, feeCollectionKeeper,
	)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))

	mintKeeper.SetParams(ctx, DefaultParams())
	mintKeeper.SetMinter(ctx, DefaultInitialMinter())

	return testInput{ctx, cdc, mintKeeper}
}
