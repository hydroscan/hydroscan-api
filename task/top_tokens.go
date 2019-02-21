package task

import (
	"time"

	"github.com/hydroscan/hydroscan-api/models"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func UpdateTokenVolume24h() {
	log.Info("UpdateTokenVolume24h")
	var tokens []models.Token
	models.DB.Find(&tokens)

	timeNow := time.Now()
	time24hAgo := time.Now().Add(-24 * time.Hour)
	time48hAgo := time.Now().Add(-48 * time.Hour)

	for _, token := range tokens {
		var volume24h decimal.Decimal
		var volume24hLast decimal.Decimal
		models.DB.Model(&models.Trade{}).Where("date > ? and date <= ? and (base_token_address = ? or quote_token_address = ?)", time24hAgo, timeNow, token.Address, token.Address).Select("sum(volume_usd)").Scan(&volume24h)
		models.DB.Model(&models.Trade{}).Where("date > ? and date <= ? and (base_token_address = ? or quote_token_address = ?)", time48hAgo, time24hAgo, token.Address, token.Address).Select("sum(volume_usd)").Scan(&volume24hLast)

		var change float32 = 0
		if !volume24hLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := volume24h.Sub(volume24hLast).Div(volume24hLast).Float64()
			change = float32(changeFloat64)
		}
		models.DB.Model(&token).Updates(models.Token{Volume24h: volume24h, Volume24hChange: change})
	}
}

func UpdateTokenVolume7d() {
	log.Info("UpdateTokenVolume7d")
	var tokens []models.Token
	models.DB.Find(&tokens)

	timeNow := time.Now()
	time7dAgo := time.Now().Add(-7 * 24 * time.Hour)
	time14dAgo := time.Now().Add(-14 * 24 * time.Hour)

	for _, token := range tokens {
		var volume7d decimal.Decimal
		var volume7dLast decimal.Decimal
		models.DB.Model(&models.Trade{}).Where("date > ? and date <= ? and (base_token_address = ? or quote_token_address = ?)", time7dAgo, timeNow, token.Address, token.Address).Select("sum(volume_usd)").Scan(&volume7d)
		models.DB.Model(&models.Trade{}).Where("date > ? and date <= ? and (base_token_address = ? or quote_token_address = ?)", time14dAgo, time7dAgo, token.Address, token.Address).Select("sum(volume_usd)").Scan(&volume7dLast)

		var change float32 = 0
		if volume7dLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := volume7d.Sub(volume7dLast).Div(volume7dLast).Float64()
			change = float32(changeFloat64)
		}
		models.DB.Model(&token).Updates(models.Token{Volume7d: volume7d, Volume7dChange: change})
	}
}

func UpdateTokenVolumeAll() {
	log.Info("UpdateTokenVolume7d")
	var tokens []models.Token

	models.DB.Find(&tokens)
	for _, token := range tokens {
		var volumeAll decimal.Decimal
		models.DB.Model(&models.Trade{}).Where("base_token_address = ? or quote_token_address = ?", token.Address, token.Address).Select("sum(volume_usd)").Scan(&volumeAll)
		models.DB.Model(&token).Updates(models.Token{VolumeAll: volumeAll})
	}
}
