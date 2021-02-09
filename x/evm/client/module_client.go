package client

import (
	"github.com/orientwalt/htdf/client"
	htdfservicecmd "github.com/orientwalt/htdf/x/core/client/cli"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group htdfservice queries under a subcommand
	htdfsvcQueryCmd := &cobra.Command{
		Use:   "hs",
		Short: "Querying commands for the htdfservice module",
	}

	htdfsvcQueryCmd.AddCommand(client.GetCommands()...)

	return htdfsvcQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	htdfsvcTxCmd := &cobra.Command{
		Use:   "hs",
		Short: "HtdfService transactions subcommands",
	}

	htdfsvcTxCmd.AddCommand(client.PostCommands(
		//htdfservicecmd.GetCmdAdd(mc.cdc),
		//htdfservicecmd.GetCmdIssue(mc.cdc),
		htdfservicecmd.GetCmdSend(mc.cdc),
		htdfservicecmd.GetCmdCreate(mc.cdc),
		htdfservicecmd.GetCmdSign(mc.cdc),
		htdfservicecmd.GetCmdBroadCast(mc.cdc),
	)...)

	return htdfsvcTxCmd
}
