package slashing

import (
	"fmt"
	"os"
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/params"
	"github.com/orientwalt/htdf/x/slashing/types"
	stake "github.com/orientwalt/htdf/x/staking/types"
	log "github.com/sirupsen/logrus"
)

func init() {
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

// Keeper of the slashing store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	validatorSet sdk.ValidatorSet
	paramspace   params.Subspace

	// codespace
	codespace sdk.CodespaceType
	// metrics
	metrics *Metrics
}

// NewKeeper creates a slashing keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, vs sdk.ValidatorSet, paramspace params.Subspace, codespace sdk.CodespaceType, metrics *Metrics) Keeper {
	keeper := Keeper{
		storeKey:     key,
		cdc:          cdc,
		validatorSet: vs,
		paramspace:   paramspace.WithKeyTable(types.ParamKeyTable()),
		codespace:    codespace,
		metrics:      metrics,
	}
	return keeper
}

// handle a validator signing two blocks at the same height
// power: power of the double-signing validator at the height of infraction
func (k Keeper) handleDoubleSign(ctx sdk.Context, addr crypto.Address, infractionHeight int64, timestamp time.Time, power int64) {
	logger := ctx.Logger().With("module", "x/slashing")

	// calculate the age of the evidence
	time := ctx.BlockHeader().Time
	age := time.Sub(timestamp)

	// fetch the validator public key
	consAddr := sdk.ConsAddress(addr)
	pubkey, err := k.getPubkey(ctx, addr)
	if err != nil {
		// Ignore evidence that cannot be handled.
		// NOTE:
		// We used to panic with:
		// `panic(fmt.Sprintf("Validator consensus-address %v not found", consAddr))`,
		// but this couples the expectations of the app to both Tendermint and
		// the simulator.  Both are expected to provide the full range of
		// allowable but none of the disallowed evidence types.  Instead of
		// getting this coordination right, it is easier to relax the
		// constraints and ignore evidence that cannot be handled.
		return
	}

	// Reject evidence if the double-sign is too old
	if age > k.MaxEvidenceAge(ctx) {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, age of %d past max age of %d",
			pubkey.Address(), infractionHeight, age, k.MaxEvidenceAge(ctx)))
		return
	}

	// Get validator and signing info
	validator := k.validatorSet.ValidatorByConsAddr(ctx, consAddr)
	if validator == nil || validator.GetStatus() == sdk.Unbonded {
		// Defensive.
		// Simulation doesn't take unbonding periods into account, and
		// Tendermint might break this assumption at some point.
		return
	}

	// fetch the validator signing info
	signInfo, found := k.getValidatorSigningInfo(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", consAddr))
	}

	// validator is already tombstoned
	if signInfo.Tombstoned {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, validator already tombstoned", pubkey.Address(), infractionHeight))
		return
	}

	// double sign confirmed
	logger.Info(fmt.Sprintf("Confirmed double sign from %s at height %d, age of %d", pubkey.Address(), infractionHeight, age))

	// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height.
	// Note that this *can* result in a negative "distributionHeight", up to -ValidatorUpdateDelay,
	// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
	// That's fine since this is just used to filter unbonding delegations & redelegations.
	distributionHeight := infractionHeight - sdk.ValidatorUpdateDelay

	// get the percentage slash penalty fraction
	fraction := k.SlashFractionDoubleSign(ctx)

	// Slash validator
	// `power` is the int64 power of the validator as provided to/by
	// Tendermint. This value is validator.Tokens as sent to Tendermint via
	// ABCI, and now received as evidence.
	// The fraction is passed in to separately to slash unbonding and rebonding delegations.
	k.validatorSet.Slash(ctx, consAddr, distributionHeight, power, fraction)

	// Jail validator if not already jailed
	// begin unbonding validator if not already unbonding (tombstone)
	if !validator.IsJailed() {
		k.validatorSet.Jail(ctx, consAddr)
	}

	// Set tombstoned to be true
	signInfo.Tombstoned = true

	// Set jailed until to be forever (max time)
	signInfo.JailedUntil = types.DoubleSignJailEndTime

	// Set validator signing info
	k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
}

// handle a validator signature, must be called once per validator per block
// TODO refactor to take in a consensus address, additionally should maybe just take in the pubkey too
func (k Keeper) handleValidatorSignature(ctx sdk.Context, addr crypto.Address, power int64, signed bool) {
	logger := ctx.Logger().With("module", "x/slashing")
	height := ctx.BlockHeight()
	consAddr := sdk.ConsAddress(addr)
	pubkey, err := k.getPubkey(ctx, addr)
	if err != nil {
		panic(fmt.Sprintf("Validator consensus-address %v not found", consAddr))
	}

	// fetch signing info
	signInfo, found := k.getValidatorSigningInfo(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", consAddr))
	}

	// this is a relative index, so it counts blocks the validator *should* have signed
	// will use the 0-value default signing info if not present, except for start height
	index := signInfo.IndexOffset % k.SignedBlocksWindow(ctx)
	signInfo.IndexOffset++

	// Update signed block bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.getValidatorMissedBlockBitArray(ctx, consAddr, index)
	missed := !signed
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed, increment counter
		k.setValidatorMissedBlockBitArray(ctx, consAddr, index, true)
		signInfo.MissedBlocksCounter++
	case previous && !missed:
		// Array value has changed from missed to not missed, decrement counter
		k.setValidatorMissedBlockBitArray(ctx, consAddr, index, false)
		signInfo.MissedBlocksCounter--
	default:
		// Array value at this index has not changed, no need to update counter
	}

	if missed {
		logger.Info(fmt.Sprintf("Absent validator %s (%v) at height %d, %d missed, threshold %d", addr, pubkey, height, signInfo.MissedBlocksCounter, k.MinSignedPerWindow(ctx)))
	}

	minHeight := signInfo.StartHeight + k.SignedBlocksWindow(ctx)
	maxMissed := k.SignedBlocksWindow(ctx) - k.MinSignedPerWindow(ctx)

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		validator := k.validatorSet.ValidatorByConsAddr(ctx, consAddr)
		if validator != nil && !validator.IsJailed() {

			// Downtime confirmed: slash and jail the validator
			logger.Info(fmt.Sprintf("Validator %s past min height of %d and below signed blocks threshold of %d",
				pubkey.Address(), minHeight, k.MinSignedPerWindow(ctx)))

			// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height,
			// and subtract an additional 1 since this is the LastCommit.
			// Note that this *can* result in a negative "distributionHeight" up to -ValidatorUpdateDelay-1,
			// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
			// That's fine since this is just used to filter unbonding delegations & redelegations.

			// move from app/v1/slashing to here , yqq , 2021-04-08
			// Disable Slashing for ValidatorSignatureMissing. junying-todo, 2020-05-26
			// distributionHeight := height - sdk.ValidatorUpdateDelay - 1
			// k.validatorSet.Slash(ctx, consAddr, distributionHeight, power, k.SlashFractionDowntime(ctx))

			log.Infof("No Slashing For ValidatorSignatureMissing(Jailed:%s)\n", pubkey.Address())
			k.validatorSet.Jail(ctx, consAddr)
			signInfo.JailedUntil = ctx.BlockHeader().Time.Add(k.DowntimeJailDuration(ctx))

			// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
			signInfo.MissedBlocksCounter = 0
			signInfo.IndexOffset = 0
			k.clearValidatorMissedBlockBitArray(ctx, consAddr)
		} else {
			// Validator was (a) not found or (b) already jailed, don't slash
			logger.Info(fmt.Sprintf("Validator %s would have been slashed for downtime, but was either not found in store or already jailed",
				pubkey.Address()))
		}
	}

	// Set the updated signing info
	k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
}

// Punish proposer censorship by slashing malefactor's stake
func (k Keeper) handleProposerCensorship(ctx sdk.Context, addr crypto.Address, infractionHeight int64) (tags sdk.Tags) {
	logger := ctx.Logger()
	time := ctx.BlockHeader().Time
	consAddr := sdk.ConsAddress(addr)
	_, err := k.getPubkey(ctx, addr)
	if err != nil {
		panic(fmt.Sprintf("Validator consensus-address %v not found", consAddr))
	}

	// Get validator.
	validator := k.validatorSet.ValidatorByConsAddr(ctx, consAddr)
	if validator == nil || validator.GetStatus() == sdk.Unbonded {
		// Defensive.
		// Simulation doesn't take unbonding periods into account, and
		// Tendermint might break this assumption at some point.
		return
	}
	logger.Info("The malefactor proposer proposed a invalid block",
		"proposer_address", validator.GetOperator().String(),
		"block_height", ctx.BlockHeight(), "consensus_address", consAddr.String())

	distributionHeight := infractionHeight - stake.ValidatorUpdateDelay
	// Slash validator
	// `power` is the int64 power of the validator as provided to/by
	// Tendermint. This value is validator.Tokens as sent to Tendermint via
	// ABCI, and now received as evidence.
	// The revisedFraction (which is the new fraction to be slashed) is passed
	// in separately to separately slash unbonding and rebonding delegations.
	tags = k.validatorSet.Slash(ctx, consAddr, distributionHeight, validator.GetBondedTokens().Int64(), k.SlashFractionCensorship(ctx))

	// Jail validator if not already jailed
	if !validator.IsJailed() {
		k.validatorSet.Jail(ctx, consAddr)
	}

	// Set or updated validator jail duration
	signInfo, found := k.getValidatorSigningInfo(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", consAddr))
	}
	signInfo.JailedUntil = time.Add(k.CensorshipJailDuration(ctx))
	k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
	return
}

func (k Keeper) addPubkey(ctx sdk.Context, pubkey crypto.PubKey) {
	addr := pubkey.Address()
	k.setAddrPubkeyRelation(ctx, addr, pubkey)
}

func (k Keeper) getPubkey(ctx sdk.Context, address crypto.Address) (crypto.PubKey, error) {
	store := ctx.KVStore(k.storeKey)
	var pubkey crypto.PubKey
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.GetAddrPubkeyRelationKey(address)), &pubkey)
	if err != nil {
		return nil, fmt.Errorf("address %v not found", address)
	}
	return pubkey, nil
}

func (k Keeper) setAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address, pubkey crypto.PubKey) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pubkey)
	store.Set(types.GetAddrPubkeyRelationKey(addr), bz)
}

func (k Keeper) deleteAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAddrPubkeyRelationKey(addr))
}

// MaxEvidenceAge - max age for evidence
func (k Keeper) MaxEvidenceAge(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, types.KeyMaxEvidenceAge, &res)
	return
}

// SignedBlocksWindow - sliding window for downtime slashing
func (k Keeper) SignedBlocksWindow(ctx sdk.Context) (res int64) {
	k.paramspace.Get(ctx, types.KeySignedBlocksWindow, &res)
	return
}

// Downtime slashing threshold
func (k Keeper) MinSignedPerWindow(ctx sdk.Context) int64 {
	var minSignedPerWindow sdk.Dec
	k.paramspace.Get(ctx, types.KeyMinSignedPerWindow, &minSignedPerWindow)
	signedBlocksWindow := k.SignedBlocksWindow(ctx)

	// NOTE: RoundInt64 will never panic as minSignedPerWindow is
	//       less than 1.
	return minSignedPerWindow.MulInt64(signedBlocksWindow).RoundInt64()
}

// Downtime unbond duration
func (k Keeper) DowntimeJailDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, types.KeyDowntimeJailDuration, &res)
	return
}

// SlashFractionDoubleSign
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, types.KeySlashFractionDoubleSign, &res)
	return
}

// Censorship jail duration
func (k Keeper) CensorshipJailDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, types.KeyCensorshipJailDuration, &res)
	return
}

// SlashFractionDowntime
func (k Keeper) SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, types.KeySlashFractionDowntime, &res)
	return
}

// Slash fraction for Censorship
func (k Keeper) SlashFractionCensorship(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, types.KeySlashFractionCensorship, &res)
	return
}

// GetParams returns the total set of slashing parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

/////////////////////
// Stored by *validator* address (not operator address)
func (k Keeper) getValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress) (info types.ValidatorSigningInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorSigningInfoKey(address))
	if bz == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
	found = true
	return
}

// Stored by *validator* address (not operator address)
func (k Keeper) IterateValidatorSigningInfos(ctx sdk.Context, handler func(address sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSigningInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		address := types.GetValidatorSigningInfoAddress(iter.Key())
		var info types.ValidatorSigningInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &info)
		if handler(address, info) {
			break
		}
	}
}

// Stored by *validator* address (not operator address)
func (k Keeper) SetValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress, info types.ValidatorSigningInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(info)
	store.Set(types.GetValidatorSigningInfoKey(address), bz)
}

// Stored by *validator* address (not operator address)
func (k Keeper) getValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64) (missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorMissedBlockBitArrayKey(address, index))
	if bz == nil {
		// lazy: treat empty key as not missed
		missed = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
	return
}

// Stored by *validator* address (not operator address)
func (k Keeper) IterateValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, handler func(index int64, missed bool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	index := int64(0)
	// Array may be sparse
	for ; index < k.SignedBlocksWindow(ctx); index++ {
		var missed bool
		bz := store.Get(types.GetValidatorMissedBlockBitArrayKey(address, index))
		if bz == nil {
			continue
		}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
		if handler(index, missed) {
			break
		}
	}
}

// Stored by *validator* address (not operator address)
func (k Keeper) setValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(missed)
	store.Set(types.GetValidatorMissedBlockBitArrayKey(address, index), bz)
}

// Stored by *validator* address (not operator address)
func (k Keeper) clearValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorMissedBlockBitArrayPrefixKey(address))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
