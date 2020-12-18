package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/rpc"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	sdkRest "github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/utils/unit_convert"

	"github.com/orientwalt/htdf/x/auth"
	htdfservice "github.com/orientwalt/htdf/x/core"
	distrTypes "github.com/orientwalt/htdf/x/distribution/types"
	stakingTypes "github.com/orientwalt/htdf/x/staking/types"
)

func AccountListRequestHandlerFn(w http.ResponseWriter, r *http.Request) {

	ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())

	accounts, err := ksw.Accounts()
	if err != nil {
		return
	}

	bJsonFormat := false
	vars := r.URL.Query()
	_, ok := vars["jsonformat"]
	if ok {
		bJsonFormat = true

		//fmt.Printf("jsonformat=%s\n", jsonformat)
	}

	var index int
	if bJsonFormat == false {
		for _, account := range accounts {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("Account #%d: {%s}\n", index, account.Address)))
			index++
		}
	} else {
		data, err := json.Marshal(accounts)
		if err != nil {
			fmt.Printf("Marshal error|err=%s\n", err)
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Write(data)
	}
}

type AccountTxsReq struct {
	Address    string `json:"address"`
	FromHeight int64  `json:"fromHeight"`
	EndHeight  int64  `json:"endHeight"`
	Flag       int64  `json:"flag"`
}

func parseTx(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {
	var tx auth.StdTx

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

//GetAccountTxsFn
//@param Address:  address of account
//@param FromHeight: query tx in the block height range [fromHeight,endHeight];  when FromHeight and endHeight is 0,  range is [CHAIN_HEIGHT-800, CHAIN_HEIGHT]
//@param EndHeight:  query tx in the block height range [fromHeight,endHeight]   when FromHeight and endHeight is 0,  range is [CHAIN_HEIGHT-800, CHAIN_HEIGHT]
//@param Flag:   query flag; 0, address appears both in fromAddress and toAddress; 1, address appears in fromAddress;  2, address appears in toAddress
func GetAccountTxsFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accountTxsReq AccountTxsReq
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&accountTxsReq)
		if err != nil {
			fmt.Printf("Decode error|err=%s\n", err)
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		//fmt.Printf("accountTxsReq=%v\n", accountTxsReq)

		//check address
		_, err = sdk.AccAddressFromBech32(accountTxsReq.Address)
		if err != nil {
			fmt.Printf("AccAddressFromBech32 error|err=%s\n", err)
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		fromHeight := accountTxsReq.FromHeight
		endHeight := accountTxsReq.EndHeight

		chainHeight, err := rpc.GetChainHeight(cliCtx)
		if err != nil {
			fmt.Printf("GetChainHeight error|err=%s\n", err)
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		//fmt.Printf("chainHeight=%d\n", chainHeight)

		//from the chain height
		scope := int64(1000)
		if fromHeight == 0 && endHeight == 0 {
			endHeight = chainHeight
			fromHeight = chainHeight - scope
		}

		//correct height parameter
		if fromHeight <= 0 {
			fromHeight = 1
		}

		if endHeight > chainHeight || endHeight == 0 {
			endHeight = chainHeight
		}

		//fmt.Printf("final|fromHeight=%d|endHeight=%d\n", fromHeight, endHeight)

		// get the node
		node, err := cliCtx.GetNode()
		if err != nil {
			fmt.Printf("getNode error|err=%s\n", err)
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result ResultAccountTxs
		result.ChainHeight = chainHeight
		result.FromHeight = fromHeight
		result.EndHeight = endHeight

		for height := fromHeight; height <= endHeight; height++ {
			//get Block info
			resultBlock, err := node.Block(&height)
			if err != nil {
				fmt.Printf("get block error|err=%s\n", err)
				sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			//fmt.Printf("currHeight=%d\n", height)
			//fmt.Printf("txTotal=%d\n", len(resultBlock.Block.Txs))

			for _, tx := range resultBlock.Block.Txs {
				sdkTx, err := parseTx(cdc, tx)
				if err != nil {
					fmt.Printf("parseTx error|err=%s\n", err)
					sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}

				switch iMsg := sdkTx.(type) {
				case auth.StdTx:
					var displayTx DisplayTx
					for _, msg := range iMsg.GetMsgs() {
						//fmt.Printf("msg|route=%s|type=%s\n", msg.Route(), msg.Type())

						switch msg := msg.(type) {
						case htdfservice.MsgSend:
							//fmt.Printf("msg|from=%s|to=%s\n", msg.From, msg.To)

							if (accountTxsReq.Flag == 0 && (msg.From.String() == accountTxsReq.Address || msg.To.String() == accountTxsReq.Address)) ||
								(accountTxsReq.Flag == 1 && msg.From.String() == accountTxsReq.Address) ||
								(accountTxsReq.Flag == 2 && msg.To.String() == accountTxsReq.Address) {
								//fmt.Printf("msg found|from=%s|to=%s\n", msg.From, msg.To)

								displayTx.From = msg.From.String()
								displayTx.To = msg.To.String()
								displayTx.Amount = unit_convert.DefaultCoinsToBigCoins(msg.Amount)
								displayTx.Hash = hex.EncodeToString(tx.Hash())
								displayTx.Height = height
								displayTx.Time = resultBlock.Block.Header.Time.Local().Format("2006-01-02 15:04:05")
								displayTx.Memo = iMsg.Memo
								displayTx.Data = msg.Data

								if len(displayTx.Data) > 0 {
									displayTx.TxClassify = rpc.TxClassify_Contract
								} else {
									displayTx.TxClassify = rpc.TxClassify_Normal
								}
								displayTx.TypeName = msg.Type()

								result.ArrTx = append(result.ArrTx, displayTx)
							}

						case stakingTypes.MsgDelegate:
							if msg.DelegatorAddress.String() == accountTxsReq.Address {
								displayTx.From = msg.DelegatorAddress.String()
								displayTx.To = msg.ValidatorAddress.String()
								displayTx.Hash = hex.EncodeToString(tx.Hash())
								displayTx.Amount = unit_convert.DefaultCoinsToBigCoins([]sdk.Coin{msg.Amount})
								displayTx.Memo = iMsg.Memo
								displayTx.TxClassify = rpc.TxClassify_Extension
								displayTx.TypeName = msg.Type()
								result.ArrTx = append(result.ArrTx, displayTx)
							}

						case distrTypes.MsgWithdrawDelegatorReward:
							if msg.DelegatorAddress.String() == accountTxsReq.Address {
								displayTx.From = msg.DelegatorAddress.String()
								displayTx.To = msg.ValidatorAddress.String()
								displayTx.Hash = hex.EncodeToString(tx.Hash())
								displayTx.Memo = iMsg.Memo
								displayTx.TxClassify = rpc.TxClassify_Extension
								displayTx.TypeName = msg.Type()
								result.ArrTx = append(result.ArrTx, displayTx)
							}

						default:
							fmt.Printf("ignore type|type=%s|route=%s\n", msg.Type(), msg.Route())
							continue
						}
					}

				default:
					fmt.Printf("unknown type: %+v\n", iMsg)
				}
			}
		}

		sdkRest.PostProcessResponse(w, cdc, &result, cliCtx.Indent)

	}
}
