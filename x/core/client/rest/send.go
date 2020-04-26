package rest

import (
	"fmt"
	"net/http"

	"github.com/orientwalt/htdf/accounts/keystore"
	htdfRest "github.com/orientwalt/htdf/accounts/rest"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/crypto/keys/keyerror"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/utils/unit_convert"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	htdfservice "github.com/orientwalt/htdf/x/core"
	hscorecli "github.com/orientwalt/htdf/x/core/client/cli"
)

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	To      string       `json:"to"`
	Amount  sdk.Coins    `json:"amount"`
	Data    string       `json:"data"`
	// GasPrice  string       `json:"gas_price"`  // uint: HTDF/gallon
	// GasWanted string       `json:"gas_wanted"` // unit: gallon
}

// var msgCdc = codec.New()

// func init() {
// 	bank.RegisterCodec(msgCdc)
// }

// SendTxRequestHandlerFn - http request handler to send coins to a address.
func SendTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req SendReq
		var mreq htdfRest.SendShiftReq
		if !rest.ReadRESTReq(w, r, cdc, &mreq) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq.From = mreq.BaseReq.From
		req.BaseReq.Memo = mreq.BaseReq.Memo
		req.BaseReq.ChainID = mreq.BaseReq.ChainID
		req.BaseReq.AccountNumber = mreq.BaseReq.AccountNumber
		req.BaseReq.Sequence = mreq.BaseReq.Sequence
		req.BaseReq.GasPrice = mreq.BaseReq.GasPrice
		req.BaseReq.GasWanted = mreq.BaseReq.GasWanted
		req.BaseReq.GasAdjustment = mreq.BaseReq.GasAdjustment
		req.BaseReq.Simulate = mreq.BaseReq.Simulate
		req.To = mreq.To
		req.Data = mreq.Data
		// req.GasPrice = mreq.GasPrice
		// req.GasWanted = mreq.GasWanted

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {

			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		toAddr, err := sdk.AccAddressFromBech32(req.To)
		if err != nil {

			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		_, gasWanted, err := client.ParseGas(req.BaseReq.GasWanted)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var gasPrice uint64
		gasPrice, err = client.ParseGasPrice(req.BaseReq.GasPrice)
		// when access smart contract, extract gas field
		// var gasPrice, gasWanted uint64
		// if len(req.Data) > 0 {
		// 	gasPrice, err = strconv.ParseUint(unit_convert.BigAmountToDefaultAmount(req.GasPrice), 10, 64)
		// 	if err != nil {
		// 		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		// 		return
		// 	}

		// 	gasWanted, err = strconv.ParseUint(req.GasWanted, 10, 64)
		// 	if err != nil {
		// 		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		// 		return
		// 	}
		// }

		fmt.Printf("gasPrice=%d|gasWanted=%d\n", gasPrice, gasWanted)

		msg := htdfservice.NewMsgSendForData(fromAddr, toAddr, unit_convert.BigCoinsToDefaultCoins(mreq.Amount), req.Data, gasPrice, gasWanted)
		CompleteAndBroadcastTxREST(w, cliCtx, req.BaseReq, mreq.BaseReq.Password, []sdk.Msg{msg}, cdc)

	}
}

//-----------------------------------------------------------------------------
// Building / Sending utilities

// CompleteAndBroadcastTxREST implements a utility function that facilitates
// sending a series of messages in a signed tx. In addition, it will handle
// tx gas simulation and estimation.
//
// NOTE: Also see CompleteAndBroadcastTxCLI.
func CompleteAndBroadcastTxREST(w http.ResponseWriter, cliCtx context.CLIContext,
	baseReq rest.BaseReq, password string, msgs []sdk.Msg, cdc *codec.Codec) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, baseReq.GasAdjustment, client.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gasWanted, err := client.ParseGas(baseReq.GasWanted)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var gasPrice uint64
	gasPrice, err = client.ParseGasPrice(baseReq.GasPrice)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := authtxb.NewTxBuilder(
		utils.GetTxEncoder(cdc), baseReq.AccountNumber,
		baseReq.Sequence, gasWanted, gasAdj, baseReq.Simulate,
		baseReq.ChainID, baseReq.Memo, gasPrice,
	)

	// get fromaddr
	fromaddr := msgs[0].(htdfservice.MsgSend).GetSigners()[0]

	txBldr, err = hscorecli.PrepareTxBuilder(txBldr, cliCtx, fromaddr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if baseReq.Simulate || simAndExec {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, client.ErrInvalidGasAdjustment.Error())
			return
		}

		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if baseReq.Simulate {
			rest.WriteSimulationResponse(w, cdc, txBldr.GasWanted())
			return
		}
	}

	bech32 := sdk.AccAddress.String(fromaddr)

	if err != nil {
		return
	}

	ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())
	txBytes, err := ksw.BuildAndSign(txBldr, bech32, password, msgs)
	if keyerror.IsErrKeyNotFound(err) {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	} else if keyerror.IsErrWrongPassword(err) {
		rest.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
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
