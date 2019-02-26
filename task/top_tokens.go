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

	type QueryResult struct {
		Volume24h     decimal.Decimal
		Volume24hLast decimal.Decimal
	}

	for _, token := range tokens {
		result := QueryResult{}

		models.DB.Raw("SELECT sum(volume_usd) as volume24h FROM trades where date > ? and date <= ? and (base_token_address = ? or quote_token_address = ?)",
			time24hAgo, timeNow, token.Address, token.Address).Scan(&result)
		models.DB.Raw("SELECT sum(volume_usd) as volume24h_last FROM trades where date > ? and date <= ? and (base_token_address = ? or quote_token_address = ?)",
			time48hAgo, time24hAgo, token.Address, token.Address).Scan(&result)

		var change float32 = 0
		if !result.Volume24hLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := result.Volume24h.Sub(result.Volume24hLast).Div(result.Volume24hLast).Float64()
			change = float32(changeFloat64)
		}
		// Using map[string]interface{} instead of models.Token since
		// when update with struct, GORM will only update those fields that with non blank value
		// nothing will be updated as "", 0, false are blank values of their types
		models.DB.Model(&token).Updates(map[string]interface{}{"volume_24h": result.Volume24h, "volume_24h_change": change})
	}
}

func UpdateTokenVolume7d() {
	log.Info("UpdateTokenVolume7d")
	var tokens []models.Token
	models.DB.Find(&tokens)

	timeNow := time.Now()
	time7dAgo := time.Now().Add(-7 * 24 * time.Hour)
	time14dAgo := time.Now().Add(-14 * 24 * time.Hour)

	type QueryResult struct {
		Volume7d     decimal.Decimal
		Volume7dLast decimal.Decimal
	}

	for _, token := range tokens {
		result := QueryResult{}

		models.DB.Raw("SELECT sum(volume_usd) as volume7d FROM trades where date > ? and date <= ? and (base_token_address = ? or quote_token_address = ?)",
			time7dAgo, timeNow, token.Address, token.Address).Scan(&result)
		models.DB.Raw("SELECT sum(volume_usd) as volume7d_last FROM trades where date > ? and date <= ? and (base_token_address = ? or quote_token_address = ?)",
			time14dAgo, time7dAgo, token.Address, token.Address).Scan(&result)

		var change float32 = 0
		if !result.Volume7dLast.Equal(decimal.NewFromFloat32(0)) {
			changeFloat64, _ := result.Volume7d.Sub(result.Volume7dLast).Div(result.Volume7dLast).Float64()
			change = float32(changeFloat64)
		}
		models.DB.Model(&token).Updates(map[string]interface{}{"volume_7d": result.Volume7d, "volume_7d_change": change})
	}
}

func UpdateTokenVolumeAll() {
	log.Info("UpdateTokenVolume7d")
	var tokens []models.Token
	models.DB.Find(&tokens)

	type QueryResult struct {
		VolumeAll decimal.Decimal
	}

	for _, token := range tokens {
		result := QueryResult{}
		models.DB.Raw("SELECT sum(volume_usd) as volume_all FROM trades WHERE base_token_address = ? or quote_token_address = ?", token.Address, token.Address).Scan(&result)
		models.DB.Model(&token).Updates(map[string]interface{}{"volume_all": result.VolumeAll})
	}
}
