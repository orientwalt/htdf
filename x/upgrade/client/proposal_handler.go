package client

import (
	govclient "github.com/orientwalt/htdf/x/gov/client"
	"github.com/orientwalt/htdf/x/upgrade/client/cli"
	"github.com/orientwalt/htdf/x/upgrade/client/rest"
)

var ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal, rest.ProposalRESTHandler)
var CancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal, rest.ProposalCancelRESTHandler)
