// Package rpc contains RPC handler methods and utilities to start
// Ethermint's Web3-compatibly JSON-RPC server.
package rpc

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/orientwalt/htdf/client/context"
)

// RPC namespaces and API version
const (
	Web3Namespace     = "web3"
	EthNamespace      = "eth"
	PersonalNamespace = "personal"
	NetNamespace      = "net"

	apiVersion = "1.0"
)

// GetRPCAPIs returns the list of all APIs
func GetRPCAPIs(cliCtx context.CLIContext /*, keys []secp256k1.PrivKeySecp256k1*/) []rpc.API {
	nonceLock := new(AddrLocker)
	backend := NewEthermintBackend(cliCtx)

	return []rpc.API{
		{
			Namespace: Web3Namespace,
			Version:   apiVersion,
			Service:   NewPublicWeb3API(),
			Public:    true,
		},
		// {
		// 	Namespace: EthNamespace,
		// 	Version:   apiVersion,
		// 	Service:   NewPublicEthAPI(cliCtx, backend, nonceLock, keys),
		// 	Public:    true,
		// },
		{
			Namespace: PersonalNamespace,
			Version:   apiVersion,
			Service:   NewPersonalEthAPI(cliCtx, nonceLock),
			Public:    false,
		},
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   NewPublicFilterAPI(cliCtx, backend),
			Public:    true,
		},
		{
			Namespace: NetNamespace,
			Version:   apiVersion,
			Service:   NewPublicNetAPI(cliCtx),
			Public:    true,
		},
	}
}
