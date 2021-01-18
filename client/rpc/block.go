package rpc

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	distrTypes "github.com/orientwalt/htdf/x/distribution/types"
	stakingTypes "github.com/orientwalt/htdf/x/staking/types"

	"github.com/gorilla/mux"
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/types/rest"
	sdkRest "github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/utils/unit_convert"
	"github.com/orientwalt/htdf/x/auth"
	coretypes "github.com/orientwalt/htdf/x/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

//BlockCommand returns the verified block data for a given heights
func BlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [height]",
		Short: "Get verified data for a the block at given height",
		Args:  cobra.MaximumNArgs(1),
		RunE:  printBlock,
	}
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	viper.BindPFlag(client.FlagNode, cmd.Flags().Lookup(client.FlagNode))
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	viper.BindPFlag(client.FlagTrustNode, cmd.Flags().Lookup(client.FlagTrustNode))
	return cmd
}

func getBlock(cliCtx context.CLIContext, height *int64) ([]byte, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(height)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		check, err := cliCtx.Verify(res.Block.Height)
		if err != nil {
			return nil, err
		}

		if err := tmliteProxy.ValidateHeader(&res.Block.Header, check); err != nil {
			return nil, err
		}

		err = tmliteProxy.ValidateBlock(res.Block, check)
		if err != nil {
			return nil, err
		}
	}

	if cliCtx.Indent {
		return cdc.MarshalJSONIndent(res, "", "  ")
	}
	return cdc.MarshalJSON(res)
}

// get the current blockchain height
func GetChainHeight(cliCtx context.CLIContext) (int64, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return -1, err
	}
	status, err := node.Status()
	if err != nil {
		return -1, err
	}
	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// CMD

func printBlock(cmd *cobra.Command, args []string) error {
	var height *int64
	// optional height
	if len(args) > 0 {
		h, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		if h > 0 {
			tmp := int64(h)
			height = &tmp
		}
	}

	output, err := getBlock(context.NewCLIContext(), height)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// REST

// REST handler to get a block
func BlockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				"ERROR: Couldn't parse block height. Assumed format is '/block/{height}'.")
			return
		}
		chainHeight, err := GetChainHeight(cliCtx)
		if height > chainHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound,
				"ERROR: Requested block height is bigger then the chain length.")
			return
		}
		output, err := getBlock(cliCtx, &height)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

// REST handler to get the latest block
func LatestBlockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		height, err := GetChainHeight(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		output, err := getBlock(cliCtx, &height)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

//BlockDetails struct
const (
	TxClassify_Normal    uint32 = 0
	TxClassify_Contract  uint32 = 1
	TxClassify_Extension uint32 = 2
)

type DisplayTx struct {
	From       string
	To         string
	Amount     []sdk.BigCoin
	Hash       string
	Memo       string
	Data       string
	TxClassify uint32
	TypeName   string
}

type DisplayBlock struct {
	Txs        []DisplayTx          `json:"txs"`
	Evidence   tmTypes.EvidenceData `json:"evidence"`
	LastCommit *tmTypes.Commit      `json:"last_commit"`
}

type BlockInfo struct {
	BlockMeta *tmTypes.BlockMeta `json:"block_meta"`
	Block     DisplayBlock       `json:"block"`
	Time      string             `json:"time"`
}

type DisplayFee struct {
	// Amount    []sdk.BigCoin `json:"amount"` // junying-todo, 2019-12-07
	GasWanted string `json:"gas_wanted"`
	GasPrice  string `json:"gas_price"`
}

type StdTx struct {
	Msgs       []DisplayTx         `json:"msg"`
	Fee        DisplayFee          `json:"fee"`
	Signatures []auth.StdSignature `json:"signatures"`
	Memo       string              `json:"memo"`
}

type GetTxResponse struct {
	Height    int64               `json:"height"`
	TxHash    string              `json:"txhash"`
	Code      uint32              `json:"code,omitempty"`
	Log       sdk.ABCIMessageLogs `json:"log,omitempty"`
	Info      string              `json:"info,omitempty"`
	GasWanted int64               `json:"gas_wanted,omitempty"`
	GasUsed   int64               `json:"gas_used,omitempty"`
	Events    sdk.StringEvents    `json:"events,omitempty"`
	Codespace string              `json:"codespace,omitempty"`
	Tx        StdTx               `json:"tx,omitempty"`
	// Data      string              `json:"data,omitempty"`

}

// GetBlockDetailFn
func GetBlockDetailFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ERROR: Couldn't parse block height. Assumed format is '/block/{height}'."))
			return
		}
		chainHeight, err := GetChainHeight(cliCtx)
		if height > chainHeight {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR: Requested block height is bigger then the chain length."))
			return
		}

		// get the node
		node, err := cliCtx.GetNode()
		if err != nil {
			fmt.Printf("getNode error|err=%s\n", err)
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var blockInfo BlockInfo
		//get Block info
		resultBlock, err := node.Block(&height)
		if err != nil {
			fmt.Printf("get block error|err=%s\n", err)
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// blockInfo.BlockMeta = resultBlock.BlockMeta
		blockInfo.Block.Evidence = resultBlock.Block.Evidence
		blockInfo.Block.LastCommit = resultBlock.Block.LastCommit
		blockInfo.Time = resultBlock.Block.Header.Time.Local().Format("2006-01-02 15:04:05")

		for _, tx := range resultBlock.Block.Data.Txs {
			sdkTx, err := parseTx(cdc, tx)
			if err != nil {
				fmt.Printf("parseTx error|err=%s\n", err)
				sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			switch currTx := sdkTx.(type) {
			case auth.StdTx:
				var displayTx DisplayTx
				for _, msg := range currTx.GetMsgs() {
					//fmt.Printf("msg|route=%s|type=%s\n", msg.Route(), msg.Type())
					switch msg := msg.(type) {
					case coretypes.MsgSend:

						displayTx.From = msg.From.String()
						displayTx.To = msg.To.String()
						displayTx.Hash = hex.EncodeToString(tx.Hash())
						displayTx.Amount = unit_convert.DefaultCoinsToBigCoins(msg.Amount)
						displayTx.Memo = currTx.Memo
						displayTx.Data = msg.Data

						if len(displayTx.Data) > 0 {
							displayTx.TxClassify = TxClassify_Contract
						} else {
							displayTx.TxClassify = TxClassify_Normal
						}
						displayTx.TypeName = msg.Type()

						blockInfo.Block.Txs = append(blockInfo.Block.Txs, displayTx)

						//fmt.Printf("msg|from=%s|to=%s\n", msg.From, msg.To)

					case stakingTypes.MsgDelegate:
						displayTx.From = msg.DelegatorAddress.String()
						displayTx.To = msg.ValidatorAddress.String()
						displayTx.Hash = hex.EncodeToString(tx.Hash())
						displayTx.Amount = unit_convert.DefaultCoinsToBigCoins([]sdk.Coin{msg.Amount})
						displayTx.Memo = currTx.Memo
						displayTx.TxClassify = TxClassify_Extension
						displayTx.TypeName = msg.Type()
						blockInfo.Block.Txs = append(blockInfo.Block.Txs, displayTx)

					case distrTypes.MsgWithdrawDelegatorReward:
						displayTx.From = msg.DelegatorAddress.String()
						displayTx.To = msg.ValidatorAddress.String()
						displayTx.Hash = hex.EncodeToString(tx.Hash())
						displayTx.Memo = currTx.Memo
						displayTx.TxClassify = TxClassify_Extension
						displayTx.TypeName = msg.Type()
						blockInfo.Block.Txs = append(blockInfo.Block.Txs, displayTx)

					default:
						fmt.Printf("ignore type|type=%s|route=%s\n", msg.Type(), msg.Route())
						continue
					}
				}

			default:
				fmt.Printf("unknown type: %+v\n", currTx)
			}
		}

		sdkRest.PostProcessResponse(w, cdc, &blockInfo, cliCtx.Indent)
	}
}

func parseTx(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {
	var tx auth.StdTx

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func formatTxResult(cdc *codec.Codec, res *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (sdk.TxResponse, error) {
	tx, err := parseTx(cdc, res.Tx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return sdk.NewResponseResultTx(res, tx, resBlock.Block.Time.Format(time.RFC3339)), nil
}

//
func GetTxFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		hash, err := hex.DecodeString(hashHexStr)
		if err != nil {
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// get the node
		node, err := cliCtx.GetNode()
		if err != nil {
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resultTx, err := node.Tx(hash, !cliCtx.TrustNode)
		if err != nil {
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBlocks, err := getBlocksForTxResults(cliCtx, []*ctypes.ResultTx{resultTx})
		if err != nil {
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		txResp, err := formatTxResult(cdc, resultTx, resBlocks[resultTx.Height])
		if err != nil {
			sdkRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var getTxResponse GetTxResponse
		getTxResponse.Height = txResp.Height
		getTxResponse.TxHash = txResp.TxHash
		getTxResponse.Code = txResp.Code
		// getTxResponse.Data = txResp.Data
		getTxResponse.Log = txResp.Logs
		getTxResponse.Info = txResp.Info
		getTxResponse.GasWanted = txResp.GasWanted
		getTxResponse.GasUsed = txResp.GasUsed
		getTxResponse.Events = txResp.Events
		getTxResponse.Codespace = txResp.Codespace

		// switch currTx := txResp.Tx.(type) {
		// case auth.StdTx:
		// 	// getTxResponse.Tx.Fee.Amount = unit_convert.DefaultCoinsToBigCoins(currTx.Fee.Amount())
		// 	// getTxResponse.Tx.Fee.GasWanted = unit_convert.DefaultAmoutToBigAmount(strconv.FormatUint(currTx.Fee.GasWanted, 10))
		// 	getTxResponse.Tx.Signatures = currTx.Signatures
		// 	getTxResponse.Tx.Memo = currTx.Memo

		// 	var displayTx DisplayTx
		// 	for _, msg := range currTx.GetMsgs() {
		// 		//fmt.Printf("msg|route=%s|type=%s\n", msg.Route(), msg.Type())
		// 		switch msg := msg.(type) {
		// 		case coretypes.MsgSend:
		// 			displayTx.From = msg.From
		// 			displayTx.To = msg.To
		// 			displayTx.Amount = unit_convert.DefaultCoinsToBigCoins(msg.Amount)
		// 			getTxResponse.Tx.Msgs = append(getTxResponse.Tx.Msgs, displayTx)

		// 		default:
		// 			fmt.Printf("ignore type|type=%s|route=%s\n", msg.Type(), msg.Route())
		// 			continue
		// 		}
		// 	}

		// default:
		// 	fmt.Printf("unknown type: %+v\n", currTx)
		// }

		sdkRest.PostProcessResponse(w, cdc, getTxResponse, cliCtx.Indent)
	}
}

func getBlocksForTxResults(cliCtx context.CLIContext, resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	resBlocks := make(map[int64]*ctypes.ResultBlock)

	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := node.Block(&resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}
