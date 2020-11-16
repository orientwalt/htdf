package subspace

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
)

// Keys for parameter access
const (
	TestParamStore = "ParamsTest"
)

// Returns components for testing
func DefaultTestComponents(t *testing.T) (sdk.Context, Subspace, func() sdk.CommitID) {
	cdc := codec.New()
	key := sdk.NewKVStoreKey(StoreKey)
	tkey := sdk.NewTransientStoreKey(TStoreKey)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.SetTracer(os.Stdout)
	ms.SetTracingContext(sdk.TraceContext{})
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))
	subspace := NewSubspace(cdc, key, tkey, TestParamStore)

	return ctx, subspace, ms.Commit
}
