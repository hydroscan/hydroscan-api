package apiv1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/hydroscan/hydroscan-api/utils"
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
	Transaction       string `form:"transaction"`
}

type TradesChartQuery struct {
	Filter         string `form:"filter"`
	TokenAddress   string `form:"tokenAddress"`
	TraderAddress  string `form:"traderAddress"`
	RelayerAddress string `form:"relayerAddress"`
}

func GetTrades(c *gin.Context) {
	query := TradesQuery{1, 25, "", "", "", "", "", ""}
	c.BindQuery(&query)

	page := query.Page
	pageSize := query.PageSize
	offset := (page - 1) * pageSize

	var trades []models.Trade
	statment := models.DB.Table("trades").Order("block_number desc").Order("log_index desc")
	if query.Transaction != "" {
		statment = statment.Where("transaction_hash = ?", query.Transaction)
	} else if query.BaseTokenAddress != "" && query.QuoteTokenAddress != "" {
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

func GetTrader(c *gin.Context) {
	address := c.Params.ByName("address")
	type TopToken struct {
		Address      string          `json:"address"`
		Name         string          `json:"name"`
		Symbol       string          `json:"symbol"`
		Volume       decimal.Decimal `json:"volume"`
		VolumeLast   decimal.Decimal `json:"volumeLast"`
		VolumeChange float32         `json:"volumeChange"`
	}

	var res struct {
		Address          string          `json:"address"`
		Volume24h        decimal.Decimal `json:"volume24h"`
		Volume24hLast    decimal.Decimal `json:"volume24hLast"`
		Volume24hChange  float32         `json:"volume24hChange"`
		Trades24h        decimal.Decimal `json:"trades24h"`
		Trades24hLast    decimal.Decimal `json:"trades24hLast"`
		Trades24hChange  float32         `json:"trades24hChange"`
		TotalMakerRabate decimal.Decimal `json:"totalMakerRabate"`
		TopTokens        []TopToken      `json:"topTokens"`
	}

	timeNow := time.Now()
	time24hAgo := time.Now().Add(-24 * time.Hour)
	time48hAgo := time.Now().Add(-48 * time.Hour)

	models.DB.Raw(`SELECT sum(trades.volume_usd) AS volume24h, count(*) AS trades24h
		FROM trades WHERE (trades.maker_address = ? OR trades.taker_address = ?) AND date >= ? AND date < ?`,
		address, address, time24hAgo, timeNow).Scan(&res)

	models.DB.Raw(`SELECT sum(trades.volume_usd) AS volume24h_last, count(*) AS trades24h_last
		FROM trades WHERE (trades.maker_address = ? OR trades.taker_address = ?) AND date >= ? AND date < ?`,
		address, address, time48hAgo, time24hAgo).Scan(&res)

	models.DB.Raw(`SELECT sum(maker_rebate) FROM trades WHERE (trades.maker_address = ? OR trades.taker_address = ?)`,
		address, address).Scan(&res)

	var topTokens []TopToken

	const baseFields = "t.address, t.name, t.symbol"
	models.DB.Raw(`SELECT `+baseFields+`, t.volume, sum(trades.volume_usd) AS volume_last
		FROM (
			SELECT `+baseFields+`, sum(trades.volume_usd) AS volume FROM tokens AS t, trades
			WHERE (trades.base_token_address = t.address OR trades.quote_token_address = t.address)
			AND (trades.maker_address = ? OR trades.taker_address = ?)
			AND trades.date >= ? AND trades.date < ?
			GROUP BY `+baseFields+`
			ORDER BY volume DESC LIMIT 3 OFFSET 0
		) AS t LEFT JOIN trades
		ON (t.address = trades.base_token_address OR t.address = trades.quote_token_address)
		AND (trades.maker_address = ? OR trades.taker_address = ?)
		AND trades.date >= ? AND trades.date < ?
		GROUP BY `+baseFields+`, t.volume
		ORDER BY t.volume DESC`,
		address, address, time24hAgo, timeNow,
		address, address, time48hAgo, time24hAgo).Scan(&topTokens)

	for i, token := range topTokens {
		if !token.VolumeLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := token.Volume.Sub(token.VolumeLast).Div(token.VolumeLast).Float64()
			topTokens[i].VolumeChange = float32(changeFloat64)
		}
	}

	res.TopTokens = topTokens

	if !res.Volume24hLast.Equal(decimal.NewFromFloat32(0)) {
		changeFloat64, _ := res.Volume24h.Sub(res.Volume24hLast).Div(res.Volume24hLast).Float64()
		res.Volume24hChange = float32(changeFloat64)
	}

	if !res.Trades24hLast.Equal(decimal.NewFromFloat32(0)) {
		changeFloat64, _ := res.Trades24h.Sub(res.Trades24hLast).Div(res.Trades24hLast).Float64()
		res.Trades24hChange = float32(changeFloat64)
	}

	res.Address = address
	c.JSON(200, res)
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

func GetTradesSearch(c *gin.Context) {
	keyword := c.Query("keyword")

	var res struct {
		SearchType string `json:"searchType"`
		SearchKey  string `json:"searchKey"`
	}

	if utils.IsAddress(keyword) {
		if isTrue, searchKey := utils.IsToken(keyword); isTrue {
			res.SearchType = "TOKEN"
			res.SearchKey = searchKey

		} else if isTrue, searchKey := utils.IsRelayer(keyword); isTrue {
			res.SearchType = "RELAYER"
			res.SearchKey = searchKey

		} else if isTrue, searchKey := utils.IsTrader(keyword); isTrue {
			res.SearchType = "TRADER"
			res.SearchKey = searchKey

		}
	} else if utils.IsTransaction(keyword) {
		res.SearchType = "TRANSACTION"
		res.SearchKey = keyword

	} else {
		token := models.Token{}
		if err := models.DB.Where("name ILIKE ? or symbol ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").Order("volume_24h desc").First(&token).Error; err == nil {
			res.SearchType = "TOKENS"
			res.SearchKey = keyword
		}
	}

	c.JSON(200, res)
}
