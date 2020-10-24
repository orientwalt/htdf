package v039

import (
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/codec"
	cryptocodec "github.com/orientwalt/htdf/crypto/codec"
	v038auth "github.com/orientwalt/htdf/x/auth/legacy/v038"
	v039auth "github.com/orientwalt/htdf/x/auth/legacy/v039"
	"github.com/orientwalt/htdf/x/genutil/types"
)

// Migrate migrates exported state from v0.38 to a v0.39 genesis state.
//
// NOTE: No actual migration occurs since the types do not change, but JSON
// serialization of accounts do change.
func Migrate(appState types.AppMap, _ client.Context) types.AppMap {
	v038Codec := codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(v038Codec)
	v038auth.RegisterLegacyAminoCodec(v038Codec)

	v039Codec := codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(v039Codec)
	v039auth.RegisterLegacyAminoCodec(v039Codec)

	// migrate x/auth state (JSON serialization only)
	if appState[v038auth.ModuleName] != nil {
		var authGenState v038auth.GenesisState
		v038Codec.MustUnmarshalJSON(appState[v038auth.ModuleName], &authGenState)

		delete(appState, v038auth.ModuleName) // delete old key in case the name changed
		appState[v039auth.ModuleName] = v039Codec.MustMarshalJSON(v039auth.Migrate(authGenState))
	}

	return appState
}
