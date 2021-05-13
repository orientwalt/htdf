package types

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/params"
)

const (
	// Default parameter namespace
	DefaultParamSpace = sdk.DefaultDenom // junying-todo, 2019-08-27, "stake" to satoshi
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.

	// 2021-04-08 yqq, modify 14 days to 3 days
	DefaultUnbondingTime time.Duration = time.Second * 60 * 60 * 24 * 3

	// Default maximum number of bonded validators
	DefaultMaxValidators uint16 = 50

	// Default maximum entries in a UBD/RED pair
	DefaultMaxEntries uint16 = 7

	// Delay, in blocks, between when validator updates are returned to Tendermint and when they are applied
	// For example, if this is 0, the validator set at the end of a block will sign the next block, or
	// if this is 1, the validator set at the end of a block will sign the block after the next.
	// Constant as this should not change without a hard fork.
	ValidatorUpdateDelay int64 = 1
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime = []byte("UnbondingTime")
	KeyMaxValidators = []byte("MaxValidators")
	KeyMaxEntries    = []byte("KeyMaxEntries")
	KeyBondDenom     = []byte("BondDenom")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	UnbondingTime time.Duration `json:"unbonding_time"` // time duration of unbonding
	MaxValidators uint16        `json:"max_validators"` // maximum number of validators (max uint16 = 65535)
	MaxEntries    uint16        `json:"max_entries"`    // max entries for either unbonding delegation or redelegation (per pair/trio)
	// note: we need to be a bit careful about potential overflow here, since this is user-determined
	BondDenom string `json:"bond_denom"` // bondable coin denomination
}

func NewParams(unbondingTime time.Duration, maxValidators, maxEntries uint16,
	bondDenom string) Params {

	return Params{
		UnbondingTime: unbondingTime,
		MaxValidators: maxValidators,
		MaxEntries:    maxEntries,
		BondDenom:     bondDenom,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyUnbondingTime, Value: &p.UnbondingTime},
		{Key: KeyMaxValidators, Value: &p.MaxValidators},
		{Key: KeyMaxEntries, Value: &p.MaxEntries},
		{Key: KeyBondDenom, Value: &p.BondDenom},
	}
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := MsgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := MsgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultUnbondingTime, DefaultMaxValidators, DefaultMaxEntries, sdk.DefaultBondDenom)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time:    %s
  Max Validators:    %d
  Max Entries:       %d
  Bonded Coin Denom: %s`, p.UnbondingTime,
		p.MaxValidators, p.MaxEntries, p.BondDenom)
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// // validate a set of params
// func (p Params) Validate() error {
// 	if p.BondDenom == "" {
// 		return fmt.Errorf("staking parameter BondDenom can't be an empty string")
// 	}
// 	if p.MaxValidators == 0 {
// 		return fmt.Errorf("staking parameter MaxValidators must be a positive integer")
// 	}
// 	return nil
// }

//
func ValidateParams(p Params) error {
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}

	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}
	return nil
}

//
func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyUnbondingTime):
		unbondingTime, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateUnbondingTime(unbondingTime); err != nil {
			return nil, err
		}
		return unbondingTime, nil
	case string(KeyMaxValidators):
		maxValidators, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMaxValidators(uint16(maxValidators)); err != nil {
			return nil, err
		}
		return uint16(maxValidators), nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyUnbondingTime):
		err := cdc.UnmarshalJSON(bytes, &p.UnbondingTime)
		return p.UnbondingTime.String(), err
	case string(KeyMaxValidators):
		err := cdc.UnmarshalJSON(bytes, &p.MaxValidators)
		return strconv.Itoa(int(p.MaxValidators)), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

func validateUnbondingTime(v time.Duration) sdk.Error {
	// planed to be modified by junying, 2019-11-20
	// because two-week unbond time may be too long
	// if v < 2*sdk.Week

	// 2021-05-12, yqq
	// We make a flag for test environment.
	// In test environment, we don't check unboding time.
	flagTest := false
	_, ok := os.LookupEnv("HTDF_TEST_ENV")
	if !ok {
		flagTest = true
	}

	if !flagTest && v < 3*sdk.Day {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidUnbondingTime,
			fmt.Sprintf("Invalid UnbondingTime [%s] should be greater than or equal to 3 days", v.String()))
	}
	return nil
}

func validateMaxValidators(v uint16) sdk.Error {
	if v < 50 || v > 200 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxValidators, fmt.Sprintf("Invalid MaxValidators [%d] should be between [50, 200]", v))
	}
	return nil
}
