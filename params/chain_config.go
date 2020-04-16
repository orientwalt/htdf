package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

var (
	MainnetChainConfig = &params.ChainConfig{
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(1),
		DAOForkBlock:        big.NewInt(2),
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(3),
		EIP150Hash:          common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:         big.NewInt(4),
		EIP158Block:         big.NewInt(5),
		ByzantiumBlock:      big.NewInt(6),
		ConstantinopleBlock: big.NewInt(7),
	}
)
