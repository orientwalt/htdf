package types

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

var _ auth.Account = (*Account)(nil)

// ----------------------------------------------------------------------------
// Main Ethermint account
// ----------------------------------------------------------------------------

// BaseAccount implements the auth.Account interface and embeds an
// auth.BaseAccount type. It is compatible with the auth.AccountMapper.
type Account struct {
	auth.Account

	// merkle root of the storage trie
	//
	// TODO: good chance we may not need this
	//Root ethcmn.Hash

	CodeHash []byte
}

// ProtoBaseAccount defines the prototype function for BaseAccount used for an
// account mapper.
func ProtoBaseAccount() auth.Account {
	return &Account{Account: &auth.BaseAccount{}}
}

// Balance returns the balance of an account.
func (acc Account) Balance() sdk.Int {
	return acc.GetCoins().AmountOf(sdk.DefaultDenom)
}

// SetBalance sets an account's balance.
func (acc Account) SetBalance(amt sdk.Int) {
	oldCoins := acc.GetCoins()
	var newCoins sdk.Coins

	bUpdateDenom := false
	for _, coin := range oldCoins {
		if coin.Denom == sdk.DefaultDenom {
			coin = sdk.NewCoin(sdk.DefaultDenom, amt)
			bUpdateDenom = true
		}
		newCoins = append(newCoins, coin)
	}

	//insert new denom;  evm use "satoshi"
	if bUpdateDenom == false {
		newCoins = append(newCoins, sdk.NewCoin(sdk.DefaultDenom, amt))
	}

	acc.SetCoins(newCoins)
}

func NewAccount(account auth.Account) *Account {
	return &Account{Account: account}
}

// ----------------------------------------------------------------------------
// Code & Storage
// ----------------------------------------------------------------------------

// Account code and storage type aliases.
type (
	Code    []byte
	Storage map[ethcmn.Hash]ethcmn.Hash
)

func (c Code) String() string {
	return string(c)
}

func (c Storage) String() (str string) {
	for key, value := range c {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

// Copy returns a copy of storage.
func (c Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range c {
		cpy[key] = value
	}

	return cpy
}
