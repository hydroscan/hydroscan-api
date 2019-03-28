package task

import (
	"time"

	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/shopspring/decimal"

	log "github.com/sirupsen/logrus"
)

type RelayerTrades24hData struct {
	Trades24h        uint64
	Trades24hChange  float32
	Traders24h       uint64
	Traders24hChange float32
}

func UpdateRelayerVolume24h() {
	log.Info("UpdateRelayerVolume24h")
	var relayers []models.Relayer
	models.DB.Find(&relayers)

	timeNow := time.Now()
	time24hAgo := time.Now().Add(-24 * time.Hour)
	time48hAgo := time.Now().Add(-48 * time.Hour)

	type QueryResult struct {
		Volume24h     decimal.Decimal
		Volume24hLast decimal.Decimal
	}

	for _, relayer := range relayers {
		result := QueryResult{}

		models.DB.Raw("SELECT sum(volume_usd) AS volume24h FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?",
			time24hAgo, timeNow, relayer.Address).Scan(&result)
		models.DB.Raw("SELECT sum(volume_usd) AS volume24h_last FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?",
			time48hAgo, time24hAgo, relayer.Address).Scan(&result)

		// the minimum change is -1 (-100%), so using -2 represent not a number
		var change float32 = -2
		if !result.Volume24hLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := result.Volume24h.Sub(result.Volume24hLast).Div(result.Volume24hLast).Float64()
			change = float32(changeFloat64)
		}
		// Using map[string]interface{} instead of models.Relayer since
		// when update with struct, GORM will only update those fields that with non blank value
		// nothing will be updated as "", 0, false are blank values of their types
		models.DB.Model(&relayer).Updates(map[string]interface{}{"volume_24h": result.Volume24h, "volume_24h_change": change})
	}
}

func UpdateRelayerTrades24h() {
	log.Info("UpdateRelayerTrades")
	var relayers []models.Relayer
	models.DB.Find(&relayers)

	timeNow := time.Now()
	time24hAgo := time.Now().Add(-24 * time.Hour)
	time48hAgo := time.Now().Add(-48 * time.Hour)

	type QueryResult struct {
		Trades24h      uint64
		Trades24hLast  uint64
		Traders24h     uint64
		Traders24hLast uint64
	}

	for _, relayer := range relayers {
		result := QueryResult{}

		models.DB.Raw("SELECT count(*) AS trades24h FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?",
			time24hAgo, timeNow, relayer.Address).Scan(&result)
		models.DB.Raw("SELECT count(*) AS trades24h_last FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?",
			time48hAgo, time24hAgo, relayer.Address).Scan(&result)

		models.DB.Raw(`SELECT count(*) AS traders24h
			FROM (
			SELECT maker_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
     		UNION
     		SELECT taker_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
     		) AS traders`,
			time24hAgo, timeNow, relayer.Address,
			time24hAgo, timeNow, relayer.Address).Scan(&result)
		models.DB.Raw(`SELECT count(*) AS traders24h_last
			FROM (
			SELECT maker_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
     		UNION
     		SELECT taker_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
     		) AS traders`,
			time48hAgo, time24hAgo, relayer.Address,
			time48hAgo, time24hAgo, relayer.Address).Scan(&result)

		var trades24hChange float32 = -2
		if result.Trades24hLast != 0 {
			trades24hChange = float32((float64(result.Trades24h) - float64(result.Trades24hLast)) / float64(result.Trades24hLast))
		}

		var traders24hChange float32 = -2
		if result.Traders24hLast != 0 {
			traders24hChange = float32((float64(result.Traders24h) - float64(result.Traders24hLast)) / float64(result.Traders24hLast))
		}

		data := RelayerTrades24hData{result.Trades24h, trades24hChange, result.Traders24h, traders24hChange}
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		log.Info(string(b))
		err = redis.Client.HSet("TOKENS_TRADES_24H_DATA", relayer.Address, string(b)).Err()
		if err != nil {
			panic(err)
		}
	}
}

func GetRelayerTrades24hData(address string) RelayerTrades24hData {
	data := RelayerTrades24hData{}
	res, err := redis.Client.HGet("TOKENS_TRADES_24H_DATA", address).Result()
	if err != nil {
		return data
	}
	json.Unmarshal([]byte(res), &data)

	return data
}
