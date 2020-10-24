package tx

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/orientwalt/htdf/codec"
	codectypes "github.com/orientwalt/htdf/codec/types"
	"github.com/orientwalt/htdf/std"
	"github.com/orientwalt/htdf/testutil/testdata"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth/testutil"
)

func TestGenerator(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &testdata.TestMsg{})
	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	suite.Run(t, testutil.NewTxConfigTestSuite(NewTxConfig(protoCodec, DefaultSignModes)))
}
