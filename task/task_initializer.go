package task

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

const ProtocolV1 = "1"
const ProtocolV1_1 = "1.1"
const ProtocolV1_2 = "1.2"

const HydroExchangeAddressV1 = "0x2cB4B49C0d6E9db2164d94Ce48853BF77C4D883E"
const HydroExchangeABIV1 = `[{"constant":false,"inputs":[{"name":"delegate","type":"address"}],"name":"approveDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newConfig","type":"bytes32"}],"name":"changeDiscountConfig","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"proxyAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"filled","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"cancelled","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"relayerDelegates","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"exitIncentiveSystem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"DOMAIN_SEPARATOR","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"discountConfig","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"relayer","type":"address"}],"name":"canMatchOrdersFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"}],"name":"getDiscountedRate","outputs":[{"name":"result","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"EIP712_ORDER_TYPE","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"FEE_RATE_BASE","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"renounceOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"DISCOUNT_RATE_BASE","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"trader","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"takerOrderParam","type":"tuple"},{"components":[{"name":"trader","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"makerOrderParams","type":"tuple[]"},{"components":[{"name":"baseToken","type":"address"},{"name":"quoteToken","type":"address"},{"name":"relayer","type":"address"}],"name":"orderAddressSet","type":"tuple"}],"name":"matchOrders","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"isOwner","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"relayer","type":"address"}],"name":"isParticipant","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"trader","type":"address"},{"name":"relayer","type":"address"},{"name":"baseToken","type":"address"},{"name":"quoteToken","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"}],"name":"order","type":"tuple"}],"name":"cancelOrder","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"joinIncentiveSystem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"EIP712_DOMAIN_TYPEHASH","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"hotTokenAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"delegate","type":"address"}],"name":"revokeDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"_proxyAddress","type":"address"},{"name":"hotTokenAddress","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"orderHash","type":"bytes32"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"baseToken","type":"address"},{"indexed":false,"name":"quoteToken","type":"address"},{"indexed":false,"name":"relayer","type":"address"},{"indexed":false,"name":"maker","type":"address"},{"indexed":false,"name":"taker","type":"address"},{"indexed":false,"name":"baseTokenAmount","type":"uint256"},{"indexed":false,"name":"quoteTokenAmount","type":"uint256"},{"indexed":false,"name":"makerFee","type":"uint256"},{"indexed":false,"name":"takerFee","type":"uint256"},{"indexed":false,"name":"makerGasFee","type":"uint256"},{"indexed":false,"name":"makerRebate","type":"uint256"},{"indexed":false,"name":"takerGasFee","type":"uint256"}],"name":"Match","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"},{"indexed":true,"name":"delegate","type":"address"}],"name":"RelayerApproveDelegate","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"},{"indexed":true,"name":"delegate","type":"address"}],"name":"RelayerRevokeDelegate","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"}],"name":"RelayerExit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"}],"name":"RelayerJoin","type":"event"}]`
const HydroStartBlockNumberV1 = 6885289

const HydroExchangeAddressV1_1 = "0xE2a0BFe759e2A4444442Da5064ec549616FFF101"
const HydroExchangeABIV1_1 = `[{"constant":false,"inputs":[{"name":"delegate","type":"address"}],"name":"approveDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newConfig","type":"bytes32"}],"name":"changeDiscountConfig","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"proxyAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"filled","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"cancelled","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"relayerDelegates","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"exitIncentiveSystem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"DOMAIN_SEPARATOR","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"discountConfig","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"relayer","type":"address"}],"name":"canMatchOrdersFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"}],"name":"getDiscountedRate","outputs":[{"name":"result","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"EIP712_ORDER_TYPE","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"FEE_RATE_BASE","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"REBATE_RATE_BASE","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"renounceOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"DISCOUNT_RATE_BASE","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"trader","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"takerOrderParam","type":"tuple"},{"components":[{"name":"trader","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"makerOrderParams","type":"tuple[]"},{"name":"baseTokenFilledAmounts","type":"uint256[]"},{"components":[{"name":"baseToken","type":"address"},{"name":"quoteToken","type":"address"},{"name":"relayer","type":"address"}],"name":"orderAddressSet","type":"tuple"}],"name":"matchOrders","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"isOwner","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"relayer","type":"address"}],"name":"isParticipant","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"trader","type":"address"},{"name":"relayer","type":"address"},{"name":"baseToken","type":"address"},{"name":"quoteToken","type":"address"},{"name":"baseTokenAmount","type":"uint256"},{"name":"quoteTokenAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"}],"name":"order","type":"tuple"}],"name":"cancelOrder","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"joinIncentiveSystem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"SUPPORTED_ORDER_VERSION","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"EIP712_DOMAIN_TYPEHASH","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"hotTokenAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"delegate","type":"address"}],"name":"revokeDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"_proxyAddress","type":"address"},{"name":"hotTokenAddress","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"orderHash","type":"bytes32"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"components":[{"name":"baseToken","type":"address"},{"name":"quoteToken","type":"address"},{"name":"relayer","type":"address"}],"indexed":false,"name":"addressSet","type":"tuple"},{"components":[{"name":"maker","type":"address"},{"name":"taker","type":"address"},{"name":"buyer","type":"address"},{"name":"makerFee","type":"uint256"},{"name":"makerRebate","type":"uint256"},{"name":"takerFee","type":"uint256"},{"name":"makerGasFee","type":"uint256"},{"name":"takerGasFee","type":"uint256"},{"name":"baseTokenFilledAmount","type":"uint256"},{"name":"quoteTokenFilledAmount","type":"uint256"}],"indexed":false,"name":"result","type":"tuple"}],"name":"Match","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"},{"indexed":true,"name":"delegate","type":"address"}],"name":"RelayerApproveDelegate","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"},{"indexed":true,"name":"delegate","type":"address"}],"name":"RelayerRevokeDelegate","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"}],"name":"RelayerExit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"relayer","type":"address"}],"name":"RelayerJoin","type":"event"}]`
const HydroStartBlockNumberV1_1 = 7454912

const HydroExchangeAddressV1_2 = "0x241e82C79452F51fbfc89Fac6d912e021dB1a3B7"
const HydroExchangeABIV1_2 = `[{"payable":true,"stateMutability":"payable","type":"fallback"},{"constant":false,"inputs":[{"components":[{"name":"actionType","type":"uint8"},{"name":"encodedParams","type":"bytes"}],"name":"actions","type":"tuple[]"}],"name":"batch","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"hash","type":"bytes32"},{"name":"signerAddress","type":"address"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"isValidSignature","outputs":[{"name":"isValid","type":"bool"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"getAllMarketsCount","outputs":[{"name":"count","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"assetAddress","type":"address"}],"name":"getAsset","outputs":[{"components":[{"name":"lendingPoolToken","type":"address"},{"name":"priceOracle","type":"address"},{"name":"interestModel","type":"address"}],"name":"asset","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"assetAddress","type":"address"}],"name":"getAssetOraclePrice","outputs":[{"name":"price","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"marketID","type":"uint16"}],"name":"getMarket","outputs":[{"components":[{"name":"baseAsset","type":"address"},{"name":"quoteAsset","type":"address"},{"name":"liquidateRate","type":"uint256"},{"name":"withdrawRate","type":"uint256"},{"name":"auctionRatioStart","type":"uint256"},{"name":"auctionRatioPerBlock","type":"uint256"}],"name":"market","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"},{"name":"marketID","type":"uint16"}],"name":"isAccountLiquidatable","outputs":[{"name":"isLiquidatable","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"},{"name":"marketID","type":"uint16"}],"name":"getAccountDetails","outputs":[{"components":[{"name":"liquidatable","type":"bool"},{"name":"status","type":"uint8"},{"name":"debtsTotalUSDValue","type":"uint256"},{"name":"balancesTotalUSDValue","type":"uint256"}],"name":"details","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAuctionsCount","outputs":[{"name":"count","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getCurrentAuctions","outputs":[{"name":"","type":"uint32[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"auctionID","type":"uint32"}],"name":"getAuctionDetails","outputs":[{"components":[{"name":"borrower","type":"address"},{"name":"marketID","type":"uint16"},{"name":"debtAsset","type":"address"},{"name":"collateralAsset","type":"address"},{"name":"leftDebtAmount","type":"uint256"},{"name":"leftCollateralAmount","type":"uint256"},{"name":"ratio","type":"uint256"},{"name":"price","type":"uint256"},{"name":"finished","type":"bool"}],"name":"details","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"auctionID","type":"uint32"},{"name":"amount","type":"uint256"}],"name":"fillAuctionWithAmount","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"user","type":"address"},{"name":"marketID","type":"uint16"}],"name":"liquidateAccount","outputs":[{"name":"isLiquidatable","type":"bool"},{"name":"auctionID","type":"uint32"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"asset","type":"address"}],"name":"getTotalBorrow","outputs":[{"name":"amount","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"asset","type":"address"}],"name":"getTotalSupply","outputs":[{"name":"amount","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"asset","type":"address"},{"name":"user","type":"address"},{"name":"marketID","type":"uint16"}],"name":"getAmountBorrowed","outputs":[{"name":"amount","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"asset","type":"address"},{"name":"user","type":"address"}],"name":"getAmountSupplied","outputs":[{"name":"amount","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"asset","type":"address"},{"name":"extraBorrowAmount","type":"uint256"}],"name":"getInterestRates","outputs":[{"name":"borrowInterestRate","type":"uint256"},{"name":"supplyInterestRate","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"asset","type":"address"}],"name":"getInsuranceBalance","outputs":[{"name":"amount","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"delegate","type":"address"}],"name":"approveDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"delegate","type":"address"}],"name":"revokeDelegate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"joinIncentiveSystem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"exitIncentiveSystem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"relayer","type":"address"}],"name":"canMatchOrdersFrom","outputs":[{"name":"canMatch","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"relayer","type":"address"}],"name":"isParticipant","outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"asset","type":"address"},{"name":"user","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"marketID","type":"uint16"},{"name":"asset","type":"address"},{"name":"user","type":"address"}],"name":"marketBalanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"marketID","type":"uint16"},{"name":"asset","type":"address"},{"name":"user","type":"address"}],"name":"getMarketTransferableAmount","outputs":[{"name":"amount","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"trader","type":"address"},{"name":"relayer","type":"address"},{"name":"baseAsset","type":"address"},{"name":"quoteAsset","type":"address"},{"name":"baseAssetAmount","type":"uint256"},{"name":"quoteAssetAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"}],"name":"order","type":"tuple"}],"name":"cancelOrder","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"orderHash","type":"bytes32"}],"name":"isOrderCancelled","outputs":[{"name":"isCancelled","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"components":[{"name":"trader","type":"address"},{"name":"baseAssetAmount","type":"uint256"},{"name":"quoteAssetAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"takerOrderParam","type":"tuple"},{"components":[{"name":"trader","type":"address"},{"name":"baseAssetAmount","type":"uint256"},{"name":"quoteAssetAmount","type":"uint256"},{"name":"gasTokenAmount","type":"uint256"},{"name":"data","type":"bytes32"},{"components":[{"name":"config","type":"bytes32"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"signature","type":"tuple"}],"name":"makerOrderParams","type":"tuple[]"},{"name":"baseAssetFilledAmounts","type":"uint256[]"},{"components":[{"name":"baseAsset","type":"address"},{"name":"quoteAsset","type":"address"},{"name":"relayer","type":"address"}],"name":"orderAddressSet","type":"tuple"}],"name":"params","type":"tuple"}],"name":"matchOrders","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"}],"name":"getDiscountedRate","outputs":[{"name":"rate","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getHydroTokenAddress","outputs":[{"name":"hydroTokenAddress","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"orderHash","type":"bytes32"}],"name":"getOrderFilledAmount","outputs":[{"name":"amount","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"components":[{"name":"baseToken","type":"address"},{"name":"quoteToken","type":"address"},{"name":"relayer","type":"address"}],"indexed":false,"name":"addressSet","type":"tuple"},{"components":[{"name":"maker","type":"address"},{"name":"taker","type":"address"},{"name":"buyer","type":"address"},{"name":"makerFee","type":"uint256"},{"name":"makerRebate","type":"uint256"},{"name":"takerFee","type":"uint256"},{"name":"makerGasFee","type":"uint256"},{"name":"takerGasFee","type":"uint256"},{"name":"baseTokenFilledAmount","type":"uint256"},{"name":"quoteTokenFilledAmount","type":"uint256"}],"indexed":false,"name":"result","type":"tuple"}],"name":"Match","type":"event"}]`
const HydroStartBlockNumberV1_2 = 8399662

var ResourcePath = "/resource"
var EthClient *ethclient.Client
var err error
var contractABIV1 abi.ABI
var contractABIV1_1 abi.ABI
var contractABIV1_2 abi.ABI

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

	contractABIV1_2, err = abi.JSON(strings.NewReader(string(HydroExchangeABIV1_2)))
	if err != nil {
		panic(err)
	}
}
