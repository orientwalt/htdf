package transient

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/orientwalt/htdf/store/dbadapter"
	"github.com/orientwalt/htdf/store/types"
)

var _ types.Committer = (*Store)(nil)
var _ types.KVStore = (*Store)(nil)

// Store is a wrapper for a MemDB with Commiter implementation
type Store struct {
	dbadapter.Store
}

// Constructs new MemDB adapter
func NewStore() *Store {
	return &Store{Store: dbadapter.Store{DB: dbm.NewMemDB()}}
}

// Implements CommitStore
// Commit cleans up Store.
func (ts *Store) Commit() (id types.CommitID) {
	ts.Store = dbadapter.Store{DB: dbm.NewMemDB()}
	return
}

// Implements CommitStore
func (ts *Store) SetPruning(pruning types.PruningOptions) {
}

// Implements CommitStore
func (ts *Store) LastCommitID() (id types.CommitID) {
	return
}

func (ts *Store) SetLastCommitID(int64) {}

// Implements Store.
func (ts *Store) GetStoreType() types.StoreType {
	return types.StoreTypeTransient
}
