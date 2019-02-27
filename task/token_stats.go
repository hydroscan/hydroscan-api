package task

import (
	"time"

	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/shopspring/decimal"

	log "github.com/sirupsen/logrus"
)

type TokenTrades24hData struct {
	Trades24h        uint64
	Trades24hChange  float32
	Traders24h       uint64
	Traders24hChange float32
}

func UpdateTokenVolume24h() {
	log.Info("UpdateTokenVolume24h")
	var tokens []models.Token
	models.DB.Find(&tokens)

	timeNow := time.Now()
	time24hAgo := time.Now().Add(-24 * time.Hour)
	time48hAgo := time.Now().Add(-48 * time.Hour)

	type QueryResult struct {
		Volume24h     decimal.Decimal
		Volume24hLast decimal.Decimal
	}

	for _, token := range tokens {
		result := QueryResult{}

		models.DB.Raw("SELECT sum(volume_usd) AS volume24h FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)",
			time24hAgo, timeNow, token.Address, token.Address).Scan(&result)
		models.DB.Raw("SELECT sum(volume_usd) AS volume24h_last FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)",
			time48hAgo, time24hAgo, token.Address, token.Address).Scan(&result)

		var change float32 = 0
		if !result.Volume24hLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := result.Volume24h.Sub(result.Volume24hLast).Div(result.Volume24hLast).Float64()
			change = float32(changeFloat64)
		}
		// Using map[string]interface{} instead of models.Token since
		// when update with struct, GORM will only update those fields that with non blank value
		// nothing will be updated as "", 0, false are blank values of their types
		models.DB.Model(&token).Updates(map[string]interface{}{"volume_24h": result.Volume24h, "volume_24h_change": change})
	}
}

func UpdateTokenVolume7d() {
	log.Info("UpdateTokenVolume7d")
	var tokens []models.Token
	models.DB.Find(&tokens)

	timeNow := time.Now()
	time7dAgo := time.Now().Add(-7 * 24 * time.Hour)
	time14dAgo := time.Now().Add(-14 * 24 * time.Hour)

	type QueryResult struct {
		Volume7d     decimal.Decimal
		Volume7dLast decimal.Decimal
	}

	for _, token := range tokens {
		result := QueryResult{}

		models.DB.Raw("SELECT sum(volume_usd) AS volume7d FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)",
			time7dAgo, timeNow, token.Address, token.Address).Scan(&result)
		models.DB.Raw("SELECT sum(volume_usd) AS volume7d_last FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)",
			time14dAgo, time7dAgo, token.Address, token.Address).Scan(&result)

		var change float32 = 0
		if !result.Volume7dLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := result.Volume7d.Sub(result.Volume7dLast).Div(result.Volume7dLast).Float64()
			change = float32(changeFloat64)
		}
		models.DB.Model(&token).Updates(map[string]interface{}{"volume_7d": result.Volume7d, "volume_7d_change": change})
	}
}

func UpdateTokenVolume1m() {
	log.Info("UpdateTokenVolume1m")
	var tokens []models.Token
	models.DB.Find(&tokens)

	timeNow := time.Now()
	time1mAgo := time.Now().Add(-30 * 24 * time.Hour)
	time2mAgo := time.Now().Add(-60 * 24 * time.Hour)

	type QueryResult struct {
		Volume1m     decimal.Decimal
		Volume1mLast decimal.Decimal
	}

	for _, token := range tokens {
		result := QueryResult{}

		models.DB.Raw("SELECT sum(volume_usd) AS volume1m FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)",
			time1mAgo, timeNow, token.Address, token.Address).Scan(&result)
		models.DB.Raw("SELECT sum(volume_usd) AS volume1m_last FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)",
			time2mAgo, time1mAgo, token.Address, token.Address).Scan(&result)

		var change float32 = 0
		if !result.Volume1mLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := result.Volume1m.Sub(result.Volume1mLast).Div(result.Volume1mLast).Float64()
			change = float32(changeFloat64)
		}
		models.DB.Model(&token).Updates(map[string]interface{}{"volume_1m": result.Volume1m, "volume_1m_change": change})
	}
}

func UpdateTokenVolumeAll() {
	log.Info("UpdateTokenVolume7d")
	var tokens []models.Token
	models.DB.Find(&tokens)

	type QueryResult struct {
		VolumeAll decimal.Decimal
	}

	for _, token := range tokens {
		result := QueryResult{}
		models.DB.Raw("SELECT sum(volume_usd) AS volume_all FROM trades WHERE base_token_address = ? OR quote_token_address = ?",
			token.Address, token.Address).Scan(&result)
		models.DB.Model(&token).Updates(map[string]interface{}{"volume_all": result.VolumeAll})
	}
}

func UpdateTokenTrades24h() {
	log.Info("UpdateTokenTrades")
	var tokens []models.Token
	models.DB.Find(&tokens)

	timeNow := time.Now()
	time24hAgo := time.Now().Add(-24 * time.Hour)
	time48hAgo := time.Now().Add(-48 * time.Hour)

	type QueryResult struct {
		Trades24h      uint64
		Trades24hLast  uint64
		Traders24h     uint64
		Traders24hLast uint64
	}

	for _, token := range tokens {
		result := QueryResult{}

		models.DB.Raw("SELECT count(*) AS trades24h FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)",
			time24hAgo, timeNow, token.Address, token.Address).Scan(&result)
		models.DB.Raw("SELECT count(*) AS trades24h_last FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)",
			time48hAgo, time24hAgo, token.Address, token.Address).Scan(&result)

		models.DB.Raw(`SELECT count(*) AS traders24h
			FROM (
			SELECT maker_address FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)
     		UNION
     		SELECT taker_address FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)
     		) AS traders`,
			time24hAgo, timeNow, token.Address, token.Address,
			time24hAgo, timeNow, token.Address, token.Address).Scan(&result)
		models.DB.Raw(`SELECT count(*) AS traders24h_last
			FROM (
			SELECT maker_address FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)
     		UNION
     		SELECT taker_address FROM trades WHERE date > ? AND date <= ? AND (base_token_address = ? OR quote_token_address = ?)
     		) AS traders`,
			time48hAgo, time24hAgo, token.Address, token.Address,
			time48hAgo, time24hAgo, token.Address, token.Address).Scan(&result)

		var trades24hChange float32 = 0
		if result.Trades24hLast != 0 {
			trades24hChange = float32((float64(result.Trades24h) - float64(result.Trades24hLast)) / float64(result.Trades24hLast))
		}

		var traders24hChange float32 = 0
		if result.Traders24hLast != 0 {
			traders24hChange = float32((float64(result.Traders24h) - float64(result.Traders24hLast)) / float64(result.Traders24hLast))
		}

		data := TokenTrades24hData{result.Trades24h, trades24hChange, result.Traders24hLast, traders24hChange}
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		log.Info(string(b))
		err = redis.Client.HSet("TOKENS_TRADES_24H_DATA", token.Address, string(b)).Err()
		if err != nil {
			panic(err)
		}
	}
}

func GetTrades24hData(address string) TokenTrades24hData {
	res, err := redis.Client.HGet("TOKENS_TRADES_24H_DATA", address).Result()
	if err != nil {
		panic(err)
	}
	data := TokenTrades24hData{}
	json.Unmarshal([]byte(res), &data)

	return data
}
