package client

import (
	"github.com/orientwalt/htdf/x/distribution/client/cli"
	"github.com/orientwalt/htdf/x/distribution/client/rest"
	govclient "github.com/orientwalt/htdf/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
