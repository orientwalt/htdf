package rpc

import (
	"github.com/spf13/viper"

	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
)

// PublicNetAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicNetAPI struct {
	networkVersion string
}

// NewPublicNetAPI creates an instance of the public Net Web3 API.
func NewPublicNetAPI(_ context.CLIContext) *PublicNetAPI {
	chainID := viper.GetString(client.FlagChainID)
	// parse the chainID from a integer string
	// intChainID, err := strconv.ParseUint(chainID, 0, 64)
	// if err != nil {
	// 	panic(fmt.Sprintf("invalid chainID: %s, must be integer format", chainID))
	// }

	return &PublicNetAPI{
		networkVersion: chainID,
	}
}

// Version returns the current ethereum protocol version.
func (s *PublicNetAPI) Version() string {
	return s.networkVersion
}
