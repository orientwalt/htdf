package cli

import (
	"fmt"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/keys"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	evmtypes "github.com/orientwalt/htdf/x/evm/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// junying-todo-20190327
// GetCmdSign is the CLI command for signing unsigned transaction
/*
	inspired by
	hscli tx sign unsigned.json --name junying >> signed.json
	hscli tx sign --validate-signatures signed.json
	hscli tx sign --signature-only  test.json --name junying
*/
func GetCmdSign(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign [unsignedtransaction]",
		Short: "sign a transaction",
		Long:  "hscli tx sign 7b0a202...23 --sequence 1 --account-number 0 --offline=true --encode=false",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			// load sign tx from string
			stdTx, err := evmtypes.ReadStdTxFromRawData(cliCtx.Codec, args[0])
			if err != nil {
				return err
			}

			// if no signers
			if len(stdTx.GetSigners()) == 0 {
				return err //err.
			}

			passphrase, err := keys.ReadShortPassphraseFromStdin(sdk.AccAddress.String(stdTx.GetSigners()[0]))
			if err != nil {
				return err
			}

			offlineflag := viper.GetBool(evmtypes.FlagOffline)

			// sign
			res, err := SignTransaction(authtxb.NewTxBuilderFromCLI(), cliCtx, stdTx, passphrase, offlineflag)
			if err != nil {
				return err
			}

			// print
			encodeflag := viper.GetBool(evmtypes.FlagEncode)
			if !encodeflag {
				fmt.Printf("%s\n", res)
			} else {
				fmt.Printf("%s\n", evmtypes.Encode_Hex(res))
			}
			return nil
		},
	}
	cmd.Flags().Bool(evmtypes.FlagEncode, true, "encode enabled")
	cmd.Flags().Bool(evmtypes.FlagOffline, false, "offline disabled")
	return client.PostCommands(cmd)[0]
}

func populateAccountFromState(txBldr authtxb.TxBuilder, cliCtx context.CLIContext,
	addr sdk.AccAddress) (authtxb.TxBuilder, error) {
	if txBldr.AccountNumber() == 0 {
		accNum, err := cliCtx.GetAccountNumber(addr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithAccountNumber(accNum)
	}

	if txBldr.Sequence() == 0 {
		accSeq, err := cliCtx.GetAccountSequence(addr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithSequence(accSeq)
	}
	return txBldr, nil
}

//
func SignStdTx(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, stdTx auth.StdTx, passphrase string, offline bool) (signedTx auth.StdTx, err error) {
	// from address
	if len(stdTx.GetSigners()) == 0 {
		return signedTx, nil
	}
	fromaddr := stdTx.GetSigners()[0]
	// accountnumber, accountsequence check
	if !offline {
		txBldr, err = populateAccountFromState(txBldr, cliCtx, fromaddr)
		if err != nil {
			return signedTx, err
		}
	}

	ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())

	// signature
	return ksw.SignStdTx(txBldr, stdTx, sdk.AccAddress.String(fromaddr), passphrase)
}

//
func SignTransaction(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, stdTx auth.StdTx, passphrase string, offline bool) (res []byte, err error) {
	// signature
	signedTx, err := SignStdTx(txBldr, cliCtx, stdTx, passphrase, offline)
	if err != nil {
		return []byte("signing failed"), err
	}

	switch cliCtx.Indent {
	case true:
		res, err = cliCtx.Codec.MarshalJSONIndent(signedTx, "", "  ")
	default:
		res, err = cliCtx.Codec.MarshalJSON(signedTx)
	}

	if err != nil {
		return []byte("json creating failed"), err
	}
	return res, err
}
