package task

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const RetryNumber = 3

func SubscribeLogs() {
	FetchHistoricalLogs()

	var client *ethclient.Client
	var err error
	dialRetry := RetryNumber

	for dialRetry == RetryNumber || (err != nil && dialRetry > 0) {
		client, err = ethclient.Dial(os.Getenv("WEB3_WS"))
		dialRetry -= 1
	}
	if err != nil {
		panic(err)
	}

	contractAddress := common.HexToAddress(os.Getenv("HYDRO_EXCHANGE_ADDRESS"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	eventLogs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, eventLogs)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Warn(err)

			dialRetry = RetryNumber
			for err != nil && dialRetry > 0 {
				client, err = ethclient.Dial(os.Getenv("WEB3_WS"))
				dialRetry -= 1
			}
			if err != nil {
				panic(err)
			}

			sub, err = client.SubscribeFilterLogs(context.Background(), query, eventLogs)
			if err != nil {
				panic(err)
			}

		case eventLog := <-eventLogs:
			log.Info("recieve log: ", eventLog.BlockNumber, eventLog.Index)
			checkMissingBlocks(eventLog.BlockNumber)
			saveEventLog(eventLog)
		}
	}
}

func checkMissingBlocks(blockNumber uint64) {
	fromBlock := getFromBlockNumber()
	if blockNumber-fromBlock > 1 {
		fetchLogs(int64(fromBlock), int64(blockNumber))
	}
}
