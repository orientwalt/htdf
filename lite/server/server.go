package server

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite/proxy"

	// rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcclient "github.com/tendermint/tendermint/lite/client"
	"github.com/tendermint/tmlibs/cli"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

type Server struct {
	listenAddr         string
	nodeAddr           string
	chainID            string
	home               string
	maxOpenConnections int
	cacheSize          int
}

func getFilePath() string {
	rootDir := viper.GetString(cli.HomeFlag)
	defaultKeyStoreHome := filepath.Join(rootDir, "litenode-data")
	return defaultKeyStoreHome
}

// func init() {
// 	srv := &Server{
// 		listenAddr:         "tcp://localhost:8888",
// 		nodeAddr:           "tcp://localhost:26657",
// 		chainID:            "testchain",
// 		home:               getFilePath(),
// 		maxOpenConnections: 900,
// 		cacheSize:          10,
// 	}

// }

func NewServer(lAddr string,
	nAddr string,
	cID string,
	dir string,
	maxConnections int,
	cSize int) *Server {
	return &Server{
		listenAddr:         lAddr,
		nodeAddr:           nAddr,
		chainID:            cID,
		home:               dir,
		maxOpenConnections: maxConnections,
		cacheSize:          cSize,
	}
}

func (s *Server) ensureAddrHasSchemeOrDefaultToTCP(addr string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "tcp", "unix":
	case "":
		u.Scheme = "tcp"
	default:
		return "", fmt.Errorf("unknown scheme %q, use either tcp or unix", u.Scheme)
	}
	return u.String(), nil
}

func (s *Server) RunProxy() error {
	nodeAddr, err := s.ensureAddrHasSchemeOrDefaultToTCP(s.nodeAddr)
	if err != nil {
		return err
	}
	listenAddr, err := s.ensureAddrHasSchemeOrDefaultToTCP(s.listenAddr)
	if err != nil {
		return err
	}

	// First, connect a client
	logger.Info("Connecting to source HTTP client...")
	node := rpcclient.NewHTTP(nodeAddr, "/websocket")

	logger.Info("Constructing Verifier...")
	cert, err := proxy.NewVerifier(s.chainID, s.home, node, logger, s.cacheSize)
	if err != nil {
		return errors.Wrap(err, "constructing Verifier")
	}
	cert.SetLogger(logger)
	sc := proxy.SecureClient(node, cert)

	logger.Info("Starting proxy...")
	err = proxy.StartProxy(sc, listenAddr, logger, s.maxOpenConnections)
	if err != nil {
		return errors.Wrap(err, "starting proxy")
	}

	return nil
}
