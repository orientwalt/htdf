package rest

import (
	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
}
