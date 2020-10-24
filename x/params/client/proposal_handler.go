package client

import (
	govclient "github.com/orientwalt/htdf/x/gov/client"
	"github.com/orientwalt/htdf/x/params/client/cli"
	"github.com/orientwalt/htdf/x/params/client/rest"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd, rest.ProposalRESTHandler)
