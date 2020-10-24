package simulation

import (
	"math/rand"

	simtypes "github.com/orientwalt/htdf/types/simulation"
	"github.com/orientwalt/htdf/x/simulation"

	"github.com/orientwalt/htdf/modules/token/types"
)

const (
	keyTokenTaxRate      = "TokenTaxRate"
	keyIssueTokenBaseFee = "IssueTokenBaseFee"
	keyMintTokenFeeRatio = "MintTokenFeeRatio"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyTokenTaxRate,
			func(r *rand.Rand) string {
				return RandomDec(r).String()
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyIssueTokenBaseFee,
			func(r *rand.Rand) string {
				return RandomInt(r).String()
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMintTokenFeeRatio,
			func(r *rand.Rand) string {
				return RandomDec(r).String()
			},
		),
	}
}
