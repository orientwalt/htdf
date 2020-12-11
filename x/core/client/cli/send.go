package cli

import (
	"fmt"
	"os"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/keys"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	htdfservice "github.com/orientwalt/htdf/x/core"
	"github.com/spf13/cobra"
)

// junying-todo-20190325
// GetCmdSend is the CLI command for sending a Send transaction
/*
	inspired by
	hscli send cosmos1yqgv2rhxcgrf5jqrxlg80at5szzlarlcy254re 5htdftoken --from junying
*/
func GetCmdSend(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [fromaddr] [toaddr] [amount]",
		Short: "create & send transaction",
		Long: `hscli tx send htdf1qn38r8re3lwlf5t6zgrdycrerd5w0 \
							 htdf1yujjc5yptpphtt665u2u6zp6gl04enlg55fajp \
							 5satoshi \
							 --gas=30000 \
							 --gas-price=100`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			fmt.Println("GetCmdSend:txBldr.GasWanted()", txBldr.GasWanted())

			fromaddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			toaddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			if txBldr.GasPrice() == 0 {
				return sdk.ErrTxDecode("no gasprice")
			}

			gas := txBldr.GasWanted()
			fmt.Println("GetCmdSend:txBldr.GasPrices():", txBldr.GasPrice())
			msg := htdfservice.NewMsgSend(fromaddr, toaddr, coins, txBldr.GasPrice(), gas)

			cliCtx.PrintResponse = true

			return CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, fromaddr) //not completed yet, need account name
		},
	}
	return client.PostCommands(cmd)[0]
}

// PrepareTxBuilder populates a TxBuilder in preparation for the build of a Tx.
func PrepareTxBuilder(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, fromaddr sdk.AccAddress) (authtxb.TxBuilder, error) {

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txBldr.AccountNumber() == 0 {
		accNum, err := cliCtx.GetAccountNumber(fromaddr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithAccountNumber(accNum)
	}

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if txBldr.Sequence() == 0 {
		accSeq, err := cliCtx.GetAccountSequence(fromaddr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithSequence(accSeq)
	}
	return txBldr, nil
}

// CompleteAndBroadcastTxCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
//
// NOTE: Also see CompleteAndBroadcastTxREST.
func CompleteAndBroadcastTxCLI(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg, fromaddr sdk.AccAddress) error {
	//
	txBldr, err := PrepareTxBuilder(txBldr, cliCtx, fromaddr)
	if err != nil {
		return err
	}

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err := utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.GasWanted()}
		fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	passphrase, err := keys.ReadShortPassphraseFromStdin(sdk.AccAddress.String(fromaddr))
	if err != nil {
		return err
	}
	addr := sdk.AccAddress.String(fromaddr)
	ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())
	txBytes, err := ksw.BuildAndSign(txBldr, addr, passphrase, msgs)
	if err != nil {
		return err
	}
	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	cliCtx.PrintOutput(res)
	return err
}
