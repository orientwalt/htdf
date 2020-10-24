package simulation

// DONTCOVER

import (
	"math/rand"

	"github.com/orientwalt/htdf/x/simulation"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{}
}
