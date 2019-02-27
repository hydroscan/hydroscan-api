package task

import (
	"time"

	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/shopspring/decimal"

	log "github.com/sirupsen/logrus"
)

type Indicators struct {
	Volume24h       decimal.Decimal `json:"volume24h"`
	Trades24h       decimal.Decimal `json:"trades24h"`
	Traders24h      decimal.Decimal `json:"traders24h"`
	MarketRabate24h decimal.Decimal `json:"marketRabate24h"`
}

func UpdateIndicators() {
	log.Info("UpdateIndicators")
	time24hAgo := time.Now().Add(-24 * time.Hour)
	indicators := Indicators{}

	models.DB.Raw("SELECT sum(volume_usd) AS volume24h FROM trades WHERE date > ?", time24hAgo).Scan(&indicators)
	models.DB.Raw("SELECT count(*) AS trades24h FROM trades WHERE date > ?", time24hAgo).Scan(&indicators)
	models.DB.Raw("SELECT count(*) AS traders24h FROM (SELECT maker_address FROM trades WHERE date > ? UNION SELECT taker_address FROM trades WHERE date > ? ) AS traders",
		time24hAgo, time24hAgo).Scan(&indicators)
	models.DB.Raw("SELECT sum(maker_rebate) AS market_rabate24h FROM trades WHERE date > ?", time24hAgo).Scan(&indicators)
	log.Info(indicators)
	b, err := json.Marshal(indicators)
	if err != nil {
		panic(err)
	}
	log.Info(string(b))

	err = redis.Client.Set("indicators", string(b), 0).Err()
	if err != nil {
		panic(err)
	}
}
