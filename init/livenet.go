package init

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/orientwalt/htdf/app"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/codec"
	srvconfig "github.com/orientwalt/htdf/server/config"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	authtx "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/orientwalt/htdf/x/mint"
	"github.com/orientwalt/htdf/x/staking"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/rand"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/orientwalt/htdf/accounts/keystore"
	v0 "github.com/orientwalt/htdf/app/v0"
	"github.com/orientwalt/htdf/server"
	hsutils "github.com/orientwalt/htdf/utils"
)

// get cmd to initialize all files for tendermint testnet and application
func LiveNetFilesCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "livenet",
		Short: "Initialize files for a hsd testnet with distributing scripts",
		Long: `livenet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
hsd livenet --chain-id testchain --v 4 -o output --validator-ip-addresses ip.list --minimum-gas-prices 100satoshi --issuer-bech-address htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3 --password-from-file password.list
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			return initLiveNet(config, cdc)
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with",
	)
	cmd.Flags().StringP(flagOutputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet",
	)
	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)",
	)
	cmd.Flags().String(flagNodeDaemonHome, ".hsd",
		"Home directory of the node's daemon configuration",
	)
	cmd.Flags().String(flagNodeCliHome, ".hscli",
		"Home directory of the node's cli configuration",
	)
	cmd.Flags().String(flagStartingIPAddress, "192.168.0.1",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)",
	)
	cmd.Flags().String(flagValidatorIPAddressList, "",
		"Read Validators IP Addresses from an ip address configuration file",
	)
	cmd.Flags().String(
		client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created",
	)
	cmd.Flags().String(
		server.FlagMinGasPrices, fmt.Sprintf("100%s", sdk.DefaultBondDenom),
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 100satoshi)",
	)
	cmd.Flags().String(flagIssuerBechAddress, "", "issuer bech address")
	// cmd.Flags().String(flagStakerBechAddress, "", "staker bech address") // blocked by junying, 2019-08-27, reason: stake doesn't exist anymore
	cmd.Flags().String(flagPasswordFromFile, "", "input password from file")
	return cmd
}

var (
	// validator stake alloc amount
	validatorStakingTokens = sdk.TokensFromTendermintPower(int64(mint.ValidatorProvisions)) // 100(*10**8)
	// issuer allocation amount
	// issuerAccTokens = sdk.TokensFromTendermintPower(int64(mint.UserProvisions)).Sub(validatorStakingTokens.Mul(sdk.NewInt(mint.ValidatorNumbers))) // 59999400(*10**8)
	issuerAccTokens = sdk.TokensFromTendermintPower(int64(mint.CurrentProvisions)) // 61000000(*10**8)

)

func initLiveNet(config *tmconfig.Config, cdc *codec.Codec) error {
	var chainID string
	outDir := viper.GetString(flagOutputDir)
	numValidators := viper.GetInt(flagNumValidators)

	chainID = viper.GetString(client.FlagChainID)
	if chainID == "" {
		chainID = "chain-" + cmn.Str(6)
	}

	monikers := make([]string, numValidators)
	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]crypto.PubKey, numValidators)

	hsConfig := srvconfig.DefaultConfig()
	hsConfig.MinGasPrices = viper.GetString(server.FlagMinGasPrices)

	var (
		accs     []v0.GenesisAccount
		genFiles []string
	)
	// add issuer address
	issuerBechAddr := viper.GetString(flagIssuerBechAddress)
	var err error
	if issuerBechAddr == "" {
		buffer := client.BufferStdin()
		issuerBechAddr, err = client.GetString("Issuer Address: ", buffer)
	}

	issuerAccAddr, err := sdk.AccAddressFromBech32(issuerBechAddr)
	if err != nil {
		return err
	}

	accs = append(accs, v0.GenesisAccount{
		Address: issuerAccAddr,
		Coins: sdk.Coins{
			sdk.NewCoin(DefaultDenom, issuerAccTokens),
		},
	})

	// blocked by junying, 2019-10-16
	// reason: stake doesn't exist anymore
	// add staker address
	// blocked by junying, 2019-09-11
	// reasone: no stake at present
	// stakerBechAddr := viper.GetString(flagStakerBechAddress)
	// if stakerBechAddr == "" {
	// 	buffer := client.BufferStdin()
	// 	stakerBechAddr, err = client.GetString("Staker Address: ", buffer)
	// }

	// blocked by junying, 2019-08-27
	// reason: stake doesn't exist anymore
	// stakerAccAddr, err := sdk.AccAddressFromBech32(stakerBechAddr)
	// if err != nil {
	// 	return err
	// }

	// blocked by junying, 2019-08-27
	// reason: stake doesn't exist anymore
	// accs = append(accs, v0.GenesisAccount{
	// 	Address: stakerAccAddr,
	// 	Coins: sdk.Coins{
	// 		sdk.NewCoin(sdk.DefaultBondDenom, stakerStakingTokens),
	// 	},
	// })

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)
		nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
		nodeCliHomeName := viper.GetString(flagNodeCliHome)
		nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
		clientDir := filepath.Join(outDir, nodeDirName, nodeCliHomeName)
		gentxsDir := filepath.Join(outDir, "gentxs")

		config.SetRoot(nodeDir)

		err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		err = os.MkdirAll(clientDir, nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		monikers = append(monikers, nodeDirName)
		config.Moniker = nodeDirName

		ip := viper.GetString(flagStartingIPAddress)
		ipconfigfile := viper.GetString(flagValidatorIPAddressList)
		if ipconfigfile != "" {
			// buffer := client.BufferStdin()
			// prompt := fmt.Sprintf("IP Address for account '%s':", nodeDirName)
			// ip, err = client.GetString(prompt, buffer)
			ip, _, err = hsutils.ReadString(ipconfigfile, i+1)
		} else {
			ip, err = getIP(i, ip)
		}
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		// write ip
		hsutils.WriteString(filepath.Join(nodeDir, "config/ip.conf"), ip)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		nodeIDs[i], valPubKeys[i], err = InitializeNodeValidatorFiles(config)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		// write node id
		hsutils.WriteString(filepath.Join(nodeDir, "config/node.conf"), nodeIDs[i])
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		memo := ""
		genFiles = append(genFiles, config.GenesisFile())

		buf := client.BufferStdin()
		prompt := fmt.Sprintf(
			"Password for account '%s(%s)' (default %s):", nodeDirName, ip, app.DefaultKeyPass,
		)

		passfile := viper.GetString(flagPasswordFromFile)
		var keyPass string
		if passfile == "" {
			keyPass, err := client.GetPassword(prompt, buf)
			if err != nil && keyPass != "" {
				// An error was returned that either failed to read the password from
				// STDIN or the given password is not empty but failed to meet minimum
				// length requirements.
				_ = os.RemoveAll(outDir)
				return err
			}

			if keyPass == "" {
				keyPass = app.DefaultKeyPass
			}
		} else {
			keyPass, _, err = hsutils.ReadString(passfile, i+1)
			if err != nil {
				_ = os.RemoveAll(outDir)
				return err
			}
			keyPass = strings.TrimSuffix(keyPass, "\n")
		}

		accaddr, secret, err := server.GenerateSaveCoinKeyEx(clientDir, keyPass)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		err = writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint)
		if err != nil {
			return err
		}

		accs = append(accs, v0.GenesisAccount{
			Address: accaddr,
			Coins: sdk.Coins{
				sdk.NewCoin(sdk.DefaultBondDenom, validatorStakingTokens),
			},
		})

		// junying-todo, 2020-01-13
		// commission rate change
		rate, err := sdk.NewDecFromStr("0.1")
		if err != nil {
			return err
		}
		maxRate, err := sdk.NewDecFromStr("0.2")
		if err != nil {
			return err
		}
		maxChangeRrate, err := sdk.NewDecFromStr("0.01")
		if err != nil {
			return err
		}

		msg := staking.NewMsgCreateValidator(
			sdk.ValAddress(accaddr),
			valPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, validatorStakingTokens),
			staking.NewDescription(nodeDirName, "", "", ""),
			staking.NewCommissionMsg(rate, maxRate, maxChangeRrate), // junying-todo, 2020-01-13, commission rate change
			sdk.OneInt(),
		)
		// make unsigned transaction
		unsignedTx := auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{}, memo)
		txBldr := authtx.NewTxBuilderFromCLI().WithChainID(chainID) //.WithMemo(memo)

		addr := sdk.AccAddress.String(accaddr)
		ksw := keystore.NewKeyStoreWallet(filepath.Join(clientDir, "keystores"))
		signedTx, err := ksw.SignStdTx(txBldr, unsignedTx, addr, keyPass)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		txBytes, err := cdc.MarshalJSON(signedTx)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		// gather gentxs folder
		err = writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBytes)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		hsConfigFilePath := filepath.Join(nodeDir, "config/hsd.toml")
		srvconfig.WriteConfigFile(hsConfigFilePath, hsConfig)
	}

	if err := initGenFiles(cdc, chainID, accs, genFiles, numValidators); err != nil {
		return err
	}

	err = collectGenFilesEx(
		cdc, config, chainID, monikers, nodeIDs, valPubKeys, numValidators,
		outDir, viper.GetString(flagNodeDirPrefix), viper.GetString(flagNodeDaemonHome),
	)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully initialized %d node directories\n", numValidators)
	return nil
}

func collectGenFilesEx(
	cdc *codec.Codec, config *tmconfig.Config, chainID string,
	monikers, nodeIDs []string, valPubKeys []crypto.PubKey,
	numValidators int, outDir, nodeDirPrefix, nodeDaemonHomeName string,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
		gentxsDir := filepath.Join(outDir, "gentxs")
		moniker := monikers[i]
		config.Moniker = nodeDirName

		config.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := newInitConfig(chainID, gentxsDir, moniker, nodeID, valPubKey)

		genDoc, err := LoadGenesisDoc(cdc, config.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := genAppStateFromConfigEx(cdc, config, initCfg, genDoc)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}
		genFile := config.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		err = ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime)
		if err != nil {
			return err
		}
	}

	return nil
}
