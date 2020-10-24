package rest

import (
	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client"
)

func RegisterHandlers(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}
