package v2

import (
	"encoding/json"
	"log"

	"github.com/orientwalt/htdf/app/protocol"
	"github.com/orientwalt/htdf/app/v1/slashing"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	"github.com/orientwalt/htdf/x/crisis"
	distr "github.com/orientwalt/htdf/x/distribution"
	"github.com/orientwalt/htdf/x/gov"
	"github.com/orientwalt/htdf/x/guardian"
	"github.com/orientwalt/htdf/x/mint"
	"github.com/orientwalt/htdf/x/service"
	stake "github.com/orientwalt/htdf/x/staking"
	"github.com/orientwalt/htdf/x/upgrade"
	tmtypes "github.com/tendermint/tendermint/types"
)

// export the state of htdf for a genesis file
func (p *ProtocolV2) ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (
	appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	if forZeroHeight {
		p.prepForZeroHeightGenesis(ctx, jailWhiteList)
	}

	// iterate to get the accounts
	accounts := []GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	p.accountMapper.IterateAccounts(ctx, appendAccount)
	// fileAccounts := []GenesisAccount{}
	// for _, acc := range accounts {
	// 	if acc.Coins == nil {
	// 		continue
	// 	}
	// 	var coinsString []string
	// 	for _, coin := range acc.Coins {
	// 		coinsString = append(coinsString, coin.String())
	// 	}
	// 	fileAccounts = append(fileAccounts,
	// 		GenesisAccount{
	// 			Address:       acc.Address,
	// 			Coins:         acc.Coins,
	// 			Sequence:      acc.Sequence,
	// 			AccountNumber: acc.AccountNumber,
	// 		})
	// }

	genState := NewGenesisFileStateEx(
		accounts,
		auth.ExportGenesis(ctx, p.accountMapper, p.feeCollectionKeeper),
		stake.ExportGenesis(ctx, p.StakeKeeper),
		mint.ExportGenesis(ctx, p.mintKeeper),
		distr.ExportGenesis(ctx, p.distrKeeper),
		gov.ExportGenesis(ctx, p.govKeeper),
		upgrade.ExportGenesis(ctx),
		service.ExportGenesis(ctx, p.serviceKeeper),
		guardian.ExportGenesis(ctx, p.guardianKeeper),
		slashing.ExportGenesis(ctx, p.slashingKeeper),
		crisis.ExportGenesis(ctx, p.crisisKeeper),
	)
	appState, err = codec.MarshalJSONIndent(p.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = stake.WriteValidators(ctx, p.StakeKeeper)
	return appState, validators, nil
}

// prepare for fresh start at zero height
func (p *ProtocolV2) prepForZeroHeightGenesis(ctx sdk.Context, jailWhiteList []string) {

	applyWhiteList := false

	//Check if there is a whitelist
	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, addr := range jailWhiteList {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			log.Fatal(err)
		}
		whiteListMap[addr] = true
	}

	/* Just to be safe, assert the invariants on current state. */
	p.assertRuntimeInvariants(ctx)

	/* Handle fee distribution state. */

	// withdraw all validator commission
	p.StakeKeeper.IterateValidators(ctx, func(_ int64, val sdk.Validator) (stop bool) {
		_, _ = p.distrKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		return false
	})

	// withdraw all delegator rewards
	dels := p.StakeKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		_, _ = p.distrKeeper.WithdrawDelegationRewards(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	}

	// clear validator slash events
	p.distrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	p.distrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators
	p.StakeKeeper.IterateValidators(ctx, func(_ int64, val sdk.Validator) (stop bool) {

		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps := p.distrKeeper.GetValidatorOutstandingRewards(ctx, val.GetOperator())
		feePool := p.distrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps)
		p.distrKeeper.SetFeePool(ctx, feePool)

		p.distrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	// reinitialize all delegations
	for _, del := range dels {
		p.distrKeeper.Hooks().BeforeDelegationCreated(ctx, del.DelegatorAddress, del.ValidatorAddress)
		p.distrKeeper.Hooks().AfterDelegationModified(ctx, del.DelegatorAddress, del.ValidatorAddress)
	}

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	p.StakeKeeper.IterateRedelegations(ctx, func(_ int64, red stake.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		p.StakeKeeper.SetRedelegation(ctx, red)
		return false
	})

	// iterate through unbonding delegations, reset creation height
	p.StakeKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stake.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		p.StakeKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(protocol.KeyStake)
	iter := sdk.KVStoreReversePrefixIterator(store, stake.ValidatorsKey)
	counter := int16(0)

	var valConsAddrs []sdk.ConsAddress
	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])
		validator, found := p.StakeKeeper.GetValidator(ctx, addr)
		if !found {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		valConsAddrs = append(valConsAddrs, validator.ConsAddress())
		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		p.StakeKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	_ = p.StakeKeeper.ApplyAndReturnValidatorSetUpdates(ctx)

	/* Handle slashing state. */

	// reset start height on signing infos
	p.slashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashing.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			p.slashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
	/* Handle gov state. */

	gov.PrepForZeroHeightGenesis(ctx, p.govKeeper)

	/* Handle service state. */
	service.PrepForZeroHeightGenesis(ctx, p.serviceKeeper)
}
