package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/orientwalt/htdf/types"
)

func TestMsgSendRoute(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("atom", 10))
	var msg = NewMsgSendDefault(addr1, addr2, coins)

	require.Equal(t, msg.Route(), "htdfservice")
	require.Equal(t, msg.Type(), "send")
}

func TestMsgSendValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	atom123 := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	atom0 := sdk.NewCoins(sdk.NewInt64Coin("atom", 0))
	atom123eth123 := sdk.NewCoins(sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 123))
	atom123eth0 := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 0)}

	var emptyAddr sdk.AccAddress

	cases := []struct {
		valid bool
		tx    MsgSend
	}{
		{true, NewMsgSendDefault(addr1, addr2, atom123)},       // valid send
		{true, NewMsgSendDefault(addr1, addr2, atom123eth123)}, // valid send with multiple coins
		{false, NewMsgSendDefault(addr1, addr2, atom0)},        // non positive coin
		{false, NewMsgSendDefault(addr1, addr2, atom123eth0)},  // non positive coin in multicoins
		{false, NewMsgSendDefault(emptyAddr, addr2, atom123)},  // empty from addr
		{false, NewMsgSendDefault(addr1, emptyAddr, atom123)},  // empty to addr
		{true, NewMsgSend(addr1, addr2, atom123, 0, 30000)},    // gas below MinGasPrice(100)
		{false, NewMsgSend(addr1, addr2, atom123, 100, 10000)}, // gas below MinGas(30000)
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgSendGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("input"))
	addr2 := sdk.AccAddress([]byte("output"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("atom", 10))
	var msg = NewMsgSendDefault(addr1, addr2, coins)
	res := string(msg.GetSignBytes())
	expected := `{"Amount":[{"amount":"10","denom":"atom"}],"Data":"","From":"htdf1d9h8qat5gn84g8","GasPrice":100,"GasWanted":30000,"To":"htdf1da6hgur4wsj5g5jq"}`
	require.Equal(t, expected, res)
}

func TestMsgSendGetSigners(t *testing.T) {
	var msg = NewMsgSendDefault(sdk.AccAddress([]byte("input1")), sdk.AccAddress{}, sdk.NewCoins())
	res := msg.GetSigners()
	// TODO: fix this !
	require.Equal(t, fmt.Sprintf("%v", res), "[696E70757431]")
}

/*
// what to do w/ this test?
func TestMsgSendSigners(t *testing.T) {
	signers := []sdk.AccAddress{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	someCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	inputs := make([]Input, len(signers))
	for i, signer := range signers {
		inputs[i] = bank.NewInput(signer, someCoins)
	}
	tx := NewMsgSendDefault(inputs, nil)

	require.Equal(t, signers, tx.Signers())
}
*/
