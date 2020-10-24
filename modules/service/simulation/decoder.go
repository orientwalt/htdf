package simulation

import (
	"bytes"
	"fmt"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/types/kv"

	"github.com/orientwalt/htdf/modules/service/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding service type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ServiceDefinitionKey):
			var definition1, definition2 types.ServiceDefinition
			cdc.MustUnmarshalBinaryBare(kvA.Value, &definition1)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &definition2)
			return fmt.Sprintf("%v\n%v", definition1, definition2)

		default:
			panic(fmt.Sprintf("invalid service key prefix %X", kvA.Key[:1]))
		}
	}
}
