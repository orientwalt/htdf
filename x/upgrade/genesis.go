package upgrade

import (
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/version"
)
const defaultProtocolVersion = version.ProtocolVersion


// GenesisState - all upgrade state that must be provided at genesis
type GenesisState struct {
	GenesisVersion VersionInfo //`json:"genesis_version"`
}

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	genesisVersion := data.GenesisVersion

	k.AddNewVersionInfo(ctx, genesisVersion)
	k.protocolKeeper.ClearUpgradeConfig(ctx)
	k.protocolKeeper.SetCurrentVersion(ctx, genesisVersion.UpgradeInfo.Protocol.Version)
}

// WriteGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context) GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.DefaultUpgradeConfig(defaultProtocolVersion, "https://github.com/orientwalt/htdf/releases/tag/v"+version.Version), true),
	}
}

// ExportGenesisEx  only use for exporting data to htdf2.0
func ExportGenesisEx(ctx sdk.Context) GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.DefaultUpgradeConfig(0, "https://github.com/orientwalt/htdf/releases/tag/v"+version.Version), true),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.DefaultUpgradeConfig(defaultProtocolVersion, "https://github.com/orientwalt/htdf/releases/tag/v"+version.Version), true),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.DefaultUpgradeConfig(defaultProtocolVersion, "https://github.com/orientwalt/htdf/releases/tag/v"+version.Version), true),
	}
}
