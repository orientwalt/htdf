package types

import (
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	GetAllAccounts(ctx sdk.Context) (accounts []auth.Account)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	SetAccount(ctx sdk.Context, account auth.Account)
	RemoveAccount(ctx sdk.Context, account auth.Account)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	// InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error)

	DelegateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	UndelegateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)

	IncreaseLoosenToken(ctx sdk.Context, amt sdk.Coins)
	DecreaseLoosenToken(ctx sdk.Context, amt sdk.Coins)
	BurnCoinsFromAddr(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	BurnCoinsFromPool(ctx sdk.Context, pool string, amt sdk.Coins) (sdk.Tags, sdk.Error)
	// GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	// SetBalance(ctx sdk.Context, addr sdk.AccAddress, balance sdk.Coin) error
}

// expected fee collection keeper interface
type FeeCollectionKeeper interface {
	AddCollectedFees(sdk.Context, sdk.Coins) sdk.Coins
}
