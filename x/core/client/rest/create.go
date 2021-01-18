package rest

import (
	"net/http"

	htdfRest "github.com/orientwalt/htdf/accounts/rest"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/utils/unit_convert"
	"github.com/orientwalt/htdf/x/auth"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	coretypes "github.com/orientwalt/htdf/x/core/types"
)

// CreateReq defines the properties of a send request's body.
type CreateReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	To      string       `json:"to"`
	Amount  sdk.Coins    `json:"amount"`
	Encode  bool         `json:"encode"`
}

// CreateTxRequestHandlerFn - http request handler to send coins to a address.
func CreateTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load Request
		var req CreateReq
		var mreq htdfRest.CreateShiftReq
		if !rest.ReadRESTReq(w, r, cdc, &mreq) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq.From = mreq.BaseReq.From
		req.BaseReq.Memo = mreq.BaseReq.Memo
		req.BaseReq.ChainID = mreq.BaseReq.ChainID
		req.BaseReq.AccountNumber = mreq.BaseReq.AccountNumber
		req.BaseReq.Sequence = mreq.BaseReq.Sequence
		// req.BaseReq.Fees = unit_convert.BigCoinsToDefaultCoins(mreq.BaseReq.Fees)
		req.BaseReq.GasPrice = mreq.BaseReq.GasPrice
		req.BaseReq.GasWanted = mreq.BaseReq.GasWanted
		req.BaseReq.GasAdjustment = mreq.BaseReq.GasAdjustment
		req.BaseReq.Simulate = mreq.BaseReq.Simulate
		req.To = mreq.To
		req.Encode = mreq.Encode

		// Santize
		BaseReq := req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// When generate only is supplied, the from field must be a valid Bech32
		// address.
		fromAddr, err := sdk.AccAddressFromBech32(BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		toAddr, err := sdk.AccAddressFromBech32(req.To)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := coretypes.NewMsgSendDefault(fromAddr, toAddr, unit_convert.BigCoinsToDefaultCoins(mreq.Amount))
		WriteGenerateStdTxResponse(w, cdc, cliCtx, BaseReq, []sdk.Msg{msg}, req.Encode)

		return
	}
}

// junying-todo-20190330
// WriteGenerateStdTxResponse writes response for the generate only mode.
/*
	github.com/cosmos/cosmos-sdk/client/rest
	WriteGenerateStdTxResponse
*/
func WriteGenerateStdTxResponse(w http.ResponseWriter, cdc *codec.Codec,
	cliCtx context.CLIContext, br rest.BaseReq, msgs []sdk.Msg, encodeflag bool) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, client.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gasWanted, err := client.ParseGas(br.GasWanted)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var gasPrice uint64
	gasPrice, err = client.ParseGasPrice(br.GasPrice)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := authtxb.NewTxBuilder(
		utils.GetTxEncoder(cdc), br.AccountNumber, br.Sequence, gasWanted, gasAdj,
		br.Simulate, br.ChainID, br.Memo, gasPrice,
	)

	if simAndExec {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, client.ErrInvalidGasAdjustment.Error())
			return
		}

		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	stdMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := cdc.MarshalJSON(auth.NewStdTx(stdMsg.Msgs, stdMsg.Fee, nil, stdMsg.Memo))
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !encodeflag {
		w.Write(output)
	} else {
		encoded := coretypes.Encode_Hex(output)
		w.Write([]byte(encoded))
	}

	return
}
