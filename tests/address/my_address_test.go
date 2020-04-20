package address

import (
	"fmt"
	sdk "github.com/orientwalt/htdf/types"

	"github.com/magiconair/properties/assert"
	"github.com/orientwalt/htdf/params"
	"testing"
)

func init() {
	// set address prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(params.Bech32PrefixAccAddr, params.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(params.Bech32PrefixValAddr, params.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(params.Bech32PrefixConsAddr, params.Bech32PrefixConsPub)
	config.Seal()
}

func TestBech32Address(t *testing.T) {

	//bech32 to binary
	bech32Contract := "htdf1nkkc48lfchy92ahg50akj2384v4yfqpm4hsq6y"
	binContractAddr, err := sdk.AccAddressFromBech32(bech32Contract)
	assert.Equal(t, err, nil)

	//binary to bech32
	assert.Equal(t, bech32Contract == binContractAddr.String(), true)
	fmt.Printf("binContractAddr=%x|binAddress.String=%s\n", binContractAddr, binContractAddr.String())

	bech32Minter := "htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3"
	accaddrTmp, err := sdk.AccAddressFromHex("85ced8ddf399d75c9e381e01f3bddcefb9132fe9")
	assert.Equal(t, err, nil)
	assert.Equal(t, bech32Minter == accaddrTmp.String(), true)

	fmt.Printf("accaddrTmp.String=%s\n", accaddrTmp.String())

	//bech32 to binary
	bech32BitKeep := "htdf17qarupfh9gee0yvywhxfy2zv39fjttracvgapx"
	binBitKeepAddr, err := sdk.AccAddressFromBech32(bech32BitKeep)
	assert.Equal(t, err, nil)
	fmt.Printf("binBitKeepAddr=%x|\n", binBitKeepAddr)

}
