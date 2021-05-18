package init

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/orientwalt/htdf/app"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/codec"
	srvconfig "github.com/orientwalt/htdf/server/config"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	authtx "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/orientwalt/htdf/x/guardian"
	"github.com/orientwalt/htdf/x/staking"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/os"
	rand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/orientwalt/htdf/accounts/keystore"
	v0 "github.com/orientwalt/htdf/app/v0"
	"github.com/orientwalt/htdf/server"
	hsutils "github.com/orientwalt/htdf/utils"
)

var (
	flagNodeDirPrefix          = "node-dir-prefix"
	flagNumValidators          = "v"
	flagOutputDir              = "output-dir"
	flagNodeDaemonHome         = "node-daemon-home"
	flagNodeCliHome            = "node-cli-home"
	flagStartingIPAddress      = "starting-ip-address"
	flagValidatorIPAddressList = "validator-ip-addresses"
	flagIssuerBechAddress      = "issuer-bech-address"
	flagAccountsFilePath       = "accounts-file-path"
	flagStakerBechAddress      = "staker-bech-address"
	flagPasswordFromFile       = "password-from-file"
)

const (
	//
	nodeDirPerm = 0755
	//
	DefaultDenom = sdk.DefaultDenom
)

// get cmd to initialize all files for tendermint testnet and application
func TestnetFilesCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a hsd testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	hsd testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			return initTestnet(config, cdc)
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
		server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)",
	)

	return cmd
}

func initTestnet(config *tmconfig.Config, cdc *codec.Codec) error {
	var chainID string
	outDir := viper.GetString(flagOutputDir)
	numValidators := viper.GetInt(flagNumValidators)

	chainID = viper.GetString(client.FlagChainID)
	if chainID == "" {
		chainID = "chain-" + rand.Str(6)
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
		ipconf := viper.GetString(flagValidatorIPAddressList)
		if ipconf != "" {
			// buffer := client.BufferStdin()
			// prompt := fmt.Sprintf("IP Address for account '%s':", nodeDirName)
			// ip, err = client.GetString(prompt, buffer)
			ip, _, err = hsutils.ReadString(ipconf, i+1)
		} else {
			ip, err = getIP(i, ip)
		}
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}
		nodeIDs[i], valPubKeys[i], err = InitializeNodeValidatorFiles(config)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, config.GenesisFile())

		buf := client.BufferStdin()
		prompt := fmt.Sprintf(
			"Password for account '%s(%s)' (default %s):", nodeDirName, ip, app.DefaultKeyPass,
		)

		keyPass, err := client.GetPassword(prompt, buf)
		if err != nil && keyPass != "" {
			// An error was returned that either failed to read the password from
			// STDIN or the given password is not empty but failed to meet minimum
			// length requirements.
			return err
		}

		if keyPass == "" {
			keyPass = app.DefaultKeyPass
		}
		fmt.Println(clientDir)
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

		accTokens := sdk.TokensFromTendermintPower(30000000)
		// accStakingTokens := sdk.TokensFromTendermintPower(250000000)
		accs = append(accs, v0.GenesisAccount{
			Address: accaddr,
			Coins: sdk.Coins{
				sdk.NewCoin(DefaultDenom, accTokens),
				// sdk.NewCoin(sdk.DefaultBondDenom, accStakingTokens), //
			},
		})

		// commission rate change
		rate, _ := sdk.NewDecFromStr("0.1")
		maxRate, _  := sdk.NewDecFromStr("0.2")
		maxChangeRrate, _ := sdk.NewDecFromStr("0.01")

		valTokens := sdk.TokensFromTendermintPower(1)
		msg := staking.NewMsgCreateValidator(
			sdk.ValAddress(accaddr),
			valPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			staking.NewDescription(nodeDirName, "", "", ""),

			// 2021-05-18, yqq, 
			// keep same commission configuration as mainnet
			// staking.NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			staking.NewCommissionMsg(rate, maxRate, maxChangeRrate),
			sdk.OneInt(),
		)
		// make unsigned transaction
		unsignedTx := auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{}, memo)
		txBldr := authtx.NewTxBuilderFromCLI().WithChainID(chainID).WithMemo(memo)

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

	// yqq, 2021-04-27 , we set accs[0] as default guardian
	defaultGuardian := accs[0].Address
	if err := initGenFiles(cdc, chainID, accs, genFiles, numValidators, defaultGuardian); err != nil {
		return err
	}

	err := collectGenFiles(
		cdc, config, chainID, monikers, nodeIDs, valPubKeys, numValidators,
		outDir, viper.GetString(flagNodeDirPrefix), viper.GetString(flagNodeDaemonHome),
	)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully initialized %d node directories\n", numValidators)
	return nil
}

//
func initGenFiles(
	cdc *codec.Codec, chainID string, accs []v0.GenesisAccount,
	genFiles []string, numValidators int, guadianAddr sdk.AccAddress,
) error {

	appGenState := v0.NewDefaultGenesisState()
	appGenState.Accounts = accs
	gd := guardian.NewGuardian("genesis", guardian.Genesis, guadianAddr, guadianAddr)
	appGenState.GuardianData.Profilers[0] = gd
	appGenState.GuardianData.Trustees[0] = gd

	appGenStateJSON, err := codec.MarshalJSONIndent(cdc, appGenState)
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

func collectGenFiles(
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

		nodeAppState, err := genAppStateFromConfig(cdc, config, initCfg, genDoc)
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

func getIP(i int, startingIPAddr string) (string, error) {
	var (
		ip  string
		err error
	)

	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
	} else {
		ip, err = calculateIP(startingIPAddr, i)
		if err != nil {
			return "", err
		}
	}

	return ip, nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := cmn.EnsureDir(writePath, 0700) //0700-junying-todo-20190422
	if err != nil {
		return err
	}

	err = cmn.WriteFile(file, contents, 0600) //0600-junying-todo-20190422
	if err != nil {
		return err
	}

	return nil
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

// func addGenesisAccount(
// 	cdc *codec.Codec, appState v0.GenesisState, addr sdk.AccAddress,
// ) (v0.GenesisState, error) {
// 	var genAcc sdk.AccAddress
// 	for _, stateAcc := range appState.GuardianData.Profilers {
// 		if stateAcc.Address.Equals(addr) {
// 			return appState, fmt.Errorf("the application state already contains account %v", addr)
// 		}
// 		genAcc = stateAcc.Address
// 	}

// 	guardian := guardian.NewGuardian("genesis", guardian.Genesis, addr, addr)

// 	if genAcc.Empty() {
// 		appState.GuardianData.Profilers[0] = guardian
// 		appState.GuardianData.Trustees[0] = guardian
// 	} else {
// 		appState.GuardianData.Profilers = append(appState.GuardianData.Profilers, guardian)
// 		appState.GuardianData.Trustees = append(appState.GuardianData.Trustees, guardian)
// 	}

// 	return appState, nil
// }
