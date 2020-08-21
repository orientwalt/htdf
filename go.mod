module github.com/orientwalt/htdf

go 1.12

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20191023202215-f096da5361bb
	github.com/astaxie/beego v1.12.0
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/cosmos/ethermint v0.0.0-20190802135314-3f32f9ba8a1f
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
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jinzhu/gorm v1.9.11
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-isatty v0.0.10
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pelletier/go-toml v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/rakyll/statik v0.1.6
	github.com/rjeczalik/notify v0.9.2
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
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
	gopkg.in/fatih/set.v0 v0.1.0 // indirect
	gopkg.in/karalabe/cookiejar.v2 v2.0.0-20150724131613-8dcd6a7f4951 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/resty.v1 v1.12.0
)

replace (
	github.com/tendermint/iavl v0.12.4 => github.com/orientwalt/iavl v0.12.4
	github.com/tendermint/tendermint v0.31.5 => github.com/orientwalt/tendermint v90.0.7+incompatible
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
)
