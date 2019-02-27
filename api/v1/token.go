package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/jinzhu/gorm"
)

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

func GetToken(c *gin.Context) {
	address := c.Params.ByName("address")
	token := models.Token{}
	if err := models.DB.Where("address = ?", address).First(&token).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, token)
	}
}
