package tx

import (
	"encoding/hex"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// MempoolTxsResponse for query mempool txs
// NOTE: because of Tendermint, the maximum count of quering is 100
type MempoolTxsResponse struct {
	TotalTxs   int              `json:"mempool_txs_count"`
	TotalBytes int64            `json:"mempool_total_bytes"`
	Count      int              `json:"txs_count"`
	MempoolTxs []sdk.TxResponse `json:"txs"`
}

func NewMempoolTxsResponse(totalTxs, count int, totalBytes int64, mempoolTxs []sdk.TxResponse) MempoolTxsResponse {
	return MempoolTxsResponse{
		TotalTxs:   totalTxs,
		Count:      count,
		TotalBytes: totalBytes,
		MempoolTxs: mempoolTxs,
	}
}

func (mp MempoolTxsResponse) Empty() bool {
	return len(mp.MempoolTxs) == 0
}

func (mp *MempoolTxsResponse) ParseUnconfirmedTxs(cdc *codec.Codec, result *ctypes.ResultUnconfirmedTxs) error {
	for _, txBytes := range result.Txs {
		sdktx, err := parseTx(cdc, txBytes)
		if err != nil {
			return err
		}
		mp.MempoolTxs = append(mp.MempoolTxs, sdk.TxResponse{TxHash: hex.EncodeToString(txBytes.Hash()), Height: 0, Tx: sdktx})
	}
	mp.Count = result.Count
	mp.TotalBytes = result.TotalBytes
	mp.TotalTxs = result.Total
	return nil
}

// MempoolTxNumResponse for query mempool txs count
type MempoolTxNumResponse struct {
	Total      int   `json:"mempool_txs_count"`
	TotalBytes int64 `json:"mempool_total_bytes"`
}

func NewMempoolTxNumResponse(total int, totalBytes int64) MempoolTxNumResponse  {
	return MempoolTxNumResponse{
		Total: total,
		TotalBytes: totalBytes,
	}
}
