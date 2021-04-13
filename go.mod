module github.com/orientwalt/htdf

go 1.14

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cespare/cp v1.1.1 // indirect
	github.com/confio/ics23-iavl v0.6.0
	github.com/confio/ics23-tendermint v0.6.1
	github.com/confio/ics23/go v0.6.3
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/ledger-cosmos-go v0.11.1
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1
	github.com/emicklei/proto v1.8.0
	github.com/ethereum/go-ethereum v1.10.2
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/uint256 v1.1.1
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/magiconair/properties v1.8.4
	github.com/mattn/go-isatty v0.0.12
	github.com/pelletier/go-toml v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/prometheus/common v0.15.0 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/shopspring/decimal v0.0.0-20191009025716-f1972eb1d1f5
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/iavl v0.14.0
	github.com/tendermint/tendermint v0.34.1
	github.com/tendermint/tm-db v0.5.1
	github.com/tendermint/tmlibs v0.9.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/ethereum/go-ethereum v1.10.2 => github.com/orientwalt/go-ethereum v1.10.2
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.1
	github.com/tendermint/iavl => github.com/tendermint/iavl v0.13.3
	github.com/tendermint/tendermint => github.com/orientwalt/tendermint v0.99.4
)
