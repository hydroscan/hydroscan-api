package task

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type MatchEventV1 struct {
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

type OrderAddressSet struct {
	BaseToken  common.Address
	QuoteToken common.Address
	Relayer    common.Address
}

type MatchResult struct {
	Maker                  common.Address
	Taker                  common.Address
	Buyer                  common.Address
	MakerFee               *big.Int
	MakerRebate            *big.Int
	TakerFee               *big.Int
	MakerGasFee            *big.Int
	TakerGasFee            *big.Int
	BaseTokenFilledAmount  *big.Int
	QuoteTokenFilledAmount *big.Int
}

type MatchEventV1_1 struct {
	AddressSet OrderAddressSet
	Result     MatchResult
}

func FetchHistoricalLogs(fetchAll bool) {
	log.Info("FetchHistoricalLogs")
	fromBlock := getFromBlockNumber()

	if fetchAll {
		// fromBlock = HydroStartBlockNumberV1
		fromBlock = HydroStartBlockNumberV1_1
	}

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

func FetchRecentLogs() {
	fromBlock := getFromBlockNumber()
	lastBlock := getLastBlockNumber()
	if lastBlock-100 < fromBlock {
		fromBlock = lastBlock - 100
	}
	fetchLogs(int64(fromBlock), int64(lastBlock))
}

func fetchLogs(fromBlock int64, toBlock int64) {
	log.Info("fetchLogs: ", fromBlock, " - ", toBlock)

	contractAddressV1 := common.HexToAddress(HydroExchangeAddressV1)
	contractAddressV1_1 := common.HexToAddress(HydroExchangeAddressV1_1)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: []common.Address{
			contractAddressV1,
			contractAddressV1_1,
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
	log.Info("saveEventLog: ", eventLog.BlockNumber, eventLog.Index)

	mTrade := models.Trade{}

	if err := models.DB.Where("block_number = ? AND log_index = ?",
		eventLog.BlockNumber, eventLog.Index).First(&mTrade).Error; gorm.IsRecordNotFoundError(err) {

		if eventLog.Removed {
			log.Info("Event Log Removed: ", eventLog.BlockNumber, eventLog.Index)
			return
		}

		if eventLog.BlockNumber >= HydroStartBlockNumberV1_1 {
			saveEventLogV1_1(eventLog)
		} else {
			saveEventLogV1(eventLog)
		}

	} else {
		if eventLog.Removed {
			models.DB.Delete(&mTrade)
			log.Info("Event Log Found and Removed ", eventLog.BlockNumber, eventLog.Index)
			return
		}
	}
}

func saveEventLogV1(eventLog types.Log) {
	log.Info("saveEventLogV1: ", eventLog.BlockNumber, eventLog.Index)

	mTrade := models.Trade{}
	match := MatchEventV1{}
	err = contractABIV1.Unpack(&match, "Match", eventLog.Data)
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
		ProtocolVersion:    ProtocolV1,
	}

	if err = models.DB.Where("block_number = ? AND log_index = ?", eventLog.BlockNumber, eventLog.Index).First(&mTrade).Error; gorm.IsRecordNotFoundError(err) {
		models.DB.Create(&mTrade)
		log.Info("Saved Event Log: ", eventLog.BlockNumber, eventLog.Index)
	}
}

func saveEventLogV1_1(eventLog types.Log) {
	log.Info("saveEventLogV1_1: ", eventLog.BlockNumber, eventLog.Index)

	mTrade := models.Trade{}
	match := MatchEventV1_1{}
	err = contractABIV1_1.Unpack(&match, "Match", eventLog.Data)
	if err != nil {
		log.Warn(err)
		return
	}

	baseToken := GetToken(match.AddressSet.BaseToken.Hex())
	quoteToken := GetToken(match.AddressSet.QuoteToken.Hex())
	blockTime := getBlockTime(eventLog.BlockNumber)

	quoteTokenAmount := decimal.NewFromBigInt(match.Result.QuoteTokenFilledAmount, int32(-quoteToken.Decimals))
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
		RelayerAddress:     match.AddressSet.Relayer.Hex(),
		MakerAddress:       match.Result.Maker.Hex(),
		TakerAddress:       match.Result.Taker.Hex(),
		BuyerAddress:       match.Result.Buyer.Hex(),
		BaseTokenAmount:    decimal.NewFromBigInt(match.Result.BaseTokenFilledAmount, int32(-baseToken.Decimals)),
		QuoteTokenAmount:   quoteTokenAmount,
		MakerFee:           decimal.NewFromBigInt(match.Result.MakerFee, int32(-quoteToken.Decimals)),
		TakerFee:           decimal.NewFromBigInt(match.Result.TakerFee, int32(-quoteToken.Decimals)),
		MakerGasFee:        decimal.NewFromBigInt(match.Result.MakerGasFee, int32(-quoteToken.Decimals)),
		MakerRebate:        decimal.NewFromBigInt(match.Result.MakerRebate, int32(-quoteToken.Decimals)),
		TakerGasFee:        decimal.NewFromBigInt(match.Result.TakerGasFee, int32(-quoteToken.Decimals)),
		ProtocolVersion:    ProtocolV1_1,
	}

	if err = models.DB.Where("block_number = ? AND log_index = ?", eventLog.BlockNumber, eventLog.Index).First(&mTrade).Error; gorm.IsRecordNotFoundError(err) {
		models.DB.Create(&mTrade)
		log.Info("Saved Event Log: ", eventLog.BlockNumber, eventLog.Index)
	}
}

func getFromBlockNumber() uint64 {
	var number uint64
	mTrade := models.Trade{}
	if err := models.DB.Order("block_number desc").Take(&mTrade).Error; gorm.IsRecordNotFoundError(err) {
		number = uint64(HydroStartBlockNumberV1)
	} else {
		number = mTrade.BlockNumber
	}
	log.Info("getFromBlockNumber ", number)
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
	log.Info("getBlockTime ", blockNumber)
	blockNumberBigInt := big.NewInt(int64(blockNumber))

	var block *types.Block
	var err error
	dialRetries := MaxReties

	for dialRetries == MaxReties || (err != nil && dialRetries > 0) {
		if dialRetries != MaxReties {
			time.Sleep(1000 * time.Millisecond)
		}
		log.Info("getBlockTime dialRetries ", dialRetries)
		block, err = EthClient.BlockByNumber(context.Background(), blockNumberBigInt)
		dialRetries -= 1
	}
	if err != nil {
		log.Warn("getBlockTime err ")
		panic(err)
	}

	return block.Time().Uint64()
}
