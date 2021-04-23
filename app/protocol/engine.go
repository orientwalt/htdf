package protocol

import (
	"fmt"
	"os"

	sdk "github.com/orientwalt/htdf/types"
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

type ProtocolEngine struct {
	protocols      map[uint64]Protocol
	current        uint64
	next           uint64
	ProtocolKeeper sdk.ProtocolKeeper
}

func NewProtocolEngine(protocolKeeper sdk.ProtocolKeeper) ProtocolEngine {
	engine := ProtocolEngine{
		make(map[uint64]Protocol),
		0,
		0,
		protocolKeeper,
	}
	return engine
}

func (pe *ProtocolEngine) LoadCurrentProtocol(kvStore sdk.KVStore) (bool, uint64) {
	// find the current version from store
	current := pe.ProtocolKeeper.GetCurrentVersionByStore(kvStore)
	p, flag := pe.protocols[current]
	log.Infoln("/---------flag1----------/", flag)
	if flag == true {
		p.Load()
		pe.current = current
	}
	return flag, current
}

// To be used for Protocol with version > 0
func (pe *ProtocolEngine) Activate(version uint64) bool {
	p, flag := pe.protocols[version]
	log.Infoln("/---------flag2----------/", flag)
	if flag == true {
		p.Load()
		p.Init()
		pe.current = version
	}
	return flag
}

func (pe *ProtocolEngine) GetCurrentProtocol() Protocol {
	return pe.protocols[pe.current]
}

func (pe *ProtocolEngine) GetCurrentVersion() uint64 {
	return pe.current
}

func (pe *ProtocolEngine) Add(p Protocol) Protocol {
	if p.GetVersion() != pe.next {
		panic(fmt.Errorf("Wrong version being added to the protocol engine: %d; Expecting %d", p.GetVersion(), pe.next))
	}
	pe.protocols[pe.next] = p
	pe.next++
	return p
}

func (pe *ProtocolEngine) GetByVersion(v uint64) (Protocol, bool) {
	p, flag := pe.protocols[v]
	return p, flag
}

func (pe *ProtocolEngine) GetKVStoreKeys() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		KeyMain,
		KeyAccount,
		KeyStake,
		KeyMint,
		KeyDistr,
		KeySlashing,
		KeyGov,
		KeyFee,
		KeyParams,
		KeyUpgrade,
		KeyService,
		KeyGuardian,
		KeyStorage,
		KeyCode}
}

func (pe *ProtocolEngine) GetTransientStoreKeys() []*sdk.TransientStoreKey {
	return []*sdk.TransientStoreKey{
		TkeyStake,
		TkeyDistr,
		TkeyParams}
}
