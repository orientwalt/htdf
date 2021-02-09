package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/x/auth"
	coretypes "github.com/orientwalt/htdf/x/evm/types"
)

type broadcastBody struct {
	Tx auth.StdTx `json:"tx"`
}

// BroadcastTxRequestHandlerFn returns the broadcast tx REST handler
func BroadcastTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m broadcastBody
		if ok := unmarshalBodyOrReturnBadRequest(cliCtx, w, r, &m); !ok {
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(m.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

type broadcastBodyRaw struct {
	Tx string `json:"tx"`
}

// BroadcastTxRequestHandlerFn returns the broadcast tx REST handler
func BroadcastTxRawRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req broadcastBodyRaw
		if ok := unmarshalBodyOrReturnBadRequest(cliCtx, w, r, &req); !ok {
			return
		}

		stdTx, err := coretypes.ReadStdTxFromRawData(cliCtx.Codec, req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(stdTx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithBroadcastMode("sync")

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func unmarshalBodyOrReturnBadRequest(cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request, m interface{}) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return false
	}
	err = cliCtx.Codec.UnmarshalJSON(body, m)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

type broadcastRawBody struct {
	Tx string `json:"rawdata"`
}

// BroadcastTxRequestHandlerFn returns the broadcast tx REST handler
func BroadcastRawTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m broadcastRawBody
		if ok := unmarshalBodyOrReturnBadRequest(cliCtx, w, r, &m); !ok {
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(m.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
