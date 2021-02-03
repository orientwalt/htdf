package slashing

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/orientwalt/htdf/types"
	slashingtypes "github.com/orientwalt/htdf/x/slashing/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := slashingtypes.NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(t, string(bytes), `{"address":"htdfvaloper1v93xxeqdmnmfc"}`)
}
