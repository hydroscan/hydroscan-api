package task

import (
	"context"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type MatchEvent struct {
	BaseToken        common.Address
	QuoteToken       common.Address
	Relayer          common.Address
	Maker            common.Address
	Taker            common.Address
	BaseTokenAmount  *big.Int
	QuoteTokenAmount *big.Int
	MakerFee         *big.Int
	TakerFee         *big.Int
	MakerGasFee      *big.Int
	MakerRebate      *big.Int
	TakerGasFee      *big.Int
}

func FetchHistoricalLogs() {
	log.Info("FetchHistoricalLogs")
	fromBlock := getFromBlockNumber()
	lastBlock := getLastBlockNumber()

	if fromBlock > lastBlock {
		return
	}

	pageSize := uint64(1000)
	for fromBlock+pageSize < lastBlock {
		toBlock := fromBlock + pageSize
		fetchLogs(int64(fromBlock), int64(toBlock))
		fromBlock = toBlock
		lastBlock = getLastBlockNumber()
	}

	if fromBlock+pageSize >= lastBlock {
		fetchLogs(int64(fromBlock), int64(lastBlock))
	}

	UpdateHistoryTradePrice()
}

func fetchLogs(fromBlock int64, toBlock int64) {
	log.Info("fetchLogs: ", fromBlock, " - ", toBlock)

	contractAddress := common.HexToAddress(os.Getenv("HYDRO_EXCHANGE_ADDRESS"))
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	eventLogs, err := EthClient.FilterLogs(context.Background(), query)
	if err != nil {
		panic(err)
	}

	for _, eventLog := range eventLogs {
		saveEventLog(eventLog)
	}
}

func saveEventLog(eventLog types.Log) {
	log.Info("saveEventLog", eventLog.BlockNumber, eventLog.Index)
	if eventLog.Removed {
		log.Info("event log Removed ")
		return
	}

	mTrade := models.Trade{}

	if err := models.DB.Where("block_number = ? and log_index = ?", eventLog.BlockNumber, eventLog.Index).First(&mTrade).Error; gorm.IsRecordNotFoundError(err) {
		match := MatchEvent{}
		err = contractABI.Unpack(&match, "Match", eventLog.Data)
		if err != nil {
			log.Warn(err)
			return
		}

		baseToken := GetToken(match.BaseToken.Hex())
		quoteToken := GetToken(match.QuoteToken.Hex())
		blockTime := getBlockTime(eventLog.BlockNumber)

		quoteTokenAmount := decimal.NewFromBigInt(match.QuoteTokenAmount, int32(-quoteToken.Decimals))
		quoteTokenPriceUSD := quoteToken.PriceUSD
		date := time.Unix(int64(blockTime), 0)
		// if duration is too long unset price now. fetch history price later.
		if time.Now().Sub(date) > time.Hour {
			quoteTokenPriceUSD = decimal.New(0, 0)
		}

		mTrade = models.Trade{
			BlockNumber:        eventLog.BlockNumber,
			BlockHash:          eventLog.BlockHash.Hex(),
			TransactionHash:    eventLog.TxHash.Hex(),
			LogIndex:           eventLog.Index,
			Date:               date,
			QuoteTokenPriceUSD: quoteTokenPriceUSD,
			VolumeUSD:          quoteTokenPriceUSD.Mul(quoteTokenAmount),
			BaseTokenAddress:   baseToken.Address,
			QuoteTokenAddress:  quoteToken.Address,
			RelayerAddress:     match.Relayer.Hex(),
			MakerAddress:       match.Maker.Hex(),
			TakerAddress:       match.Taker.Hex(),
			BaseTokenAmount:    decimal.NewFromBigInt(match.BaseTokenAmount, int32(-baseToken.Decimals)),
			QuoteTokenAmount:   quoteTokenAmount,
			MakerFee:           decimal.NewFromBigInt(match.MakerFee, int32(-quoteToken.Decimals)),
			TakerFee:           decimal.NewFromBigInt(match.TakerFee, int32(-quoteToken.Decimals)),
			MakerGasFee:        decimal.NewFromBigInt(match.MakerGasFee, int32(-quoteToken.Decimals)),
			MakerRebate:        decimal.NewFromBigInt(match.MakerRebate, int32(-quoteToken.Decimals)),
			TakerGasFee:        decimal.NewFromBigInt(match.TakerGasFee, int32(-quoteToken.Decimals)),
		}

		if err = models.DB.Where("block_number = ? and log_index = ?", eventLog.BlockNumber, eventLog.Index).First(&mTrade).Error; gorm.IsRecordNotFoundError(err) {
			models.DB.Create(&mTrade)
		}
	}
}

func getFromBlockNumber() uint64 {
	var number uint64
	mTrade := models.Trade{}
	if err := models.DB.Order("block_number desc").Take(&mTrade).Error; gorm.IsRecordNotFoundError(err) {
		i, err := strconv.ParseInt(os.Getenv("HYDRO_START_BLOCK_NUMBER"), 10, 64)
		if err != nil {
			panic(err)
		}
		number = uint64(i)
	} else {
		number = mTrade.BlockNumber
	}
	log.Info(number)
	return number
}

func getLastBlockNumber() uint64 {
	header, err := EthClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return header.Number.Uint64()
}

func getBlockTime(blockNumber uint64) uint64 {
	blockNumberBigInt := big.NewInt(int64(blockNumber))
	block, err := EthClient.BlockByNumber(context.Background(), blockNumberBigInt)
	if err != nil {
		panic(err)
	}
	return block.Time().Uint64()
}
