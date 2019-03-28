package task

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func SubscribeLogs() {
	FetchHistoricalLogs()

	var client *ethclient.Client
	var err error
	dialRetries := MaxReties

	log.Info("dial: ", viper.GetString("web3_ws"))
	for dialRetries == MaxReties || (err != nil && dialRetries > 0) {
		if dialRetries != MaxReties {
			time.Sleep(1000 * time.Millisecond)
		}
		client, err = ethclient.Dial(viper.GetString("web3_ws"))
		dialRetries -= 1
	}
	if err != nil {
		panic(err)
	}

	contractAddress := common.HexToAddress(HydroExchangeAddressV1_1)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	log.Info("SubscribeFilterLogs")
	eventLogs := make(chan types.Log)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sub, err := client.SubscribeFilterLogs(ctx, query, eventLogs)
	if err != nil {
		panic(err)
	}

	log.Info("for select run")
	for {
		select {
		case err := <-sub.Err():
			log.Warn("select sub err ", err)

			dialRetries = MaxReties
			for err != nil && dialRetries > 0 {
				log.Info("dialRetries: ", dialRetries)
				if dialRetries != MaxReties {
					time.Sleep(1000 * time.Millisecond)
				}
				client, err = ethclient.Dial(viper.GetString("web3_ws"))
				dialRetries -= 1
			}
			if err != nil {
				panic(err)
			}

			log.Info("retry SubscribeFilterLogs")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			sub, err = client.SubscribeFilterLogs(ctx, query, eventLogs)
			if err != nil {
				panic(err)
			}
			log.Info("subscribed")

		case eventLog := <-eventLogs:
			log.Info("receive log: ", eventLog.BlockNumber, eventLog.Index)
			checkMissingBlocks(eventLog.BlockNumber)
			saveEventLog(eventLog)

		case <-time.After(60 * time.Second):
			log.Info("timeout 1min retry dial")
			// https://github.com/ethereum/go-ethereum/blob/245f3146c26698193c4b479e7bc5825b058c444a/rpc/subscription.go#L243
			sub.Unsubscribe()
		}
	}
}

func checkMissingBlocks(blockNumber uint64) {
	fromBlock := getFromBlockNumber()
	if blockNumber-fromBlock > 1 {
		fetchLogs(int64(fromBlock), int64(blockNumber))
	}
}
