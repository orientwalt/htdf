package types

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	sdkerrors "github.com/orientwalt/htdf/types/errors"
	"github.com/orientwalt/htdf/x/auth"
	"github.com/tendermint/go-amino"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func Decode_Hex(str string) ([]byte, error) {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		//panic(fmt.Sprintf("invalid hex string: %q", str))
		return nil, err
	}
	return b, nil
}

//
func Encode_Hex(str []byte) string {
	return hex.EncodeToString(str)
}

// Read and decode a StdTx from rawdata
func ReadStdTxFromRawData(cdc *amino.Codec, str string) (stdTx auth.StdTx, err error) {
	bytes, err := Decode_Hex(str)
	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return stdTx, err
	}
	return stdTx, err
}

// Read and decode a StdTx from rawdata
func ReadStdTxFromString(cdc *amino.Codec, str string) (stdTx auth.StdTx, err error) {
	bytes := []byte(str)
	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return stdTx, err
	}
	return stdTx, err
}

// ValidateSigner attempts to validate a signer for a given slice of bytes over
// which a signature and signer is given. An error is returned if address
// derived from the signature and bytes signed does not match the given signer.
func ValidateSigner(signBytes, sig []byte, signer ethcmn.Address) error {
	pk, err := ethcrypto.SigToPub(signBytes, sig)

	if err != nil {
		return errors.Wrap(err, "failed to derive public key from signature")
	} else if ethcrypto.PubkeyToAddress(*pk) != signer {
		return fmt.Errorf("invalid signature for signer: %s", signer)
	}

	return nil
}

func rlpHash(x interface{}) (hash ethcmn.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	//nolint:gosec,errcheck
	rlp.Encode(hasher, x)
	hasher.Sum(hash[:0])

	return hash
}

// EncodeResultData takes all of the necessary data from the EVM execution
// and returns the data as a byte slice encoded with amino
func EncodeResultData(data ResultData) ([]byte, error) {
	return codec.New().MarshalBinaryLengthPrefixed(data)
}

// DecodeResultData decodes an amino-encoded byte slice into ResultData
func DecodeResultData(in []byte) (ResultData, error) {
	var data ResultData
	err := codec.New().UnmarshalBinaryLengthPrefixed(in, &data)
	if err != nil {
		return ResultData{}, err
	}
	return data, nil
}

// ----------------------------------------------------------------------------
// Auxiliary

// TxDecoder returns an sdk.TxDecoder that can decode both auth.StdTx and
// MsgEthereumTx transactions.
func TxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var tx sdk.Tx

		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		// sdk.Tx is an interface. The concrete message types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryBare(txBytes, &tx)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		return tx, nil
	}
}

// recoverEthSig recovers a signature according to the Ethereum specification and
// returns the sender or an error.
//
// Ref: Ethereum Yellow Paper (BYZANTIUM VERSION 69351d5) Appendix F
// nolint: gocritic
func recoverEthSig(R, S, Vb *big.Int, sigHash ethcmn.Hash) (ethcmn.Address, error) {
	if Vb.BitLen() > 8 {
		return ethcmn.Address{}, errors.New("invalid signature")
	}

	V := byte(Vb.Uint64() - 27)
	if !ethcrypto.ValidateSignatureValues(V, R, S, true) {
		return ethcmn.Address{}, errors.New("invalid signature")
	}

	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)

	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V

	// recover the public key from the signature
	pub, err := ethcrypto.Ecrecover(sigHash[:], sig)
	if err != nil {
		return ethcmn.Address{}, err
	}

	if len(pub) == 0 || pub[0] != 4 {
		return ethcmn.Address{}, errors.New("invalid public key")
	}

	var addr ethcmn.Address
	copy(addr[:], ethcrypto.Keccak256(pub[1:])[12:])

	return addr, nil
}
