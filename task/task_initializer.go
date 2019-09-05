package task

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

var ResourcePath = "/resource"
var EthClient *ethclient.Client
var err error
var contractABIV1 abi.ABI
var contractABIV1_1 abi.ABI

// var contractABIV2 abi.ABI

const MaxReties = 5

func InitEthClient() {
	if viper.GetString("resource_path") != "" {
		ResourcePath = viper.GetString("resource_path")
	}

	EthClient, err = ethclient.Dial(viper.GetString("web3_url"))
	if err != nil {
		log.Panic(err)
	}

	contractABIV1, err = abi.JSON(strings.NewReader(string(HydroExchangeABIV1)))
	if err != nil {
		panic(err)
	}

	contractABIV1_1, err = abi.JSON(strings.NewReader(string(HydroExchangeABIV1_1)))
	if err != nil {
		panic(err)
	}

	// contractABIV2, err = abi.JSON(strings.NewReader(string(HydroExchangeABIV2)))
	// if err != nil {
	// 	panic(err)
	// }
}
