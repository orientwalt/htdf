package store

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/orientwalt/htdf/store/cache"
	"github.com/orientwalt/htdf/store/rootmulti"
	"github.com/orientwalt/htdf/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewCommitKVStoreCacheManager() types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize)
}
