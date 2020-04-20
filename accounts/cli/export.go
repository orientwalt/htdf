package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/orientwalt/htdf/crypto/keys/mintkey"

	"github.com/orientwalt/htdf/accounts/keystore"
	"github.com/spf13/cobra"
)

func GetExportPivKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Export all private key list",
		Long:  "export private key from .hscli/keystores",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())

			priv, err := GetPrivateKey(ksw, args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("%s	%s\n", args[0], priv)
			return nil
		},
	}
}

func GetPrivateKey(ksw *keystore.KeyStoreWallet, addr string, passphrase string) (string, error) {
	privKeyArmor, err := ksw.GetPrivKey(addr)
	if err != nil {
		return "", err
	}
	privKey, err := mintkey.UnarmorDecryptPrivKey(privKeyArmor, passphrase)
	if err != nil {
		return "", err
	}
	strPrivKey := hex.EncodeToString(privKey.Bytes())
	priv := strPrivKey[10:]
	return priv, err
}
