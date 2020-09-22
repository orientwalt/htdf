package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	common "github.com/tendermint/tendermint/libs/bytes"

	"github.com/orientwalt/htdf/app"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
)

var (
	priv = ed25519.GenPrivKey()
	addr = sdk.AccAddress(priv.PubKey().Address())
)

func TestParseQueryResponse(t *testing.T) {
	cdc := app.MakeLatestCodec()
	sdkResBytes := cdc.MustMarshalBinaryLengthPrefixed(sdk.Result{GasUsed: 10})
	gas, err := parseQueryResponse(cdc, sdkResBytes)
	assert.Equal(t, gas, uint64(10))
	assert.Nil(t, err)
	gas, err = parseQueryResponse(cdc, []byte("fuzzy"))
	assert.Equal(t, gas, uint64(0))
	assert.NotNil(t, err)
}

func TestCalculateGas(t *testing.T) {
	cdc := app.MakeLatestCodec()
	makeQueryFunc := func(gasUsed uint64, wantErr bool) func(string, common.HexBytes) ([]byte, error) {
		return func(string, common.HexBytes) ([]byte, error) {
			if wantErr {
				return nil, errors.New("")
			}
			return cdc.MustMarshalBinaryLengthPrefixed(sdk.Result{GasUsed: gasUsed}), nil
		}
	}
	type args struct {
		queryFuncGasUsed uint64
		queryFuncWantErr bool
		adjustment       float64
	}
	tests := []struct {
		name         string
		args         args
		wantEstimate uint64
		wantAdjusted uint64
		wantErr      bool
	}{
		{"error", args{0, true, 1.2}, 0, 0, true},
		{"adjusted gas", args{10, false, 1.2}, 10, 12, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryFunc := makeQueryFunc(tt.args.queryFuncGasUsed, tt.args.queryFuncWantErr)
			gotEstimate, gotAdjusted, err := CalculateGas(queryFunc, cdc, []byte(""), tt.args.adjustment)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, gotEstimate, tt.wantEstimate)
			assert.Equal(t, gotAdjusted, tt.wantAdjusted)
		})
	}
}

func TestDefaultTxEncoder(t *testing.T) {
	cdc := makeCodec()

	defaultEncoder := auth.DefaultTxEncoder(cdc)
	encoder := GetTxEncoder(cdc)

	compareEncoders(t, defaultEncoder, encoder)
}

func TestConfiguredTxEncoder(t *testing.T) {
	cdc := makeCodec()

	customEncoder := func(tx sdk.Tx) ([]byte, error) {
		return json.Marshal(tx)
	}

	config := sdk.GetConfig()
	config.SetTxEncoder(customEncoder)

	encoder := GetTxEncoder(cdc)

	compareEncoders(t, customEncoder, encoder)
}

func compareEncoders(t *testing.T, expected sdk.TxEncoder, actual sdk.TxEncoder) {
	msgs := []sdk.Msg{sdk.NewTestMsg(addr)}
	tx := auth.NewStdTx(msgs, auth.StdFee{}, []auth.StdSignature{}, "")

	defaultEncoderBytes, err := expected(tx)
	require.NoError(t, err)
	encoderBytes, err := actual(tx)
	require.NoError(t, err)
	require.Equal(t, defaultEncoderBytes, encoderBytes)
}

func makeCodec() *codec.Codec {
	cdc := app.MakeLatestCodec()
	cdc.RegisterConcrete(sdk.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}

func TestReadStdTxFromFile(t *testing.T) {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)

	// Build a test transaction
	fee := auth.NewStdFee(50000, 200000)
	stdTx := auth.NewStdTx([]sdk.Msg{}, fee, []auth.StdSignature{}, "foomemo")

	// Write it to the file
	encodedTx, _ := cdc.MarshalJSON(stdTx)
	jsonTxFile := writeToNewTempFile(t, string(encodedTx))
	defer os.Remove(jsonTxFile.Name())

	// Read it back
	decodedTx, err := ReadStdTxFromFile(cdc, jsonTxFile.Name())
	require.Nil(t, err)
	require.Equal(t, decodedTx.Memo, "foomemo")
}

func writeToNewTempFile(t *testing.T, data string) *os.File {
	fp, err := ioutil.TempFile(os.TempDir(), "client_tx_test")
	require.Nil(t, err)

	_, err = fp.WriteString(data)
	require.Nil(t, err)

	return fp
}
