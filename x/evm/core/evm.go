// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	// "github.com/ethereum/go-ethereum/core/vm"

	evm "github.com/orientwalt/htdf/x/evm/core/vm"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	//Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}

type IMessage interface {
	FromAddress() common.Address
}

// FakeChainContext impl interface ChainContext
type FakeChainContext struct {
}

func (self FakeChainContext) GetHeader(hash common.Hash, number uint64) *ethtypes.Header {
	return nil
}

// NewEVMBlockContext creates a new context for use in the EVM.
func NewEVMBlockContext(header abci.Header, chainCtx ChainContext, author *common.Address, height uint64) evm.BlockContext {

	// If we don't have an explicit author (i.e. not mining), extract from the header
	// var beneficiary common.Address
	// if author == nil {
	// 	beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	// } else {
	// 	beneficiary = *author
	// }

	beneficiary := *author

	curBlockHeader := &types.Header{
		ParentHash: common.BytesToHash(header.LastBlockId.Hash),
		Number:     big.NewInt(int64(header.Height)),
		Difficulty: big.NewInt(1),
		GasLimit:   0,
		GasUsed:    0,
		Time:       uint64(header.Time.Unix()), // time should be deterministic
		Extra:      nil,
	}

	return evm.BlockContext{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(curBlockHeader, chainCtx),         // for evm opCode BLOCKHASH
		Coinbase:    beneficiary,                                 // for evm opCode COINBASE
		BlockNumber: new(big.Int).Set(curBlockHeader.Number),     // for evm opCode BLOCKNUMBER
		Time:        new(big.Int).SetUint64(curBlockHeader.Time), // for evm opCode BLOCKTIME
		Difficulty:  new(big.Int).Set(curBlockHeader.Difficulty), // for evm opCode DIFFICULTY
		GasLimit:    curBlockHeader.GasLimit,                     // for evm opCode GASLIMIT
	}
}

// NewEVMTxContext creates a new transaction context for a single transaction.
func NewEVMTxContext(msg IMessage) evm.TxContext {
	return evm.TxContext{
		Origin: msg.FromAddress(),
		// GasPrice: new(big.Int).Set(msg.GetGasPrice()),
		GasPrice: new(big.Int).SetInt64(0),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	// Cache will initially contain [refHash.parent],
	// Then fill up with [refHash.p, refHash.pp, refHash.ppp, ...]
	var cache []common.Hash

	return func(n uint64) common.Hash {
		// If there's no hash cache yet, make one
		if len(cache) == 0 {
			cache = append(cache, ref.ParentHash)
		}
		if idx := ref.Number.Uint64() - n - 1; idx < uint64(len(cache)) {
			return cache[idx]
		}
		// No luck in the cache, but we can start iterating from the last element we already know
		lastKnownHash := cache[len(cache)-1]
		lastKnownNumber := ref.Number.Uint64() - uint64(len(cache))

		for {
			header := chain.GetHeader(lastKnownHash, lastKnownNumber)
			if header == nil {
				break
			}
			cache = append(cache, header.ParentHash)
			lastKnownHash = header.ParentHash
			lastKnownNumber = header.Number.Uint64() - 1
			if n == lastKnownNumber {
				return lastKnownHash
			}
		}
		return common.Hash{}
	}
}

// CanTransfer checks wether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db evm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db evm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}
