package app

import (
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth/ante"
	authkeeper "github.com/orientwalt/htdf/x/auth/keeper"
	"github.com/orientwalt/htdf/x/auth/signing"
	bankkeeper "github.com/orientwalt/htdf/x/bank/keeper"

	oraclekeeper "github.com/orientwalt/htdf/modules/oracle/keeper"
	oracletypes "github.com/orientwalt/htdf/modules/oracle/types"
	tokenkeeper "github.com/orientwalt/htdf/modules/token/keeper"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	tk tokenkeeper.Keeper,
	ok oraclekeeper.Keeper,
	oak oracletypes.AuthKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, bk),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
		tokenkeeper.NewValidateTokenFeeDecorator(tk, bk),
		oraclekeeper.NewValidateOracleAuthDecorator(ok, oak),
	)
}
