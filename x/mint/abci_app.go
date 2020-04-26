package mint

import (
	"os"

	sdk "github.com/orientwalt/htdf/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	// This decodes a valid hex string into a sepc256k1Pubkey for use in transaction simulation
	// junying-todo,2020-01-17
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info" //trace/debug/info/warn/error/parse/fatal/panic
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.FatalLevel //TraceLevel/DebugLevel/InfoLevel/WarnLevel/ErrorLevel/ParseLevel/FatalLevel/PanicLevel
	}
	// set global log level
	log.SetLevel(ll)
	log.SetFormatter(&log.TextFormatter{}) //&log.JSONFormatter{})
}

// junying-todo, 2019-07-17
//	BlksPerRound = 100
//	rewards+commission+community-pool
//	hscli query distr rewards htdf1zulqmaqlsgrgmagenaqf02p8kfgsuqkdwgwj80
//	121793749706.0satoshi
//	* 4 = 487174998824
//	not true becasue proper get more rewards,that's, different rewards on every node.
//  hscli query distr commission cosmosvaloper1lwjmdnks33xwnmfayc64ycprww49n33mtm92ne
// 	hscli query distr community-pool

// junying-todo, 2019-07-15
// single node: 88.2 for delegators, 11.8 for validator(commission)
// commission is validating fee
// commission rate changes?
func calcParams(ctx sdk.Context, k Keeper) (sdk.Dec, sdk.Dec, sdk.Dec) {
	// fetch params
	totalSupply := k.sk.TotalTokens(ctx)
	log.Infoln("totalSupply", totalSupply)
	// block index
	curBlkHeight := ctx.BlockHeight()
	log.Infoln("curBlkHeight:", curBlkHeight)

	// check terminate condition, junying-todo, 2019-12-05
	if totalSupply.GT(sdk.NewInt(TotalLiquidAsSatoshi)) { // || curBlkHeight > TotalMineableBlks {
		return sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0)
	}

	// sine params
	curAmplitude := k.sk.Amplitude(ctx)
	curCycle := k.sk.Cycle(ctx)
	curLastIndex := k.sk.LastIndex(ctx)
	// check if it's time for new cycle
	if curBlkHeight >= (curLastIndex + curCycle) {
		k.sk.SetAmplitude(ctx, randomAmplitude(curBlkHeight))
		k.sk.SetCycle(ctx, randomCycle(curAmplitude))
		k.sk.SetLastIndex(ctx, curBlkHeight)
	}

	AnnualProvisionsDec, Inflation, BlockProvision := GetMineToken(curBlkHeight, totalSupply, curAmplitude, curCycle, curLastIndex)
	// junying-todo, 2020-02-04
	k.SetReward(ctx, curBlkHeight, BlockProvision.TruncateInt64())
	return AnnualProvisionsDec, Inflation, BlockProvision
}

// GetMineToken...
func GetMineToken(curBlkHeight int64, totalSupply sdk.Int, curAmplitude int64, curCycle int64, curLastIndex int64) (sdk.Dec, sdk.Dec, sdk.Dec) {

	// block index
	// log.Infoln("curBlkHeight:", curBlkHeight)

	AnnualProvisionsDec := sdk.NewDec(AnnualProvisionAsSatoshi)
	// Inflation = AnnualProvisions / totalSupply
	Inflation := AnnualProvisionsDec.Quo(sdk.NewDecFromInt(totalSupply))

	BlockReward := calcRewardAsSatoshi(curAmplitude, curCycle, curBlkHeight-curLastIndex)
	if BlockReward < 0 {
		panic(0)
	}
	BlockProvision := sdk.NewDec(BlockReward)
	log.Infoln("BlockProvision:", BlockReward)
	log.Infoln("curAmplitude:", curAmplitude)
	log.Infoln("curCycle:", curCycle)
	log.Infoln("curLastIndex:", curLastIndex)

	return AnnualProvisionsDec, Inflation, BlockProvision
}

// Inflate every block, update inflation parameters once per hour
func BeginBlocker(ctx sdk.Context, k Keeper) {

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	//params := k.GetParams(ctx)

	// recalculate inflation rate
	var provisionAmt sdk.Dec
	minter.AnnualProvisions, minter.Inflation, provisionAmt = calcParams(ctx, k)

	k.SetMinter(ctx, minter)

	// mint coins, add to collected fees, update supply
	//fmt.Printf("AnnualProvisions: %s, Inflation: %s, provisionAmt: %s\n", minter.AnnualProvisions.String(), minter.Inflation.String(), provisionAmt.TruncateInt().String())
	mintedCoin := sdk.NewCoin(sdk.DefaultDenom, provisionAmt.TruncateInt())
	k.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoin})
	k.sk.InflateSupply(ctx, mintedCoin.Amount)

}
