package htdfservice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	v2auth "github.com/orientwalt/htdf/app/v2/auth"
	"github.com/tendermint/tendermint/crypto"
)

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool) {
	_, result, abort := anteHandler(ctx, tx, simulate)
	require.False(t, abort)
	require.Equal(t, sdk.CodeOK, result.Code)
	require.True(t, result.IsOK())
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool, code sdk.CodeType) {
	_, result, abort := anteHandler(ctx, tx, simulate)
	require.True(t, abort, "abort, expected: true, got: false")

	require.Equal(t, code, result.Code, fmt.Sprintf("Expected %v, got %v", code, result))
	require.Equal(t, sdk.CodespaceRoot, result.Codespace, "code not match")

	// if code == sdk.CodeOutOfGas {
	// stdTx, ok := tx.(auth.StdTx)
	// require.True(t, ok, "tx must be in form auth.StdTx")
	// GasWanted set correctly
	// require.Equal(t, stdTx.Fee.GasWanted, result.GasWanted, "Gas wanted not set correctly")
	// require.True(t, result.GasUsed > result.GasWanted, "GasUsed not greated than GasWanted")
	// Check that context is set correctly
	// require.Equal(t, result.GasUsed, newCtx.GasMeter().GasConsumed(), "Context not updated correctly")
	// }
}

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

// add by yqq 2020-11-17
func TestSendMsg_GenesisBlock(t *testing.T) {

	// setup
	input := setupTestInput()
	ctx := input.ctx
	anteHandler := auth.NewAnteHandler(input.ak, input.fck)

	// keys and addresses
	priv1, _, addr1 := keyPubAddr()

	// set the accounts
	acc1 := input.ak.NewAccountWithAddress(ctx, addr1)
	input.ak.SetAccount(ctx, acc1)

	// msg and signatures
	var tx sdk.Tx
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	// account's balance is not enough for paying fee
	var blkNum int64 = 0
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(1500000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 2)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSendForData(addr1, addr2, sendAmount, "aabbccddeeff", sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		// checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFunds)
		checkValidTx(t, anteHandler, ctx, tx, false)
	}

	// account's balance is not enough for sendAmount
	{
		ctx = ctx.WithBlockHeight(blkNum)

		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(30000, 100)

		// balance only enough for paying fee
		acc1.SetCoins(sendFee.Amount())
		input.ak.SetAccount(ctx, acc1)

		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 30000*10000)) // greater than balance
		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkValidTx(t, anteHandler, ctx, tx, false)
	}

	// sendAmount is 0 satoshi , and data is empty
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(30000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 100000000000000)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkValidTx(t, anteHandler, ctx, tx, false)
	}

	// gasWanted is lower than defaultGasWanted
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(3000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 1))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 100000000000000)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkValidTx(t, anteHandler, ctx, tx, false)
	}

}

// add by yqq 2020-11-17
func TestSendMsg_NonGenesisBlock(t *testing.T) {

	// setup
	input := setupTestInput()
	ctx := input.ctx
	anteHandler := auth.NewAnteHandler(input.ak, input.fck)

	// keys and addresses
	priv1, _, addr1 := keyPubAddr()

	// set the accounts
	acc1 := input.ak.NewAccountWithAddress(ctx, addr1)
	input.ak.SetAccount(ctx, acc1)

	// msg and signatures
	var tx sdk.Tx
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	// account's balance is not enough for paying fee
	var blkNum int64 = 100
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(1500000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 2)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSendForData(addr1, addr2, sendAmount, "aabbccddeeff", sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		
		// auth  be v2 auth , not check balance for fee, so this ok for old version 
		checkValidTx(t, anteHandler, ctx, tx, false) 
	}

	// Data is very long , it will cause Out of gas
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(1500000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 2)))
		input.ak.SetAccount(ctx, acc1)

		contractData := ""
		for i := 0; i < 1000000; i++ {
			contractData += fmt.Sprintf("%02X", i*i)
			if len(contractData) > 200000 {
				break
			}
		}
		t.Logf("contractData' length : %d\n", len(contractData))
		msgSend := NewMsgSendForData(addr1, addr2, sendAmount, contractData, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeOutOfGas)
	}

	//  account's balance is not enough for sendAmount
	{
		ctx = ctx.WithBlockHeight(blkNum)

		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(30000, 100)

		// balance only enough for paying fee
		acc1.SetCoins(sendFee.Amount())
		input.ak.SetAccount(ctx, acc1)

		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 30000*10000)) // greater than balance
		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkValidTx(t, anteHandler, ctx, tx, false)
		// checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientCoins)
	}

	// sendAmount is 0 satoshi , and data is empty
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(30000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 100000000000000)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientCoins)
	}

	//  gasWanted is lower than defaultGasWanted
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(3000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 1))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 100000000000000)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInvalidGas)
	}

}




// add by yqq 2020-11-17
func TestSendMsg_GenesisBlock_V2(t *testing.T) {

	// setup
	input := setupTestInput()
	ctx := input.ctx
	anteHandler := v2auth.NewAnteHandler(input.ak, input.fck)

	// keys and addresses
	priv1, _, addr1 := keyPubAddr()

	// set the accounts
	acc1 := input.ak.NewAccountWithAddress(ctx, addr1)
	input.ak.SetAccount(ctx, acc1)

	// msg and signatures
	var tx sdk.Tx
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	// account's balance is not enough for paying fee
	var blkNum int64 = 0
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(1500000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 2)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSendForData(addr1, addr2, sendAmount, "aabbccddeeff", sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		// checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFunds)
		checkValidTx(t, anteHandler, ctx, tx, false)
	}

	// account's balance is not enough for sendAmount
	{
		ctx = ctx.WithBlockHeight(blkNum)

		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(30000, 100)

		// balance only enough for paying fee
		acc1.SetCoins(sendFee.Amount())
		input.ak.SetAccount(ctx, acc1)

		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 30000*10000)) // greater than balance
		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkValidTx(t, anteHandler, ctx, tx, false)
	}

	// sendAmount is 0 satoshi , and data is empty
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(30000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 100000000000000)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkValidTx(t, anteHandler, ctx, tx, false)
	}

	// gasWanted is lower than defaultGasWanted
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(3000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 1))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 100000000000000)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkValidTx(t, anteHandler, ctx, tx, false)
	}

}

// add by yqq 2020-11-17
func TestSendMsg_NonGenesisBlock_V2(t *testing.T) {

	// setup
	input := setupTestInput()
	ctx := input.ctx
	anteHandler := v2auth.NewAnteHandler(input.ak, input.fck)

	// keys and addresses
	priv1, _, addr1 := keyPubAddr()

	// set the accounts
	acc1 := input.ak.NewAccountWithAddress(ctx, addr1)
	input.ak.SetAccount(ctx, acc1)

	// msg and signatures
	var tx sdk.Tx
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	// account's balance is not enough for paying fee
	var blkNum int64 = 100
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(1500000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 2)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSendForData(addr1, addr2, sendAmount, "aabbccddeeff", sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFunds)
	}

	// Data is very long , it will cause Out of gas
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(1500000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 2)))
		input.ak.SetAccount(ctx, acc1)

		contractData := ""
		for i := 0; i < 1000000; i++ {
			contractData += fmt.Sprintf("%02X", i*i)
			if len(contractData) > 200000 {
				break
			}
		}
		t.Logf("contractData' length : %d\n", len(contractData))
		msgSend := NewMsgSendForData(addr1, addr2, sendAmount, contractData, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeOutOfGas)
	}

	//  account's balance is not enough for sendAmount
	{
		ctx = ctx.WithBlockHeight(blkNum)

		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(30000, 100)

		// balance only enough for paying fee
		acc1.SetCoins(sendFee.Amount())
		input.ak.SetAccount(ctx, acc1)

		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 30000*10000)) // greater than balance
		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkValidTx(t, anteHandler, ctx, tx, false)
		// checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientCoins)
	}

	// sendAmount is 0 satoshi , and data is empty
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(30000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 0))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 100000000000000)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientCoins)
	}

	//  gasWanted is lower than defaultGasWanted
	{
		ctx = ctx.WithBlockHeight(blkNum) // test for non-genesis
		_, _, addr2 := keyPubAddr()
		sendFee := auth.NewStdFee(3000, 100)
		sendAmount := sdk.NewCoins(sdk.NewInt64Coin("satoshi", 1))

		acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("satoshi", 100000000000000)))
		input.ak.SetAccount(ctx, acc1)

		msgSend := NewMsgSend(addr1, addr2, sendAmount, sendFee.GasPrice, sendFee.GasWanted)
		msgSends := []sdk.Msg{msgSend}
		seqs[0] = acc1.GetSequence()
		tx = newTestTx(ctx, msgSends, privs, accnums, seqs, sendFee)
		checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInvalidGas)
	}

}
