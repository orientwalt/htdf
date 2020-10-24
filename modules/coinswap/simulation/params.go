package simulation

import (
	"math/rand"

	simtypes "github.com/orientwalt/htdf/types/simulation"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return nil
}
