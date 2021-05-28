package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	common "github.com/tendermint/tendermint/libs/os"

	"github.com/orientwalt/htdf/app"
	v0 "github.com/orientwalt/htdf/app/v0"
	"github.com/orientwalt/htdf/client/keys"
	"github.com/orientwalt/htdf/codec"
	i "github.com/orientwalt/htdf/init"
	"github.com/orientwalt/htdf/server"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/guardian"
)

const (
	flagOverwrite  = "overwrite"
	flagClientHome = "home-client"
)

// AddGuardianAccountCmd returns add-guardian-account cobra Command.
func AddGuardianAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-guardian-account [address]",
		Short: "Add guardian account to genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
				if err != nil {
					return err
				}

				info, err := kb.Get(args[0])
				if err != nil {
					return err
				}

				addr = info.GetAddress()
			}

			genFile := config.GenesisFile()
			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `hsd init` first", genFile)
			}

			genDoc, err := i.LoadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState v0.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			appState, err = addGenesisAccount(cdc, appState, addr)
			if err != nil {
				return err
			}

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			return i.ExportGenesisFile(genFile, genDoc.ChainID, genDoc.InitialHeight, nil, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, app.DefaultCLIHome, "client's home directory")

	return cmd
}

func addGenesisAccount(
	cdc *codec.Codec, appState v0.GenesisState, addr sdk.AccAddress,
) (v0.GenesisState, error) {
	var genAcc sdk.AccAddress
	for _, stateAcc := range appState.GuardianData.Profilers {
		if stateAcc.Address.Equals(addr) {
			return appState, fmt.Errorf("the application state already contains account %v", addr)
		}
		genAcc = stateAcc.Address
	}

	guardian := guardian.NewGuardian("genesis", guardian.Genesis, addr, addr)

	if genAcc.Empty() {
		appState.GuardianData.Profilers[0] = guardian
		appState.GuardianData.Trustees[0] = guardian
	} else {
		appState.GuardianData.Profilers = append(appState.GuardianData.Profilers, guardian)
		appState.GuardianData.Trustees = append(appState.GuardianData.Trustees, guardian)
	}

	return appState, nil
}
