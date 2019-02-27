package apiv1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type TokensQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Filter   string `form:"filter"`
}

func GetTokens(c *gin.Context) {
	query := TokensQuery{1, 25, "24H"}
	c.BindQuery(&query)

	order := "volume_24h desc"
	switch query.Filter {
	case "24H":
		order = "volume_24h desc"
	case "7D":
		order = "volume_7d desc"
	case "1M":
		order = "volume_1m desc"
	case "ALL":
		order = "volume_all desc"
	default:
		c.AbortWithStatus(404)
		return
	}

	page := query.Page
	pageSize := query.PageSize
	offset := (page - 1) * pageSize

	var tokens []models.Token
	if err := models.DB.Offset(offset).Limit(pageSize).Order(order).Find(&tokens).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		for i, _ := range tokens {
			tradesData := task.GetTrades24hData(tokens[i].Address)
			// tokens[i].Trades24h = tradesData.Trades24h
			// tokens[i].Trades24hChange = tradesData.Trades24hChange
			tokens[i].Traders24h = tradesData.Traders24h
			// tokens[i].Traders24hChange = tradesData.Traders24hChange
			tokens[i].Amount24h = tradesData.Amount24h
		}

		type resType struct {
			Page     int            `json:"page"`
			PageSize int            `json:"pageSize"`
			Count    uint64         `json:"count"`
			Tokens   []models.Token `json:"tokens"`
		}
		res := resType{page, pageSize, 0, tokens}
		models.DB.Table("tokens").Count(&res.Count)

		c.JSON(200, res)
	}
}

// // Example select tokens with volume, change only using SQL instead of cache
// // this is little complex and may be slow when data is enough

// select t.address, t.name, t.symbol, t.decimals, t.price_usd, t.price_updated_at, t.volume, sum(trades.volume_usd) as volume_last
// from (
//   select t.address, t.name, t.symbol, t.decimals, t.price_usd, t.price_updated_at, sum(trades.volume_usd) as volume
//   from tokens as t, trades
//   where (trades.base_token_address = t.address or trades.quote_token_address = t.address) and trades.date >= '2019-02-25T00:00:00+08:00' and trades.date < '2019-02-26T00:00:00+08:00'
//   group by t.address, t.name, t.symbol, t.decimals, t.price_usd, t.price_updated_at
//   order by volume desc limit 25 offset 0
// ) as t, trades
// where (t.address = trades.base_token_address or t.address = trades.quote_token_address) and trades.date >= '2019-02-24T00:00:00+08:00' and trades.date < '2019-02-25T00:00:00+08:00'
// group by t.address, t.name, t.symbol, t.decimals, t.price_usd, t.price_updated_at, t.volume
// order by t.volume desc;

// const queryTokensColumns = "t.address, t.name, t.symbol, t.decimals, t.price_usd, t.price_updated_at"
// const queryTokensSQL = `select ` + queryTokensColumns + `, t.volume, sum(trades.volume_usd) as volume_last
// from (
//   select ` + queryTokensColumns + `, sum(trades.volume_usd) as volume
//   from tokens as t, trades
//   where (trades.base_token_address = t.address or trades.quote_token_address = t.address) and trades.date >= ? and trades.date < ?
//   group by ` + queryTokensColumns + `
//   order by volume desc limit ? offset ?
// ) as t, trades
// where (t.address = trades.base_token_address or t.address = trades.quote_token_address) and trades.date >= ? and trades.date < ?
// group by ` + queryTokensColumns + `, t.volume
// order by t.volume desc`

func GetToken(c *gin.Context) {
	address := c.Params.ByName("address")
	token := models.Token{}
	if err := models.DB.Where("address = ?", address).First(&token).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		tradesData := task.GetTrades24hData(token.Address)
		token.Trades24h = tradesData.Trades24h
		token.Trades24hChange = tradesData.Trades24hChange
		token.Traders24h = tradesData.Traders24h
		token.Traders24hChange = tradesData.Traders24hChange
		token.Amount24h = tradesData.Amount24h

		c.JSON(200, token)
	}
}

func GetTokenChart(c *gin.Context) {
	address := c.Params.ByName("address")
	token := models.Token{}
	if err := models.DB.Where("address = ?", address).First(&token).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	}

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
	models.DB.Raw(`SELECT date_trunc(?, date) AS dt, sum(volume_usd), count(*) AS trades_count
		FROM trades WHERE date >= ? AND (base_token_address = ? OR quote_token_address = ?)
		GROUP BY dt ORDER BY dt`, trunc, from, address, address).Scan(&res)

	var resTraders []struct {
		TradersCount uint64 `json:"traders"`
	}
	// select traders
	// SELECT dt, count(*) FROM (SELECT date_trunc('hour', date) AS dt, maker_address FROM trades WHERE date > '2019-02-26t00:00:00+08:00'UNION SELECT date_trunc('hour', date) AS dt, taker_address FROM trades WHERE date > '2019-02-26t00:00:00+08:00' ) AS traders GROUP BY dt ORDER BY dt;
	models.DB.Raw(`SELECT dt, count(*) AS traders_count
		FROM (
		SELECT date_trunc(?, date) AS dt, maker_address FROM trades WHERE date > ? AND (base_token_address = ? OR quote_token_address = ?)
		UNION
		SELECT date_trunc(?, date) AS dt, taker_address FROM trades WHERE date > ? AND (base_token_address = ? OR quote_token_address = ?)
		) AS traders GROUP BY dt ORDER BY dt`,
		trunc, from, address, address,
		trunc, from, address, address).Scan(&resTraders)
	for i, _ := range res {
		res[i].TradersCount = resTraders[i].TradersCount
	}
	c.JSON(200, res)
}
