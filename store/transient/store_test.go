package transient

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	// "github.com/orientwalt/htdf/store/transient"
	"github.com/orientwalt/htdf/store/types"
)

var k, v = []byte("hello"), []byte("world")

func TestTransientStore(t *testing.T) {
	tstore := NewStore()

	require.Nil(t, tstore.Get(k))

	tstore.Set(k, v)

	require.Equal(t, v, tstore.Get(k))

	tstore.Commit()

	require.Nil(t, tstore.Get(k))

	// no-op
	tstore.SetPruning(types.PruningOptions{})

	emptyCommitID := tstore.LastCommitID()
	require.Equal(t, emptyCommitID.Version, int64(0))
	require.True(t, bytes.Equal(emptyCommitID.Hash, nil))
	require.Equal(t, types.StoreTypeTransient, tstore.GetStoreType())
}
