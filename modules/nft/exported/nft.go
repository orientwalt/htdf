package exported

import (
	sdk "github.com/orientwalt/htdf/types"
)

// NFT non fungible token interface
type NFT interface {
	GetID() string
	GetName() string
	GetOwner() sdk.AccAddress
	GetURI() string
	GetData() string
}
