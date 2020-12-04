package mint

import (
	"fmt"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/params"
)

// mint parameters
type Params struct {
	MintDenom           string  `json:"mint_denom"`            // type of coin to mint
	InflationRateChange sdk.Dec `json:"inflation_rate_change"` // maximum annual change in inflation rate
	InflationMax        sdk.Dec `json:"inflation_max"`         // maximum inflation rate
	InflationMin        sdk.Dec `json:"inflation_min"`         // minimum inflation rate
	GoalBonded          sdk.Dec `json:"goal_bonded"`           // goal of percent bonded atoms
	BlocksPerYear       uint64  `json:"blocks_per_year"`       // expected blocks per year
}

func NewParams(mintDenom string, inflationRateChange, inflationMax,
	inflationMin, goalBonded sdk.Dec, blocksPerYear uint64) Params {

	return Params{
		MintDenom:           mintDenom,
		InflationRateChange: inflationRateChange,
		InflationMax:        inflationMax,
		InflationMin:        inflationMin,
		GoalBonded:          goalBonded,
		BlocksPerYear:       blocksPerYear,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:           sdk.DefaultBondDenom, //junying-todo-07102019//"satoshi", //
		InflationRateChange: sdk.NewDecWithPrec(13, 2),
		InflationMax:        sdk.NewDecWithPrec(20, 2),
		InflationMin:        sdk.NewDecWithPrec(7, 2),
		GoalBonded:          sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
	}
}

func validateParams(params Params) error {
	if params.GoalBonded.LT(sdk.ZeroDec()) {
		return fmt.Errorf("mint parameter GoalBonded should be positive, is %s ", params.GoalBonded.String())
	}
	if params.GoalBonded.GT(sdk.OneDec()) {
		return fmt.Errorf("mint parameter GoalBonded must be <= 1, is %s", params.GoalBonded.String())
	}
	if params.InflationMax.LT(params.InflationMin) {
		return fmt.Errorf("mint parameter Max inflation must be greater than or equal to min inflation")
	}
	if params.MintDenom == "" {
		return fmt.Errorf("mint parameter MintDenom can't be an empty string")
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:             %s
  Inflation Rate Change:  %s
  Inflation Max:          %s
  Inflation Min:          %s
  Goal Bonded:            %s
  Blocks Per Year:        %d
`,
		p.MintDenom, p.InflationRateChange, p.InflationMax,
		p.InflationMin, p.GoalBonded, p.BlocksPerYear,
	)
}

// Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamspace
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: ParamStoreKeyParams, Value: &p.InflationRateChange},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(ParamStoreKeyParams):
		inflation, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if inflation.GT(sdk.NewDecWithPrec(2, 1)) || inflation.LT(sdk.NewDecWithPrec(0, 0)) {
			return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidMintInflation, fmt.Sprintf("Mint Inflation [%s] should be between [0, 0.2] ", value))
		}
		return inflation, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(ParamStoreKeyParams):
		err := cdc.UnmarshalJSON(bytes, &p.InflationRateChange)
		return p.InflationRateChange.String(), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

//______________________________________________________________________

// get inflation params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) Params {
	var params Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set inflation params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
