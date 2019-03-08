package apiv1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/utils"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func GetRelayers(c *gin.Context) {
	var relayers []models.Relayer
	if err := models.DB.Find(&relayers).Error; gorm.IsRecordNotFoundError(err) {
		c.JSON(404, relayers)
	} else {
		c.JSON(200, relayers)
	}
}

func GetRelayer(c *gin.Context) {
	slug := c.Params.ByName("slug")
	relayer := models.Relayer{}

	statment := models.DB.Where("slug = ?", slug)
	if utils.IsAddress(slug) {
		if isTrue, key := utils.IsRelayer(slug); isTrue {
			statment = models.DB.Where("address = ?", key)
		} else {
			c.JSON(404, relayer)
		}
	}

	if err := statment.First(&relayer).Error; gorm.IsRecordNotFoundError(err) {
		c.JSON(404, relayer)
	} else {
		address := relayer.Address
		type TopToken struct {
			Address      string          `json:"address"`
			Name         string          `json:"name"`
			Symbol       string          `json:"symbol"`
			Volume       decimal.Decimal `json:"volume"`
			VolumeLast   decimal.Decimal `json:"volumeLast"`
			VolumeChange float32         `json:"volumeChange"`
		}

		var res struct {
			Name             string          `json:"name"`
			Url              string          `json:"url"`
			Slug             string          `json:"slug"`
			Address          string          `json:"address"`
			Volume24h        decimal.Decimal `json:"volume24h"`
			Volume24hLast    decimal.Decimal `json:"volume24hLast"`
			Volume24hChange  float32         `json:"volume24hChange"`
			Trades24h        decimal.Decimal `json:"trades24h"`
			Trades24hLast    decimal.Decimal `json:"trades24hLast"`
			Trades24hChange  float32         `json:"trades24hChange"`
			Traders24h       decimal.Decimal `json:"traders24h"`
			Traders24hLast   decimal.Decimal `json:"traders24hLast"`
			Traders24hChange float32         `json:"traders24hChange"`
			TopTokens        []TopToken      `json:"topTokens"`
		}

		timeNow := time.Now()
		time24hAgo := time.Now().Add(-24 * time.Hour)
		time48hAgo := time.Now().Add(-48 * time.Hour)

		models.DB.Raw(`SELECT sum(trades.volume_usd) AS volume24h, count(*) AS trades24h
		FROM trades WHERE relayer_address = ? AND date >= ? AND date < ?`,
			address, time24hAgo, timeNow).Scan(&res)

		models.DB.Raw(`SELECT sum(trades.volume_usd) AS volume24h_last, count(*) AS trades24h_last
		FROM trades WHERE relayer_address = ? AND date >= ? AND date < ?`,
			address, time48hAgo, time24hAgo).Scan(&res)

		models.DB.Raw(`SELECT count(*) AS traders24h
			FROM (
			SELECT maker_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
     		UNION
     		SELECT taker_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
     		) AS traders`,
			time24hAgo, timeNow, address,
			time24hAgo, timeNow, address).Scan(&res)
		models.DB.Raw(`SELECT count(*) AS traders24h_last
			FROM (
			SELECT maker_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
     		UNION
     		SELECT taker_address FROM trades WHERE date > ? AND date <= ? AND relayer_address = ?
     		) AS traders`,
			time48hAgo, time24hAgo, address,
			time48hAgo, time24hAgo, address).Scan(&res)

		models.DB.Raw(`SELECT sum(maker_rebate) FROM trades WHERE relayer_address = ?`,
			address).Scan(&res)

		var topTokens []TopToken
		const baseFields = "t.address, t.name, t.symbol"
		models.DB.Raw(`SELECT `+baseFields+`, t.volume, sum(trades.volume_usd) AS volume_last
		FROM (
			SELECT `+baseFields+`, sum(trades.volume_usd) AS volume FROM tokens AS t, trades
			WHERE (trades.base_token_address = t.address OR trades.quote_token_address = t.address)
			AND relayer_address = ?
			AND trades.date >= ? AND trades.date < ?
			GROUP BY `+baseFields+`
			ORDER BY volume DESC LIMIT 3 OFFSET 0
		) AS t LEFT JOIN trades
		ON (t.address = trades.base_token_address OR t.address = trades.quote_token_address)
		AND relayer_address = ?
		AND trades.date >= ? AND trades.date < ?
		GROUP BY `+baseFields+`, t.volume
		ORDER BY t.volume DESC`,
			address, time24hAgo, timeNow,
			address, time48hAgo, time24hAgo).Scan(&topTokens)

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

		if !res.Traders24hLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := res.Traders24h.Sub(res.Traders24hLast).Div(res.Traders24hLast).Float64()
			res.Traders24hChange = float32(changeFloat64)
		}

		res.Name = relayer.Name
		res.Url = relayer.Url
		res.Slug = relayer.Slug
		res.Address = relayer.Address

		c.JSON(200, res)
	}
}
