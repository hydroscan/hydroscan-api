package apiv1

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type TokensQuery struct {
	Page           int    `form:"page"`
	PageSize       int    `form:"pageSize"`
	Filter         string `form:"filter"`
	Keyword        string `form:"keyword"`
	TraderAddress  string `form:"traderAddress"`
	RelayerAddress string `form:"relayerAddress"`
}

type ResToken struct {
	Address         string          `json:"address"`
	Name            string          `json:"name"`
	Symbol          string          `json:"symbol"`
	Volume24h       decimal.Decimal `gorm:"column:volume_24h" json:"volume24h"`
	Amount24h       decimal.Decimal `json:"amount24h"`
	Volume24hLast   decimal.Decimal `gorm:"column:volume_24h_last" json:"volume24hLast"`
	Volume24hChange float32         `gorm:"column:volume_24h_change" json:"volume24hChange"`
	Traders24h      uint64          `json:"traders24h"`
	PriceUSD        decimal.Decimal `gorm:"column:price_usd" json:"priceUSD"`
	PriceUpdatedAt  time.Time       `gorm:"column:price_updated_at" json:"priceUpdatedAt"`
}

type ResType struct {
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Count    uint64     `json:"count"`
	Tokens   []ResToken `json:"tokens"`
}

func GetTokens(c *gin.Context) {
	query := TokensQuery{1, 25, "24H", "", "", ""}
	c.BindQuery(&query)
	log.Info("query ", query)

	page := query.Page
	pageSize := query.PageSize
	offset := (page - 1) * pageSize

	var res ResType
	if query.Keyword != "" {
		res = getTokensByKeyword(page, pageSize, offset, query.Keyword)
	} else if query.RelayerAddress != "" {
		res = getTokensByRelayer(page, pageSize, offset, query.RelayerAddress)
	} else if query.TraderAddress != "" {
		res = getTokensByTrader(page, pageSize, offset, query.TraderAddress)
	} else {
		res = getTokensDefault(page, pageSize, offset, query.Filter)
	}

	c.JSON(200, res)

}

func getTokensDefault(page int, pageSize int, offset int, filter string) ResType {
	var tokens []ResToken
	res := ResType{page, pageSize, 0, tokens}

	orderField := "volume_24h"
	switch filter {
	case "24H":
		orderField = "volume_24h"
	case "7D":
		orderField = "volume_7d"
	case "1M":
		orderField = "volume_1m"
	case "ALL":
		orderField = "volume_all"
	default:
		orderField = "volume_24h"
	}
	models.DB.Raw("SELECT * FROM tokens ORDER BY "+orderField+" DESC LIMIT ? OFFSET ?", pageSize, offset).Scan(&tokens)

	for i, _ := range tokens {
		tradesData := task.GetTrades24hData(tokens[i].Address)
		tokens[i].Traders24h = tradesData.Traders24h
		tokens[i].Amount24h = tradesData.Amount24h
	}
	models.DB.Raw("SELECT COUNT(*) FROM tokens").Scan(&res)

	res.Tokens = tokens
	return res
}

func getTokensByKeyword(page int, pageSize int, offset int, keyword string) ResType {
	var tokens []ResToken
	res := ResType{page, pageSize, 0, tokens}

	keyword = strings.TrimSpace(keyword)
	models.DB.Raw("SELECT * FROM tokens WHERE name ILIKE ? OR symbol ILIKE ? ORDER BY volume_24h DESC LIMIT ? OFFSET ?",
		"%"+keyword+"%", "%"+keyword+"%", pageSize, offset).Scan(&tokens)

	for i, _ := range tokens {
		tradesData := task.GetTrades24hData(tokens[i].Address)
		tokens[i].Traders24h = tradesData.Traders24h
		tokens[i].Amount24h = tradesData.Amount24h
	}

	models.DB.Raw("SELECT COUNT(*) FROM tokens WHERE name ILIKE ? OR symbol ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").Scan(&res)

	res.Tokens = tokens
	return res
}

func getTokensByTrader(page int, pageSize int, offset int, traderAddress string) ResType {
	var tokens []ResToken
	res := ResType{page, pageSize, 0, tokens}

	timeNow := time.Now()
	time24hAgo := time.Now().Add(-24 * time.Hour)
	time48hAgo := time.Now().Add(-48 * time.Hour)

	const baseFields = "t.address, t.name, t.symbol, t.price_usd, t.price_updated_at"
	models.DB.Raw(`SELECT `+baseFields+`, t.volume AS volume_24h, sum(trades.volume_usd) AS volume_24h_last
		FROM (
			SELECT `+baseFields+`, sum(trades.volume_usd) AS volume FROM tokens AS t, trades
			WHERE (trades.base_token_address = t.address OR trades.quote_token_address = t.address)
			AND (trades.maker_address = ? OR trades.taker_address = ?)
			AND trades.date >= ? AND trades.date < ?
			GROUP BY `+baseFields+`
			ORDER BY volume DESC LIMIT ? OFFSET ?
		) AS t LEFT JOIN trades
		ON (t.address = trades.base_token_address OR t.address = trades.quote_token_address)
		AND (trades.maker_address = ? OR trades.taker_address = ?)
		AND trades.date >= ? AND trades.date < ?
		GROUP BY `+baseFields+`, t.volume
		ORDER BY t.volume DESC`,
		traderAddress, traderAddress, time24hAgo, timeNow, pageSize, offset,
		traderAddress, traderAddress, time48hAgo, time24hAgo).Scan(&tokens)

	type QueryResult struct {
		AsBaseTokenAmount24h  decimal.Decimal
		AsQuoteTokenAmount24h decimal.Decimal
	}

	for i, token := range tokens {
		result := QueryResult{}

		models.DB.Raw(`SELECT sum(base_token_amount) AS as_base_token_amount24h FROM trades
			WHERE date > ? AND date <= ?
			AND base_token_address = ?
			AND (maker_address = ? OR taker_address = ?)`,
			time24hAgo, timeNow, token.Address, traderAddress, traderAddress).Scan(&result)
		models.DB.Raw(`SELECT sum(quote_token_amount) AS as_quote_token_amount24h FROM trades
			WHERE date > ? AND date <= ?
			AND quote_token_address = ?
			AND (maker_address = ? OR taker_address = ?)`,
			time24hAgo, timeNow, token.Address, traderAddress, traderAddress).Scan(&result)

		tokens[i].Amount24h = result.AsBaseTokenAmount24h.Add(result.AsQuoteTokenAmount24h)

		if !token.Volume24hLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := token.Volume24h.Sub(token.Volume24hLast).Div(token.Volume24hLast).Float64()
			tokens[i].Volume24hChange = float32(changeFloat64)
		}
	}

	models.DB.Raw(`SELECT count(*)
		FROM (
		SELECT base_token_address FROM trades WHERE date > ? AND date <= ? AND (maker_address = ? OR taker_address = ?)
 		UNION
 		SELECT quote_token_address FROM trades WHERE date > ? AND date <= ? AND (maker_address = ? OR taker_address = ?)
 		) AS traders`,
		time24hAgo, timeNow, traderAddress, traderAddress,
		time24hAgo, timeNow, traderAddress, traderAddress).Scan(&res)

	res.Tokens = tokens
	return res
}

func getTokensByRelayer(page int, pageSize int, offset int, relayerAddress string) ResType {
	var tokens []ResToken
	res := ResType{page, pageSize, 0, tokens}

	timeNow := time.Now()
	time24hAgo := time.Now().Add(-24 * time.Hour)
	time48hAgo := time.Now().Add(-48 * time.Hour)

	const baseFields = "t.address, t.name, t.symbol, t.price_usd, t.price_updated_at"
	models.DB.Raw(`SELECT `+baseFields+`, t.volume AS volume_24h, sum(trades.volume_usd) AS volume_24h_last
		FROM (
			SELECT `+baseFields+`, sum(trades.volume_usd) AS volume FROM tokens AS t, trades
			WHERE (trades.base_token_address = t.address OR trades.quote_token_address = t.address)
			AND trades.relayer_address = ?
			AND trades.date >= ? AND trades.date < ?
			GROUP BY `+baseFields+`
			ORDER BY volume DESC LIMIT ? OFFSET ?
		) AS t LEFT JOIN trades
		ON (t.address = trades.base_token_address OR t.address = trades.quote_token_address)
		AND trades.relayer_address = ?
		AND trades.date >= ? AND trades.date < ?
		GROUP BY `+baseFields+`, t.volume
		ORDER BY t.volume DESC`,
		relayerAddress, time24hAgo, timeNow, pageSize, offset,
		relayerAddress, time48hAgo, time24hAgo).Scan(&tokens)

	type QueryResult struct {
		Traders24h            uint64
		AsBaseTokenAmount24h  decimal.Decimal
		AsQuoteTokenAmount24h decimal.Decimal
	}

	for i, token := range tokens {
		result := QueryResult{}

		models.DB.Raw(`SELECT sum(base_token_amount) AS as_base_token_amount24h FROM trades
			WHERE date > ? AND date <= ?
			AND base_token_address = ?
			AND relayer_address = ?`,
			time24hAgo, timeNow, token.Address, relayerAddress).Scan(&result)
		models.DB.Raw(`SELECT sum(quote_token_amount) AS as_quote_token_amount24h FROM trades
			WHERE date > ? AND date <= ?
			AND quote_token_address = ?
			AND relayer_address = ?`,
			time24hAgo, timeNow, token.Address, relayerAddress).Scan(&result)

		models.DB.Raw(`SELECT count(*) AS traders24h
			FROM (
				SELECT maker_address FROM trades WHERE date > ? AND date <= ?
					AND (base_token_address = ? OR quote_token_address = ?)
					AND relayer_address = ?
     			UNION
     			SELECT taker_address FROM trades WHERE date > ? AND date <= ?
     				AND (base_token_address = ? OR quote_token_address = ?)
     				AND relayer_address = ?
     		) AS traders`,
			time24hAgo, timeNow, token.Address, token.Address, relayerAddress,
			time24hAgo, timeNow, token.Address, token.Address, relayerAddress).Scan(&result)

		tokens[i].Amount24h = result.AsBaseTokenAmount24h.Add(result.AsQuoteTokenAmount24h)
		tokens[i].Traders24h = result.Traders24h

		if !token.Volume24hLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := token.Volume24h.Sub(token.Volume24hLast).Div(token.Volume24hLast).Float64()
			tokens[i].Volume24hChange = float32(changeFloat64)
		}
	}

	models.DB.Raw(`SELECT count(*)
		FROM (
		SELECT base_token_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
 		UNION
 		SELECT quote_token_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
 		) AS traders`,
		time24hAgo, timeNow, relayerAddress,
		time24hAgo, timeNow, relayerAddress).Scan(&res)

	res.Tokens = tokens
	return res
}

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
