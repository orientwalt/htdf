package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/tx"
	"github.com/orientwalt/htdf/types/rest"

	"github.com/orientwalt/htdf/modules/record/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/record/records", recordPostHandlerFn(cliCtx)).Methods("POST")
}

func recordPostHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RecordCreateReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgCreateRecord(req.Contents, req.Creator)

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
