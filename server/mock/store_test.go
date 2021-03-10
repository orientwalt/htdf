package mock

import (
	"testing"

	"github.com/stretchr/testify/require"

	dbm "github.com/tendermint/tm-db"

	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
)

func TestStore(t *testing.T) {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	key := sdk.NewKVStoreKey("test")
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	require.Nil(t, err)

	store := cms.GetKVStore(key)
	require.NotNil(t, store)

	k := []byte("hello")
	v := []byte("world")
	require.False(t, store.Has(k))
	store.Set(k, v)
	require.True(t, store.Has(k))
	require.Equal(t, v, store.Get(k))
	store.Delete(k)
	require.False(t, store.Has(k))
}
