package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec" // XXX fix
	sdk "github.com/orientwalt/htdf/types"
	slashingtypes "github.com/orientwalt/htdf/x/slashing/types"
)

// GetCmdQuerySigningInfo implements the command to query signing info.
func GetCmdQuerySigningInfo(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "signing-info [validator-conspub]",
		Short: "Query a validator's signing information",
		Long: strings.TrimSpace(`Use a validators' consensus public key to find the signing-info for that validator:

$ hscli query slashing signing-info cosmosvalconspub1zcjduepqfhvwcmt7p06fvdgexxhmz0l8c7sgswl7ulv7aulk364x4g5xsw7sr0k2g5
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pk, err := sdk.GetConsPubKeyBech32(args[0])
			if err != nil {
				return err
			}

			consAddr := sdk.ConsAddress(pk.Address())
			key := slashingtypes.GetValidatorSigningInfoKey(consAddr)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("Validator %s not found in slashing store", consAddr)
			}

			var signingInfo slashingtypes.ValidatorSigningInfo
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &signingInfo)
			return cliCtx.PrintOutput(signingInfo)
		},
	}
}

// GetCmdQueryParams implements a command to fetch slashing parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current slashing parameters",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query genesis parameters for the slashing module:

$ hscli query slashing params
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/parameters", slashingtypes.QuerierRoute)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params slashingtypes.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}
