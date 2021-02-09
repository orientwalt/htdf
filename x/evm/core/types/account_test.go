package types

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"

	"testing"

	"github.com/magiconair/properties/assert"
)

const (
	strAddr = "htdf1gh8yeqhxx7n29fx3fuaksqpdsau7gxc0rteptt"
)

func init() {

	// set address prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()
}

func TestAccountUpdate(t *testing.T) {

	accAddr, err := sdk.AccAddressFromBech32(strAddr)
	assert.Equal(t, err, nil)

	baseAccount := auth.NewBaseAccountWithAddress(accAddr)
	coins := sdk.Coins{sdk.NewCoin(sdk.DefaultDenom, sdk.NewInt(12345678)), sdk.NewCoin("stake", sdk.NewInt(2345678))}
	baseAccount.Coins = coins

	fmt.Printf("baseAccount=%v\n", baseAccount)

	account := NewAccount(&baseAccount)

	fmt.Printf("account|coins=%v\n", account.GetCoins())
	account.SetBalance(sdk.NewInt(30000))

	fmt.Printf("account|coins=%v\n", account.GetCoins())

}

func TestAccountInsert(t *testing.T) {

	accAddr, err := sdk.AccAddressFromBech32(strAddr)
	assert.Equal(t, err, nil)

	baseAccount := auth.NewBaseAccountWithAddress(accAddr)
	account := NewAccount(&baseAccount)

	account.SetBalance(sdk.NewInt(3500350))
	fmt.Printf("account|coins=%v\n", account.GetCoins())
	assert.Equal(t, account.Balance().String() == "3500350", true)

}
