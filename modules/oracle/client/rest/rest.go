package rest

import (
	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client"
)

// RegisterHandlers registers oracle REST handlers to a router
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
