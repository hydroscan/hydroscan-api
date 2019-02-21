package task

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

var coinmarketcapSlugs map[string]string

type MissingTime struct {
	Min time.Time
	Max time.Time
}

type CoinmarketcapHistoryPrice struct {
	PriceUSD [][]decimal.Decimal `json:"price_usd"`
}

func UpdateHistoryTradePrice() {
	log.Info("UpdateHistoryTradePrice")
	getCoinmarketcapSlugs()
	for address, slug := range coinmarketcapSlugs {
		missingTime := MissingTime{}
		log.Info(address)
		models.DB.Raw("select min(date), max(date) from trades where quote_token_price_usd = 0 and quote_token_address = ?", address).Scan(&missingTime)
		if missingTime.Min.IsZero() { // no missing price
			continue
		}

		fetchHistoryTradePriceAndSave(address, slug, missingTime)
	}
}

func fetchHistoryTradePriceAndSave(address string, slug string, missingTime MissingTime) {
	log.Info(missingTime)
	from := missingTime.Min.Unix() - 3600
	end := missingTime.Max.Unix() + 3600

	for from < end {
		to := from + 30*24*3600 // fetch 1 month data once time
		fetchHistoryIntervalAndSave(address, slug, from, to)
		from = to
	}
}

func fetchHistoryIntervalAndSave(address string, slug string, from int64, to int64) {
	url := "https://graphs2.coinmarketcap.com/currencies/" + slug + "/" + strconv.FormatInt(from*1000, 10) + "/" + strconv.FormatInt(to*1000, 10)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	cmcHistoryPrice := CoinmarketcapHistoryPrice{}
	json.Unmarshal([]byte(body), &cmcHistoryPrice)

	lastTime := time.Unix(from, 0)
	for _, timePrice := range cmcHistoryPrice.PriceUSD {
		timeFloat64, _ := timePrice[0].Float64()
		nextTime := time.Unix(int64(timeFloat64/1000), 0)
		trades := []models.Trade{}
		models.DB.Where("quote_token_price_usd = 0 and quote_token_address = ? and date >= ? and date <= ?", address, lastTime, nextTime).Find(&trades)
		for _, trade := range trades {
			models.DB.Model(&trade).Updates(models.Trade{QuoteTokenPriceUSD: timePrice[1], VolumeUSD: timePrice[1].Mul(trade.QuoteTokenAmount)})
		}
		lastTime = nextTime
	}
}

func UpdateOnlyVolumeUSD() {
	models.DB.Exec("UPDATE trades SET volume_usd = quote_token_price_usd * quote_token_amount WHERE volume_usd = 0 and quote_token_price_usd != 0")
	log.Info("UpdateOnlyVolumeUSD")
}

func getCoinmarketcapSlugs() {
	jsonFile, err := os.Open("resource/coinmarketcap_slugs.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &coinmarketcapSlugs)
}
