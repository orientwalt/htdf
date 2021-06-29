package init

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	common "github.com/tendermint/tendermint/libs/rand"

	"github.com/orientwalt/htdf/app"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/server"
)

const (
	flagOverwrite    = "overwrite"
	flagClientHome   = "home-client"
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

type printInfo struct {
	Moniker    string          `json:"moniker"`
	ChainID    string          `json:"chain_id"`
	NodeID     string          `json:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message"`
}

func displayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s\n", string(out)) // nolint: errcheck
	return nil
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command { // nolint: golint
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.Str(6))
			}
			initialHeight := viper.GetInt64(flagInitialHeight)
			if initialHeight < 0 {
				initialHeight = 0
			}

			dbBackend := viper.GetString(flagDBBackend)
			if(dbBackend == "") {
				dbBackend = "cleveldb"
			}

			nodeID, _, err := InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			var appState json.RawMessage
			genFile := config.GenesisFile()

			if appState, err = initializeEmptyGenesis(cdc, genFile, chainID,
				viper.GetBool(flagOverwrite)); err != nil {
				return err
			}

			if err = ExportGenesisFile(genFile, chainID, initialHeight, nil, appState); err != nil {
				return err
			}

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)
			config.BaseConfig.DBBackend = dbBackend
			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Int64(flagInitialHeight, 0, "genesis block's initial height")
	cmd.Flags().String(flagDBBackend, "cleveldb", "Database backend: goleveldb or cleveldb")

	return cmd
}
