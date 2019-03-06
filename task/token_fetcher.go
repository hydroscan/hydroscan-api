package task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TokenPrice struct {
	Rate decimal.Decimal `json:"rate"`
	TS   int64           `json:"ts"`
}

type TokenInfo struct {
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	Decimals          string
	DecimalsInterface interface{} `json:"decimals"` // ethplorer api return decimals type can be string or number
	TotalSupply       string      `json:"totalSupply"`
	HoldersCount      uint64      `json:"holdersCount"`
	Price             TokenPrice  `json:"price"`
}

func GetToken(address string) models.Token {
	mToken := models.Token{}
	if err := models.DB.Where("address = ?", address).First(&mToken).Error; gorm.IsRecordNotFoundError(err) {
		tokenInfo := GetTokenInfo(address)

		decimals, err := strconv.ParseUint(tokenInfo.Decimals, 10, 64)
		if err != nil {
			panic(err)
		}

		mToken = models.Token{
			Address:        address,
			Decimals:       uint(decimals),
			Name:           tokenInfo.Name,
			Symbol:         tokenInfo.Symbol,
			TotalSupply:    tokenInfo.TotalSupply,
			HoldersCount:   tokenInfo.HoldersCount,
			PriceUSD:       tokenInfo.Price.Rate,
			PriceUpdatedAt: time.Unix(tokenInfo.Price.TS, 0),
		}
		models.DB.Create(&mToken)
	}
	return mToken
}

func UpdateTokenPrices() {
	log.Info("UpdateTokenPrices")

	mTokens := []models.Token{}
	models.DB.Find(&mTokens)

	for _, mToken := range mTokens {
		tokenInfo := GetTokenInfo(mToken.Address)
		models.DB.Model(&mToken).Updates(models.Token{
			Name:           tokenInfo.Name,
			Symbol:         tokenInfo.Symbol,
			TotalSupply:    tokenInfo.TotalSupply,
			HoldersCount:   tokenInfo.HoldersCount,
			PriceUSD:       tokenInfo.Price.Rate,
			PriceUpdatedAt: time.Unix(tokenInfo.Price.TS, 0),
		})
	}
}

func GetTokenInfo(address string) TokenInfo {
	log.Info("GetTokenInfo ", address)
	url := "http://api.ethplorer.io/getTokenInfo/" + address + "?apiKey=" + viper.GetString("ethplorer_apikey")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	tokenInfo := TokenInfo{}
	json.Unmarshal([]byte(body), &tokenInfo)
	// ethplorer api return decimals type can be string or number
	if tokenInfo.Decimals == "" {
		tokenInfo.Decimals = fmt.Sprint(tokenInfo.DecimalsInterface)
	}
	return tokenInfo
}
