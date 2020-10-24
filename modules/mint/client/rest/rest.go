package rest

import (
	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client"
)

// RegisterHandlers registers minting module REST handlers on the provided router.
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}
