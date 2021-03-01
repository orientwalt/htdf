package rpc

import (
	"github.com/orientwalt/htdf/client"
	// "github.com/orientwalt/htdf/client/flags"

	"github.com/orientwalt/htdf/client/lcd"
	"github.com/orientwalt/htdf/codec"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// flagUnlockKey = "unlock-key"
	flagWebsocket = "wsport"
)

// EmintServeCmd creates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func EmintServeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := lcd.ServeCommand(cdc, registerRoutes)
	// cmd.Flags().String(flagUnlockKey, "", "Select a key to unlock on the RPC server")
	cmd.Flags().String(flagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().StringP(client.FlagBroadcastMode, "b", client.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")
	return cmd
}

// registerRoutes creates a new server and registers the `/rpc` endpoint.
// Rpc calls are enabled based on their associated module (eg. "eth").
func registerRoutes(rs *lcd.RestServer) {
	s := rpc.NewServer()
	// accountName := viper.GetString(flagUnlockKey)
	// accountNames := strings.Split(accountName, ",")

	// var emintKeys []secp256k1.PrivKeySecp256k1
	// if len(accountName) > 0 {
	// 	var err error
	// 	inBuf := bufio.NewReader(os.Stdin)

	// 	keyringBackend := viper.GetString(client.FlagKeyringBackend)
	// 	passphrase := ""
	// 	switch keyringBackend {
	// 	case keyring.BackendOS:
	// 		break
	// 	case keyring.BackendFile:
	// 		passphrase, err = client.GetPassword(
	// 			"Enter password to unlock key for RPC API: ",
	// 			inBuf)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}

	// 	emintKeys, err = unlockKeyFromNameAndPassphrase(accountNames, passphrase)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	apis := GetRPCAPIs(rs.CliCtx /*, emintKeys*/)

	// TODO: Allow cli to configure modules https://github.com/ChainSafe/ethermint/issues/74
	whitelist := make(map[string]bool)

	// Register all the APIs exposed by the services
	for _, api := range apis {
		if whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := s.RegisterName(api.Namespace, api.Service); err != nil {
				panic(err)
			}
		}
	}

	// Web3 RPC API route
	rs.Mux.HandleFunc("/", s.ServeHTTP).Methods("POST", "OPTIONS")

	// Register all other Cosmos routes
	// client.RegisterRoutes(rs.CliCtx, rs.Mux)
	// authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	// app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)

	// start websockets server
	websocketAddr := viper.GetString(flagWebsocket)
	ws := NewWebsocketsServer(rs.CliCtx, websocketAddr)
	ws.Start()
}

// func unlockKeyFromNameAndPassphrase(accountNames []string, passphrase string) (emintKeys []secp256k1.PrivKeySecp256k1, err error) {
// 	keybase, err := keyring.NewKeyring(
// 		sdk.KeyringServiceName(),
// 		viper.GetString(client.FlagKeyringBackend),
// 		viper.GetString(client.FlagHome),
// 		os.Stdin,
// 	)
// 	if err != nil {
// 		return
// 	}

// 	// try the for loop with array []string accountNames
// 	// run through the bottom code inside the for loop
// 	for _, acc := range accountNames {
// 		// With keyring keybase, password is not required as it is pulled from the OS prompt
// 		privKey, err := keybase.ExportPrivateKeyObject(acc, passphrase)
// 		if err != nil {
// 			return nil, err
// 		}

// 		var ok bool
// 		emintKey, ok := privKey.(secp256k1.PrivKeySecp256k1)
// 		if !ok {
// 			panic(fmt.Sprintf("invalid private key type: %T", privKey))
// 		}
// 		emintKeys = append(emintKeys, emintKey)
// 	}

// 	return
// }
