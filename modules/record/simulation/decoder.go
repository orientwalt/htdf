package simulation

import (
	"bytes"
	"fmt"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/types/kv"

	"github.com/orientwalt/htdf/modules/record/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding slashing type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.RecordKey):
			var recordA, recordB types.Record
			cdc.MustUnmarshalBinaryBare(kvA.Value, &recordA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &recordB)
			return fmt.Sprintf("%v\n%v", recordA, recordB)
		default:
			panic(fmt.Sprintf("invalid record key prefix %X", kvA.Key[:1]))
		}
	}
}
