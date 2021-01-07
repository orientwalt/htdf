package transient

import (
	dbm "github.com/tendermint/tendermint/libs/db"

	"github.com/orientwalt/htdf/store/types"

	"github.com/orientwalt/htdf/store/dbadapter"
)

var _ types.Committer = (*Store)(nil)
var _ types.KVStore = (*Store)(nil)

// Store is a wrapper for a MemDB with Commiter implementation
type Store struct {
	dbadapter.Store
}

// Constructs new MemDB adapter
func NewStore() *Store {
	return &Store{dbadapter.Store{dbm.NewMemDB()}}
}

// Implements CommitStore
// Commit cleans up Store.
func (ts *Store) Commit([]*types.KVStoreKey) (id types.CommitID) {
	ts.Store = dbadapter.Store{dbm.NewMemDB()}
	return
}

// Implements Committer/CommitStore.
func (ts *Store) CommitWithVersion(KVStoreList []*types.KVStoreKey, _ int64) types.CommitID {
	return ts.Commit(KVStoreList)
}

// Implements CommitStore
func (ts *Store) SetPruning(pruning types.PruningOptions) {
}

// Implements CommitStore
func (ts *Store) LastCommitID() (id types.CommitID) {
	return
}

// Implements Store.
func (ts *Store) GetStoreType() types.StoreType {
	return types.StoreTypeTransient
}
