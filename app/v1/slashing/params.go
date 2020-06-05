package slashing

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace                         = ModuleName
	DefaultMaxEvidenceAge       int64         = 60 * 2 * 60
	DefaultSignedBlocksWindow   int64         = 100
	DefaultDowntimeJailDuration time.Duration = 60 * 10 * time.Second
	BlocksPerMinute                           = 12                        // 5 seconds a block
	BlocksPerDay                              = BlocksPerMinute * 60 * 24 // 17280
)

// The Double Sign Jail period ends at Max Time supported by Amino (Dec 31, 9999 - 23:59:59 GMT)
var (
	DoubleSignJailEndTime          = time.Unix(253402300799, 0)
	DefaultMinSignedPerWindow      = sdk.NewDecWithPrec(5, 1)
	DefaultSlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
	DefaultSlashFractionDowntime   = sdk.NewDec(1).Quo(sdk.NewDec(100))
)

// Parameter store keys
var (
	KeyMaxEvidenceAge          = []byte("MaxEvidenceAge")
	KeySignedBlocksWindow      = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow      = []byte("MinSignedPerWindow")
	KeyDoubleSignJailDuration  = []byte("DoubleSignJailDuration")
	KeyCensorshipJailDuration  = []byte("CensorshipJailDuration")
	KeyDowntimeJailDuration    = []byte("DowntimeJailDuration")
	KeySlashFractionDoubleSign = []byte("SlashFractionDoubleSign")
	KeySlashFractionDowntime   = []byte("SlashFractionDowntime")
	KeySlashFractionCensorship = []byte("SlashFractionCensorship")
)

// ParamKeyTable for slashing module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for slashing at genesis
type Params struct {
	MaxEvidenceAge          int64         `json:"max_evidence_age"`
	SignedBlocksWindow      int64         `json:"signed_blocks_window"`
	MinSignedPerWindow      sdk.Dec       `json:"min_signed_per_window"`
	DoubleSignJailDuration  time.Duration `json:"double_sign_jail_duration"`
	CensorshipJailDuration  time.Duration `json:"censorship_jail_duration"`
	DowntimeJailDuration    time.Duration `json:"downtime_jail_duration"`
	SlashFractionDoubleSign sdk.Dec       `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   sdk.Dec       `json:"slash_fraction_downtime"`
	SlashFractionCensorship sdk.Dec       `json:"slash_fraction_censorship"`
}

func (p Params) String() string {
	return fmt.Sprintf(`Slashing Params:
  MaxEvidenceAge:          %d
  SignedBlocksWindow:      %d
  MinSignedPerWindow:      %s
  DowntimeJailDuration:    %s
  SlashFractionDoubleSign: %d
  SlashFractionDowntime:   %d`, p.MaxEvidenceAge,
		p.SignedBlocksWindow, p.MinSignedPerWindow,
		p.DowntimeJailDuration, p.SlashFractionDoubleSign,
		p.SlashFractionDowntime)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMaxEvidenceAge, &p.MaxEvidenceAge},
		{KeySignedBlocksWindow, &p.SignedBlocksWindow},
		{KeyMinSignedPerWindow, &p.MinSignedPerWindow},
		{KeyDowntimeJailDuration, &p.DowntimeJailDuration},
		{KeySlashFractionDoubleSign, &p.SlashFractionDoubleSign},
		{KeySlashFractionDowntime, &p.SlashFractionDowntime},
		{KeySlashFractionCensorship, &p.SlashFractionCensorship},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyMaxEvidenceAge):
		maxEvidenceAge, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMaxEvidenceAge(maxEvidenceAge); err != nil {
			return nil, err
		}
		return maxEvidenceAge, nil
	case string(KeySignedBlocksWindow):
		signedBlocksWindow, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSignedBlocksWindow(signedBlocksWindow); err != nil {
			return nil, err
		}
		return signedBlocksWindow, nil
	case string(KeyMinSignedPerWindow):
		minSignedPerWindow, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMinSignedPerWindow(minSignedPerWindow); err != nil {
			return nil, err
		}
		return minSignedPerWindow, nil
	case string(KeyDoubleSignJailDuration):
		doubleSignJailDuration, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateDoubleSignJailDuration(doubleSignJailDuration); err != nil {
			return nil, err
		}
		return doubleSignJailDuration, nil
	case string(KeyDowntimeJailDuration):
		downtimeJailDuration, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateDowntimeJailDuration(downtimeJailDuration); err != nil {
			return nil, err
		}
		return downtimeJailDuration, nil
	case string(KeyCensorshipJailDuration):
		censorshipJailDuration, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateCensorshipJailDuration(censorshipJailDuration); err != nil {
			return nil, err
		}
		return censorshipJailDuration, nil
	case string(KeySlashFractionDoubleSign):
		slashFractionDoubleSign, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSlashFractionDoubleSign(slashFractionDoubleSign); err != nil {
			return nil, err
		}
		return slashFractionDoubleSign, nil
	case string(KeySlashFractionDowntime):
		slashFractionDowntime, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSlashFractionDowntime(slashFractionDowntime); err != nil {
			return nil, err
		}
		return slashFractionDowntime, nil
	case string(KeySlashFractionCensorship):
		slashFractionCensorship, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSlashFractionCensorship(slashFractionCensorship); err != nil {
			return nil, err
		}
		return slashFractionCensorship, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func validateMaxEvidenceAge(p int64) sdk.Error {
	if p < 1*60 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash MaxEvidenceAge [%d] should be between [1 minute,) ", p))
	}
	return nil
}

func validateSignedBlocksWindow(p int64) sdk.Error {
	if p < 10 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash SignedBlocksWindow [%d] should be between [100, 140000] ", p))
	}
	return nil
}

func validateMinSignedPerWindow(p sdk.Dec) sdk.Error {
	if p.IsNegative() || p.GT(sdk.OneDec()) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Min signed per window should be less than or equal to one and greater than zero, is %s ", p.String()))
	}
	return nil
}

func validateDoubleSignJailDuration(p time.Duration) sdk.Error {
	if p <= 0 || p >= 2*sdk.Week {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash DoubleSignJailDuration [%s] should be between (0, 2weeks) ", p.String()))
	}
	return nil
}

func validateDowntimeJailDuration(p time.Duration) sdk.Error {
	if p < 1*time.Minute {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash DowntimeJailDuration [%s] should be between (0, 1minute) ", p.String()))
	}
	return nil
}

func validateCensorshipJailDuration(p time.Duration) sdk.Error {
	if p <= 0 || p >= 2*sdk.Week {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash CensorshipJailDuration [%s] should be between (0, 2weeks) ", p.String()))
	}
	return nil
}

func validateSlashFractionDoubleSign(p sdk.Dec) sdk.Error {
	if p.IsNegative() || p.GT(sdk.OneDec()) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, "Slashing fraction double sign should be less than or equal to one and greater than zero, is %s", p.String())
	}
	return nil
}

func validateSlashFractionDowntime(p sdk.Dec) sdk.Error {
	if p.IsNegative() || p.GT(sdk.OneDec()) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, "Slashing fraction downtime should be less than or equal to one and greater than zero, is %s", p.String())
	}
	return nil
}

func validateSlashFractionCensorship(p sdk.Dec) sdk.Error {
	if p.LT(sdk.ZeroDec()) || p.GT(sdk.NewDecWithPrec(1, 1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash SlashFractionCensorship [%s] should be between [0, 0.1] ", p.String()))
	}
	return nil
}

func (p *Params) GetParamSpace() string {
	return DefaultParamspace
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyMaxEvidenceAge):
		err := cdc.UnmarshalJSON(bytes, &p.MaxEvidenceAge)
		return strconv.FormatInt(p.MaxEvidenceAge, 10), err
	case string(KeySignedBlocksWindow):
		err := cdc.UnmarshalJSON(bytes, &p.SignedBlocksWindow)
		return strconv.FormatInt(p.SignedBlocksWindow, 10), err
	case string(KeyMinSignedPerWindow):
		err := cdc.UnmarshalJSON(bytes, &p.MinSignedPerWindow)
		return p.MinSignedPerWindow.String(), err
	case string(KeyDoubleSignJailDuration):
		err := cdc.UnmarshalJSON(bytes, &p.DoubleSignJailDuration)
		return p.DoubleSignJailDuration.String(), err
	case string(KeyDowntimeJailDuration):
		err := cdc.UnmarshalJSON(bytes, &p.DowntimeJailDuration)
		return p.DowntimeJailDuration.String(), err
	case string(KeyCensorshipJailDuration):
		err := cdc.UnmarshalJSON(bytes, &p.CensorshipJailDuration)
		return p.CensorshipJailDuration.String(), err
	case string(KeySlashFractionDoubleSign):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFractionDoubleSign)
		return p.SlashFractionDoubleSign.String(), err
	case string(KeySlashFractionDowntime):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFractionDowntime)
		return p.SlashFractionDowntime.String(), err
	case string(KeySlashFractionCensorship):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFractionCensorship)
		return p.SlashFractionCensorship.String(), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// Default parameters for this module
func DefaultParams() Params {
	return Params{
		MaxEvidenceAge:          DefaultMaxEvidenceAge,
		SignedBlocksWindow:      DefaultSignedBlocksWindow,
		MinSignedPerWindow:      DefaultMinSignedPerWindow,
		DowntimeJailDuration:    DefaultDowntimeJailDuration,
		SlashFractionDoubleSign: DefaultSlashFractionDoubleSign,
		SlashFractionDowntime:   DefaultSlashFractionDowntime,
	}
}

func validateParams(p Params) sdk.Error {
	if err := validateMaxEvidenceAge(p.MaxEvidenceAge); err != nil {
		return err
	}
	if err := validateSignedBlocksWindow(p.SignedBlocksWindow); err != nil {
		return err
	}
	if err := validateMinSignedPerWindow(p.MinSignedPerWindow); err != nil {
		return err
	}
	if err := validateDoubleSignJailDuration(p.DoubleSignJailDuration); err != nil {
		return err
	}
	if err := validateDowntimeJailDuration(p.DowntimeJailDuration); err != nil {
		return err
	}
	if err := validateCensorshipJailDuration(p.CensorshipJailDuration); err != nil {
		return err
	}
	if err := validateSlashFractionDoubleSign(p.SlashFractionDoubleSign); err != nil {
		return err
	}
	if err := validateSlashFractionDowntime(p.SlashFractionDowntime); err != nil {
		return err
	}
	if err := validateSlashFractionCensorship(p.SlashFractionCensorship); err != nil {
		return err
	}
	return nil
}

// MaxEvidenceAge - max age for evidence
func (k Keeper) MaxEvidenceAge(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyMaxEvidenceAge, &res)
	return
}

// SignedBlocksWindow - sliding window for downtime slashing
func (k Keeper) SignedBlocksWindow(ctx sdk.Context) (res int64) {
	k.paramspace.Get(ctx, KeySignedBlocksWindow, &res)
	return
}

// Downtime slashing threshold
func (k Keeper) MinSignedPerWindow(ctx sdk.Context) int64 {
	var minSignedPerWindow sdk.Dec
	k.paramspace.Get(ctx, KeyMinSignedPerWindow, &minSignedPerWindow)
	signedBlocksWindow := k.SignedBlocksWindow(ctx)

	// NOTE: RoundInt64 will never panic as minSignedPerWindow is
	//       less than 1.
	return minSignedPerWindow.MulInt64(signedBlocksWindow).RoundInt64()
}

// Downtime unbond duration
func (k Keeper) DowntimeJailDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyDowntimeJailDuration, &res)
	return
}

// SlashFractionDoubleSign
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, KeySlashFractionDoubleSign, &res)
	return
}

// Censorship jail duration
func (k Keeper) CensorshipJailDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyCensorshipJailDuration, &res)
	return
}

// SlashFractionDowntime
func (k Keeper) SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, KeySlashFractionDowntime, &res)
	return
}

// Slash fraction for Censorship
func (k Keeper) SlashFractionCensorship(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, KeySlashFractionCensorship, &res)
	return
}

// GetParams returns the total set of slashing parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}
