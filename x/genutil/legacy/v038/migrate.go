package v038

import (
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/codec"
	cryptocodec "github.com/orientwalt/htdf/crypto/codec"
	v036auth "github.com/orientwalt/htdf/x/auth/legacy/v036"
	v038auth "github.com/orientwalt/htdf/x/auth/legacy/v038"
	v036distr "github.com/orientwalt/htdf/x/distribution/legacy/v036"
	v038distr "github.com/orientwalt/htdf/x/distribution/legacy/v038"
	v036genaccounts "github.com/orientwalt/htdf/x/genaccounts/legacy/v036"
	"github.com/orientwalt/htdf/x/genutil/types"
	v036staking "github.com/orientwalt/htdf/x/staking/legacy/v036"
	v038staking "github.com/orientwalt/htdf/x/staking/legacy/v038"
)

// Migrate migrates exported state from v0.36/v0.37 to a v0.38 genesis state.
func Migrate(appState types.AppMap, _ client.Context) types.AppMap {
	v036Codec := codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(v036Codec)

	v038Codec := codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(v038Codec)
	v038auth.RegisterLegacyAminoCodec(v038Codec)

	if appState[v036genaccounts.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var authGenState v036auth.GenesisState
		v036Codec.MustUnmarshalJSON(appState[v036auth.ModuleName], &authGenState)

		var genAccountsGenState v036genaccounts.GenesisState
		v036Codec.MustUnmarshalJSON(appState[v036genaccounts.ModuleName], &genAccountsGenState)

		// delete deprecated genaccounts genesis state
		delete(appState, v036genaccounts.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v038auth.ModuleName] = v038Codec.MustMarshalJSON(v038auth.Migrate(authGenState, genAccountsGenState))
	}

	// migrate staking state
	if appState[v036staking.ModuleName] != nil {
		var stakingGenState v036staking.GenesisState
		v036Codec.MustUnmarshalJSON(appState[v036staking.ModuleName], &stakingGenState)

		delete(appState, v036staking.ModuleName) // delete old key in case the name changed
		appState[v038staking.ModuleName] = v038Codec.MustMarshalJSON(v038staking.Migrate(stakingGenState))
	}

	// migrate distribution state
	if appState[v036distr.ModuleName] != nil {
		var distrGenState v036distr.GenesisState
		v036Codec.MustUnmarshalJSON(appState[v036distr.ModuleName], &distrGenState)

		delete(appState, v036distr.ModuleName) // delete old key in case the name changed
		appState[v038distr.ModuleName] = v038Codec.MustMarshalJSON(v038distr.Migrate(distrGenState))
	}

	return appState
}
