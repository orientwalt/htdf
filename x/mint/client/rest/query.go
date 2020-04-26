package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/rpc"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/x/mint"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/minting/parameters",
		queryParamsHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/minting/inflation",
		queryInflationHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/minting/annual-provisions",
		queryAnnualProvisionsHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/minting/total-provisions",
		queryTotalProvisionsHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/minting/rewards/{height}",
		queryBlockRewardHandlerFn(cdc, cliCtx),
	).Methods("GET")
}

func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryParameters)

		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryInflationHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryInflation)

		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryAnnualProvisionsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryAnnualProvisions)

		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryTotalProvisionsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryTotalProvisions)

		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryBlockRewardHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				"ERROR: Couldn't parse block height. Assumed format is '/rewards/{height}'.")
			return
		}
		chainHeight, err := rpc.GetChainHeight(cliCtx)
		if height > chainHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound,
				"ERROR: Requested block height is bigger then the chain length.")
			return
		}

		bz, err := cliCtx.Codec.MarshalJSON(mint.NewQueryBlockRewardParams(height))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound,
				"ERROR: height must be integer.")
			return
		}

		route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryBlockRewards)

		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
