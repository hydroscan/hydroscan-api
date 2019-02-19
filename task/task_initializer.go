package task

import (
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

const HydroExchangeABI = `[{"constant":false,"inputs":[{"name":"delegate","type":"address"}],"name":"approveDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newConfig","type":"bytes32"}],"name":"changeDiscountConfig","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"proxyAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"filled","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"cancelled","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"relayerDelegates","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"exitIncentiveSystem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"DOMAIN_SEPARATOR","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"discountConfig","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"relayer","type":"address"}],"name":"canMatchOrdersFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"}],"name":"getDiscountedRate","outputs":[{"name":"result","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"EIP712_ORDER_TYPE","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"FEE_RATE_BASE","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"renounceOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"DISCOUNT_RATE_BASE","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"trader","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"takerOrderParam","type":"tuple"},{"components":[{"name":"trader","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"makerOrderParams","type":"tuple[]"},{"components":[{"name":"baseToken","type":"address"},{"name":"quoteToken","type":"address"},{"name":"relayer","type":"address"}],"name":"orderAddressSet","type":"tuple"}],"name":"matchOrders","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"isOwner","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"relayer","type":"address"}],"name":"isParticipant","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"trader","type":"address"},{"name":"relayer","type":"address"},{"name":"baseToken","type":"address"},{"name":"quoteToken","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"}],"name":"order","type":"tuple"}],"name":"cancelOrder","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"joinIncentiveSystem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"EIP712_DOMAIN_TYPEHASH","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"hotTokenAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"delegate","type":"address"}],"name":"revokeDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"_proxyAddress","type":"address"},{"name":"hotTokenAddress","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"orderHash","type":"bytes32"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"baseToken","type":"address"},{"indexed":false,"name":"quoteToken","type":"address"},{"indexed":false,"name":"relayer","type":"address"},{"indexed":false,"name":"maker","type":"address"},{"indexed":false,"name":"taker","type":"address"},{"indexed":false,"name":"baseTokenAmount","type":"uint256"},{"indexed":false,"name":"quoteTokenAmount","type":"uint256"},{"indexed":false,"name":"makerFee","type":"uint256"},{"indexed":false,"name":"takerFee","type":"uint256"},{"indexed":false,"name":"makerGasFee","type":"uint256"},{"indexed":false,"name":"makerRebate","type":"uint256"},{"indexed":false,"name":"takerGasFee","type":"uint256"}],"name":"Match","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"},{"indexed":true,"name":"delegate","type":"address"}],"name":"RelayerApproveDelegate","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"},{"indexed":true,"name":"delegate","type":"address"}],"name":"RelayerRevokeDelegate","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"}],"name":"RelayerExit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"}],"name":"RelayerJoin","type":"event"}]`

var EthClient *ethclient.Client
var contractABI abi.ABI
var err error

func InitEthClient() {
	EthClient, err = ethclient.Dial(os.Getenv("WEB3_URL"))
	if err != nil {
		log.Panic(err)
	}

	contractABI, err = abi.JSON(strings.NewReader(string(HydroExchangeABI)))
	if err != nil {
		panic(err)
	}
}
