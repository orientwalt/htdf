package rest

import (
	"fmt"
	"net/http"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/types/rest"

	svrConfig "github.com/orientwalt/htdf/server/config"

	"github.com/gorilla/mux"
	"github.com/orientwalt/htdf/x/auth"
	htdfservice "github.com/orientwalt/htdf/x/core"
)

const (
	restName = "custom"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {

	if svrConfig.ApiSecurityLevel == svrConfig.ValueSecurityLevel_Low {
		r.HandleFunc(fmt.Sprintf("/%s/send", storeName), SendTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
		r.HandleFunc(fmt.Sprintf("/%s/create", storeName), CreateTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
		r.HandleFunc(fmt.Sprintf("/%s/sign", storeName), SignTxRawRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	}
	r.HandleFunc(fmt.Sprintf("/%s/broadcast", storeName), BroadcastTxRawRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	//
	r.HandleFunc(
		fmt.Sprintf("/%s/contract/{address}/{code}", storeName),
		QueryContractRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")
}

// query contractREST Handler
func QueryContractRequestHandlerFn(
	storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		code := vars["code"]

		contractaddr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		bz, err := cliCtx.Codec.MarshalJSON(htdfservice.NewQueryContractParams(contractaddr, code))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound,
				"ERROR: height must be integer.")
			return
		}
		//
		route := fmt.Sprintf("custom/%s/%s", htdfservice.QuerierRoute, htdfservice.QueryContract)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
