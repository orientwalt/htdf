package rest

import (
	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/types/rest"

	"github.com/orientwalt/htdf/modules/record/types"
)

// Rest variable names
// nolint
const (
	RestRecordID = "record-id"
)

// RegisterHandlers defines routes that get registered by the main application
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type RecordCreateReq struct {
	BaseReq  rest.BaseReq    `json:"base_req" yaml:"base_req"` // base req
	Contents []types.Content `json:"contents" yaml:"contents"`
	Creator  sdk.AccAddress  `json:"creator" yaml:"creator"`
}
