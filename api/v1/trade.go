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
)

type TradesQuery struct {
	Page              int    `form:"page"`
	PageSize          int    `form:"pageSize"`
	BaseTokenAddress  string `form:"baseTokenAddress"`
	QuoteTokenAddress string `form:"quoteTokenAddress"`
}

func GetTrades(c *gin.Context) {
	query := TradesQuery{1, 25, "", ""}
	c.BindQuery(&query)

	page := query.Page
	pageSize := query.PageSize
	offset := (page - 1) * pageSize

	var trades []models.Trade
	if err := models.DB.Order("block_number desc").Order("log_index desc").Offset(offset).Limit(pageSize).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").Find(&trades).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		type resType struct {
			Page     int            `json:"page"`
			PageSize int            `json:"pageSize"`
			Count    uint64         `json:"count"`
			Trades   []models.Trade `json:"trades"`
		}
		res := resType{page, pageSize, 0, trades}
		models.DB.Table("trades").Count(&res.Count)

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
	filter := c.DefaultQuery("filter", "1M")
	var res []struct {
		Dt           time.Time       `json:"date"`
		Sum          decimal.Decimal `json:"volume"`
		TradesCount  uint64          `json:"trades"`
		TradersCount uint64          `json:"traders"`
	}
	trunc := "day"
	from := time.Now().Add(-30 * 24 * time.Hour)
	switch filter {
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
	models.DB.Raw(`select date_trunc(?, date) as dt, sum(volume_usd), count(*) as trades_count
		from trades where date >= ? group by dt order by dt`, trunc, from).Scan(&res)
	// select dt, count(*) from ( select date_trunc('hour', date) as dt, maker_address from trades where date > '2019-02-26t00:00:00+08:00' union select date_trunc('hour', date) as dt, taker_address from trades where date > '2019-02-26t00:00:00+08:00' ) as traders group by dt order by dt;
	// select traders

	var resTraders []struct {
		TradersCount uint64 `json:"traders"`
	}
	models.DB.Raw(`select dt, count(*) as traders_count from (
		select date_trunc(?, date) as dt, maker_address from trades WHERE date > ?
		union
		select date_trunc(?, date) as dt, taker_address from trades WHERE date > ?
		) as traders group by dt order by dt`, trunc, from, trunc, from).Scan(&resTraders)

	for i, _ := range res {
		res[i].TradersCount = resTraders[i].TradersCount
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
