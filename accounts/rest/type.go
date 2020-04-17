package rest

import (
	sdk "github.com/orientwalt/htdf/types"
	"github.com/tendermint/tendermint/crypto"
)

//accinfo
type NewAccInfo struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []sdk.BigCoin  `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

type AccountBody struct {
	Type  string     `json:"type"`
	Value NewAccInfo `json:"value"`
}

//send
type SendShiftReq struct {
	BaseReq   SendDataBaseReq `json:"base_req"`
	To        string          `json:"to"`
	Amount    []sdk.BigCoin   `json:"amount"`
	Data      string          `json:"data"`
	GasPrice  string          `json:"gas_price"`  // unit: HTDF/gallon
	GasWanted string          `json:"gas_wanted"` // unit: gallon
}

type SendDataBaseReq struct {
	From          string `json:"from"`
	Password      string `json:"password"`
	Memo          string `json:"memo"`
	ChainID       string `json:"chain_id"`
	AccountNumber uint64 `json:"account_number"`
	Sequence      uint64 `json:"sequence"`
	GasPrice      string `json:"gas_price"`
	GasWanted     string `json:"gas_wanted"`
	GasAdjustment string `json:"gas_adjustment"`
	GenerateOnly  bool   `json:"generate_only"`
	Simulate      bool   `json:"simulate"`
}

//create
type CreateShiftReq struct {
	//BaseReq    rest.BaseReq     `json:"base_req"`
	BaseReq SendDataBaseReq `json:"base_req"`
	To      string          `json:"to"`
	Amount  []sdk.BigCoin   `json:"amount"`
	Encode  bool            `json:"encode"`
}

type DisplayTx struct {
	From       string
	To         string
	Amount     []sdk.BigCoin
	Hash       string
	Height     int64
	Time       string
	Memo       string
	Data       string
	TxClassify uint32
	TypeName   string
}

type ResultAccountTxs struct {
	ChainHeight int64
	FromHeight  int64
	EndHeight   int64
	ArrTx       []DisplayTx
}
