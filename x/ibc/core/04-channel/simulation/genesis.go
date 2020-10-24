package simulation

import (
	"math/rand"

	simtypes "github.com/orientwalt/htdf/types/simulation"
	"github.com/orientwalt/htdf/x/ibc/core/04-channel/types"
)

// GenChannelGenesis returns the default channel genesis state.
func GenChannelGenesis(_ *rand.Rand, _ []simtypes.Account) types.GenesisState {
	return types.DefaultGenesisState()
}
