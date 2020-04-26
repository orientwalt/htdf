package keystore

import (
	"fmt"

	"path/filepath"

	"github.com/orientwalt/htdf/accounts"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"

	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

var _ accounts.KeyStoreWallets = (*KeyStoreWallet)(nil)

type DefaultKeyStorePath func() string

var DefaultKeyStoreHome = defaultKeyStoreHome()

func defaultKeyStoreHome() DefaultKeyStorePath {
	return func() string {
		rootDir := viper.GetString(cli.HomeFlag)
		defaultKeyStoreHome := filepath.Join(rootDir, "keystores")
		return defaultKeyStoreHome
	}
}

type KeyStoreWallet struct {
	keyStore *KeyStore
	scan     *scaner
}

func NewKeyStoreWallet(path string) *KeyStoreWallet {
	ksw := &KeyStoreWallet{
		keyStore: NewKeyStore(path),
		scan:     newScaner(path),
	}

	return ksw
}

func (ksw *KeyStoreWallet) Accounts() ([]accounts.Account, error) {
	accounts, err := ksw.scan.accounts()
	if err != nil {
		return nil, err
	}

	return accounts, err
}

func (ksw *KeyStoreWallet) GetPrivKey(addr string) (string, error) {
	key, err := ksw.scan.getSigner(addr)
	if err != nil {
		return "", err
	}

	return key.PrivKey, err
}

func (ksw *KeyStoreWallet) Update(addr string, passphrase, newPassphrase string) error {
	key, err := ksw.scan.getSigner(addr)
	if err != nil {
		return err
	}

	ksw.keyStore.key = key

	account := accounts.Account{Address: addr}
	acc, err := ksw.scan.find(account)
	if err != nil {
		return err
	}

	err = ksw.keyStore.update(acc, passphrase, newPassphrase)
	return err
}

func (ksw *KeyStoreWallet) Drop(addr string) error {

	found := ksw.scan.hasAddress(addr)
	if found {
		account := accounts.Account{Address: addr}
		acc, err := ksw.scan.find(account)
		if err != nil {
			return err
		}

		err = ksw.keyStore.drop(acc)
		return err
	}

	return ErrNoMatch
}

func (ksw *KeyStoreWallet) BuildAndSign(txbuilder authtxb.TxBuilder, addr string, passphrase string, msgs []sdk.Msg) ([]byte, error) {
	msg, err := BuildSignMsg(txbuilder, msgs)
	if err != nil {
		return nil, err
	}

	return ksw.Sign(txbuilder, addr, passphrase, msg)
}

func BuildSignMsg(txbuilder authtxb.TxBuilder, msgs []sdk.Msg) (authtxb.StdSignMsg, error) {
	fmt.Println("BuildSignMsg:txbuilder.GasWanted()", txbuilder.GasWanted())
	chainID := txbuilder.ChainID()
	if chainID == "" {
		return authtxb.StdSignMsg{}, fmt.Errorf("chain ID required but not specified")
	}
	// junying-todo, 2019-11-08
	// converted from fee based to gas*gasprice based
	// if txbuilder.GasPrices().IsZero() {
	// 	return authtxb.StdSignMsg{}, errors.New("gasprices can't not be zero")
	// }
	// if txbuilder.GasWanted() <= 0 {
	// 	return authtxb.StdSignMsg{}, errors.New("gasWanted must be provided")
	// }
	fmt.Println("BuildSignMsg:Fee", auth.NewStdFee(txbuilder.GasWanted(), txbuilder.GasPrice()), txbuilder.GasWanted())
	return authtxb.StdSignMsg{
		ChainID:       txbuilder.ChainID(),
		AccountNumber: txbuilder.AccountNumber(),
		Sequence:      txbuilder.Sequence(),
		Memo:          txbuilder.Memo(),
		Msgs:          msgs,
		Fee:           auth.NewStdFee(txbuilder.GasWanted(), txbuilder.GasPrice()), // auth.NewStdFee(txbuilder.GasWanted(), fees),
	}, nil
}

func (ksw *KeyStoreWallet) Sign(txbuilder authtxb.TxBuilder, addr string, passphrase string, msg authtxb.StdSignMsg) ([]byte, error) {
	key, err := ksw.scan.getSigner(addr)
	if err != nil {
		return nil, err
	}

	ksw.keyStore.key = key
	// fmt.Println(msg.Bytes())
	sig, err := ksw.makeSignature(passphrase, msg)
	if err != nil {
		return nil, err
	}

	en := txbuilder.TxEncoder()

	tx := auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo)

	return en(tx)
}

func (ksw *KeyStoreWallet) SignStdTx(txbuilder authtxb.TxBuilder, stdTx auth.StdTx, addr string, passphrase string) (signedStdTx auth.StdTx, err error) {
	key, err := ksw.scan.getSigner(addr)
	if err != nil {
		return
	}

	ksw.keyStore.key = key

	stdSignature, err := ksw.makeSignature(passphrase, authtxb.StdSignMsg{
		ChainID:       txbuilder.ChainID(),
		AccountNumber: txbuilder.AccountNumber(),
		Sequence:      txbuilder.Sequence(),
		Fee:           stdTx.Fee,
		Msgs:          stdTx.GetMsgs(),
		Memo:          stdTx.GetMemo(),
	})
	if err != nil {
		return
	}

	sigs := stdTx.GetSignatures()

	if len(sigs) == 0 {
		sigs = []auth.StdSignature{stdSignature}
	} else {
		sigs = append(sigs, stdSignature)
	}

	signedStdTx = auth.NewStdTx(stdTx.GetMsgs(), stdTx.Fee, sigs, stdTx.GetMemo())
	return
}

// MakeSignature builds a StdSignature given keybase, key name, passphrase, and a StdSignMsg.
func (ksw *KeyStoreWallet) makeSignature(passphrase string, msg authtxb.StdSignMsg) (sig auth.StdSignature, err error) {
	sigBytes, pubkey, err := ksw.keyStore.key.Sign(passphrase, msg.Bytes())
	if err != nil {
		return
	}

	return auth.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	}, nil
}
