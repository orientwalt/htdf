package types

import (
	"fmt"
	"time"
)

// Signing info for a validator
type ValidatorSigningInfo struct {
	StartHeight         int64     `json:"start_height"`          // height at which validator was first a candidate OR was unjailed
	IndexOffset         int64     `json:"index_offset"`          // index offset into signed block bit array
	JailedUntil         time.Time `json:"jailed_until"`          // timestamp validator cannot be unjailed until
	Tombstoned          bool      `json:"tombstoned"`            // whether or not a validator has been tombstoned (killed out of validator set)
	MissedBlocksCounter int64     `json:"missed_blocks_counter"` // missed blocks counter (to avoid scanning the array every time)
}

// Construct a new `ValidatorSigningInfo` struct
func NewValidatorSigningInfo(startHeight, indexOffset int64, jailedUntil time.Time,
	tombstoned bool, missedBlocksCounter int64) ValidatorSigningInfo {

	return ValidatorSigningInfo{
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
	}
}

// Return human readable signing info
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Start Height:          %d
Index Offset:          %d
Jailed Until:          %v
Tombstoned:            %t
Missed Blocks Counter: %d`,
		i.StartHeight, i.IndexOffset, i.JailedUntil,
		i.Tombstoned, i.MissedBlocksCounter)
}
