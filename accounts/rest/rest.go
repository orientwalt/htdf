package rest

import (
	"github.com/gorilla/mux"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	svrConfig "github.com/orientwalt/htdf/server/config"
)

const (
	restName = "custom"
)

// resgister REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	//r.HandleFunc("/keys", QueryKeysRequestHandler(indent)).Methods("GET")
	//r.HandleFunc("/keys", AddNewKeyRequestHandler(indent)).Methods("POST")
	//r.HandleFunc("/keys/seed", SeedRequestHandler).Methods("GET")
	//r.HandleFunc("/keys/{name}/recover", RecoverRequestHandler(indent)).Methods("POST")
	//r.HandleFunc("/keys/{name}", GetKeyRequestHandler(indent)).Methods("GET")
	//r.HandleFunc("/keys/{name}", UpdateKeyRequestHandler).Methods("PUT")
	//r.HandleFunc("/keys/{name}", DeleteKeyRequestHandler).Methods("DELETE")

	if svrConfig.ApiSecurityLevel == svrConfig.ValueSecurityLevel_Low {
		r.HandleFunc("/accounts/newaccount", NewAccountRequestHandlerFn).Methods("POST")
		r.HandleFunc("/accounts/export", ExportAccountRequestHandlerFn(cliCtx, cdc)).Methods("POST")
	}

	r.HandleFunc("/accounts/list", AccountListRequestHandlerFn).Methods("GET")
	r.HandleFunc("/accounts/transactions", GetAccountTxsFn(cliCtx, cdc)).Methods("POST")
}

// register REST route
func RegisterRoute(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(
		"/auth/accounts/{address}",
		QueryAccountRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/bank/balances/{address}",
		QueryBalancesRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")

}
