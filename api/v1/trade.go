package apiv1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type TradesQuery struct {
	Page              int    `form:"page"`
	PageSize          int    `form:"pageSize"`
	BaseTokenAddress  string `form:"baseTokenAddress"`
	QuoteTokenAddress string `form:"quoteTokenAddress"`
	TokenAddress      string `form:"tokenAddress"`
	TraderAddress     string `form:"traderAddress"`
	RelayerAddress    string `form:"relayerAddress"`
}

type TradesChartQuery struct {
	Filter         string `form:"filter"`
	TokenAddress   string `form:"tokenAddress"`
	TraderAddress  string `form:"traderAddress"`
	RelayerAddress string `form:"relayerAddress"`
}

func GetTrades(c *gin.Context) {
	query := TradesQuery{1, 25, "", "", "", "", ""}
	c.BindQuery(&query)

	page := query.Page
	pageSize := query.PageSize
	offset := (page - 1) * pageSize

	var trades []models.Trade
	statment := models.DB.Table("trades").Order("block_number desc").Order("log_index desc")
	if query.BaseTokenAddress != "" && query.QuoteTokenAddress != "" {
		statment = statment.Where("base_token_address = ? AND quote_token_address = ?", query.BaseTokenAddress, query.QuoteTokenAddress)
	} else if query.TokenAddress != "" {
		statment = statment.Where("base_token_address = ? OR quote_token_address = ?", query.TokenAddress, query.TokenAddress)
	} else if query.TraderAddress != "" {
		statment = statment.Where("maker_address = ? OR taker_address = ?", query.TraderAddress, query.TraderAddress)
	} else if query.RelayerAddress != "" {
		statment = statment.Where("relayer_address = ?", query.RelayerAddress)
	}

	if err := statment.Offset(offset).Limit(pageSize).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").Find(&trades).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		type resType struct {
			Page     int            `json:"page"`
			PageSize int            `json:"pageSize"`
			Count    uint64         `json:"count"`
			Trades   []models.Trade `json:"trades"`
		}
		res := resType{page, pageSize, 0, trades}
		statment.Count(&res.Count)

		c.JSON(200, res)
	}
}

func GetTrade(c *gin.Context) {
	uuid := c.Params.ByName("uuid")
	trade := models.Trade{}
	if err := models.DB.Where("uuid = ?", uuid).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").First(&trade).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, trade)
	}
}

func GetTradesChart(c *gin.Context) {
	query := TradesChartQuery{"1M", "", "", ""}
	c.BindQuery(&query)

	log.Info(query)
	trunc := "day"
	from := time.Now().Add(-30 * 24 * time.Hour)
	switch query.Filter {
	case "24H":
		trunc = "hour"
		from = time.Now().Add(-24 * time.Hour)
	case "7D":
		trunc = "hour"
		from = time.Now().Add(-7 * 24 * time.Hour)
	case "1M":
		trunc = "day"
		from = time.Now().Add(-30 * 24 * time.Hour)
	case "1Y":
		trunc = "day"
		from = time.Now().Add(-365 * 24 * time.Hour)
	case "ALL":
		trunc = "day"
		from = time.Now().Add(-1000 * 24 * time.Hour)
	default:
		c.AbortWithStatus(404)
		return
	}

	var res []struct {
		Dt           time.Time       `json:"date"`
		Sum          decimal.Decimal `json:"volume"`
		TradesCount  uint64          `json:"trades"`
		TradersCount uint64          `json:"traders"`
	}
	var resTraders []struct {
		TradersCount uint64 `json:"traders"`
	}

	if query.TokenAddress != "" {

		models.DB.Raw(`SELECT date_trunc(?, date) AS dt, sum(volume_usd), count(*) AS trades_count
		FROM trades WHERE date >= ? AND (base_token_address = ? OR quote_token_address = ?)
		GROUP BY dt ORDER BY dt`, trunc, from, query.TokenAddress, query.TokenAddress).Scan(&res)

		models.DB.Raw(`SELECT dt, count(*) AS traders_count
		FROM (
		SELECT date_trunc(?, date) AS dt, maker_address FROM trades WHERE date > ? AND (base_token_address = ? OR quote_token_address = ?)
		UNION
		SELECT date_trunc(?, date) AS dt, taker_address FROM trades WHERE date > ? AND (base_token_address = ? OR quote_token_address = ?)
		) AS traders GROUP BY dt ORDER BY dt`,
			trunc, from, query.TokenAddress, query.TokenAddress,
			trunc, from, query.TokenAddress, query.TokenAddress).Scan(&resTraders)

	} else if query.TraderAddress != "" {

		models.DB.Raw(`SELECT date_trunc(?, date) AS dt, sum(volume_usd), count(*) AS trades_count
		FROM trades WHERE date >= ? AND (maker_address = ? OR taker_address = ?)
		GROUP BY dt ORDER BY dt`, trunc, from, query.TraderAddress, query.TraderAddress).Scan(&res)

	} else if query.RelayerAddress != "" {

		models.DB.Raw(`SELECT date_trunc(?, date) AS dt, sum(volume_usd), count(*) AS trades_count
		FROM trades WHERE date >= ? AND relayer_address = ?
		GROUP BY dt ORDER BY dt`, trunc, from, query.RelayerAddress).Scan(&res)

		models.DB.Raw(`SELECT dt, count(*) AS traders_count
		FROM (
		SELECT date_trunc(?, date) AS dt, maker_address FROM trades WHERE date > ? AND relayer_address = ?
		UNION
		SELECT date_trunc(?, date) AS dt, taker_address FROM trades WHERE date > ? AND relayer_address = ?
		) AS traders GROUP BY dt ORDER BY dt`,
			trunc, from, query.RelayerAddress,
			trunc, from, query.RelayerAddress).Scan(&resTraders)

	} else {

		models.DB.Raw(`SELECT date_trunc(?, date) AS dt, sum(volume_usd), count(*) AS trades_count
		FROM trades WHERE date >= ? GROUP BY dt ORDER BY dt`, trunc, from).Scan(&res)

		// select traders
		// SELECT dt, count(*) FROM (SELECT date_trunc('hour', date) AS dt, maker_address FROM trades WHERE date > '2019-02-26t00:00:00+08:00'UNION SELECT date_trunc('hour', date) AS dt, taker_address FROM trades WHERE date > '2019-02-26t00:00:00+08:00' ) AS traders GROUP BY dt ORDER BY dt;
		models.DB.Raw(`SELECT dt, count(*) AS traders_count
		FROM (
		SELECT date_trunc(?, date) AS dt, maker_address FROM trades WHERE date > ?
		UNION
		SELECT date_trunc(?, date) AS dt, taker_address FROM trades WHERE date > ?
		) AS traders GROUP BY dt ORDER BY dt`, trunc, from, trunc, from).Scan(&resTraders)

	}

	if len(resTraders) > 0 {
		for i, _ := range res {
			res[i].TradersCount = resTraders[i].TradersCount
		}
	}

	c.JSON(200, res)
}

func GetTradesIndicators(c *gin.Context) {
	res, err := redis.Client.Get("indicators").Result()
	if err != nil {
		panic(err)
	}
	indicators := task.Indicators{}
	json.Unmarshal([]byte(res), &indicators)
	c.JSON(200, indicators)
}
