package auth

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/params"
)

// DefaultParamspace defines the default auth module parameter subspace
const DefaultParamspace = "auth"
const DefaultFeeCollection = "fee"

// Default parameter values
const (
	DefaultMaxMemoCharacters uint64 = 256
	DefaultTxSigLimit        uint64 = 7
	DefaultInitialHeight     uint64 = 1
	//DefaultTxSizeLimit            uint64 = 1000
	DefaultTxSizeCostPerByte      uint64 = 10
	DefaultSigVerifyCostED25519   uint64 = 590
	DefaultSigVerifyCostSecp256k1 uint64 = 4000 // modified by junying, 2019-09-04, 1000 to 4000
)

var (
	MinimumGasPrice      = sdk.ZeroInt()
	MaximumGasPrice      = sdk.NewIntWithDecimal(1, 18) //1htdf, 10^18satoshi
	MinimumTxSizeLimit   = uint64(500)
	MaximumTxSizeLimit   = uint64(1500)
	MinimumInitialHeight = uint64(1)
)

// Parameter keys
var (
	KeygasPriceThreshold      = []byte("GasPriceThreshold")
	KeyMaxMemoCharacters      = []byte("MaxMemoCharacters")
	KeyTxSigLimit             = []byte("TxSigLimit")
	KeyTxSizeLimit            = []byte("TxSizeLimit")
	KeyTxSizeCostPerByte      = []byte("TxSizeCostPerByte")
	KeySigVerifyCostED25519   = []byte("SigVerifyCostED25519")
	KeySigVerifyCostSecp256k1 = []byte("SigVerifyCostSecp256k1")
)

var (
	KeyInitialHeight = []byte("InitialHeight")
)

// Params defines the parameters for the auth module.
type Params struct {
	GasPriceThreshold sdk.Int `json:"gas_price_threshold"`
	MaxMemoCharacters uint64  `json:"max_memo_characters"`
	TxSigLimit        uint64  `json:"tx_sig_limit"`
	InitialHeight     uint64  `json:"initial_height"`
	//TxSizeLimit            uint64  `json:"tx_size"` // tx size limit
	TxSizeCostPerByte      uint64 `json:"tx_size_cost_per_byte"`
	SigVerifyCostED25519   uint64 `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256k1 uint64 `json:"sig_verify_cost_secp256k1"`
}

var _ params.ParamSet = &Params{}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeygasPriceThreshold, Value: &p.GasPriceThreshold},
		{Key: KeyMaxMemoCharacters, Value: &p.MaxMemoCharacters},
		{Key: KeyTxSigLimit, Value: &p.TxSigLimit},
		{Key: KeyInitialHeight, Value: &p.InitialHeight},
		//{KeyTxSizeLimit, &p.TxSizeLimit},
		{Key: KeyTxSizeCostPerByte, Value: &p.TxSizeCostPerByte},
		{Key: KeySigVerifyCostED25519, Value: &p.SigVerifyCostED25519},
		{Key: KeySigVerifyCostSecp256k1, Value: &p.SigVerifyCostSecp256k1},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := msgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := msgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		GasPriceThreshold: sdk.NewIntWithDecimal(6, 12),
		MaxMemoCharacters: DefaultMaxMemoCharacters,
		TxSigLimit:        DefaultTxSigLimit,
		InitialHeight:     DefaultInitialHeight,
		//TxSizeLimit:            DefaultTxSizeLimit,
		TxSizeCostPerByte:      DefaultTxSizeCostPerByte,
		SigVerifyCostED25519:   DefaultSigVerifyCostED25519,
		SigVerifyCostSecp256k1: DefaultSigVerifyCostSecp256k1,
	}
}

// Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamspace
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeygasPriceThreshold):
		threshold, ok := sdk.NewIntFromString(value)
		if !ok {
			return nil, params.ErrInvalidString(value)
		}
		if !threshold.GT(MinimumGasPrice) || threshold.GT(MaximumGasPrice) {
			return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidGasPriceThreshold, fmt.Sprintf("Gas price threshold (%s) should be (0, 10^18satoshi]", value))
		}
		return threshold, nil
	case string(KeyTxSizeLimit):
		txsize, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if txsize < MinimumTxSizeLimit || txsize > MaximumTxSizeLimit {
			return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidTxSizeLimit, fmt.Sprintf("Tx size limit (%s) should be [500, 1500]", value))
		}
		return txsize, nil

	case string(KeyInitialHeight):
		initialHeight, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if initialHeight < MinimumInitialHeight {
			return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidInitialHeight, fmt.Sprintf("Initial Height (%s) should be same or greater than 1", value))
		}
		return initialHeight, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeygasPriceThreshold):
		err := cdc.UnmarshalJSON(bytes, &p.GasPriceThreshold)
		return p.GasPriceThreshold.String(), err
	case string(KeyTxSizeLimit):
		err := cdc.UnmarshalJSON(bytes, &p.TxSigLimit)
		return strconv.FormatUint(uint64(p.TxSigLimit), 10), err
	case string(KeyInitialHeight):
		err := cdc.UnmarshalJSON(bytes, &p.InitialHeight)
		return strconv.FormatUint(uint64(p.InitialHeight), 10), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("GasPriceThreshold: %d\n", p.GasPriceThreshold))
	sb.WriteString(fmt.Sprintf("MaxMemoCharacters: %d\n", p.MaxMemoCharacters))
	sb.WriteString(fmt.Sprintf("TxSigLimit: %d\n", p.TxSigLimit))
	sb.WriteString(fmt.Sprintf("InitialHeight: %d\n", p.InitialHeight))
	// sb.WriteString(fmt.Sprintf("TxSizeLimit: %d\n", p.TxSizeLimit))
	sb.WriteString(fmt.Sprintf("TxSizeCostPerByte: %d\n", p.TxSizeCostPerByte))
	sb.WriteString(fmt.Sprintf("SigVerifyCostED25519: %d\n", p.SigVerifyCostED25519))
	sb.WriteString(fmt.Sprintf("SigVerifyCostSecp256k1: %d\n", p.SigVerifyCostSecp256k1))
	return sb.String()
}
