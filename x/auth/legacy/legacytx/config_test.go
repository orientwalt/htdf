package legacytx_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/orientwalt/htdf/codec"
	cryptoAmino "github.com/orientwalt/htdf/crypto/codec"
	"github.com/orientwalt/htdf/testutil/testdata"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth/legacy/legacytx"
	"github.com/orientwalt/htdf/x/auth/testutil"
)

func testCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptoAmino.RegisterCrypto(cdc)
	cdc.RegisterConcrete(&testdata.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}

func TestStdTxConfig(t *testing.T) {
	cdc := testCodec()
	txGen := legacytx.StdTxConfig{Cdc: cdc}
	suite.Run(t, testutil.NewTxConfigTestSuite(txGen))
}
