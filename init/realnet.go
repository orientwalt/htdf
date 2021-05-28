package init

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/orientwalt/htdf/app"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/codec"
	srvconfig "github.com/orientwalt/htdf/server/config"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	authtx "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/orientwalt/htdf/x/staking"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/rand"

	"github.com/orientwalt/htdf/accounts/keystore"
	v0 "github.com/orientwalt/htdf/app/v0"
	"github.com/orientwalt/htdf/server"
	hsutils "github.com/orientwalt/htdf/utils"
)

// get cmd to initialize all files for tendermint testnet and application
func RealNetFilesCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "realnet",
		Short: "Initialize files for a hsd mainnet with exported accounts",
		Long: `realnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	hsd realnet --chain-id testchain --v 4 -o output --validator-ip-addresses ip.list --minimum-gas-prices 100satoshi --accounts-file-path accounts.list --password-from-file password.list
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			return initRealNet(config, cdc)
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with",
	)
	cmd.Flags().Int64(flagInitialHeight, 0,
		"Genesis Block's Initial Height",
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
	cmd.Flags().String(flagAccountsFilePath, "", "issuer bech address")
	cmd.Flags().String(flagPasswordFromFile, "", "input password from file")
	return cmd
}

func initRealNet(config *tmconfig.Config, cdc *codec.Codec) error {
	var chainID string
	outDir := viper.GetString(flagOutputDir)
	numValidators := viper.GetInt(flagNumValidators)

	chainID = viper.GetString(client.FlagChainID)
	if chainID == "" {
		chainID = "chain-" + cmn.Str(6)
	}
	// junying-added, initial-height
	initialHeight := viper.GetInt64(flagInitialHeight)
	if initialHeight < 0 {
		initialHeight = 0
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
	// read accounts from account.list
	accFilePath := viper.GetString(flagAccountsFilePath)
	accounts, balances, err := hsutils.ReadAccounts(accFilePath)

	for index, acc := range accounts {
		issuerAccAddr, err := sdk.AccAddressFromBech32(acc)
		if err != nil {
			return err
		}

		balance, ok := sdk.NewIntFromString(strconv.Itoa(balances[index]))
		if !ok {
			continue
		}
		accs = append(accs, v0.GenesisAccount{
			Address: issuerAccAddr,
			Coins: sdk.Coins{
				sdk.NewCoin(DefaultDenom, balance),
			},
		})
	}

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

		valTokens := validatorStakingTokens //sdk.TokensFromTendermintPower(100)
		msg := staking.NewMsgCreateValidator(
			sdk.ValAddress(accaddr),
			valPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			staking.NewDescription(nodeDirName, "", "", ""),
			staking.NewCommissionMsg(rate, maxRate, maxChangeRrate), // junying-todo, 2020-01-13, commission rate change
			sdk.OneInt(),
		)
		// make unsigned transaction
		unsignedTx := auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{}, memo)
		txBldr := authtx.NewTxBuilderFromCLI().WithChainID(chainID) //.WithMemo(memo)

		addr := sdk.AccAddress.String(accaddr)
		defaultKeyStoreHome := filepath.Join(clientDir, "keystores")
		ksw := keystore.NewKeyStoreWallet(defaultKeyStoreHome)
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

	// yqq, 2021-04-27 , we set accs[0] as default guardian
	defaultGuardian := accs[0].Address
	if err := initGenFiles(cdc, chainID, accs, genFiles, numValidators, initialHeight, defaultGuardian); err != nil {
		return err
	}

	err = collectGenFilesEx(
		cdc, config, chainID, initialHeight, monikers, nodeIDs, valPubKeys, numValidators,
		outDir, viper.GetString(flagNodeDirPrefix), viper.GetString(flagNodeDaemonHome),
	)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully initialized %d node directories\n", numValidators)
	return nil
}
