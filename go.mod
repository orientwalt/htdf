module github.com/orientwalt/htdf

go 1.12

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20191023202215-f096da5361bb // indirect
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/cosmos/ledger-cosmos-go v0.11.1
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/emicklei/proto v1.8.0
	github.com/ethereum/go-ethereum v1.8.27
	github.com/go-kit/kit v0.9.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-isatty v0.0.10
	github.com/pelletier/go-toml v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/rakyll/statik v0.1.6
	github.com/shopspring/decimal v0.0.0-20191009025716-f1972eb1d1f5
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.3.0
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.12.1
	github.com/tendermint/tendermint v0.31.5
	github.com/tendermint/tmlibs v0.9.0
	golang.org/x/crypto v0.0.0-20191112222119-e1110fd1c708
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)

replace (
	github.com/tendermint/iavl v0.12.1 => github.com/orientwalt/iavl v0.12.92
	github.com/tendermint/tendermint v0.31.5 => github.com/orientwalt/tendermint v0.90.10
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
)
