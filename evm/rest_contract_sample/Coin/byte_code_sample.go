package main

import (
	"encoding/hex"
	"github.com/orientwalt/htdf/params"
	"github.com/orientwalt/htdf/types"
	"math/big"

	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"

	sdk "github.com/orientwalt/htdf/types"

	"io/ioutil"
	"os"
)

var (
	strTestContractToAddress = "htdf1vms0n5t80acapjnvr4t9xeelucujq58zml4kg2"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func loadBin(filename string) []byte {
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(errors.New("loadBin error"))
	}

	fmt.Printf("code=%x\n", code)

	return hexutil.MustDecode("0x" + string(code))
}
func loadAbi(filename string) abi.ABI {
	abiFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("open error|err=%s\n", err)
		panic(errors.New("loadBin error"))
	}
	defer abiFile.Close()
	abiObj, err := abi.JSON(abiFile)
	if err != nil {
		panic(errors.New("loadBin error"))
	}

	return abiObj
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usageg:  %s abiFileName  binFilename minterAddress \n", os.Args[0])
		fmt.Printf("    ##minterAddress  : \"nil\" means have no minterAddress\n")
		os.Exit(1)
	}

	// set address prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(params.Bech32PrefixAccAddr, params.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(params.Bech32PrefixValAddr, params.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(params.Bech32PrefixConsAddr, params.Bech32PrefixConsPub)
	config.Seal()

	abiFileName := os.Args[1]
	binFileName := os.Args[2]
	strMinterAddress := os.Args[3]

	//create contract
	data := loadBin(binFileName)
	fmt.Printf("contractCode, create contract|Code=%s\n", hex.EncodeToString(data))

	//minter
	abiObj := loadAbi(abiFileName)
	contractByteCode, err := abiObj.Pack("minter")
	must(err)
	fmt.Printf("contractCode, minter|Code=%s\n", hex.EncodeToString(contractByteCode))

	//==================access created contract=====================================
	if strMinterAddress == "nil" {
		fmt.Printf("have no strMinterAddress\n")
		os.Exit(0)
	}

	//address convert
	testContractToAddress, err := sdk.AccAddressFromBech32(strTestContractToAddress)
	must(err)
	eaTestContractToAddress := types.ToEthAddress(testContractToAddress)

	minterAddress, err := sdk.AccAddressFromBech32(strMinterAddress)
	must(err)
	eaMinterAddress := types.ToEthAddress(minterAddress)

	//mint
	contractByteCode, err = abiObj.Pack("mint", eaMinterAddress, big.NewInt(1000000))
	must(err)
	fmt.Printf("contractCode, mint|minterAddress=%s|Code=%s\n", minterAddress.String(), hex.EncodeToString(contractByteCode))

	//send
	contractByteCode, err = abiObj.Pack("send", eaTestContractToAddress, big.NewInt(30))
	must(err)
	fmt.Printf("contractCode, send|testContractToAddress=%s|Code=%s\n", testContractToAddress.String(), hex.EncodeToString(contractByteCode))

	//get balance
	contractByteCode, err = abiObj.Pack("balances", eaTestContractToAddress)
	must(err)
	fmt.Printf("contractCode, get balance|testContractToAddress=%s|Code=%s\n", testContractToAddress.String(), hex.EncodeToString(contractByteCode))

	//get minter balance
	contractByteCode, err = abiObj.Pack("balances", eaMinterAddress)
	must(err)
	fmt.Printf("contractCode, get balance|strMinterAddress=%s|Code=%s\n", strMinterAddress, hex.EncodeToString(contractByteCode))

}
