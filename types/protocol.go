package types

import (
	"fmt"
	"os"

	"github.com/orientwalt/htdf/codec"
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

const (
	AppVersionTag = "app_version"
	MainStore     = "main"
)

var (
	UpgradeConfigKey     = []byte("upgrade_config")
	CurrentVersionKey    = []byte("current_version")
	LastFailedVersionKey = []byte("last_failed_version")
	cdc                  = codec.New()
)

type ProtocolDefinition struct {
	Version   uint64 `json:"version"`
	Software  string `json:"software"`
	Height    uint64 `json:"height"`
	Threshold Dec    `json:"threshold"`
}

type UpgradeConfig struct {
	ProposalID uint64
	Protocol   ProtocolDefinition
}

func (uc UpgradeConfig) String() string {
	return fmt.Sprintf("proposalID: %v, version: %v, software: %s, height: %v, threshold: %s",
		uc.ProposalID, uc.Protocol.Version, uc.Protocol.Software, uc.Protocol.Height, uc.Protocol.Threshold.String(),
	)
}

func NewProtocolDefinition(version uint64, software string, height uint64, threshold Dec) ProtocolDefinition {
	return ProtocolDefinition{
		version,
		software,
		height,
		threshold,
	}
}

func NewUpgradeConfig(proposalID uint64, protocol ProtocolDefinition) UpgradeConfig {
	return UpgradeConfig{
		proposalID,
		protocol,
	}
}

func DefaultUpgradeConfig(protocolVersion uint64, software string) UpgradeConfig {
	return UpgradeConfig{
		ProposalID: uint64(0),
		Protocol:   NewProtocolDefinition(protocolVersion, software, uint64(1), NewDecWithPrec(9, 1)),
	}
}

type ProtocolKeeper struct {
	storeKey StoreKey
	cdc      *codec.Codec
}

func NewProtocolKeeper(key StoreKey) ProtocolKeeper {
	return ProtocolKeeper{key, cdc}
}

func (pk ProtocolKeeper) GetCurrentVersionByStore(store KVStore) uint64 {
	bz := store.Get(CurrentVersionKey)
	if bz == nil {
		return 0
	}
	var currentVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentVersion)
	return currentVersion
}

func (pk ProtocolKeeper) GetUpgradeConfigByStore(store KVStore) (upgradeConfig UpgradeConfig, found bool) {
	bz := store.Get(UpgradeConfigKey)
	if bz == nil {
		return upgradeConfig, false
	}
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

func (pk ProtocolKeeper) GetCurrentVersion(ctx Context) uint64 {
	store := ctx.KVStore(pk.storeKey)
	log.Debugln("!------1--------GetCurrentVersion---------------------!")
	bz := store.Get(CurrentVersionKey)
	log.Debugln("!------2--------GetCurrentVersion---------------------!")
	if bz == nil {
		return 0
	}
	log.Debugln("!------3--------GetCurrentVersion---------------------!")
	var currentVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &currentVersion)
	log.Debugln("!------4--------GetCurrentVersion---------------------!", currentVersion)
	return currentVersion
}

func (pk ProtocolKeeper) SetCurrentVersion(ctx Context, currentVersion uint64) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(currentVersion)
	store.Set(CurrentVersionKey, bz)
}

func (pk ProtocolKeeper) GetLastFailedVersion(ctx Context) uint64 {
	store := ctx.KVStore(pk.storeKey)
	bz := store.Get(LastFailedVersionKey)
	if bz == nil {
		return 0 // default value
	}
	var lastFailedVersion uint64
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lastFailedVersion)
	return lastFailedVersion
}

func (pk ProtocolKeeper) SetLastFailedVersion(ctx Context, lastFailedVersion uint64) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(lastFailedVersion)
	store.Set(LastFailedVersionKey, bz)
}

func (pk ProtocolKeeper) GetUpgradeConfig(ctx Context) (upgradeConfig UpgradeConfig, found bool) {
	store := ctx.KVStore(pk.storeKey)
	bz := store.Get(UpgradeConfigKey)
	if bz == nil {
		return upgradeConfig, false
	}
	pk.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &upgradeConfig)
	return upgradeConfig, true
}

func (pk ProtocolKeeper) SetUpgradeConfig(ctx Context, upgradeConfig UpgradeConfig) {
	store := ctx.KVStore(pk.storeKey)
	bz := pk.cdc.MustMarshalBinaryLengthPrefixed(upgradeConfig)
	store.Set(UpgradeConfigKey, bz)
}

func (pk ProtocolKeeper) ClearUpgradeConfig(ctx Context) {
	store := ctx.KVStore(pk.storeKey)
	store.Delete(UpgradeConfigKey)
}

func (pk ProtocolKeeper) IsValidVersion(ctx Context, version uint64) bool {
	log.Debugln("2--------------handleMsgSubmitSoftwareUpgradeProposal---------------------")
	currentVersion := pk.GetCurrentVersion(ctx)

	lastFailedVersion := pk.GetLastFailedVersion(ctx)
	return isValidVersion(currentVersion, lastFailedVersion, version)
}

func isValidVersion(currentVersion uint64, lastFailedVersion uint64, version uint64) bool {
	if currentVersion >= lastFailedVersion {
		return currentVersion+1 == version
	} else {
		return lastFailedVersion == version || lastFailedVersion+1 == version
	}
}
