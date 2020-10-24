package types

import (
	"strings"

	codectypes "github.com/orientwalt/htdf/codec/types"
	clienttypes "github.com/orientwalt/htdf/x/ibc/core/02-client/types"
	commitmenttypes "github.com/orientwalt/htdf/x/ibc/core/23-commitment/types"
	host "github.com/orientwalt/htdf/x/ibc/core/24-host"
	"github.com/orientwalt/htdf/x/ibc/core/exported"
)

// NewQueryConnectionResponse creates a new QueryConnectionResponse instance
func NewQueryConnectionResponse(
	connectionID string, connection ConnectionEnd, proof []byte, height clienttypes.Height,
) *QueryConnectionResponse {
	path := commitmenttypes.NewMerklePath(strings.Split(host.ConnectionPath(connectionID), "/"))
	return &QueryConnectionResponse{
		Connection:  &connection,
		Proof:       proof,
		ProofPath:   path.Pretty(),
		ProofHeight: height,
	}
}

// NewQueryClientConnectionsResponse creates a new ConnectionPaths instance
func NewQueryClientConnectionsResponse(
	clientID string, connectionPaths []string, proof []byte, height clienttypes.Height,
) *QueryClientConnectionsResponse {
	path := commitmenttypes.NewMerklePath(strings.Split(host.ClientConnectionsPath(clientID), "/"))
	return &QueryClientConnectionsResponse{
		ConnectionPaths: connectionPaths,
		Proof:           proof,
		ProofPath:       path.Pretty(),
		ProofHeight:     height,
	}
}

// NewQueryClientConnectionsRequest creates a new QueryClientConnectionsRequest instance
func NewQueryClientConnectionsRequest(clientID string) *QueryClientConnectionsRequest {
	return &QueryClientConnectionsRequest{
		ClientId: clientID,
	}
}

// NewQueryConnectionClientStateResponse creates a newQueryConnectionClientStateResponse instance
func NewQueryConnectionClientStateResponse(identifiedClientState clienttypes.IdentifiedClientState, proof []byte, height clienttypes.Height) *QueryConnectionClientStateResponse {
	path := commitmenttypes.NewMerklePath(strings.Split(host.FullClientPath(identifiedClientState.ClientId, host.ClientStatePath()), "/"))
	return &QueryConnectionClientStateResponse{
		IdentifiedClientState: &identifiedClientState,
		Proof:                 proof,
		ProofPath:             path.Pretty(),
		ProofHeight:           height,
	}
}

// NewQueryConnectionConsensusStateResponse creates a newQueryConnectionConsensusStateResponse instance
func NewQueryConnectionConsensusStateResponse(clientID string, anyConsensusState *codectypes.Any, consensusStateHeight exported.Height, proof []byte, height clienttypes.Height) *QueryConnectionConsensusStateResponse {
	path := commitmenttypes.NewMerklePath(strings.Split(host.FullClientPath(clientID, host.ConsensusStatePath(consensusStateHeight)), "/"))
	return &QueryConnectionConsensusStateResponse{
		ConsensusState: anyConsensusState,
		ClientId:       clientID,
		Proof:          proof,
		ProofPath:      path.Pretty(),
		ProofHeight:    height,
	}
}
