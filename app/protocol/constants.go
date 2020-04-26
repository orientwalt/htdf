package protocol

import sdk "github.com/orientwalt/htdf/types"

const (
	// all store name
	AccountStore         = "acc"
	CoreStore            = "hs"
	StakeStore           = "staking"
	StakeTransientStore  = "transient_stake"
	MintStore            = "mint"
	DistrStore           = "distr"
	DistrTransientStore  = "transient_distr"
	SlashingStore        = "slashing"
	GovStore             = "gov"
	FeeStore             = "fee"
	ParamsStore          = "params"
	ParamsTransientStore = "transient_params"
	ServiceStore         = "service"
	GuardianStore        = "guardian"
	UpgradeStore         = "upgrade"
	Storage              = "storage"
	Code                 = "code"

	// all route for query and handler
	BankRoute     = "bank"
	AccountRoute  = AccountStore
	CoreRoute     = CoreStore
	StakeRoute    = StakeStore
	DistrRoute    = DistrStore
	SlashingRoute = SlashingStore
	GovRoute      = GovStore
	ParamsRoute   = ParamsStore
	ServiceRoute  = ServiceStore
	GuardianRoute = GuardianStore
	UpgradeRoute  = UpgradeStore
	MintRoute     = MintStore // junying-todo, 2020-02-05
)

var (
	KeyMain     = sdk.NewKVStoreKey(sdk.MainStore)
	KeyAccount  = sdk.NewKVStoreKey(AccountStore)
	KeyStake    = sdk.NewKVStoreKey(StakeStore)
	TkeyStake   = sdk.NewTransientStoreKey(StakeTransientStore)
	KeyMint     = sdk.NewKVStoreKey(MintStore)
	KeyDistr    = sdk.NewKVStoreKey(DistrStore)
	TkeyDistr   = sdk.NewTransientStoreKey(DistrTransientStore)
	KeySlashing = sdk.NewKVStoreKey(SlashingStore)
	KeyGov      = sdk.NewKVStoreKey(GovStore)
	KeyFee      = sdk.NewKVStoreKey(FeeStore)
	KeyParams   = sdk.NewKVStoreKey(ParamsStore)
	TkeyParams  = sdk.NewTransientStoreKey(ParamsTransientStore)
	KeyService  = sdk.NewKVStoreKey(ServiceStore)
	KeyGuardian = sdk.NewKVStoreKey(GuardianStore)
	KeyUpgrade  = sdk.NewKVStoreKey(UpgradeStore)
	KeyStorage  = sdk.NewKVStoreKey(Storage) // junying-todo
	KeyCode     = sdk.NewKVStoreKey(Code)    // junying-todo
)
