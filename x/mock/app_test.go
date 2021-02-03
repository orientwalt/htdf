package mock

import (
	"testing"

	newevmtypes "github.com/orientwalt/htdf/evm/types"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	evmtypes "github.com/orientwalt/htdf/x/core/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type (
	expectedBalance struct {
		addr  sdk.AccAddress
		coins sdk.Coins
	}

	appTestCase struct {
		expSimPass       bool
		expPass          bool
		msgs             []sdk.Msg
		accNums          []uint64
		accSeqs          []uint64
		privKeys         []crypto.PrivKey
		expectedBalances []expectedBalance
	}
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	addr3 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	priv4 = secp256k1.GenPrivKey()
	addr4 = sdk.AccAddress(priv4.PubKey().Address())

	coins     = sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)}
	halfCoins = sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 5)}
	manyCoins = sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), sdk.NewInt64Coin("satoshi", 1)}
	freeFee   = auth.NewStdFee(100000, 0)

	sendMsg1 = evmtypes.NewMsgSendDefault(addr1, addr2, coins)
	sendMsg2 = evmtypes.NewMsgSendDefault(addr1, addr2, manyCoins)
)

// initialize the mock application for this module
func getMockApp(t *testing.T) *App {
	mapp, err := getBenchmarkMockApp()
	require.NoError(t, err)
	return mapp
}

func TestMsgSendWithAccounts(t *testing.T) {
	mapp := getMockApp(t)
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 6700000000000)},
	}

	accs := newevmtypes.NewAccount(acc)

	SetGenesis(mapp, []auth.Account{accs})

	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, accs, res1.(*newevmtypes.Account))

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{sendMsg1},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 6699996999990)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)}},
			},
		},
	}

	for _, tc := range testCases {
		SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}

	// bumping the tx nonce number without resigning should be an auth error
	mapp.BeginBlock(abci.RequestBeginBlock{})

	tx := GenTx([]sdk.Msg{sendMsg1}, []uint64{0}, []uint64{0}, priv1)
	// tx.Signatures[0].Sequence = 1

	res := mapp.Deliver(tx)

	require.Equal(t, sdk.CodeUnauthorized, res.Code, res.Log)
	require.Equal(t, sdk.CodespaceRoot, res.Codespace)

	// resigning the tx with the bumped sequence should work
	SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{sendMsg1, sendMsg2}, []uint64{0}, []uint64{1}, true, true, priv1)
}

func TestMsgSendMultipleOut(t *testing.T) {
	mapp := getMockApp(t)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42000000000)},
	}
	acc2 := &auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42000000000)},
	}

	SetGenesis(mapp, []auth.Account{acc1, acc2})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{sendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 41996999999)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42000000001)}},
			},
		},
	}

	for _, tc := range testCases {
		SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestMsgSendDependent(t *testing.T) {
	mapp := getMockApp(t)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 42000000000000)},
	}

	SetGenesis(mapp, []auth.Account{acc1})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{sendMsg1},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 41999996999990)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)}},
			},
		},
		{
			msgs:       []sdk.Msg{sendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{1},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 41999993999989)}},
			},
		},
	}

	for _, tc := range testCases {
		SignCheckDeliver(t, mapp.BaseApp, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}
