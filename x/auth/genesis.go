package auth

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	CollectedFees sdk.Coins `json:"collected_fees"`
	Params        Params    `json:"params"`
}


// GenesisStateEx only use for exporting data to htdf2.0
type GenesisStateEx struct {
	CollectedFees sdk.Coins `json:"collected_fees"`
	Params        ParamsEx    `json:"params"`
}


// NewGenesisState - Create a new genesis state
func NewGenesisState(collectedFees sdk.Coins, params Params) GenesisState {
	return GenesisState{
		Params:        params,
		CollectedFees: collectedFees,
	}
}

// NewGenesisStateEx - Create a new genesis state
func NewGenesisStateEx(collectedFees sdk.Coins, params Params) GenesisStateEx {
	newParams := ParamsEx{}
	newParams.GasPriceThreshold = params.GasPriceThreshold
	newParams.InitialHeight = 1 // TODO: yqq--------
	newParams.MaxMemoCharacters = params.MaxMemoCharacters
	newParams.SigVerifyCostED25519  = params.SigVerifyCostED25519
	newParams.SigVerifyCostSecp256k1 = params.SigVerifyCostSecp256k1
	newParams.TxSigLimit = params.TxSigLimit
	newParams.TxSizeCostPerByte = params.TxSizeCostPerByte

	return GenesisStateEx{
		Params:        newParams,
		CollectedFees: collectedFees,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(sdk.NewCoins(), DefaultParams())
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, ak AccountKeeper, fck FeeCollectionKeeper, data GenesisState) {
	ak.SetParams(ctx, data.Params)
	fck.setCollectedFees(ctx, data.CollectedFees)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak AccountKeeper, fck FeeCollectionKeeper) GenesisState {
	collectedFees := fck.GetCollectedFees(ctx)
	params := ak.GetParams(ctx)

	return NewGenesisState(collectedFees, params)
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.Params.TxSigLimit == 0 {
		return fmt.Errorf("invalid tx signature limit: %d", data.Params.TxSigLimit)
	}
	if data.Params.SigVerifyCostED25519 == 0 {
		return fmt.Errorf("invalid ED25519 signature verification cost: %d", data.Params.SigVerifyCostED25519)
	}
	if data.Params.SigVerifyCostSecp256k1 == 0 {
		return fmt.Errorf("invalid SECK256k1 signature verification cost: %d", data.Params.SigVerifyCostSecp256k1)
	}
	if data.Params.MaxMemoCharacters == 0 {
		return fmt.Errorf("invalid max memo characters: %d", data.Params.MaxMemoCharacters)
	}
	if data.Params.TxSizeCostPerByte == 0 {
		return fmt.Errorf("invalid tx size cost per byte: %d", data.Params.TxSizeCostPerByte)
	}
	return nil
}
