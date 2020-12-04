// nolint
package auth

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/params"
	xauth "github.com/orientwalt/htdf/x/auth"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	ak  xauth.AccountKeeper
	fck xauth.FeeCollectionKeeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	xauth.RegisterBaseAccount(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := xauth.NewAccountKeeper(cdc, authCapKey, pk.Subspace(xauth.DefaultParamspace), xauth.ProtoBaseAccount)
	fck := xauth.NewFeeCollectionKeeper(cdc, fckCapKey)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	ak.SetParams(ctx, xauth.DefaultParams())

	return testInput{cdc: cdc, ctx: ctx, ak: ak, fck: fck}
}

func newTestMsg(addrs ...sdk.AccAddress) *sdk.TestMsg {
	return sdk.NewTestMsg(addrs...)
}

func newStdFee() xauth.StdFee {
	return xauth.NewStdFee(2000000, 100)
}

// coins to more than cover the fee
func newCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin("satoshi", 1000000000),
	}
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func newTestTx(ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee xauth.StdFee) sdk.Tx {
	sigs := make([]xauth.StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := xauth.StdSignBytes(ctx.ChainID(), accNums[i], seqs[i], fee, msgs, "")

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = xauth.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := xauth.NewStdTx(msgs, fee, sigs, "")
	return tx
}

func newTestTxWithMemo(ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee xauth.StdFee, memo string) sdk.Tx {
	sigs := make([]xauth.StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := xauth.StdSignBytes(ctx.ChainID(), accNums[i], seqs[i], fee, msgs, memo)

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = xauth.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := xauth.NewStdTx(msgs, fee, sigs, memo)
	return tx
}

func newTestTxWithSignBytes(msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee xauth.StdFee, signBytes []byte, memo string) sdk.Tx {
	sigs := make([]xauth.StdSignature, len(privs))
	for i, priv := range privs {
		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = xauth.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := xauth.NewStdTx(msgs, fee, sigs, memo)
	return tx
}
